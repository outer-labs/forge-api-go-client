package dm

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestBucketAPI_CreateBucket(t *testing.T) {
	ctx := context.Background()

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	bucketAPI := NewBucketAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	t.Run("Create a bucket", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket(ctx, "go_testing_bucket", "transient")

		if err != nil {
			t.Fatalf("Failed to create a bucket: %s\n", err.Error())
		}
	})

	t.Run("Delete created bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(ctx, "go_testing_bucket")

		if err != nil {
			t.Fatalf("Failed to delete bucket: %s\n", err.Error())
		}
	})

	t.Run("Create a bucket with invalid name", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket(ctx, "goTestingBucket", "transient")

		if err == nil {
			t.Fatalf("Should fail creating a bucket with invalid name\n")
		}
	})

	t.Run("Create a bucket with invalid policyKey", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket(ctx, "goTestingBucket", "democracy")

		if err == nil {
			t.Fatalf("Should fail creating a bucket with invalid name\n")
		}
	})
}

func TestBucketAPI_GetBucketDetails(t *testing.T) {
	ctx := context.Background()

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	bucketAPI := NewBucketAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	testBucketKey := "my_test_bucket_key_for_go"

	t.Run("Create a bucket", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket(ctx, testBucketKey, "transient")

		if err != nil {
			t.Fatalf("Failed to create a bucket: %s\n", err.Error())
		}
	})

	t.Run("Get bucket details", func(t *testing.T) {
		_, err := bucketAPI.GetBucketDetails(ctx, testBucketKey)

		if err != nil {
			t.Fatalf("Failed to get bucket details: %s\n", err.Error())
		}
	})

	t.Run("Delete created bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(ctx, testBucketKey)

		if err != nil {
			t.Fatalf("Failed to delete bucket: %s\n", err.Error())
		}
	})

	t.Run("Get nonexistent bucket", func(t *testing.T) {
		_, err := bucketAPI.GetBucketDetails(ctx, testBucketKey+"30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting details for non-existing bucket\n")
		}
	})
}

func TestBucketAPI_ListBuckets(t *testing.T) {
	ctx := context.Background()

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	bucketAPI := NewBucketAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	t.Run("List available buckets", func(t *testing.T) {
		_, err := bucketAPI.ListBuckets(ctx, "", "", "")

		if err != nil {
			t.Fatalf("Failed to list buckets: %s\n", err.Error())
		}
	})

	t.Run("Create a bucket and find it among listed", func(t *testing.T) {

		testBucketKey := "just_for_testing"

		_, err := bucketAPI.CreateBucket(ctx, testBucketKey, "transient")

		if err != nil {
			t.Errorf("Failed to create a bucket: %s\n", err.Error())
		}

		list, err := bucketAPI.ListBuckets(ctx, "", "", "")

		if err != nil {
			t.Errorf("Failed to list buckets: %s\n", err.Error())
		}

		found := false

		for _, bucket := range list.Items {

			if bucket.BucketKey == testBucketKey {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Could not find the %s bucket\n", testBucketKey)
		}

		if err = bucketAPI.DeleteBucket(ctx, testBucketKey); err != nil {
			t.Errorf("Failed to delete bucket: %s\n", err.Error())
		}
	})

}

func ExampleBucketAPI_CreateBucket() {
	ctx := context.Background()

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	bucketAPI := NewBucketAPIWithCredentials(clientID, clientSecret, DefaultRateLimiter)

	bucket, err := bucketAPI.CreateBucket(ctx, "some_unique_name", "transient")

	if err != nil {
		// handle error
	}

	fmt.Printf("Bucket %s was created with policy %s\n",
		bucket.BucketKey,
		bucket.PolicyKey)

}
