package services

import (
	"context"
	"strconv"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type CollabService struct {
	collabRepo repos.CollabRepo
}

func NewCollabService(collabRepo repos.CollabRepo) core.CollabService {
	return &CollabService{
		collabRepo: collabRepo,
	}
}

func (s *CollabService) CreateCollaborator(ctx context.Context, userID string, req *core.CreateCollabRequest) (*core.CollabResponse, error) {
	// Create collaborator record
	collaborator, err := s.collabRepo.CreateCollaborator(ctx, userID, req.CollabType)
	if err != nil {
		return nil, core.NewInternalError("Failed to create collaborator", err)
	}

	return &core.CollabResponse{
		CollabID:   collaborator.ID,
		CollabType: collaborator.CollabType,
	}, nil
}

func (s *CollabService) GetCollaborator(ctx context.Context, collabID string) (*core.CollabResponse, error) {
	collaborator, err := s.collabRepo.GetCollaboratorByID(ctx, collabID)
	if err != nil {
		return nil, err
	}

	return &core.CollabResponse{
		CollabID:   collaborator.ID,
		CollabType: collaborator.CollabType,
	}, nil
}

func (s *CollabService) UpdateCollaborator(ctx context.Context, collabID string, req *core.UpdateCollabRequest) (*core.CollabResponse, error) {
	// Get existing collaborator
	collaborator, err := s.collabRepo.GetCollaboratorByID(ctx, collabID)
	if err != nil {
		return nil, err
	}

	// Update collaborator type if provided
	if req.CollabType != nil {
		collaborator.CollabType = *req.CollabType
		
		// Update in repository
		updatedCollab, err := s.collabRepo.UpdateCollaborator(ctx, collabID, *req.CollabType)
		if err != nil {
			return nil, core.NewInternalError("Failed to update collaborator", err)
		}
		collaborator = updatedCollab
	}

	return &core.CollabResponse{
		CollabID:   collaborator.ID,
		CollabType: collaborator.CollabType,
	}, nil
}

func (s *CollabService) ListCollaborators(ctx context.Context, page, limit string) (*core.CollabListResponse, error) {
	// Parse pagination parameters
	pageNum := 1
	limitNum := 10
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitNum = l
		}
	}
	offset := (pageNum - 1) * limitNum

	// Get collaborators
	collaborators, err := s.collabRepo.GetAllCollaborators(ctx, limitNum, offset)
	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve collaborators", err)
	}

	// Convert to response format
	collabResponses := make([]core.CollabResponse, len(collaborators))
	for i, collab := range collaborators {
		collabResponses[i] = core.CollabResponse{
			CollabID:   collab.ID,
			CollabType: collab.CollabType,
		}
	}

	return &core.CollabListResponse{
		Items: collabResponses,
		Page:  pageNum,
		Limit: limitNum,
		Total: len(collaborators), // In a real app, this would be a separate count query
	}, nil
}
