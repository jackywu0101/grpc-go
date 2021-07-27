/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/status"
)

const (
	address     = "localhost:50051"
	defaultName = "japan"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	names := []string{defaultName}
	if len(os.Args) > 1 {
		names = os.Args[1:]
	}
	for _, name := range names {
		SayHello(c, name)
	}
}

// SayHello wraps grpc call SayHello
func SayHello(c pb.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Printf("could not greet: %v", err)
		if st, ok := status.FromError(err); ok {
			log.Printf("Got error status: %s", st.Message())
			for _, detail := range st.Details() {
				if mvmErr, ok := detail.(*pb.MvmError); ok {
					log.Printf("Got MvmErr code %d, message %s", mvmErr.LibmvmError, mvmErr.Msg)
				}
			}
		}
	} else {
		log.Printf("Greeting: %s", r.GetMessage())
	}
}
