package query_server

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"

	pb "github.com/popstk/subserver/query-gateway/query"
)

var (
	configFile string
	config     Config
)

func init() {
	log.SetFlags(log.Lshortfile|log.Ltime)
	flag.StringVar(&configFile, "c", "config.json", "config file")
}


func parseConfig(root string) ([]string, error) {
	urls := make([]string, 0)
	m := make(map[string]bool)

	q := make([]string, 0, 1)
	q = append(q, root)
	for len(q) > 0 {

		uuid := q[0]
		q = q[1:]
		m[uuid] = true

		source, exist := config.Valid[uuid]
		if !exist {
			return nil, errors.New("Invalid uuid: "+uuid)
		}

		for _, s := range source {
			if s.Type == "sub" {
				_, exist := m[s.Addr]
				if !exist {
					q = append(q, s.Addr)
				}

				continue
			}

			u, err := s.Parse()
			if err != nil {
				log.Print(err)
				continue
			}
			urls = append(urls, u...)
		}
	}

	return urls, nil
}

type Server struct {}

func (s *Server) Query(ctx context.Context, request *pb.Request) (*pb.Reply, error) {
	uuid := request.Uuid
	log.Println("uuid is ", uuid)

	urls, err := parseConfig(uuid)
	if err != nil {
		return nil, err
	}

	valid := make([]string, 0, len(urls))
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u != "" {
			valid =append(valid, u)
		}
	}

	respond := strings.Join(valid, "\n")
	if request.Decode {
		respond = base64.StdEncoding.EncodeToString([]byte(respond))
	}

	return &pb.Reply{
		Message:respond,
	}, nil
}


func Serve(endpoint string) error {
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}

	s:= grpc.NewServer()
	pb.RegisterSubserverServer(s, &Server{})

	return s.Serve(lis)
}

