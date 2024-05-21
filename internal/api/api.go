package api

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"pgxrio/internal/inventory"
	"strings"
	"sync"
	"time"
)

type Server struct {
	httpAddress string
	grpcAddress string
	service     *inventory.Service

	httpServer *httpServer
	grpcServer *grpcServer

	stopFn *sync.Once
}

type httpServer struct {
	inventory  *inventory.Service
	middleware func(handler http.Handler) http.Handler
	server     *http.Server
}

func (s httpServer) Run(ctx context.Context, address string) error {

	var handler http.Handler = NewHttpServer(s.inventory)

	if s.middleware != nil {
		handler = s.middleware(handler)
	}

	s.server = &http.Server{Addr: address, Handler: handler, ReadHeaderTimeout: time.Second * 5}

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil

	return nil

}

func (s *httpServer) Shutdown(ctx context.Context) {

	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			log.Println("Gracefully shutdown http server is fail!")
		}
	}

}
func (s *grpcServer) Run(ctx context.Context, address string) error {
	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("Fail to listen server ::: %v", err)

	}

	s.server = grpc.NewServer()
	reflection.Register(s.server)

	log.Printf("grpc server listen at %s", lis.Addr())
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("fail grpc to server ::: ", err)
	}
	return nil

	return nil

}

func (s *grpcServer) Shutdown(ctx context.Context) {

	done := make(chan struct{}, 1)
	log.Println("starting shut down grpc server ...")
	go func() {
		if s.server != nil {
			s.server.GracefulStop()

		}
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		if s.server != nil {
			s.server.Stop()

		}
		log.Println("Gracefully shutdown grppc server is fail!")

	}
}

type grpcServer struct {
	inventory *inventory.Service
	server    *grpc.Server
}

func (s *Server) Run(ctx context.Context) (err error) {
	var errorChannel = make(chan error, 2)
	ctx, cancel := context.WithCancel(ctx)

	s.httpServer = &httpServer{
		inventory: s.service,
	}
	s.grpcServer = &grpcServer{inventory: s.service}
	go func() {
		err := s.httpServer.Run(ctx, s.httpAddress)
		if err != nil {
			err = fmt.Errorf("Error with HTTP server ::: %v", err)
		}
		errorChannel <- err
	}()
	go func() {
		err := s.grpcServer.Run(ctx, s.grpcAddress)
		if err != nil {
			err = fmt.Errorf("Error with GRPC server ::: %v", err)
		}
		errorChannel <- err
	}()
	var errorString []string

	for i := 0; i < cap(errorChannel); i++ {

		if err := <-errorChannel; err != nil {
			errorString = append(errorString, err.Error())

			if ctx.Err() == nil {
				s.Shutdown(context.Background())
			}
		}
	}
	if len(errorString) > 0 {
		err = errors.New(strings.Join(errorString, ", "))
	}

	cancel()
	return err
}

func (s *Server) Shutdown(ctx context.Context) {

	s.stopFn.Do(func() {
		s.httpServer.Shutdown(ctx)
		s.grpcServer.Shutdown(ctx)
	})
}
