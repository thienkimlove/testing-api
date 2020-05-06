// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
//   protoc 3.11.4
// sources:
//   health/v1/health.proto
//   sample/v1/sample.proto

package server

import (
	context "context"
	fmt "fmt"
	go_grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	grpc "google.golang.org/grpc"
	net "net"
	http "net/http"
	os "os"
	signal "os/signal"
	health "rpc.tekoapis.com/rpc/health"
	sample "rpc.tekoapis.com/rpc/sample"
	syscall "syscall"
	time "time"
)

func DefaultConfig() Config {
	return Config{
		GRPC: ServerListen{
			Host: "0.0.0.0",
			Port: 10443,
		},
		HTTP: ServerListen{
			Host: "0.0.0.0",
			Port: 10080,
		},
	}
}

type Config struct {
	GRPC ServerListen `json:"grpc"`
	HTTP ServerListen `json:"http"`
}

type ServerListen struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (l ServerListen) String() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

type Server struct {
	gRPC *grpc.Server
	mux  *runtime.ServeMux
	cfg  Config
}

func NewServer(cfg Config, opt ...grpc.ServerOption) *Server {
	go_grpc_prometheus.EnableHandlingTimeHistogram()
	opt = append(opt, grpc.ChainStreamInterceptor(go_grpc_prometheus.StreamServerInterceptor), grpc.ChainUnaryInterceptor(go_grpc_prometheus.UnaryServerInterceptor))
	return &Server{
		gRPC: grpc.NewServer(opt...),
		mux:  runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true})),
		cfg:  cfg,
	}
}

func (s *Server) Register(grpcServer ...interface{}) error {
	for _, srv := range grpcServer {
		switch _srv := srv.(type) {
		case health.HealthCheckServiceServer:
			health.RegisterHealthCheckServiceServer(s.gRPC, _srv)
			if err := health.RegisterHealthCheckServiceHandlerFromEndpoint(context.Background(), s.mux, s.cfg.GRPC.String(), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
				return err
			}
		case sample.SampleServiceServer:
			sample.RegisterSampleServiceServer(s.gRPC, _srv)
			if err := sample.RegisterSampleServiceHandlerFromEndpoint(context.Background(), s.mux, s.cfg.GRPC.String(), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unknown GPRC Service to register %#v", srv)
		}
	}
	return nil
}

func (s *Server) Serve() error {
	stop := make(chan os.Signal, 1)
	errch := make(chan error)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())
	httpMux.Handle("/", s.mux)
	httpServer := http.Server{
		Addr:    s.cfg.HTTP.String(),
		Handler: httpMux,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errch <- err
		}
	}()
	go func() {
		listener, err := net.Listen("tcp", s.cfg.GRPC.String())
		if err != nil {
			errch <- err
			return
		}
		if err := s.gRPC.Serve(listener); err != nil {
			errch <- err
		}
	}()
	for {
		select {
		case <-stop:
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			httpServer.Shutdown(ctx)
			s.gRPC.GracefulStop()
			return nil
		case err := <-errch:
			return err
		}
	}
}