/*
@author '彼时思默'
@time 2020/5/13 上午9:03
@describe:
*/
package reply

import "github.com/bishisimo/rpc_log_system/src/redux"

var DefaultReply = &redux.Reply{Status: "OK"}
var NotExistReply = &redux.Reply{Status: "do not exist"}
var MissIdReply = &redux.Reply{Status: "do not exist Id"}
