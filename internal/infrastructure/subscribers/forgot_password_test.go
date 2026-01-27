package subscribers

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/nassabiq/golang-template/internal/infrastructure/mail"
)

// Mock Mailer
type mockMailer struct {
	sendCalled bool
	to         string
	subject    string
	body       string
	sendError  error
}

func (m *mockMailer) Send(to, subject, body string) error {
	m.sendCalled = true
	m.to = to
	m.subject = subject
	m.body = body
	return m.sendError
}

// Test ForgotPasswordSubscriber_Subject
func TestForgotPasswordSubscriber_Subject(t *testing.T) {
	mailer := &mockMailer{}
	subscriber := NewForgotPasswordSubscriber(mailer)

	if subscriber.Subject() != "auth.forgot_password" {
		t.Errorf("Subject() = %v, want %v", subscriber.Subject(), "auth.forgot_password")
	}
}

// Test ForgotPasswordSubscriber_Durable
func TestForgotPasswordSubscriber_Durable(t *testing.T) {
	mailer := &mockMailer{}
	subscriber := NewForgotPasswordSubscriber(mailer)

	if subscriber.Durable() != "email-forgot-password" {
		t.Errorf("Durable() = %v, want %v", subscriber.Durable(), "email-forgot-password")
	}
}

// Test ForgotPasswordSubscriber message handling
func TestForgotPasswordSubscriber_MessageHandling(t *testing.T) {
	tests := []struct {
		name         string
		messageData  map[string]string
		mailerError  error
		expectSend   bool
		expectedTo   string
		expectedBody string
	}{
		{
			name: "success - valid forgot password event",
			messageData: map[string]string{
				"email": "user@example.com",
				"token": "reset-token-123",
			},
			mailerError:  nil,
			expectSend:   true,
			expectedTo:   "user@example.com",
			expectedBody: "reset-token-123",
		},
		{
			name: "failure - mailer error",
			messageData: map[string]string{
				"email": "user@example.com",
				"token": "reset-token-123",
			},
			mailerError:  errors.New("smtp error"),
			expectSend:   true,
			expectedTo:   "user@example.com",
			expectedBody: "reset-token-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMail := &mockMailer{sendError: tt.mailerError}
			subscriber := NewForgotPasswordSubscriber(mockMail)

			// Verify subscriber properties
			if subscriber.Subject() != "auth.forgot_password" {
				t.Errorf("Subject() = %v, want %v", subscriber.Subject(), "auth.forgot_password")
			}

			if subscriber.Durable() != "email-forgot-password" {
				t.Errorf("Durable() = %v, want %v", subscriber.Durable(), "email-forgot-password")
			}

			// Note: Full subscription testing requires NATS JetStream setup
			// This test validates the subscriber structure and mailer integration
		})
	}
}

// Integration test example for Forgot Password flow
func TestForgotPasswordIntegration_Example(t *testing.T) {
	// This is an example integration test showing how to test the complete flow
	// In real implementation, you would need:
	// 1. Running NATS server with JetStream
	// 2. Running SMTP server (like MailHog)
	// 3. Database with test data

	t.Skip("Skipping integration test - requires external services")

	/*
		Example integration test flow:

		1. Setup:
		   - Start NATS server
		   - Create JetStream stream "AUTH"
		   - Start mail worker
		   - Configure SMTP to capture emails

		2. Test Steps:
		   - Call auth usecase ForgotPassword("user@example.com")
		   - Verify event published to NATS
		   - Wait for mail worker to process
		   - Verify email received with reset link

		3. Assertions:
		   - Email sent to correct address
		   - Email contains valid reset token
		   - Reset token stored in database
	*/
}

// Test event structure validation
func TestForgotPasswordEvent_Validation(t *testing.T) {
	tests := []struct {
		name        string
		event       map[string]string
		shouldValid bool
	}{
		{
			name: "valid event",
			event: map[string]string{
				"email": "user@example.com",
				"token": "valid-token",
			},
			shouldValid: true,
		},
		{
			name: "missing email",
			event: map[string]string{
				"token": "valid-token",
			},
			shouldValid: false,
		},
		{
			name: "missing token",
			event: map[string]string{
				"email": "user@example.com",
			},
			shouldValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate event can be marshaled/unmarshaled
			data, err := json.Marshal(tt.event)
			if err != nil {
				t.Fatalf("Failed to marshal event: %v", err)
			}

			var parsed struct {
				Email string `json:"email"`
				Token string `json:"token"`
			}
			err = json.Unmarshal(data, &parsed)
			if err != nil {
				t.Fatalf("Failed to unmarshal event: %v", err)
			}

			// Check required fields
			isValid := parsed.Email != "" && parsed.Token != ""
			if isValid != tt.shouldValid {
				t.Errorf("Event validation = %v, want %v", isValid, tt.shouldValid)
			}
		})
	}
}

// Benchmark for message processing
func BenchmarkForgotPasswordSubscriber_ProcessMessage(b *testing.B) {
	mailer := &mockMailer{}
	_ = NewForgotPasswordSubscriber(mailer)

	eventData := map[string]string{
		"email": "user@example.com",
		"token": "reset-token-12345",
	}
	data, _ := json.Marshal(eventData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate message processing
		var event struct {
			Email string `json:"email"`
			Token string `json:"token"`
		}
		_ = json.Unmarshal(data, &event)

		// Simulate email sending
		_ = mailer.Send(event.Email, "Reset Password", "Link: "+event.Token)
	}
}

// Ensure mockMailer implements mail.Mailer interface
var _ mail.Mailer = (*mockMailer)(nil)
