/*
@author '彼时思默'
@time 2020/5/11 下午4:07
@describe:
*/
package entity

import (
	"sync"
)

//var MsgChanPoolStruct=sync.Pool{New: func() interface{}{return make(MsgChan, ChanSize)}}
//var MsgStrChanPoolStruct=sync.Pool{New: func() interface{}{return make(MsgStrChan, ChanSize)}}

var MsgChanPool *MsgChanPoolStruct
var MsgStrChanPool *MsgStrChanPoolStruct

func init() {
	MsgChanPool = NewMsgChanPool()
	MsgStrChanPool = NewMsgStrChanPool()
}

func NewMsgChanPool() *MsgChanPoolStruct {
	return &MsgChanPoolStruct{
		Body: sync.Pool{New: func() interface{} { return make(MsgChan, ChanSize) }},
	}
}

func NewMsgStrChanPool() *MsgStrChanPoolStruct {
	return &MsgStrChanPoolStruct{
		Body: sync.Pool{New: func() interface{} { return make(MsgStrChan, ChanSize) }},
	}
}

type MsgChanPoolStruct struct {
	Body sync.Pool
}

func (p *MsgChanPoolStruct) Put(x interface{}) {
	p.Body.Put(x)
}
func (p *MsgChanPoolStruct) Get() MsgChan {
	return p.Body.Get().(MsgChan)
}

type MsgStrChanPoolStruct struct {
	Body sync.Pool
}

func (p *MsgStrChanPoolStruct) Put(x MsgStrChan) {
	p.Body.Put(x)
}
func (p *MsgStrChanPoolStruct) Get() MsgStrChan {
	return p.Body.Get().(MsgStrChan)
}
