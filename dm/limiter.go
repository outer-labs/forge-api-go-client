package dm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/time/rate"
)

var DefaultRateLimiter = NewRateLimiter(&DefaultDataManagementLimits, &DefaultOSSLimiter, DefaultFallbackLimiter)

type ApiEndpoints map[string]map[*regexp.Regexp]*rate.Limiter

type OSSLimiter struct {
	matcher *regexp.Regexp
	limiter *rate.Limiter
}

type RateLimiter struct {
	dm       *ApiEndpoints
	oss      *OSSLimiter
	fallback *rate.Limiter
}

var DefaultDataManagementLimits = ApiEndpoints{
	"GET": {
		// Hub endpoints
		apiUrlRegexp(`hubs$`):          limitPerMinute(50),
		apiUrlRegexp(`hubs/{hub_id}$`): limitPerMinute(50),

		// Project endpoints
		apiUrlRegexp(`hubs/{hub_id}/projects\/?(\?.*)?$`):               limitPerMinute(50),
		apiUrlRegexp(`hubs/{hub_id}/projects/{project_id}$`):            limitPerMinute(50),
		apiUrlRegexp(`hubs/{hub_id}/projects/{project_id}/hub$`):        limitPerMinute(50),
		apiUrlRegexp(`hubs/{hub_id}/projects/{project_id}/topFolders$`): limitPerMinute(300),
		apiUrlRegexp(`projects/{project_id}/downloads/{download_id}$`):  limitPerMinute(300),
		apiUrlRegexp(`projects/{project_id}/jobs/{job_id}$`):            limitPerMinute(300),

		// Folder endpoints
		apiUrlRegexp(`projects/{project_id}/folders/{folder_id}$`):                     limitPerMinute(300),
		apiUrlRegexp(`projects/{project_id}/folders/{folder_id}/contents$`):            limitPerMinute(300),
		apiUrlRegexp(`projects/{project_id}/folders/{folder_id}/parent$`):              limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/folders/{folder_id}/refs$`):                limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/folders/{folder_id}/relationships/links$`): limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/folders/{folder_id}/relationships/refs$`):  limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/folders/{folder_id}/search$`):              limitPerMinute(300),

		// Item endpoints
		apiUrlRegexp(`projects/{project_id}/items/{item_id}$`):                     limitPerMinute(300),
		apiUrlRegexp(`projects/{project_id}/items/{item_id}/parent$`):              limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/items/{item_id}/refs$`):                limitPerMinute(300),
		apiUrlRegexp(`projects/{project_id}/items/{item_id}/relationships/refs$`):  limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/items/{item_id}/relationships/links$`): limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/items/{item_id}/tip$`):                 limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/items/{item_id}/versions$`):            limitPerMinute(800),

		// Version endpoints
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}$`):                     limitPerMinute(300),
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/downloadFormats$`):     limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/downloads$`):           limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/item$`):                limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/refs$`):                limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/relationships/links$`): limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/relationships/refs$`):  limitPerMinute(50),
	},
	"POST": {
		// Project endpoints
		apiUrlRegexp(`projects/{project_id}/downloads$`): limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/storage$`):   limitPerMinute(300),

		// Folder endpoints
		apiUrlRegexp(`projects/{project_id}/folders$`):                                limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/folders/{folder_id}/relationships/refs$`): limitPerMinute(50),

		// Item endpoints
		apiUrlRegexp(`projects/{project_id}/items$`):                              limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/items/{item_id}/relationships/refs$`): limitPerMinute(50),

		// Version endpoints
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/relationships/refs$`):  limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/relationships/links$`): limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/relationships/links$`): limitPerMinute(50),

		// Command endpoints
		apiUrlRegexp(`projects/{project_id}/commands$`): limitPerMinute(300),
	},
	"PATCH": {
		// Folder endpoints
		apiUrlRegexp(`projects/{project_id}/folders/{folder_id}$`): limitPerMinute(50),

		// Item endpoints
		apiUrlRegexp(`projects/{project_id}/items/{item_id}$`): limitPerMinute(50),

		// Version endpoints
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}$`):                               limitPerMinute(50),
		apiUrlRegexp(`projects/{project_id}/versions/{version_id}/relationships/links/{link_id}$`): limitPerMinute(50),
	},
}

var DefaultOSSLimiter = OSSLimiter{
	matcher: regexp.MustCompile(`^https?://developer.api.autodesk.com/oss/v2`),
	limiter: limitPerMinute(1000),
}

var DefaultFallbackLimiter = limitPerMinute(50)

func NewRateLimiter(endpoints *ApiEndpoints, oss *OSSLimiter, fallback *rate.Limiter) *RateLimiter {
	return &RateLimiter{
		dm:       endpoints,
		oss:      oss,
		fallback: fallback,
	}
}

func (r *RateLimiter) HttpRequest(
	ctx context.Context,
	method string,
	url string,
	body io.Reader,
) (*http.Request, error) {
	if err := r.limiter(method, url).Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	return http.NewRequest(method, url, body)
}

func (r *RateLimiter) limiter(method, url string) *rate.Limiter {
	if r.oss.matcher.MatchString(url) {
		return r.oss.limiter
	}

	set, ok := (*r.dm)[method]
	if !ok {
		return r.fallback
	}

	for k, v := range set {
		if k.MatchString(url) {
			return v
		}
	}

	return r.fallback
}

var variableToRegexp = regexp.MustCompile("{.+}")

func apiUrlRegexp(stub string) *regexp.Regexp {
	replaced := variableToRegexp.ReplaceAllString(stub, ".+")
	return regexp.MustCompile("^https?://developer.api.autodesk.com/data/v(1|2)/" + replaced)
}

func limitPerMinute(r time.Duration) *rate.Limiter {
	return rate.NewLimiter(rate.Every(time.Minute/r), 1)
}
