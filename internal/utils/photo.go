package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

func ValidateBase64Photo(photo string) error {
	if photo == "" {
		return nil
	}
	parts := strings.Split(photo, ",")
	if len(parts) != 2 {
		return errors.New("invalid base64 photo format")
	}
	return nil
}

func DecodeBase64Photo(photo string) ([]byte, error) {
	if err := ValidateBase64Photo(photo); err != nil {
		return nil, err
	}
	if photo == "" {
		return nil, nil
	}
	parts := strings.Split(photo, ",")
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	return data, nil
}
