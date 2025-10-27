package config

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	_ = os.Setenv("HTTP_PORT", "8080")
	_ = os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	_ = os.Setenv("KAFKA_URL", "kafka1:9092;kafka2:9092")
	_ = os.Setenv("KAFKA_REQUEST_TOPIC", "requests")
	_ = os.Setenv("KAFKA_RESPONSE_TOPIC", "responses")
	_ = os.Setenv("KAFKA_GROUP_ID", "app-group")

	defer func() {
		_ = os.Unsetenv("HTTP_PORT")
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Unsetenv("KAFKA_URL")
		_ = os.Unsetenv("KAFKA_REQUEST_TOPIC")
		_ = os.Unsetenv("KAFKA_RESPONSE_TOPIC")
		_ = os.Unsetenv("KAFKA_GROUP_ID")
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.HttpPort != "8080" {
		t.Errorf("Expected HttpPort 8080, got %s", config.HttpPort)
	}

	if config.DatabaseURL != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("Unexpected DatabaseURL: %s", config.DatabaseURL)
	}

	expectedBrokers := []string{"kafka1:9092", "kafka2:9092"}
	if !reflect.DeepEqual(config.KafkaBrokers, expectedBrokers) {
		t.Errorf("Expected brokers %v, got %v", expectedBrokers, config.KafkaBrokers)
	}

	if config.KafkaRequestTopic != "requests" {
		t.Errorf("Expected consumer topic 'requests', got %s", config.KafkaRequestTopic)
	}

	if config.KafkaResponseTopic != "responses" {
		t.Errorf("Expected producer topic 'responses', got %s", config.KafkaResponseTopic)
	}

	if config.KafkaGroupID != "app-group" {
		t.Errorf("Expected group ID 'app-group', got %s", config.KafkaGroupID)
	}
}

func TestLoadConfig_MissingRequired(t *testing.T) {
	_ = os.Setenv("HTTP_PORT", "8080")
	_ = os.Setenv("DATABASE_URL", "postgres://localhost:5432/db")
	_ = os.Unsetenv("KAFKA_URL")

	defer func() {
		_ = os.Unsetenv("HTTP_PORT")
		_ = os.Unsetenv("DATABASE_URL")
	}()

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for missing required env var, got nil")
	}

	if !strings.Contains(err.Error(), "KAFKA_URL is required") {
		t.Errorf("Error should mention missing KAFKA_URL, got: %v", err)
	}
}

func TestLoadConfig_AllMissing(t *testing.T) {
	os.Clearenv()

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for missing all env vars, got nil")
	}

	errorMsg := err.Error()
	requiredVars := []string{"HTTP_PORT", "DATABASE_URL", "KAFKA_URL", "KAFKA_REQUEST_TOPIC", "KAFKA_RESPONSE_TOPIC", "KAFKA_GROUP_ID"}

	for _, varName := range requiredVars {
		if !strings.Contains(errorMsg, varName) {
			t.Errorf("Error should mention missing %s", varName)
		}
	}
}

func TestLoadConfig_KafkaBrokersParsing(t *testing.T) {
	_ = os.Setenv("HTTP_PORT", "8080")
	_ = os.Setenv("DATABASE_URL", "postgres://localhost:5432/db")
	_ = os.Setenv("KAFKA_URL", "host1:9092;host2:9092;host3:9092")
	_ = os.Setenv("KAFKA_REQUEST_TOPIC", "requests")
	_ = os.Setenv("KAFKA_RESPONSE_TOPIC", "responses")
	_ = os.Setenv("KAFKA_GROUP_ID", "group")

	defer func() {
		_ = os.Unsetenv("HTTP_PORT")
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Unsetenv("KAFKA_URL")
		_ = os.Unsetenv("KAFKA_REQUEST_TOPIC")
		_ = os.Unsetenv("KAFKA_RESPONSE_TOPIC")
		_ = os.Unsetenv("KAFKA_GROUP_ID")
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedBrokers := []string{"host1:9092", "host2:9092", "host3:9092"}
	if !reflect.DeepEqual(config.KafkaBrokers, expectedBrokers) {
		t.Errorf("Expected brokers %v, got %v", expectedBrokers, config.KafkaBrokers)
	}
}

func TestMustGetEnv(t *testing.T) {
	_ = os.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	value, err := mustGetEnv("TEST_VAR")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != "test-value" {
		t.Errorf("Expected 'test-value', got %s", value)
	}

	_, err = mustGetEnv("NON_EXISTENT_VAR")
	if err == nil {
		t.Error("Expected error for non-existent var")
	} else if err.Error() != "NON_EXISTENT_VAR is required" {
		t.Errorf("Unexpected error message: %v", err)
	}
}
