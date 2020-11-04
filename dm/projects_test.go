package dm

import (
	"context"
	"os"
	"testing"
)

func TestProjectAPI_GetProjects(t *testing.T) {
	ctx := context.Background()

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	hubAPI := NewHubAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	testHubKey := os.Getenv("BIM_360_TEST_ACCOUNT_HUBKEY")

	t.Run("List all projects under a given hub", func(t *testing.T) {
		_, err := hubAPI.ListProjects(ctx, testHubKey)

		if err != nil {
			t.Fatalf("Failed to get project details: %s\n", err.Error())
		}
	})

	t.Run("List all projects under non-existent hub (should fail)", func(t *testing.T) {
		_, err := hubAPI.ListProjects(ctx, testHubKey+"30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting projects for non-existing hub\n")
		}
	})
}

func TestProjectAPI_GetProjectDetails(t *testing.T) {
	ctx := context.Background()

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	hubAPI := NewHubAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	testHubKey := os.Getenv("BIM_360_TEST_ACCOUNT_HUBKEY")
	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")

	t.Run("List all projects under a given hub", func(t *testing.T) {
		_, err := hubAPI.GetProjectDetails(ctx, testHubKey, testProjectKey)

		if err != nil {
			t.Fatalf("Failed to get project details: %s\n", err.Error())
		}
	})

	t.Run("List all projects under non-existent hub (should fail)", func(t *testing.T) {
		_, err := hubAPI.GetProjectDetails(ctx, testHubKey, testProjectKey+"30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting projects for non-existing hub\n")
		}
	})
}

func TestProjectAPI_GetTopFolders(t *testing.T) {
	ctx := context.Background()

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	hubAPI := NewHubAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	testHubKey := os.Getenv("BIM_360_TEST_ACCOUNT_HUBKEY")
	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")

	t.Run("List all projects under a given hub", func(t *testing.T) {
		_, err := hubAPI.GetTopFolders(ctx, testHubKey, testProjectKey)

		if err != nil {
			t.Fatalf("Failed to get project details: %s\n", err.Error())
		}
	})

	t.Run("List all projects under non-existent hub (should fail)", func(t *testing.T) {
		_, err := hubAPI.GetTopFolders(ctx, testHubKey, testProjectKey+"30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting projects for non-existing hub\n")
		}
	})
}
