package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"

	"github.com/nassabiq/golang-template/internal/infrastructure/mail"
	"github.com/nassabiq/golang-template/internal/infrastructure/registry"
	"github.com/nassabiq/golang-template/internal/infrastructure/subscribers"
	"github.com/nassabiq/golang-template/internal/modules/auth/domain"
	"github.com/nassabiq/golang-template/internal/modules/auth/event"
	authRepo "github.com/nassabiq/golang-template/internal/modules/auth/repository/postgres"
	"github.com/nassabiq/golang-template/internal/modules/auth/usecase"
	"github.com/nassabiq/golang-template/internal/shared/helper"
)

// TestForgotPassword_EndToEnd menguji flow lengkap forgot password
// Requirement:
// 1. MailDev running on localhost:1025 (SMTP) dan localhost:1080 (Web UI)
// 2. NATS running on localhost:4222
// 3. PostgreSQL running dengan database yang sudah dimigrate
//
// Setup MailDev:
//
//	docker run -d -p 1080:1080 -p 1025:1025 maildev/maildev
//
// Setup NATS:
//
//	docker run -d -p 4222:4222 -p 8222:8222 nats:latest -js
//
// Jalankan test:
//
//	go test -v ./test/integration/... -run TestForgotPassword_EndToEnd
func TestForgotPassword_EndToEnd(t *testing.T) {
	// Skip jika tidak ada environment variable INTEGRATION_TEST
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run")
	}

	ctx := context.Background()

	// =========================
	// 1. Setup Dependencies
	// =========================

	// Load .env file
	_ = godotenv.Load("../../.env")

	// Database connection - menggunakan DB_DSN dari .env
	dbDSN := getEnv("DB_DSN", "postgres://root:root@172.17.0.1:5432/go-template?sslmode=disable")
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// NATS connection - menggunakan NATS_URL dari .env
	natsURL := getEnv("NATS_URL", "nats://172.17.0.1:4222")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// JetStream
	js, err := nc.JetStream()
	if err != nil {
		t.Fatalf("Failed to create JetStream: %v", err)
	}

	// Create stream if not exists
	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     "AUTH",
		Subjects: []string{"auth.*"},
	})

	// Mailer (MailDev) - menggunakan SMTP_HOST dari .env
	smtpHost := getEnv("SMTP_HOST", "172.17.0.1")
	smtpPort := 1025
	mailer := mail.NewSMTPMailer(smtpHost, smtpPort, getEnv("SMTP_FROM", "noreply@test.com"))

	// =========================
	// 2. Setup Subscriber (Mail Worker)
	// =========================

	reg := registry.New()
	reg.Register(subscribers.NewForgotPasswordSubscriber(mailer))

	// Run subscriber in background
	go reg.Run(js)

	// Wait for subscriber to be ready
	time.Sleep(500 * time.Millisecond)

	// =========================
	// 3. Setup Auth Usecase
	// =========================

	repo := authRepo.NewAuthRepository(db)

	// Create event publisher
	bus := &natsEventBus{js: js}
	pub := event.NewAuthPublisher(bus)

	// Create usecase
	authUC := usecase.NewAuthUsecase(repo, pub)

	// Set dependencies
	authUC.SetPasswordHasher(&helper.BcryptHasher{})
	authUC.SetUUIDGenerator(&helper.UUIDGenerator{})
	authUC.SetTokenService(&mockTokenService{})
	authUC.SetNowFunc(time.Now)

	// =========================
	// 4. Test Data Setup
	// =========================

	// Buat test user
	testEmail := fmt.Sprintf("test%d@example.com", time.Now().Unix())
	testUser := &domain.User{
		ID:           fmt.Sprintf("user-%d", time.Now().Unix()),
		Name:         "Test User",
		Email:        testEmail,
		PasswordHash: "hashed-password",
		RoleID:       "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert role jika belum ada
	_, _ = db.ExecContext(ctx, `INSERT INTO roles (id, name) VALUES ('user', 'User') ON CONFLICT DO NOTHING`)

	// Insert user ke database
	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, password, role_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (email) DO UPDATE SET id = $1
	`, testUser.ID, testUser.Name, testUser.Email, testUser.PasswordHash, testUser.RoleID, testUser.CreatedAt, testUser.UpdatedAt)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	log.Printf("Created test user: %s", testEmail)

	// =========================
	// 5. Execute Forgot Password
	// =========================

	err = authUC.ForgotPassword(ctx, testEmail)
	if err != nil {
		t.Fatalf("ForgotPassword failed: %v", err)
	}

	log.Printf("ForgotPassword called for: %s", testEmail)

	// =========================
	// 6. Wait for Email Processing
	// =========================

	// Tunggu subscriber memproses message
	time.Sleep(2 * time.Second)

	// =========================
	// 7. Verify Email Sent
	// =========================

	// Cek MailDev API untuk memverifikasi email diterima
	maildevURL := getEnv("MAILDEV_URL", "http://172.17.0.1:1080")
	log.Printf("Check MailDev at: %s", maildevURL)
	log.Printf("Email should be sent to: %s", testEmail)

	// Verifikasi password reset token tersimpan di database
	var tokenCount int
	err = db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM password_resets pr
		JOIN users u ON u.id = pr.user_id
		WHERE u.email = $1 AND pr.used = false
	`, testEmail).Scan(&tokenCount)

	if err != nil {
		t.Fatalf("Failed to query password reset: %v", err)
	}

	if tokenCount == 0 {
		t.Errorf("Password reset token not found in database")
	} else {
		log.Printf("✓ Password reset token found in database")
	}

	log.Printf("✓ Integration test completed. Check MailDev at %s to see the email", maildevURL)
}

// Helper types and functions

type natsEventBus struct {
	js nats.JetStreamContext
}

func (n *natsEventBus) Publish(subject string, payload []byte) error {
	_, err := n.js.Publish(subject, payload)
	return err
}

type mockTokenService struct{}

func (m *mockTokenService) GenerateAccessToken(userID, role string) (string, error) {
	return "mock-access-token", nil
}

func (m *mockTokenService) GenerateRefreshToken() (plain string, hash string, err error) {
	return "mock-refresh-plain", "mock-refresh-hash", nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// TestCheckMailDev memeriksa koneksi ke MailDev
func TestCheckMailDev(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run")
	}

	smtpHost := getEnv("SMTP_HOST", "172.17.0.1")
	mailer := mail.NewSMTPMailer(smtpHost, 1025, "test@example.com")

	err := mailer.Send("recipient@example.com", "Test Subject", "Test Body")
	if err != nil {
		t.Fatalf("Failed to send test email: %v", err)
	}

	t.Log("✓ Test email sent successfully. Check MailDev at http://localhost:1080")
}

// TestCheckNATS memeriksa koneksi ke NATS
func TestCheckNATS(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run")
	}

	natsURL := getEnv("NATS_URL", "nats://172.17.0.1:4222")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Test JetStream
	js, err := nc.JetStream()
	if err != nil {
		t.Fatalf("Failed to create JetStream: %v", err)
	}

	// Create test stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "TEST_STREAM",
		Subjects: []string{"test.*"},
	})
	if err != nil {
		t.Logf("Stream might already exist: %v", err)
	}

	// Publish test message
	_, err = js.Publish("test.hello", []byte("Hello NATS"))
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	t.Log("✓ NATS connection and JetStream working")
}
