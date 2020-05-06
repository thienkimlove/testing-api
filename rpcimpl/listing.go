package rpcimpl

import (
	"context"
	"github.com/thienkimlove/testing-api/custom_rpc/listing"
	"go.tekoapis.com/kitchen/log/level"
)
type ListingServer struct {
     listing.UnimplementedListingServiceServer
}

func (h *ListingServer) ListingRq(ctx context.Context, message *listing.ListingRequest) (*listing.ListingResponse, error){
	level.Info(ctx).L("This is example log by grpc from listing")
	return &listing.ListingResponse{Content: "hello, reply from "}, nil
}