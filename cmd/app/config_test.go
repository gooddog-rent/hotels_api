package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {

	// load .env.test file
	err := godotenv.Load("../test.env")
	if err != nil {
		fmt.Print("Error loading test.env file. Env variables should be loaded.")
	}
}

func TestNewConfigNotEmptyData(t *testing.T) {

	cfg := NewConfig()

	if cfg.HTTP_PORT == "" || cfg.HOTELS_PATH == "" {
		t.Errorf("Config struct should not have an empty values: got %v", cfg)
	}
}

func TestNewConfigEmptyData(t *testing.T) {

	envs := []string{"HTTP_PORT", "HOTELS_PATH"}
	envsBuffer := make(map[string]string)

	// clear environments variables
	for _, env := range envs {

		// save to buffer for later restore envs
		envsBuffer[env] = os.Getenv(env)

		// clear env
		os.Setenv(env, "")
	}

	cfg := NewConfig()

	if cfg.HTTP_PORT != "" || cfg.HOTELS_PATH != "" {
		t.Errorf("Config struct should be an empty values: got %v", cfg)
	}

	// restore envs for another tests
	// loads values from .env into the system
	for k, v := range envsBuffer {
		os.Setenv(k, v)
	}
}

func TestNewConfigEnvsNotExist(t *testing.T) {

	envs := []string{"HTTP_PORT", "HOTELS_PATH"}
	envsBuffer := make(map[string]string)

	// clear environments variables
	for _, env := range envs {

		// save to buffer for later restore envs
		envsBuffer[env] = os.Getenv(env)

		// clear env
		os.Setenv(env, "")
	}

	// flush all environments
	// os.Clearenv()

	for _, env := range envs {

		// clear env
		os.Setenv(env, "")

		value := getEnv(env)
		if value != "" {
			t.Errorf("Env variable %s should not exist!", env)
		}
	}

	// restore envs for another tests
	// loads values from .env into the system
	for k, v := range envsBuffer {
		os.Setenv(k, v)
	}
}

func TestGetEnvExist(t *testing.T) {

	envs := []string{"HTTP_PORT", "HOTELS_PATH"}

	for _, env := range envs {

		value := getEnv(env)
		if value == "" {
			t.Errorf("Env variable %s does not exist!", env)
		}
	}
}

func TestGetEnvNotExist(t *testing.T) {

	envs := []string{"FAKE_HTTP_PORT", "FAKE_HOTELS_PATH"}

	for _, env := range envs {

		// clear env
		os.Setenv(env, "")

		value := getEnv(env)
		if value != "" {
			t.Errorf("Env variable %s should not exist!", env)
		}
	}
}
