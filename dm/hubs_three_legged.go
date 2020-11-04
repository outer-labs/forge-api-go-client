package dm

import (
	"context"

	"github.com/outer-labs/forge-api-go-client/oauth"
)

type HubAPI3L struct {
	Auth        oauth.ThreeLeggedAuth
	Token       TokenRefresher
	HubAPIPath  string
	RateLimiter HttpRequestLimiter
}

func NewHubAPI3LWithCredentials(
	auth oauth.ThreeLeggedAuth,
	token TokenRefresher,
	limiter HttpRequestLimiter,
) *HubAPI3L {
	return &HubAPI3L{
		Auth:        auth,
		Token:       token,
		HubAPIPath:  "/project/v1/hubs",
		RateLimiter: limiter,
	}
}

// Hub functions for use with 3legged authentication
func (a *HubAPI3L) GetHubsThreeLegged(ctx context.Context) (result ForgeResponseArray, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getHubs(ctx, a.RateLimiter, path, a.Token.Bearer().AccessToken)
}

func (a *HubAPI3L) GetHubDetailsThreeLegged(ctx context.Context, hubKey string) (result ForgeResponseObject, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getHubDetails(ctx, a.RateLimiter, path, hubKey, a.Token.Bearer().AccessToken)
}

func (a *HubAPI3L) ListProjectsThreeLegged(ctx context.Context, hubKey string) (result ForgeResponseArray, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return listProjects(ctx, a.RateLimiter, path, hubKey, "", "", "", "", a.Token.Bearer().AccessToken)
}

func (a *HubAPI3L) GetProjectDetailsThreeLegged(ctx context.Context, hubKey, projectKey string) (result ForgeResponseObject, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getProjectDetails(ctx, a.RateLimiter, path, hubKey, projectKey, a.Token.Bearer().AccessToken)
}

func (a *HubAPI3L) GetTopFoldersThreeLegged(ctx context.Context, hubKey, projectKey string) (result ForgeResponseArray, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getTopFolders(ctx, a.RateLimiter, path, hubKey, projectKey, a.Token.Bearer().AccessToken)
}
