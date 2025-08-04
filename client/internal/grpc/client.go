package grpc

import (
	"context"

	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/shared/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GophKeeperClient handles operations via gRPC
type GophKeeperClient struct {
	client pb.GophKeeperClient // gRPC generated client interface
}

// NewGophKeeperClient creates new GophKeeperClient instance
func NewGophKeeperClient(conn *grpc.ClientConn) *GophKeeperClient {
	return &GophKeeperClient{
		client: pb.NewGophKeeperClient(conn),
	}
}

// Register performs user registration via gRPC
func (c *GophKeeperClient) Register(ctx context.Context, req *models.UserRegReq) (*models.User, error) {
	res, err := c.client.Register(ctx, &pb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Salt:     req.Salt,
	})

	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:   models.UserID(res.UserId),
		JWT:  res.Token,
		Salt: res.Salt,
	}, err
}

// Login performs user authentication via gRPC
func (c *GophKeeperClient) Login(ctx context.Context, req *models.UserLoginReq) (*models.User, error) {
	res, err := c.client.Login(ctx, &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:   models.UserID(res.UserId),
		JWT:  res.Token,
		Salt: res.Salt,
	}, err
}

// Sync performs bidirectional items synchronization with server via gRPC
// Converts local items to protobuf format and back
func (c *GophKeeperClient) Sync(ctx context.Context, clientItems []models.Item, jwt string) ([]models.Item, error) {
	var reqitems = make([]*pb.Item, len(clientItems))
	for i, item := range clientItems {
		var reqitem = &pb.Item{}
		reqitem.Id = string(item.ID)
		reqitem.UserId = string(item.UserID)
		reqitem.Type = string(item.ItemType)
		reqitem.Name = item.Name
		reqitem.Metadata = item.Metadata
		reqitem.Data = item.Data
		reqitem.UpdatedAt = timestamppb.Now()
		reqitem.IsDeleted = item.IsDeleted

		reqitems[i] = reqitem
	}

	md := metadata.Pairs("authorization", "Bearer "+jwt)
	ctx = metadata.NewOutgoingContext(ctx, md)

	res, err := c.client.Sync(ctx, &pb.SyncRequest{
		Items: reqitems,
	})
	if err != nil {
		return nil, err
	}

	var serverItems = make([]models.Item, len(res.Items))
	for i, resitem := range res.Items {
		serverItems[i].ID = models.ItemID(resitem.Id)
		serverItems[i].UserID = models.UserID(resitem.UserId)
		serverItems[i].ItemType = models.ItemType(resitem.Type)
		serverItems[i].Name = resitem.Name
		serverItems[i].Metadata = resitem.Metadata
		serverItems[i].Data = resitem.Data
		serverItems[i].UpdatedAt = resitem.UpdatedAt.AsTime()
		serverItems[i].IsDeleted = resitem.IsDeleted
	}

	return serverItems, nil
}
