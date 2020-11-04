package dm

import (
	// "fmt"
	"context"
	"encoding/json"
	"net/http"

	"github.com/outer-labs/forge-api-go-client/oauth"
)

// HubAPI holds the necessary data for making calls to Forge Data Management service
type HubAPI struct {
	oauth.TwoLeggedAuth
	HubAPIPath  string
	RateLimiter HttpRequestLimiter
}

var api HubAPI

// NewHubAPIWithCredentials returns a Hub API client with default configurations
func NewHubAPIWithCredentials(ClientID, ClientSecret string, limiter HttpRequestLimiter) HubAPI {
	return HubAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/project/v1/hubs",
		limiter,
	}
}

func (api HubAPI) GetHubs(ctx context.Context) (result ForgeResponseArray, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.HubAPIPath

	return getHubs(ctx, api.RateLimiter, path, bearer.AccessToken)
}

func (api HubAPI) GetHubDetails(ctx context.Context, hubKey string) (result ForgeResponseObject, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.HubAPIPath

	return getHubDetails(ctx, api.RateLimiter, path, hubKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */

func getHubs(ctx context.Context, limiter HttpRequestLimiter, path, token string) (result ForgeResponseArray, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET", path, nil)
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

func getHubDetails(ctx context.Context, limiter HttpRequestLimiter, path, hubKey, token string) (result ForgeResponseObject, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET", path+"/"+hubKey, nil)
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
