package dm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/outer-labs/forge-api-go-client/oauth"
)

// BucketAPI holds the necessary data for making Bucket related calls to Forge Data Management service
type BucketAPI struct {
	oauth.TwoLeggedAuth
	BucketAPIPath string
	RateLimiter   HttpRequestLimiter
}

// NewBucketAPIWithCredentials returns a Bucket API client with default configurations
func NewBucketAPIWithCredentials(ClientID string, ClientSecret string, limiter HttpRequestLimiter) BucketAPI {
	return BucketAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/oss/v2/buckets",
		limiter,
	}
}

// CreateBucketRequest contains the data necessary to be passed upon bucket creation
type CreateBucketRequest struct {
	BucketKey string `json:"bucketKey"`
	PolicyKey string `json:"policyKey"`
}

// BucketDetails reflects the body content received upon creation of a bucket
type BucketDetails struct {
	BucketKey   string `json:"bucketKey"`
	BucketOwner string `json:"bucketOwner"`
	CreateDate  string `json:"createDate"`
	Permissions []struct {
		AuthID string `json:"authId"`
		Access string `json:"access"`
	} `json:"permissions"`
	PolicyKey string `json:"policyKey"`
}

// ErrorResult reflects the body content when a request failed (g.e. Bad request or key conflict)
type ErrorResult struct {
	Reason     string `json:"reason"`
	StatusCode int
}

func (e *ErrorResult) Error() string {
	return "[" + strconv.Itoa(e.StatusCode) + "] " + e.Reason
}

// ListedBuckets reflects the response when query Data Management API for buckets associated with current Forge secrets.
type ListedBuckets struct {
	Items []struct {
		BucketKey   string `json:"bucketKey"`
		CreatedDate uint64 `json:"createdDate"`
		PolicyKey   string `json:"policyKey"`
	} `json:"items"`
	Next string `json:"next"`
}

// CreateBucket creates and returns details of created bucket, or an error on failure
func (api BucketAPI) CreateBucket(ctx context.Context, bucketKey, policyKey string) (result BucketDetails, err error) {
	bearer, err := api.Authenticate("bucket:create")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath
	result, err = createBucket(ctx, api.RateLimiter, path, bucketKey, policyKey, bearer.AccessToken)

	return
}

// DeleteBucket deletes bucket given its key.
// 	WARNING: The bucket delete call is undocumented.
func (api BucketAPI) DeleteBucket(ctx context.Context, bucketKey string) error {
	bearer, err := api.Authenticate("bucket:delete")
	if err != nil {
		return err
	}
	path := api.Host + api.BucketAPIPath

	return deleteBucket(ctx, api.RateLimiter, path, bucketKey, bearer.AccessToken)
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api BucketAPI) ListBuckets(ctx context.Context, region, limit, startAt string) (result ListedBuckets, err error) {
	bearer, err := api.Authenticate("bucket:read")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath

	return listBuckets(ctx, api.RateLimiter, path, region, limit, startAt, bearer.AccessToken)
}

// GetBucketDetails returns information associated to a bucket. See BucketDetails struct.
func (api BucketAPI) GetBucketDetails(ctx context.Context, bucketKey string) (result BucketDetails, err error) {
	bearer, err := api.Authenticate("bucket:read")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath

	return getBucketDetails(ctx, api.RateLimiter, path, bucketKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */
func getBucketDetails(ctx context.Context, limiter HttpRequestLimiter, path, bucketKey, token string) (result BucketDetails, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET",
		path+"/"+bucketKey+"/details",
		nil,
	)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		err = &ErrorResult{StatusCode: response.StatusCode}
		decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}

func listBuckets(ctx context.Context, limiter HttpRequestLimiter, path, region, limit, startAt, token string) (result ListedBuckets, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET",
		path,
		nil,
	)

	if err != nil {
		return
	}

	params := req.URL.Query()
	if len(region) != 0 {
		params.Add("region", region)
	}
	if len(limit) != 0 {
		params.Add("limit", limit)
	}
	if len(startAt) != 0 {
		params.Add("startAt", startAt)
	}

	req.URL.RawQuery = params.Encode()

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		err = &ErrorResult{StatusCode: response.StatusCode}
		decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}

func createBucket(ctx context.Context, limiter HttpRequestLimiter, path, bucketKey, policyKey, token string) (result BucketDetails, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		CreateBucketRequest{
			bucketKey,
			policyKey,
		})
	if err != nil {
		return
	}

	req, err := limiter.HttpRequest(ctx, "POST",
		path,
		bytes.NewReader(body),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		err = &ErrorResult{StatusCode: response.StatusCode}
		decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}

func deleteBucket(ctx context.Context, limiter HttpRequestLimiter, path, bucketKey, token string) (err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "DELETE",
		path+"/"+bucketKey,
		nil,
	)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		err = &ErrorResult{StatusCode: response.StatusCode}
		decoder.Decode(err)
		return
	}

	return
}
