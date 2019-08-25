package query

import (
	"context"
	"encoding/base64"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
	pb "github.com/popstk/subserver/proto/query"
	"google.golang.org/grpc"
)

var (
	configFile  string
	config      Config
	configMutex sync.Mutex
)

// Server -
type Server struct{}

// Query grpc service
func (s *Server) Query(ctx context.Context, request *pb.Request) (*pb.Reply, error) {
	uuid := request.Uuid
	log.Println("uuid =>  ", uuid)

	nodes, err := ParseConfig(uuid)
	if err != nil {
		return nil, err
	}

	// filter localhost node
	var valid []string
	for _, node := range nodes {
		host, _, err := net.SplitHostPort(node.Addr())
		if err != nil {
			log.Println(err)
			continue
		}
		if host == "127.0.0.1" || host == "localhost" {
			continue
		}

		valid = append(valid, node.String())
	}

	respond := strings.Join(valid, "\n")
	if !request.Decode {
		respond = base64.StdEncoding.EncodeToString([]byte(respond))
	}

	return &pb.Reply{
		Message: respond,
	}, nil
}

// Serve serve http server with config, return listen address
func Serve(path string) (string, error) {
	configFile = path

	if err := LoadConfig(); err != nil {
		return "", errors.Wrap(err, "query: can not load config")
	}
	go ListenConfig(context.Background())

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", errors.Wrap(err, "query: can not listen")
	}

	s := grpc.NewServer()
	pb.RegisterSubscribeServer(s, &Server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Println(err)
		}
	}()

	addr := lis.Addr().String()
	log.Println("query: listen on ", addr)

	return addr, nil
}
