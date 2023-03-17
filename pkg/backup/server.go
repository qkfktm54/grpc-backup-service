package backup

import (
	"context"
	"fmt"
	"github.com/qkfktm54/grpc-backup-service/pkg/backup/proto"
	"google.golang.org/grpc"
	"net"
	"os"
	"path/filepath"
)

type Server struct {
	proto.BackupServer
	svr       *grpc.Server
	rootDir   string
	listening bool
}

func (s *Server) log(str string, args ...interface{}) {
	fmt.Printf("[Server] %s\n", fmt.Sprintf(str, args...))
}

func (s *Server) Dir(ctx context.Context, req *proto.DirectoryRequest) (*proto.DirectoryResponse, error) {
	target, err := filepath.Abs(s.rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absPath for root dir `%s`: %v", s.rootDir, err)
	}
	if req.SubDirectory != nil {
		target = filepath.Join(target, *req.SubDirectory)
	}
	s.log("[Server] Got `Dir` request for directory `%s`", target)
	stat, err := os.Stat(target)
	if err != nil {
		return nil, fmt.Errorf("directory `%s` does not exists: %v", target, err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("cannot get contents of non-directory `%s`", target)
	}
	s.log("Enumerating contents of `%s`...", target)
	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory `%s`: %v", target, err)
	}
	s.log("Enumerating contents of `%s` yielded %d results", target, len(entries))
	contents := make([]string, len(entries))
	for i, entry := range entries {
		entryName := entry.Name()
		if entry.IsDir() {
			entryName += "/"
		}
		contents[i] = entryName
	}
	return &proto.DirectoryResponse{Contents: contents}, nil
}

func (s *Server) Listen(addr string) error {
	s.log("Trying to listen on `%s`...", addr)
	if s.svr == nil {
		return fmt.Errorf("no gRPC server instance found")
	}
	if s.listening {
		return nil
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on `%s`: %v", addr, err)
	}
	s.listening = true
	defer func() { s.listening = false }()
	s.log("Server listening on `%s`", addr)
	return s.svr.Serve(lis)
}

func (s *Server) Close(force bool) {
	if force {
		s.log("Force stopping gRPC server")
		s.svr.Stop()
		return
	}
	s.log("Gracefully stopping gRPC server")
	s.svr.GracefulStop()
}

func NewServer() *Server {

	svr := grpc.NewServer()
	backupServer := &Server{
		svr:     svr,
		rootDir: "./backups",
	}
	proto.RegisterBackupServer(svr, backupServer)
	return backupServer
}
