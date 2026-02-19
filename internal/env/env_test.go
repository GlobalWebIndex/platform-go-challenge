package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	const key = "TEST_ENV_STRING"

	// Ensure clean state
	os.Unsetenv(key)

	t.Run("Returns fallback when not set", func(t *testing.T) {
		val := GetString(key, "fallback")
		assert.Equal(t, "fallback", val)
	})

	t.Run("Returns value when set", func(t *testing.T) {
		os.Setenv(key, "hello")
		defer os.Unsetenv(key)

		val := GetString(key, "fallback")
		assert.Equal(t, "hello", val)
	})
}

func TestGetInt(t *testing.T) {
	const key = "TEST_ENV_INT"

	os.Unsetenv(key)

	t.Run("Returns fallback when not set", func(t *testing.T) {
		val := GetInt(key, 42)
		assert.Equal(t, 42, val)
	})

	t.Run("Returns fallback when invalid int", func(t *testing.T) {
		os.Setenv(key, "not-a-number")
		defer os.Unsetenv(key)

		val := GetInt(key, 42)
		assert.Equal(t, 42, val)
	})

	t.Run("Returns parsed int when valid", func(t *testing.T) {
		os.Setenv(key, "123")
		defer os.Unsetenv(key)

		val := GetInt(key, 42)
		assert.Equal(t, 123, val)
	})
}

func TestGetBool(t *testing.T) {
	const key = "TEST_ENV_BOOL"

	os.Unsetenv(key)

	t.Run("Returns fallback when not set", func(t *testing.T) {
		val := GetBool(key, true)
		assert.True(t, val)
	})

	t.Run("Returns fallback when invalid bool", func(t *testing.T) {
		os.Setenv(key, "not-a-bool")
		defer os.Unsetenv(key)

		val := GetBool(key, true)
		assert.True(t, val)
	})

	t.Run("Returns parsed bool when valid true", func(t *testing.T) {
		os.Setenv(key, "true")
		defer os.Unsetenv(key)

		val := GetBool(key, false)
		assert.True(t, val)
	})

	t.Run("Returns parsed bool when valid false", func(t *testing.T) {
		os.Setenv(key, "false")
		defer os.Unsetenv(key)

		val := GetBool(key, true)
		assert.False(t, val)
	})
}
