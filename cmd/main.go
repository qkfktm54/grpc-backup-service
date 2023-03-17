package main

import (
	"context"
	"fmt"
	"github.com/qkfktm54/grpc-backup-service/pkg/backup"
	"time"
)

func main() {
	server := backup.NewServer()
	addr := "localhost:9988"
	go func() {
		fmt.Printf("Starting gRPC backup server on `%s`...\n", addr)
		started := time.Now()
		var err error
		if err = server.Listen(addr); err != nil {
			fmt.Printf("got error during server listen loop: %v\n", err)
		}
		fmt.Printf("Listening stopped after %s. Error: %v\n", time.Now().Sub(started).String(), err)
	}()
	defer server.Close(false)

	client, err := backup.NewClient(context.Background(), addr)
	if err != nil {
		fmt.Printf("failed to instantiate backup client: %v\n", err)
	}
	defer client.Close()

	for _, lsDir := range []string{".", "Documents"} {
		contents, err := client.Dir(context.Background(), lsDir)
		if err != nil {
			fmt.Printf("Got error when requesting contents of `%s`: %v\n", lsDir, err)
		} else {
			fmt.Printf("Contents of `%s`:\n", lsDir)
			for _, c := range contents {
				fmt.Println(" - ", c)
			}
		}
	}
}
