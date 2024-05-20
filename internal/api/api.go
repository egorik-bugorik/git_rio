package api

import "sync"

type Server struct {
	httpAddress string
	grpcAddress string

	httpServer *httpServer
	grpcServer *grpcServer

	fn *sync.Once
}

type httpServer struct {
}

type grpcServer struct {
}
