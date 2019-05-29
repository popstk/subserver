package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/popstk/subserver/query-gateway/query"
	query_server "github.com/popstk/subserver/query-gateway/query-server"
	"google.golang.org/grpc"
	"log"
	"net/http"
)


const (
	httpAddress = "localhost:10086"
	subserver = "localhost:10087"
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	go query_server.Serve(subserver)

	err := gw.RegisterSubserverHandlerFromEndpoint(ctx, mux, subserver, opts)

	if err != nil {
		return err
	}

	log.Println("listen: ", httpAddress)
	return http.ListenAndServe(httpAddress, mux)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

