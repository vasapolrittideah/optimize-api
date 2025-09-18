package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/vasapolrittideah/optimize-api/services/auth-service/internal/config"
	grpcHandler "github.com/vasapolrittideah/optimize-api/services/auth-service/internal/delivery/grpc"
	mongoRepo "github.com/vasapolrittideah/optimize-api/services/auth-service/internal/repository/mongo"
	"github.com/vasapolrittideah/optimize-api/services/auth-service/internal/usecase"
	"github.com/vasapolrittideah/optimize-api/shared/auth"
	"github.com/vasapolrittideah/optimize-api/shared/database"
	"github.com/vasapolrittideah/optimize-api/shared/discovery"
	"github.com/vasapolrittideah/optimize-api/shared/logger"
	"github.com/vasapolrittideah/optimize-api/shared/utilities"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger.New()

	authServiceCfg := config.NewAuthServiceConfig(logger)

	mongodb := database.NewMongoDB(logger)
	if err := mongodb.Connect(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to MongoDB")
	}
	defer func() {
		if err := mongodb.Disconnect(ctx); err != nil {
			logger.Error().Err(err).Msg("failed to disconnect from MongoDB")
		}
	}()

	consulRegistry, err := discovery.NewConsulRegistry(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create Consul registry")
	}

	serviceID := authServiceCfg.Name + "-1"
	if err := consulRegistry.Register(serviceID, authServiceCfg.Name, authServiceCfg.Address); err != nil {
		logger.Fatal().Err(err).Msg("failed to register service in Consul")
	}
	defer func() {
		if err := consulRegistry.Deregister(serviceID, authServiceCfg.Name); err != nil {
			logger.Error().Err(err).Msg("failed to deregister service from Consul")
		}
	}()

	jwtAuthenticator := auth.NewJWTAuthenticator(
		authServiceCfg.Token.Issuer,
		authServiceCfg.Token.Issuer,
	)

	identityRepo := mongoRepo.NewIdentityMongoRepository(mongodb.GetDatabase())
	sessionRepo := mongoRepo.NewSessionMongoRepository(mongodb.GetDatabase())
	userRepo := mongoRepo.NewUserMongoRepository(ctx, logger, mongodb.GetDatabase())

	authUsecase := usecase.NewAuthUsecase(identityRepo, sessionRepo, userRepo, jwtAuthenticator, authServiceCfg)

	grpcServer := grpc.NewServer()
	grpcHandler.NewAuthGRPCHandler(grpcServer, logger, authUsecase)

	utilities.RegisterHealthServer(grpcServer)

	lc := net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", authServiceCfg.Address)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create listener")
	}

	go func() {
		logger.Info().Msg("Starting gRPC server...")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal().Err(err).Msg("failed to start gRPC server")
			cancel()
		}
	}()

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan
		cancel()
	}()

	<-ctx.Done()

	logger.Info().Msg("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}
