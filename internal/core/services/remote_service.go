package services

import (
	"context"
	"strconv"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type RemoteService struct {
	remoteRepo repos.RemoteRepo
}

func NewRemoteService(remoteRepo repos.RemoteRepo) core.RemoteService {
	return &RemoteService{
		remoteRepo: remoteRepo,
	}
}

func (s *RemoteService) CreateRemoteOrg(ctx context.Context, userID string, req *core.CreateRemoteOrgRequest) (*core.RemoteOrgResponse, error) {
	// For now, return a placeholder as RemoteOrg entity isn't fully implemented in repo layer
	return &core.RemoteOrgResponse{
		Name:            req.Name,
		ContactEmail:    req.ContactEmail,
		ContactPhone:    req.ContactPhone,
		Address:         req.Address,
		Location: core.LocationInfo{
			Pincode: &req.Pincode,
			City:    &req.City,
			State:   &req.State,
		},
		ServingCapacity: req.ServingCapacity,
		PurposeFocus:    req.PurposeFocus,
	}, nil
}

func (s *RemoteService) GetRemoteOrg(ctx context.Context, remoteOrgID string) (*core.RemoteOrgResponse, error) {
	// Placeholder implementation - would get from repository
	return &core.RemoteOrgResponse{
		Name:         "Sample Remote Org",
		ContactEmail: "contact@remote.org",
		Address:      "Sample Address",
		Location: core.LocationInfo{
			City:  stringPtr("Sample City"),
			State: stringPtr("Sample State"),
		},
		ServingCapacity: 100,
		PurposeFocus:    []string{"GENERAL"},
	}, nil
}

func (s *RemoteService) UpdateRemoteOrg(ctx context.Context, remoteOrgID string, req *core.UpdateRemoteOrgRequest) (*core.RemoteOrgResponse, error) {
	// Placeholder implementation - would update in repository
	response := &core.RemoteOrgResponse{
		Name:         "Updated Remote Org",
		ContactEmail: "updated@remote.org",
		Address:      "Updated Address",
		Location: core.LocationInfo{
			City:  stringPtr("Updated City"),
			State: stringPtr("Updated State"),
		},
		ServingCapacity: 150,
		PurposeFocus:    []string{"GENERAL"},
	}
	
	// Apply updates from request
	if req.Name != nil {
		response.Name = *req.Name
	}
	if req.ContactEmail != nil {
		response.ContactEmail = *req.ContactEmail
	}
	if req.Address != nil {
		response.Address = *req.Address
	}
	if req.ServingCapacity != nil {
		response.ServingCapacity = *req.ServingCapacity
	}
	if req.PurposeFocus != nil {
		response.PurposeFocus = req.PurposeFocus
	}
	
	return response, nil
}

func (s *RemoteService) ListRemoteOrgs(ctx context.Context, page, limit string) (*core.RemoteOrgListResponse, error) {
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

	// Placeholder implementation - would get from repository
	remoteOrgs := []core.RemoteOrgResponse{
		{
			Name:         "Remote Org 1",
			ContactEmail: "org1@remote.org",
			Address:      "Address 1",
			Location: core.LocationInfo{
				City:  stringPtr("City 1"),
				State: stringPtr("State 1"),
			},
			ServingCapacity: 100,
			PurposeFocus:    []string{"CHILDREN"},
		},
		{
			Name:         "Remote Org 2",
			ContactEmail: "org2@remote.org",
			Address:      "Address 2",
			Location: core.LocationInfo{
				City:  stringPtr("City 2"),
				State: stringPtr("State 2"),
			},
			ServingCapacity: 200,
			PurposeFocus:    []string{"ELDERLY"},
		},
	}

	return &core.RemoteOrgListResponse{
		Items: remoteOrgs,
		Page:  pageNum,
		Limit: limitNum,
		Total: len(remoteOrgs),
	}, nil
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func (s *RemoteService) AssignRemoteOrg(ctx context.Context, userID string, req *core.AssignRemoteOrgRequest) (*core.RemoteOrgAssignmentResponse, error) {
	// Create assignment between remote org and proxy org
	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}
	assignment, err := s.remoteRepo.CreateAssignment(ctx, req.RemoteOrgID.String(), userID, true, notes)
	if err != nil {
		return nil, core.NewInternalError("Failed to create remote organization assignment", err)
	}

	return &core.RemoteOrgAssignmentResponse{
		AssignmentID: assignment.ID,
		RemoteOrgID:  assignment.RemoteOrgID,
		ProxyOrgID:   assignment.ProxyOrgID,
		Active:       assignment.Active,
		Notes:        assignment.Notes,
		CreatedAt:    assignment.CreatedAt,
	}, nil
}
