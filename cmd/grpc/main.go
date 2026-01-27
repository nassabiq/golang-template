package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	appConfig "github.com/nassabiq/golang-template/internal/shared/config"
	"github.com/nassabiq/golang-template/internal/shared/database"
	"github.com/nassabiq/golang-template/internal/shared/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/nassabiq/golang-template/internal/modules/auth/event"
	authHandler "github.com/nassabiq/golang-template/internal/modules/auth/handler"
	authRepository "github.com/nassabiq/golang-template/internal/modules/auth/repository/postgres"
	authUsecase "github.com/nassabiq/golang-template/internal/modules/auth/usecase"
	authctx "github.com/nassabiq/golang-template/internal/shared/middleware/auth"
	authpb "github.com/nassabiq/golang-template/proto/auth"

	userHandler "github.com/nassabiq/golang-template/internal/modules/user/handler"
	userRepository "github.com/nassabiq/golang-template/internal/modules/user/repository"
	userUsecase "github.com/nassabiq/golang-template/internal/modules/user/usecase"
	userpb "github.com/nassabiq/golang-template/proto/user"

	natsInfra "github.com/nassabiq/golang-template/internal/infrastructure/messaging/nats"
	"github.com/nassabiq/golang-template/internal/infrastructure/token"
)

func main() {
	// =========================
	// Load config
	// =========================
	cfg := appConfig.Load()

	// =========================
	// Initialize database
	// =========================
	db := database.NewPostgres(cfg.DatabaseUrl)
	defer db.Close()

	// =========================
	// JWT Middleware
	// =========================
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	verifier := authctx.NewJWTVerifier(jwtSecret)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			authctx.UnaryServerInterceptor(verifier),
		),
	)

	// =========================
	// Repository
	// =========================
	authRepo := authRepository.NewAuthRepository(db)
	userRepo := userRepository.NewUserRepository(db)

	// =========================
	// NATS / JetStream
	// =========================
	natsConn, err := natsInfra.NewNatsConnection(cfg.NatsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	jetStreamBus, err := natsInfra.NewJetStreamBus(natsConn)
	if err != nil {
		log.Fatalf("failed to create JetStream: %v", err)
	}

	// =========================
	// Event Publisher
	// =========================
	authEventPub := event.NewAuthPublisher(jetStreamBus)

	// =========================
	// Usecase
	// =========================
	passwordHasher := &helper.BcryptHasher{}
	uuidGen := &helper.UUIDGenerator{}
	tokenSvc := token.NewService(jwtSecret)

	authUC := authUsecase.NewAuthUsecase(authRepo, authEventPub)
	authUC.SetPasswordHasher(passwordHasher)
	authUC.SetUUIDGenerator(uuidGen)
	authUC.SetTokenService(tokenSvc)
	authUC.SetNowFunc(time.Now)

	userUC := userUsecase.NewUserUsecase(userRepo, passwordHasher)

	// =========================
	// GRPC Server
	// =========================
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// =========================
	// Handler
	// =========================
	authSrv := authHandler.NewAuthHandler(authUC)
	userSrv := userHandler.NewUserHandler(*userUC)

	authpb.RegisterAuthServiceServer(grpcServer, authSrv)
	userpb.RegisterUserServiceServer(grpcServer, userSrv)

	reflection.Register(grpcServer)

	// =========================
	// Graceful shutdown
	// =========================
	go func() {
		log.Printf("🚀 gRPC server running on :%s\n", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gRPC server...")
	grpcServer.GracefulStop()
}
