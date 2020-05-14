/*
@author '彼时思默'
@time 2020/5/11 上午10:22
@describe:订阅者
*/
package entity

import (
	"context"
	"github.com/bishisimo/rpc_log_system/src/redux"
)

//订阅者
type Sub struct {
	*redux.AddSubRequest
	Ctx        context.Context
	Cancel     context.CancelFunc
	MsgStrChan chan string
	*BaseInfo
}

func NewSub(sr *redux.AddSubRequest) *Sub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Sub{
		AddSubRequest: sr,
		MsgStrChan:    MsgStrChanPool.Get(),
		Ctx:           ctx,
		Cancel:        cancel,
		BaseInfo:      NewBaseInfo(),
	}
}
