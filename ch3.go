package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	ctx context.Context
	cancel func()
	https []*httpServer
}

type httpServer struct {
	addr string
	handle http.Handler
	ctx context.Context
	cancel func()
}

func NewServe(addr []string) *Server{
	ctx, cancel := context.WithCancel(context.Background())

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Hello, world!`))
	})

	var https []*httpServer
	for _, v := range addr {
		https = append(https, NewHttpServe(v,ctx,mux))
	}
	return &Server{
		ctx:    ctx,
		cancel: cancel,
		https:	https,
	}
}

func NewHttpServe(addr string, ctx context.Context,handler http.Handler) *httpServer{
	ctx, cancel := context.WithCancel(ctx)
	return &httpServer{
		addr:   addr,
		handle: handler,
		ctx:    ctx,
		cancel: cancel,
	}
}


func (s *Server) Run() error {
	g, ctx := errgroup.WithContext(s.ctx)

	for _, v := range s.https {
		g.Go(v.Serves)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-signals:
				s.cancel()
			}
		}
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (h *httpServer) Serves() error {
	s := &http.Server{
		Addr:              h.addr,
		Handler:           h.handle,
	}
	go func() {
		<-h.ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(shutCtx)
		fmt.Printf("%s server is closed\n", h.addr)
	}()

	fmt.Println("Start",h.addr)
	return  s.ListenAndServe()
}

func main(){
	s := NewServe([]string{"0.0.0.0:9090","0.0.0.0:9098"})

	if err := s.Run(); err != nil && !errors.Is(err, context.Canceled){
		log.Println(err)
	}
	fmt.Println("everything closed successfully")
}