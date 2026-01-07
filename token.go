package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// TokenSource indicates where the API token was obtained from.
type TokenSource int

const (
	TokenSourceNone TokenSource = iota
	TokenSourceFlag
	TokenSourceEnv
	TokenSourceCached
	TokenSourceManual
)

// String returns a human-readable description of the token source.
func (s TokenSource) String() string {
	switch s {
	case TokenSourceNone:
		return "no source"
	case TokenSourceFlag:
		return "command-line flag (-token)"
	case TokenSourceEnv:
		return "environment variable (YNAB_API_TOKEN)"
	case TokenSourceCached:
		return "cached token file"
	case TokenSourceManual:
		return "manual entry"
	default:
		return "unknown source"
	}
}

const (
	tokenFileName = "ynab-api-token"
	appDirName    = "ynab-export"
)

// getTokenCachePath returns the path to the token cache file.
// It prefers UserCacheDir, falling back to UserConfigDir if needed.
func getTokenCachePath() (string, error) {
	// Try cache directory first (preferred for credentials)
	cacheDir, err := os.UserCacheDir()
	if err == nil {
		return filepath.Join(cacheDir, appDirName, tokenFileName), nil
	}

	// Fall back to config directory
	configDir, configErr := os.UserConfigDir()
	if configErr != nil {
		return "", fmt.Errorf("failed to get config directory: %w", configErr)
	}

	return filepath.Join(configDir, appDirName, tokenFileName), nil
}

// LoadCachedToken attempts to load the API token from the cache file.
// Returns the token and nil error on success.
// Returns empty string and nil if no cached token exists.
// Returns empty string and an error if there was a problem reading the file.
func LoadCachedToken() (string, error) {
	tokenPath, err := getTokenCachePath()
	if err != nil {
		return "", fmt.Errorf("failed to determine cache path: %w", err)
	}

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No cached token, not an error
		}
		return "", fmt.Errorf("failed to read cached token from %s: %w", tokenPath, err)
	}

	if len(data) == 0 {
		return "", nil // Empty file, treat as no token
	}

	return string(data), nil
}

// SaveCachedToken saves the API token to the cache file.
func SaveCachedToken(token string) error {
	tokenPath, err := getTokenCachePath()
	if err != nil {
		return err
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(tokenPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Write the token with restricted permissions (owner read/write only)
	if err := os.WriteFile(tokenPath, []byte(token), 0o600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// DeleteCachedToken removes the cached token file.
func DeleteCachedToken() error {
	tokenPath, err := getTokenCachePath()
	if err != nil {
		return err
	}

	err = os.Remove(tokenPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to delete token file: %w", err)
	}
	return nil
}

// GetTokenCacheLocation returns the path where the token is/would be cached.
// Useful for displaying to the user.
func GetTokenCacheLocation() string {
	tokenPath, err := getTokenCachePath()
	if err != nil {
		return "(unable to determine cache location)"
	}
	return tokenPath
}
