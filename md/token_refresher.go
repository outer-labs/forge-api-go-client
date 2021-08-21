package md

import "github.com/ankitm123/forge-api-go-client/oauth"

type TokenRefresher interface {
	Bearer() *oauth.Bearer
	RefreshTokenIfRequired(auth oauth.ThreeLeggedAuth) error
}
