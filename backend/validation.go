package main

import (
	"fmt"
	"regexp"
	"strings"
)

var usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func validateUsername(username string) error {
	trimmed := strings.TrimSpace(username)

	switch {
	case len(trimmed) < MinUsernameLength:
		return fmt.Errorf("username too short (min %d characters)", MinUsernameLength)
	case len(trimmed) > MaxUsernameLength:
		return fmt.Errorf("username too long (max %d characters)", MaxUsernameLength)
	case !usernameRegexp.MatchString(trimmed):
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}

	return nil
}
