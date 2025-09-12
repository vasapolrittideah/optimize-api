package grpc

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vasapolrittideah/optimize-api/services/auth-service/internal/domain"
	"github.com/vasapolrittideah/optimize-api/services/auth-service/internal/usecase"
	authpbv1 "github.com/vasapolrittideah/optimize-api/shared/protos/auth/v1"
)

type authGRPCHandler struct {
	authpbv1.UnimplementedAuthServiceServer

	logger      *zerolog.Logger
	authUsecase domain.AuthUsecase
}

func NewAuthGRPCHandler(
	server *grpc.Server,
	logger *zerolog.Logger,
	authUsecase domain.AuthUsecase,
) authpbv1.AuthServiceServer {
	handler := &authGRPCHandler{
		logger:      logger,
		authUsecase: authUsecase,
	}
	authpbv1.RegisterAuthServiceServer(server, handler)

	return handler
}

func (h *authGRPCHandler) Login(ctx context.Context, req *authpbv1.LoginRequest) (*authpbv1.LoginResponse, error) {
	params := domain.LoginParams{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	tokens, err := h.authUsecase.Login(ctx, params)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to login")

		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		default:
			return nil, status.Errorf(codes.Internal, "something went wrong")
		}
	}

	return &authpbv1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *authGRPCHandler) SignUp(ctx context.Context, req *authpbv1.SignUpRequest) (*authpbv1.SignUpResponse, error) {
	params := domain.SignUpParams{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		FullName: req.GetFullName(),
	}

	tokens, err := h.authUsecase.SignUp(ctx, params)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to sign up")

		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		default:
			return nil, status.Errorf(codes.Internal, "something went wrong")
		}
	}

	return &authpbv1.SignUpResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
