/*
@author '彼时思默'
@time 2020/5/12 下午4:22
@describe:
*/
package local

import (
	"encoding/json"
	"github.com/bishisimo/rpc_log_system/entity"
	"github.com/bishisimo/rpc_log_system/utils"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

type (
	DayKey  = string //以日期字符串作为键
	PathKey = string //以路径字符串作为键
)
type LocalSub struct {
	minFormat  string                          //时间格式化
	hourFormat string                          //时间格式化
	dayFormat  string                          //时间格式化
	monFormat  string                          //时间格式化
	yearFormat string                          //时间格式化
	Dir        string                          //存放数据的文件夹
	MsgStrChan chan string                     //接受消息的chan
	Fps        map[PathKey]*os.File            //记录所有当前有效path
	FpsDaily   map[DayKey]map[PathKey]*os.File //将有效path按日期分组
}

func NewLocalSub(dir string) *LocalSub {
	utils.MakeDir(dir)
	return &LocalSub{
		minFormat:  "2006-01-02_15-04",
		hourFormat: "2006-01-02_15",
		dayFormat:  "2006-01-02",
		monFormat:  "2006-01",
		yearFormat: "2006",
		Dir:        dir,
		MsgStrChan: entity.MsgStrChanPool.Get(),
		Fps:        make(map[PathKey]*os.File),
		FpsDaily:   make(map[DayKey]map[PathKey]*os.File),
	}
}

//添加一个句柄记录
func (s *LocalSub) AddFp(filePath string, fp *os.File) {
	today := time.Now().Format(s.dayFormat)
	if s.FpsDaily[today] == nil {
		s.FpsDaily[today] = make(map[PathKey]*os.File)
	}
	s.Fps[filePath] = fp
	s.FpsDaily[today][filePath] = fp
}

//获取一个句柄
func (s *LocalSub) GetOneFp(filePath string) *os.File {
	return s.Fps[filePath]
}

//获取指定日期的所有文件句柄
func (s *LocalSub) GetOneTargetDayFp(targetDay string) map[string]*os.File {
	return s.FpsDaily[targetDay]
}

//获取当日所有句柄
func (s *LocalSub) GetTodayAllFp() map[string]*os.File {
	today := time.Now().Format(s.dayFormat)
	return s.GetOneTargetDayFp(today)
}

//获取昨日所有句柄
func (s *LocalSub) GetYesterdayAllFp() map[string]*os.File {
	d, _ := time.ParseDuration("-24h")
	yesterday := time.Now().Add(d).Format(s.dayFormat)
	return s.GetOneTargetDayFp(yesterday)
}

//清理一个路径的句柄
func (s *LocalSub) DropOneFp(filePath string) {
	today := time.Now().Format(s.dayFormat)
	fp := s.GetOneFp(filePath)
	if err := fp.Close(); err != nil {
		logrus.Errorf("close file %s error: %s", filePath, err)
	}
	delete(s.Fps, filePath)
	delete(s.FpsDaily[today], filePath)
}

//监听发布过来的消息,将消息写入文件
func (s *LocalSub) Listen() {
	go func() {
		for msgStr := range s.MsgStrChan {
			var msg entity.Msg
			err := json.Unmarshal([]byte(msgStr), &msg)
			if err != nil {
				logrus.Error("json unmarshal error:", err)
			}
			//生成文件路径
			today := time.Now().Format(s.dayFormat)
			softName := "default"
			if msg["soft_name"] != nil {
				softName = msg["soft_name"].(string)
			}
			atomId := "default"
			if msg["atom_id"] != nil {
				atomId = msg["atom_id"].(string)
			}
			dir := path.Join(s.Dir, msg["@id"].(string), softName, today)
			descPath := path.Join(dir, atomId)
			var fp *os.File
			//获取路径句柄或创建并记录
			if fp = s.GetOneFp(descPath); fp == nil {
				utils.MakeDir(dir)
				fp, err = os.OpenFile(descPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
				if err != nil {
					logrus.Errorf("open file %s failure", descPath)
				}
				s.AddFp(descPath, fp)
			}
			//写入句柄
			if fp != nil {
				_, err = fp.WriteString(msgStr + "\n")
				if err != nil {
					logrus.Errorf("write file %s failure", descPath)
				}
			}
		}
	}()
}
