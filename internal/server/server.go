package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	profilesv1 "github.com/vindosVP/snprofiles/gen/go"
	"github.com/vindosVP/snprofiles/internal/models"
	"github.com/vindosVP/snprofiles/internal/storage"
)

type ProfileProvider interface {
	CreateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error)
	GetProfile(ctx context.Context, userId int64) (*models.Profile, error)
	GetProfiles(ctx context.Context) ([]*models.Profile, error)
	UpdateProfile(ctx context.Context, userId int64, profile *models.UpdateProfile) (*models.Profile, error)
	SetProfilePhoto(ctx context.Context, userID int64, photoUUID *string) (*string, error)
}

type server struct {
	profilesv1.UnimplementedProfilesServer
	pp ProfileProvider
	l  zerolog.Logger
}

func Register(gRPCServer *grpc.Server, pp ProfileProvider, l zerolog.Logger) {
	profilesv1.RegisterProfilesServer(gRPCServer, &server{pp: pp, l: l})
}

func requestID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no metadata")
	}
	MDreqId := md.Get("requestID")
	if len(MDreqId) == 0 {
		return "", errors.New("no request id")
	}
	return MDreqId[0], nil
}

func (s *server) CreateProfile(ctx context.Context, in *profilesv1.CreateProfileRequest) (*profilesv1.CreateProfileResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to extract request ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	profile := models.ProfileFromGRPC(in.GetProfile())
	l := s.l.With().Str("requestID", reqId).Interface("profile", profile).Logger()
	l.Info().Msg("creating profile")
	cp, err := s.pp.CreateProfile(ctx, profile)
	if err != nil {
		if errors.Is(err, storage.ErrProfileAlreadyExist) {
			l.Info().Msg("profile already exists")
			return nil, status.Error(codes.AlreadyExists, "profile already exists")
		}
		l.Error().Stack().Err(err).Msg("failed to create profile")
		return nil, status.Error(codes.Internal, "failed to create profile")
	}
	l.Info().Msg("profile created successfully")
	return &profilesv1.CreateProfileResponse{Profile: cp.ToGRPC()}, nil
}

func (s *server) GetProfile(ctx context.Context, in *profilesv1.ProfileRequest) (*profilesv1.ProfileResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to extract request ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	l := s.l.With().Str("requestID", reqId).Int64("userId", in.GetUserId()).Logger()
	l.Info().Msg("getting profile")
	profile, err := s.pp.GetProfile(ctx, in.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrProfileDoesNotExist) {
			l.Info().Msg("profile does not exist")
			return nil, status.Error(codes.NotFound, "profile does not exist")
		}
		l.Error().Stack().Err(err).Msg("failed to get profile")
		return nil, status.Error(codes.Internal, "failed to get profile")
	}
	l.Info().Msg("profile retrieved successfully")
	return &profilesv1.ProfileResponse{Profile: profile.ToGRPC()}, nil
}

func (s *server) GetProfiles(ctx context.Context, _ *profilesv1.ProfilesRequest) (*profilesv1.ProfilesResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to extract request ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	l := s.l.With().Str("requestID", reqId).Logger()
	l.Info().Msg("getting profiles")
	profiles, err := s.pp.GetProfiles(ctx)
	if err != nil {
		l.Error().Stack().Err(err).Msg("failed to get profiles")
		return nil, status.Error(codes.Internal, "failed to get profiles")
	}
	profilesGRPC := make([]*profilesv1.Profile, 0, len(profiles))
	for _, profile := range profiles {
		profilesGRPC = append(profilesGRPC, profile.ToGRPC())
	}
	return &profilesv1.ProfilesResponse{Profiles: profilesGRPC}, nil
}

func (s *server) PutProfile(ctx context.Context, in *profilesv1.PutProfileRequest) (*profilesv1.PutProfileResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to extract request ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	uProfile := models.UpdateProfileFromGRPC(in.GetProfile())
	l := s.l.With().Str("requestID", reqId).Int64("userId", in.GetUserId()).Interface("fields", uProfile).Logger()
	l.Info().Msg("updating profile")
	cp, err := s.pp.UpdateProfile(ctx, in.GetUserId(), uProfile)
	if err != nil {
		if errors.Is(err, storage.ErrProfileDoesNotExist) {
			l.Info().Msg("profile does not exist")
			return nil, status.Error(codes.NotFound, "profile does not exist")
		}
		l.Error().Stack().Err(err).Msg("failed to update profile")
		return nil, status.Error(codes.Internal, "failed to update profile")
	}
	return &profilesv1.PutProfileResponse{Profile: cp.ToGRPC()}, nil
}

func (s *server) SetPhoto(ctx context.Context, in *profilesv1.SetPhotoRequest) (*profilesv1.SetPhotoResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to extract request ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	l := s.l.With().Str("requestID", reqId).Int64("userId", in.GetUserId()).Interface("photoUUID", in.PhotoUUID).Logger()
	l.Info().Msg("setting photo")
	cp, err := s.pp.SetProfilePhoto(ctx, in.GetUserId(), in.PhotoUUID)
	if err != nil {
		l.Error().Stack().Err(err).Msg("failed to set profile photo")
		return nil, status.Error(codes.Internal, "failed to set profile photo")
	}
	return &profilesv1.SetPhotoResponse{PhotoUUID: cp}, nil
}
