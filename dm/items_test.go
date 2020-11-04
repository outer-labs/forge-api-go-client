package dm

import (
	"context"
	"os"
	"testing"
)

func TestFolderAPI_GetItemDetails(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	folderAPI := NewFolderAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")
	testItemKey := os.Getenv("BIM_360_TEST_ACCOUNT_ITEMKEY")

	t.Run("List item details", func(t *testing.T) {
		_, err := folderAPI.GetItemDetails(context.Background(), testProjectKey, testItemKey)

		if err != nil {
			t.Fatalf("Failed to get item details: %s\n", err.Error())
		}
	})
}

func TestFolderAPI_GetItemTip(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	folderAPI := NewFolderAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")
	testItemKey := os.Getenv("BIM_360_TEST_ACCOUNT_ITEMKEY")

	t.Run("List item details", func(t *testing.T) {
		_, err := folderAPI.GetItemTip(context.Background(), testProjectKey, testItemKey)

		if err != nil {
			t.Fatalf("Failed to get item details: %s\n", err.Error())
		}
	})
}

func TestFolderAPI_GetItemVersions(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	folderAPI := NewFolderAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")
	testItemKey := os.Getenv("BIM_360_TEST_ACCOUNT_ITEMKEY")

	t.Run("List item details", func(t *testing.T) {
		_, err := folderAPI.GetItemVersions(context.Background(), testProjectKey, testItemKey)

		if err != nil {
			t.Fatalf("Failed to get item details: %s\n", err.Error())
		}
	})
}
