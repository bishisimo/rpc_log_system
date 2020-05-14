/*
@author '彼时思默'
@time 2020/5/12 下午4:37
@describe:
*/
package entity

import (
	"context"
	"github.com/bishisimo/rpc_log_system/src/redux"
)

type Pub struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	MsgChan
	*BaseInfo
	*redux.AddPubRequest
}

func NewPub(pr *redux.AddPubRequest) *Pub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pub{
		Ctx:           ctx,
		Cancel:        cancel,
		MsgChan:       MsgChanPool.Get(),
		BaseInfo:      NewBaseInfo(),
		AddPubRequest: pr,
	}
}
