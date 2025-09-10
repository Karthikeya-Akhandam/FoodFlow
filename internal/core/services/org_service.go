package services

import (
	"context"
	"strconv"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type OrgService struct {
	orgRepo repos.OrgRepo
}

func NewOrgService(orgRepo repos.OrgRepo) core.OrgService {
	return &OrgService{
		orgRepo: orgRepo,
	}
}

func (s *OrgService) CreateOrganization(ctx context.Context, userID string, req *core.CreateOrgRequest) (*core.OrgResponse, error) {
	// Create organization record
	org, err := s.orgRepo.CreateOrganization(ctx, userID, req.OrgType, req.PurposeFocus, req.LastMonthPeopleFed)
	if err != nil {
		return nil, core.NewInternalError("Failed to create organization", err)
	}

	return &core.OrgResponse{
		OrgID:              org.ID,
		OrgType:            org.OrgType,
		PurposeFocus:       org.PurposeFocus,
		LastMonthPeopleFed: org.LastMonthPeopleFed,
	}, nil
}

func (s *OrgService) GetOrganization(ctx context.Context, orgID string) (*core.OrgResponse, error) {
	org, err := s.orgRepo.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return &core.OrgResponse{
		OrgID:              org.ID,
		OrgType:            org.OrgType,
		PurposeFocus:       org.PurposeFocus,
		LastMonthPeopleFed: org.LastMonthPeopleFed,
	}, nil
}

func (s *OrgService) UpdateOrganization(ctx context.Context, orgID string, req *core.UpdateOrgRequest) (*core.OrgResponse, error) {
	// Get existing organization
	org, err := s.orgRepo.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.OrgType != nil {
		org.OrgType = *req.OrgType
	}
	if req.PurposeFocus != nil {
		org.PurposeFocus = req.PurposeFocus
	}
	if req.LastMonthPeopleFed != nil {
		org.LastMonthPeopleFed = *req.LastMonthPeopleFed
	}

	// Update in repository
	updatedOrg, err := s.orgRepo.UpdateOrganization(ctx, orgID, org.OrgType, org.PurposeFocus, org.LastMonthPeopleFed)
	if err != nil {
		return nil, core.NewInternalError("Failed to update organization", err)
	}

	return &core.OrgResponse{
		OrgID:              updatedOrg.ID,
		OrgType:            updatedOrg.OrgType,
		PurposeFocus:       updatedOrg.PurposeFocus,
		LastMonthPeopleFed: updatedOrg.LastMonthPeopleFed,
	}, nil
}

func (s *OrgService) ListOrganizations(ctx context.Context, page, limit string) (*core.OrgListResponse, error) {
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

	// Get organizations
	orgs, err := s.orgRepo.GetAllOrganizations(ctx, limitNum, offset)
	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve organizations", err)
	}

	// Convert to response format
	orgResponses := make([]core.OrgResponse, len(orgs))
	for i, org := range orgs {
		orgResponses[i] = core.OrgResponse{
			OrgID:              org.ID,
			OrgType:            org.OrgType,
			PurposeFocus:       org.PurposeFocus,
			LastMonthPeopleFed: org.LastMonthPeopleFed,
		}
	}

	return &core.OrgListResponse{
		Items: orgResponses,
		Page:  pageNum,
		Limit: limitNum,
		Total: len(orgs), // In a real app, this would be a separate count query
	}, nil
}
