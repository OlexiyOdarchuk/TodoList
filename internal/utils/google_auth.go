package utils

import (
	"context"
	"os"
	"todolist/internal/models"

	"google.golang.org/api/idtoken"
)

func ValidateGoogleToken(ctx context.Context, idTokenStr string) (*models.OAuthUser, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	payload, err := idtoken.Validate(ctx, idTokenStr, clientID)
	if err != nil {
		return nil, err
	}

	return &models.OAuthUser{
		ID:       payload.Subject,
		Email:    payload.Claims["email"].(string),
		Name:     payload.Claims["name"].(string),
		Provider: "google",
	}, nil
}
