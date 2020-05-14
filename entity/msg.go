/*
@author '彼时思默'
@time 2020/5/12 下午5:27
@describe:
*/
package entity

var ChanSize = 10240

type Msg = map[string]interface{}
type MsgChan = chan Msg
type MsgStrChan = chan string
