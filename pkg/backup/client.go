package backup

import (
	"context"
	"fmt"
	"github.com/qkfktm54/grpc-backup-service/pkg/backup/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	c    proto.BackupClient
	conn *grpc.ClientConn
}

func (c *Client) log(str string, args ...interface{}) {
	fmt.Printf("[Client] %s\n", fmt.Sprintf(str, args...))
}

func NewClient(ctx context.Context, serverAddr string) (*Client, error) {
	conn, err := grpc.DialContext(ctx, serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Backup server: %v", err)
	}
	return &Client{c: proto.NewBackupClient(conn), conn: conn}, nil
}

func (c *Client) Close() error {
	c.log("Closing connection")
	return c.conn.Close()
}

func (c *Client) Dir(ctx context.Context, subDir string) ([]string, error) {
	c.log("Trying to connect to get contents of `%s`", subDir)
	var sd *string
	if subDir == "" {
		sd = nil
	} else {
		sd = &subDir
	}
	req := &proto.DirectoryRequest{SubDirectory: sd}

	resp, err := c.c.Dir(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform `Dir` on Backup server: %v", err)
	}
	return resp.GetContents(), nil
}
