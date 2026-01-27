package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nassabiq/golang-template/internal/infrastructure/mail"
	natsInfra "github.com/nassabiq/golang-template/internal/infrastructure/messaging/nats"
	"github.com/nassabiq/golang-template/internal/infrastructure/registry"
	"github.com/nassabiq/golang-template/internal/infrastructure/subscribers"
	appConfig "github.com/nassabiq/golang-template/internal/shared/config"
	natsgo "github.com/nats-io/nats.go"
)

func main() {
	cfg := appConfig.Load()

	// Connect to NATS
	nc, err := natsInfra.NewNatsConnection(cfg.NatsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("failed to create JetStream: %v", err)
	}

	// Create stream if not exists
	_, err = js.AddStream(&natsgo.StreamConfig{
		Name:     "AUTH",
		Subjects: []string{"auth.*"},
	})
	if err != nil {
		// Stream might already exist, log and continue
		log.Printf("stream creation note: %v", err)
	}

	// Initialize mailer
	mailer := mail.NewSMTPMailer(
		getEnv("SMTP_HOST", "localhost"),
		getEnvAsInt("SMTP_PORT", 1025),
		getEnv("SMTP_FROM", "noreply@app.local"),
	)

	// Registry
	reg := registry.New()

	reg.Register(subscribers.NewForgotPasswordSubscriber(mailer))
	// reg.Register(subscriber.NewUserRegisteredSubscriber(mailer))
	// reg.Register(subscriber.NewPasswordChangedSubscriber(mailer))

	reg.Run(js)

	log.Println("📨 Email worker running...")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down email worker...")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return fallback
}
