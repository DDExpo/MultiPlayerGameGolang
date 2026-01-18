package main

import (
	"fmt"
	"regexp"
	"strings"
)

func validateUsername(username string) error {
	trimmed := strings.TrimSpace(username)

	if len(trimmed) < MinUsernameLength {
		return fmt.Errorf("username too short (min %d characters)", MinUsernameLength)
	}

	if len(trimmed) > MaxUsernameLength {
		return fmt.Errorf("username too long (max %d characters)", MaxUsernameLength)
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(trimmed) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}

	return nil
}
