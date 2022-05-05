package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	pb "github.com/popstk/subserver/proto/query"
	"github.com/popstk/subserver/service/query"
	"github.com/sevlyar/go-daemon"
	"google.golang.org/grpc"
)

var (
	configFile string
	gwAddr     string
	daemonMode bool
	signal     string
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.BoolVar(&daemonMode, "d", false, "daemon mode")
	flag.StringVar(&configFile, "c", "subserver.json", "config file")
	flag.StringVar(&gwAddr, "addr", "localhost:10086", "gateway address")
	flag.StringVar(&signal, "s", "", "Send signal to the daemon: stop")
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

func Run(ctx context.Context) {
	addr, err := query.Serve(configFile)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
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

		const Redirect = "redirect/"

		if strings.HasPrefix(reply.Message, Redirect) {
			http.Redirect(w, r, reply.Message[len(Redirect):], http.StatusTemporaryRedirect)
			return
		}

		_, _ = w.Write([]byte(reply.Message))
	})

	log.Println("gateway: listen on ", gwAddr)
	server := &http.Server{Addr: gwAddr, Handler: mux}
	go func() {
		<-ctx.Done()
		if err := server.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(err)
		return
	}
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	daemon.AddCommand(daemon.StringFlag(&signal, "stop"), syscall.SIGTERM, func(sig os.Signal) (err error) {
		fmt.Println("terminating...")
		cancel()
		return daemon.ErrStop
	})

	path, _ := os.Executable()
	path = filepath.Dir(path)

	if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(path, configFile)
	}

	name := filepath.Base(os.Args[0])
	cntxt := &daemon.Context{
		PidFileName: filepath.Join(path, fmt.Sprintf("%s.pid", name)),
		PidFilePerm: 0644,
		LogFileName: filepath.Join(path, fmt.Sprintf("%s.out", name)),
		LogFilePerm: 0640,
		WorkDir:     path,
		Umask:       027,
	}

	// client process signal first
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			fmt.Println("Unable send signal to the daemon: ", err)
			return
		}
		daemon.SendCommands(d)
		return
	}

	// start daemon
	if daemonMode {
		d, err := cntxt.Reborn()
		if err != nil {
			fmt.Println("Unable to run: ", err)
			os.Exit(-1)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()
	}

	go Run(ctx)

	// daemon process signal
	err := daemon.ServeSignals()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("bye bye")
}
