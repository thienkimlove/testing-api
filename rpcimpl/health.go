package rpcimpl

import (
	"context"

	"go.tekoapis.com/kitchen/log/level"
	"rpc.tekoapis.com/rpc/health"
)

type HealthServer struct {
}

func (*HealthServer) Liveness(ctx context.Context, livenessRequest *health.LivenessRequest) (*health.LivenessResponse, error) {

	level.Info(ctx).L("This is example log by grpc context")

	return &health.LivenessResponse{
		Content: "ok",
	}, nil
}

func (*HealthServer) Readiness(ctx context.Context, txt *health.ReadinessRequest) (*health.ReadinessResponse, error) {
	return &health.ReadinessResponse{
		Content: "ok",
	}, nil
}
