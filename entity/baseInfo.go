/*
@author '彼时思默'
@time 2020/5/13 上午10:38
@describe:
*/
package entity

import "time"

type BaseInfo struct {
	CreateTime string
	MsgCount   uint64
}

func NewBaseInfo() *BaseInfo {
	return &BaseInfo{
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		MsgCount:   0,
	}
}
