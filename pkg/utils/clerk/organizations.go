package clerkutils

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/organization"
)

func CreateNewOrganization(userClerkId, newOrgName string) (string, error) {

	ctx := context.Background()
	org, err := organization.Create(ctx, &organization.CreateParams{
		Name:      clerk.String(newOrgName),
		CreatedBy: clerk.String(userClerkId),
	})
	if err != nil {
		return "", err
	}

	return org.ID, nil
}
