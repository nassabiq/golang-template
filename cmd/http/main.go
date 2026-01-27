package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	httpmw "github.com/nassabiq/golang-template/cmd/http/middleware"
	"github.com/nassabiq/golang-template/internal/infrastructure/swagger"
	authpb "github.com/nassabiq/golang-template/proto/auth"
	userpb "github.com/nassabiq/golang-template/proto/user"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get ports from environment
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "8081"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	grpcAddr := fmt.Sprintf("localhost:%s", grpcPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if key == "Authorization" {
				return key, true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)

	err := authpb.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		mux,
		grpcAddr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = userpb.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		mux,
		grpcAddr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create main HTTP mux
	mainMux := http.NewServeMux()

	// Register Swagger UI handler
	swaggerFiles := map[string]string{
		"auth": "docs/swagger/proto/auth/auth.swagger.json",
		"user": "docs/swagger/proto/user/user.swagger.json",
	}
	mainMux.Handle("/swagger/", http.StripPrefix("/swagger", swagger.MultiSwaggerHandler(swaggerFiles)))

	// Register gRPC gateway handler for all other routes
	mainMux.Handle("/", mux)

	handler := httpmw.CORS(httpmw.Logging(mainMux))

	server := &http.Server{
		Addr:    ":" + httpPort,
		Handler: handler,
	}

	go func() {
		log.Printf("🌐 HTTP Gateway running on :%s", httpPort)
		log.Printf("📚 Swagger UI available at http://localhost:%s/swagger/", httpPort)
		log.Printf("🔗 Connected to gRPC server at %s", grpcAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down HTTP Gateway...")
	_ = server.Shutdown(ctx)
}
