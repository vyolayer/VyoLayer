package grpc

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/vyolayer/vyolayer/internal/apikeys/service"
	pb "github.com/vyolayer/vyolayer/proto/apikey/v1"
)

type Server struct {
	pb.UnimplementedAPIKeyServiceServer
	svc service.Service
}

func New(svc service.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) CreateAPIKey(
	ctx context.Context,
	req *pb.CreateAPIKeyRequest,
) (*pb.CreateAPIKeyResponse, error) {

	orgID, _ := uuid.Parse(req.OrganizationId)
	projID, _ := uuid.Parse(req.ProjectId)
	actorID, _ := uuid.Parse(req.ActorId)

	item, secret, err := s.svc.Create(
		ctx,
		orgID,
		projID,
		actorID,
		req.Name,
		req.Description,
		strings.ToLower(req.Environment),
		req.Scopes,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateAPIKeyResponse{
		ApiKey: toProtoAPIKey(item, req.Scopes),
		Secret: secret,
	}, nil
}

func (s *Server) ListAPIKeys(
	ctx context.Context,
	req *pb.ListAPIKeysRequest,
) (*pb.ListAPIKeysResponse, error) {

	orgID, _ := uuid.Parse(req.OrganizationId)
	projID, _ := uuid.Parse(req.ProjectId)

	items, err := s.svc.List(ctx, orgID, projID)
	if err != nil {
		return nil, err
	}

	out := make([]*pb.APIKey, 0, len(items))
	for _, it := range items {
		out = append(out, toProtoAPIKey(&it, nil))
	}

	return &pb.ListAPIKeysResponse{
		ApiKeys: out,
	}, nil
}

func (s *Server) GetAPIKey(
	ctx context.Context,
	req *pb.GetAPIKeyRequest,
) (*pb.GetAPIKeyResponse, error) {

	id, _ := uuid.Parse(req.Id)
	orgID, _ := uuid.Parse(req.OrganizationId)
	projID, _ := uuid.Parse(req.ProjectId)

	item, err := s.svc.Get(ctx, id, orgID, projID)
	if err != nil {
		return nil, err
	}

	return &pb.GetAPIKeyResponse{
		ApiKey: toProtoAPIKey(item, nil),
	}, nil
}

func (s *Server) RevokeAPIKey(
	ctx context.Context,
	req *pb.RevokeAPIKeyRequest,
) (*pb.RevokeAPIKeyResponse, error) {

	id, _ := uuid.Parse(req.Id)
	orgID, _ := uuid.Parse(req.OrganizationId)
	projID, _ := uuid.Parse(req.ProjectId)
	actorID, _ := uuid.Parse(req.ActorId)

	err := s.svc.Revoke(ctx, id, orgID, projID, actorID)
	if err != nil {
		return nil, err
	}

	return &pb.RevokeAPIKeyResponse{
		Success: true,
	}, nil
}

func (s *Server) RotateAPIKey(
	ctx context.Context,
	req *pb.RotateAPIKeyRequest,
) (*pb.RotateAPIKeyResponse, error) {

	id, _ := uuid.Parse(req.Id)
	orgID, _ := uuid.Parse(req.OrganizationId)
	projID, _ := uuid.Parse(req.ProjectId)
	actorID, _ := uuid.Parse(req.ActorId)

	item, secret, err := s.svc.Rotate(ctx, id, orgID, projID, actorID)
	if err != nil {
		return nil, err
	}

	return &pb.RotateAPIKeyResponse{
		ApiKey: toProtoAPIKey(item, nil),
		Secret: secret,
	}, nil
}

func (s *Server) ValidateAPIKey(
	ctx context.Context,
	req *pb.ValidateAPIKeyRequest,
) (*pb.ValidateAPIKeyResponse, error) {

	ok, item, scopes, err := s.svc.Validate(ctx, req.Secret)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateAPIKeyResponse{
		Valid:  ok,
		ApiKey: toProtoAPIKey(item, scopes),
		Scopes: scopes,
	}, nil
}
