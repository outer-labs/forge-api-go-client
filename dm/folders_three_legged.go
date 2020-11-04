package dm

import (
	"context"

	"github.com/outer-labs/forge-api-go-client/oauth"
)

type FolderAPI3L struct {
	Auth          oauth.ThreeLeggedAuth
	Token         TokenRefresher
	FolderAPIPath string
	RateLimiter   HttpRequestLimiter
}

func NewFolderAPI3LWithCredentials(
	auth oauth.ThreeLeggedAuth,
	token TokenRefresher,
	limiter HttpRequestLimiter,
) *FolderAPI3L {
	return &FolderAPI3L{
		Auth:          auth,
		Token:         token,
		FolderAPIPath: "/data/v1/projects",
		RateLimiter:   limiter,
	}
}

// Three legged Folder api calls
func (a FolderAPI3L) GetFolderDetailsThreeLegged(ctx context.Context, projectKey, folderKey string) (result ForgeResponseObject, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.FolderAPIPath
	return getFolderDetails(ctx, a.RateLimiter, path, projectKey, folderKey, a.Token.Bearer().AccessToken)
}

func (a FolderAPI3L) GetFolderContentsThreeLegged(ctx context.Context, projectKey, folderKey string) (result ForgeResponseArray, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.FolderAPIPath
	return getFolderContents(ctx, a.RateLimiter, path, projectKey, folderKey, a.Token.Bearer().AccessToken)
}

func (a FolderAPI3L) GetItemDetailsThreeLegged(ctx context.Context, projectKey, itemKey string) (result ForgeResponseObject, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.FolderAPIPath

	return getItemDetails(ctx, a.RateLimiter, path, projectKey, itemKey, a.Token.Bearer().AccessToken)
}
