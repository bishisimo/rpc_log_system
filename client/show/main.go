/*
@author '彼时思默'
@time 2020/5/13 上午10:26
@describe:
*/
package main

import (
	"context"
	"fmt"
	"github.com/bishisimo/rpc_log_system/src/redux"
	"github.com/bishisimo/rpc_log_system/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(utils.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := redux.NewReduxClient(conn)
	res, err := c.ShowInfo(context.Background(), &redux.BlankRequest{})
	if err != nil {
		logrus.Error("show error:", err)
	}
	fmt.Printf("info:%+v", string(res.Info))
}
