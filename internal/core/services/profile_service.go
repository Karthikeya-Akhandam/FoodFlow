package services

import (
	"context"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type ProfileService struct {
	profileRepo repos.ProfileRepo
}

func NewProfileService(profileRepo repos.ProfileRepo) core.ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*core.ProfileResponse, error) {
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &core.ProfileResponse{
		Name:        &profile.Name,
		Phone:       profile.Phone,
		Pincode:     profile.Pincode,
		City:        profile.City,
		District:    profile.District,
		State:       profile.State,
		Address:     profile.Address,
		IsRemoteOrg: profile.IsRemoteOrg,
	}, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userID string, req *core.UpdateProfileRequest) (*core.ProfileResponse, error) {
	// Handle nullable fields
	var name string
	if req.Name != nil {
		name = *req.Name
	}
	
	var isRemoteOrg bool
	if req.IsRemoteOrg != nil {
		isRemoteOrg = *req.IsRemoteOrg
	}
	
	profile, err := s.profileRepo.UpdateProfile(ctx, userID, name, req.Phone, req.Pincode, req.City, req.District, req.State, req.Address, isRemoteOrg)
	if err != nil {
		return nil, err
	}

	return &core.ProfileResponse{
		Name:        &profile.Name,
		Phone:       profile.Phone,
		Pincode:     profile.Pincode,
		City:        profile.City,
		District:    profile.District,
		State:       profile.State,
		Address:     profile.Address,
		IsRemoteOrg: profile.IsRemoteOrg,
	}, nil
}
