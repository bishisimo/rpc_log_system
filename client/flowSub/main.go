/*
@author '彼时思默'
@time 2020/5/13 上午9:37
@describe:
*/
package main

import (
	"context"
	"fmt"
	"github.com/bishisimo/rpc_log_system/src/redux"
	"github.com/bishisimo/rpc_log_system/utils"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(utils.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := redux.NewReduxClient(conn)
	res, err := c.AddSub(context.Background(), &redux.AddSubRequest{
		Id:         "std1",
		AcceptType: "std",
		IsFlow:     true,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	go func() {
		for {
			msg, err := res.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}
			fmt.Println("std:", msg.Body)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
