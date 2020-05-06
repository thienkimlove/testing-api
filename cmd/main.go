package main

import (
	"github.com/thienkimlove/testing-api/config"
	"github.com/thienkimlove/testing-api/rpcimpl"
	"github.com/thienkimlove/testing-api/rpcimpl/server"
	"log"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.tekoapis.com/kitchen/database/migrate"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	cfg, err := config.Load()

	if err != nil {
		return err
	}

	logger, err := cfg.Log.Build()

	if err != nil {
		return err
	}

	cmd := &cobra.Command{
		Use: "rpc-application",
	}

	cmd.AddCommand(&cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			grpc_zap.ReplaceGrpcLoggerV2(logger)

			s := server.NewServer(cfg.Server,
				grpc_middleware.WithUnaryServerChain(
					grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
					grpc_zap.UnaryServerInterceptor(logger),
				),
				grpc_middleware.WithStreamServerChain(
					grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
					grpc_zap.StreamServerInterceptor(logger),
				),
			)
			// Register your rpc server here
			// You may register multiple server
			// s.Register(server1, server2, server3)
			//
			healthServer := &rpcimpl.HealthServer{}
			listingServer := &rpcimpl.ListingServer{}

			if err := s.Register(healthServer, listingServer); err != nil {
				logger.Fatal("Error register servers", zap.Any("error", err))
			}

			logger.Warn("Starting server", zap.Any("grpc", cfg.Server.GRPC), zap.Any("http", cfg.Server.HTTP))

			if err := s.Serve(); err != nil {
				logger.Fatal("Error start server", zap.Any("error", err))
			}
		},
	}, migrate.Command(cfg.MigrationsFolder, cfg.MySQL.String()))
	return cmd.Execute()
}
