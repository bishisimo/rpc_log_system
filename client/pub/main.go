/*
@author '彼时思默'
@time 2020/5/13 上午9:34
@describe:
*/
package main

import (
	"context"
	"github.com/bishisimo/rpc_log_system/src/redux"
	"github.com/bishisimo/rpc_log_system/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(utils.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := redux.NewReduxClient(conn)

	_, err = c.AddPub(context.Background(), &redux.AddPubRequest{
		Id: "pub1",
	})
	if err != nil {
		logrus.Error("Add pub error:", err)
	}
	req, err := c.PubMsg(context.Background())
	if err != nil {
		logrus.Error("pub msg error:", err)
	}
	go func() {
		for true {
			err := req.Send(&redux.PubMessage{
				Id:  "pub1",
				Msg: `{"time":"` + time.Now().Format("2006-01-02 15:04:05") + `","data":"我暴饮暴食!"}`,
			})
			if err != nil {
				return
			}
			time.Sleep(time.Second)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
