package dm

import (
	"context"

	"github.com/outer-labs/forge-api-go-client/oauth"
)

// BucketAPI holds the necessary data for making Bucket related calls to Forge Data Management service
type BucketAPI3L struct {
	Auth           oauth.ThreeLeggedAuth
	Token          TokenRefresher
	BucketsAPIPath string
	RateLimiter    HttpRequestLimiter
}

// NewBucketAPIWithCredentials returns a Bucket API client with default configurations
func NewBucketAPI3LWithCredentials(auth oauth.ThreeLeggedAuth, token TokenRefresher, limiter HttpRequestLimiter) *BucketAPI3L {
	return &BucketAPI3L{
		Auth:           auth,
		Token:          token,
		BucketsAPIPath: "/oss/v2/buckets",
		RateLimiter:    limiter,
	}
}

// CreateBucket creates and returns details of created bucket, or an error on failure
func (api BucketAPI3L) CreateBucket3L(ctx context.Context, bucketKey, policyKey string) (result BucketDetails, err error) {
	if err = api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return
	}

	path := api.Auth.Host + api.BucketsAPIPath
	result, err = createBucket(ctx, api.RateLimiter, path, bucketKey, policyKey, api.Token.Bearer().AccessToken)

	return
}

// DeleteBucket deletes bucket given its key.
// 	WARNING: The bucket delete call is undocumented.
func (api BucketAPI3L) DeleteBucket3L(ctx context.Context, bucketKey string) error {
	if err := api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return err
	}

	path := api.Auth.Host + api.BucketsAPIPath

	return deleteBucket(ctx, api.RateLimiter, path, bucketKey, api.Token.Bearer().AccessToken)
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api BucketAPI3L) ListBuckets3L(ctx context.Context, region, limit, startAt string) (result ListedBuckets, err error) {
	if err = api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return
	}

	path := api.Auth.Host + api.BucketsAPIPath

	return listBuckets(ctx, api.RateLimiter, path, region, limit, startAt, api.Token.Bearer().AccessToken)
}

// GetBucketDetails returns information associated to a bucket. See BucketDetails struct.
func (api BucketAPI3L) GetBucketDetails3L(ctx context.Context, bucketKey string) (result BucketDetails, err error) {
	if err = api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return
	}

	path := api.Auth.Host + api.BucketsAPIPath
	return getBucketDetails(ctx, api.RateLimiter, path, bucketKey, api.Token.Bearer().AccessToken)
}
