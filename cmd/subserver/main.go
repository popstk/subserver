package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/popstk/subserver/proto/query"
	"github.com/popstk/subserver/service/query"
	"google.golang.org/grpc"
)

var (
	configFile   string
	gwAddr       string
	querySrvAddr string
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.StringVar(&configFile, "c", "subserver.json", "config file")
	flag.StringVar(&gwAddr, "addr", "localhost:10086", "gateway address")
}

/*
func gateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := pb.RegisterSubscribeHandlerFromEndpoint(
		ctx, mux, querySrvAddr, opts); err != nil {
		return err
	}

	log.Println("gateway: listen on ", gwAddr)
	return http.ListenAndServe(gwAddr, mux)
}
*/

func run() error {
	conn, err := grpc.Dial(querySrvAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewSubscribeClient(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		decode := false
		if r.URL.Query().Get("decode") == "true" {
			decode = true
		}
		uuid := strings.Trim(r.URL.Path, "/")

		reply, err := client.Query(context.Background(), &pb.Request{Uuid: uuid, Decode: decode})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(reply.Message))
	})

	log.Println("gateway: listen on ", gwAddr)
	return http.ListenAndServe(gwAddr, mux)
}

func main() {
	flag.Parse()

	if !filepath.IsAbs(configFile) {
		path, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		configFile = filepath.Join(filepath.Dir(path), configFile)
	}

	var err error
	querySrvAddr, err = query.Serve(configFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
