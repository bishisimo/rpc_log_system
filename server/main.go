/*
@author '彼时思默'
@time 2020/5/12 上午10:11
@describe:
*/
package main

import (
	"encoding/json"
	"fmt"
	. "github.com/bishisimo/rpc_log_system/entity"
	. "github.com/bishisimo/rpc_log_system/server/local"
	. "github.com/bishisimo/rpc_log_system/server/reply"
	"github.com/bishisimo/rpc_log_system/src/redux"
	"github.com/bishisimo/rpc_log_system/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

var Server *server

func init() {
	Server = NewServer()
}
func NewServer() *server {
	local := NewLocalSub("local_data")
	local.Listen()
	return &server{
		Local:     local,
		Pubs:      new(PubMap),
		SubsFlow:  new(SubMap),
		SubsBatch: new(SubMap),
	}
}

type server struct {
	Local     *LocalSub //map[string]*LocalSub 本地储存订阅者通道
	Pubs      *PubMap   //map[string]*Pub 发布者通道,消息体为map映射的json
	SubsFlow  *SubMap   //map[string]*Sub 流数据订阅者通道
	SubsBatch *SubMap   //map[string]*Sub 批量数据订阅者通道
}

func (s *server) AddPub(ctx context.Context, request *redux.AddPubRequest) (*redux.Reply, error) {
	if request.Id == "" {
		return MissIdReply, nil
	}
	pub := NewPub(request)
	s.Pubs.Store(request.Id, pub)
	go s.flowDistribution(request.Id)
	return DefaultReply, nil
}

func (s *server) PubMsg(pubServer redux.Redux_PubMsgServer) error {
	for {
		pubMessage, err := pubServer.Recv()
		if err == io.EOF {
			return pubServer.SendAndClose(DefaultReply)
		}
		if err != nil {
			return err
		}
		var msg Msg
		err = json.Unmarshal([]byte(pubMessage.Msg), &msg)
		if err != nil {
			logrus.Error("json error:", err)
			return err
		}
		if pub := s.Pubs.Load(pubMessage.Id); pub != nil {
			pub.MsgCount += 1
			pub.MsgChan <- msg
		} else {
			return pubServer.SendAndClose(MissIdReply)
		}
	}
}

func (s *server) RemovePub(ctx context.Context, request *redux.RemovePubRequest) (*redux.Reply, error) {
	if pub := s.Pubs.Load(request.Id); pub != nil {
		pub.Cancel()
		MsgChanPool.Put(pub.MsgChan)
		s.Pubs.Delete(request.Id)
		return DefaultReply, nil
	}
	return NotExistReply, nil
}

func (s *server) AddSub(request *redux.AddSubRequest, subServer redux.Redux_AddSubServer) error {
	sub := NewSub(request)
	if sub.IsFlow {
		s.SubsFlow.Store(sub.Id, sub)
		for {
			select {
			case <-sub.Ctx.Done():
				return nil
			case msgStr := <-sub.MsgStrChan:
				sub.MsgCount += 1
				if err := subServer.Send(&redux.FlowSubReply{
					Body: msgStr,
				}); err != nil {
					return err
				}
			}
		}
	} else {
		s.SubsBatch.Store(sub.Id, sub)
		return nil
	}
}

func (s *server) GetTargetPath(ctx context.Context, request *redux.GetTargetPathRequest) (*redux.TargetPathReply, error) {
	paths := s.Local.GetYesterdayAllFp()
	res := make([]string, 0, len(paths))
	for k := range paths {
		res = append(res, k)
	}
	return &redux.TargetPathReply{
		Paths: res,
	}, nil
}

func (s *server) RemoveSub(ctx context.Context, request *redux.RemoveSubRequest) (*redux.Reply, error) {
	if sub := s.SubsFlow.Load(request.Id); sub != nil {
		sub.Cancel()
		MsgStrChanPool.Put(sub.MsgStrChan)
		s.SubsFlow.Delete(request.Id)
		return DefaultReply, nil
	} else if sub := s.SubsBatch.Load(request.Id); sub != nil {
		sub.Cancel()
		MsgStrChanPool.Put(sub.MsgStrChan)
		s.SubsBatch.Delete(request.Id)
		return DefaultReply, nil
	}
	return NotExistReply, nil
}

//将消息分发到流订阅者
func (s *server) flowDistribution(pubId string) {
	pub := s.Pubs.Load(pubId)
	for {
		select {
		case <-pub.Ctx.Done():
			return
		case msg := <-pub.MsgChan:
			msg["@id"] = pubId
			msgByte, err := json.Marshal(&msg)
			if err != nil {
				logrus.Errorf(err.Error())
			}
			msgStr := string(msgByte)
			s.Local.MsgStrChan <- msgStr
			s.SubsFlow.Range(func(key string, sub *Sub) {
				subChan := sub.MsgStrChan
				subChan <- msgStr
			})
		}
	}
}

func (s *server) ShowInfo(ctx context.Context, request *redux.BlankRequest) (*redux.InfoReply, error) {
	info := make(map[string]map[string]map[string]string)
	pubInfo := make(map[string]map[string]string)
	flowInfo := make(map[string]map[string]string)
	batchInfo := make(map[string]map[string]string)
	s.Pubs.Range(func(id string, pub *Pub) {
		tmpPubInfo := make(map[string]string)
		tmpPubInfo["Id"] = pub.Id
		tmpPubInfo["CreateTime"] = pub.CreateTime
		tmpPubInfo["MsgCount"] = fmt.Sprintf("%d", pub.MsgCount)
		pubInfo[id] = tmpPubInfo
	})

	s.SubsFlow.Range(func(id string, sub *Sub) {
		tmpFlowInfo := make(map[string]string)
		tmpFlowInfo["Id"] = sub.Id
		tmpFlowInfo["AcceptType"] = sub.AcceptType
		tmpFlowInfo["CreateTime"] = sub.CreateTime
		tmpFlowInfo["MsgCount"] = fmt.Sprintf("%d", sub.MsgCount)
		flowInfo[id] = tmpFlowInfo
	})
	s.SubsBatch.Range(func(id string, sub *Sub) {
		tmpBatchInfo := make(map[string]string)
		tmpBatchInfo["Id"] = sub.Id
		tmpBatchInfo["AcceptType"] = sub.AcceptType
		tmpBatchInfo["MsgCount"] = fmt.Sprintf("%d", sub.MsgCount)
		batchInfo[id] = tmpBatchInfo
	})
	info["pubInfo"] = pubInfo
	info["flowInfo"] = flowInfo
	info["batchInfo"] = batchInfo
	result, err := json.Marshal(&info)
	if err != nil {
		logrus.Error("json error:", err)
	}
	return &redux.InfoReply{
		Info: result,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", utils.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	redux.RegisterReduxServer(s, Server)
	_ = s.Serve(lis)
}
