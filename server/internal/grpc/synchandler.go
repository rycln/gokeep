package grpc

import (
	"context"

	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/shared/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// syncService defines the required domain operations for sync
type syncService interface {
	SyncItems(context.Context, []models.Item) ([]models.Item, error)
}

// Sync handles sync requests
func (h *GophKeeperServer) Sync(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
	var clientitems = make([]models.Item, len(req.Items))
	for i, reqitem := range req.Items {
		clientitems[i].ID = models.ItemID(reqitem.Id)
		clientitems[i].UserID = models.UserID(reqitem.UserId)
		clientitems[i].ItemType = models.ItemType(reqitem.Type)
		clientitems[i].Name = reqitem.Name
		clientitems[i].Metadata = reqitem.Metadata
		clientitems[i].Data = reqitem.Data
		clientitems[i].UpdatedAt = reqitem.UpdatedAt.AsTime()
		clientitems[i].IsDeleted = reqitem.IsDeleted
	}

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	serveritems, err := h.sync.SyncItems(ctx, clientitems)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var resitems = make([]*pb.Item, len(serveritems))
	for i, serveritem := range serveritems {
		var resitem = &pb.Item{}
		resitem.Id = string(serveritem.ID)
		resitem.UserId = string(serveritem.UserID)
		resitem.Type = string(serveritem.ItemType)
		resitem.Name = serveritem.Name
		resitem.Metadata = serveritem.Metadata
		resitem.Data = serveritem.Data
		resitem.UpdatedAt = timestamppb.Now()
		resitem.IsDeleted = serveritem.IsDeleted

		resitems[i] = resitem
	}

	return &pb.SyncResponse{
		Items: resitems,
	}, nil
}
