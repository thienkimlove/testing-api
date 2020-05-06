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

func (h *ListingServer) ListingRq2(ctx context.Context, message *listing.ListingRequest) (*listing.ListingResponse2, error){
	level.Info(ctx).L("This is example log by grpc from listing2")
	return &listing.ListingResponse2{Content2: "hello, reply from  ListingRq2"}, nil
}