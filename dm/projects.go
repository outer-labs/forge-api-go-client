package dm

import (
	"context"
	"encoding/json"
	"net/http"
)

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api HubAPI) ListProjects(ctx context.Context, hubKey string) (result ForgeResponseArray, err error) {

	// TO DO: take in optional arguments for query params: id, ext, page, limit
	// https://forge.autodesk.com/en/docs/data/v2/reference/http/hubs-hub_id-projects-GET/
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}

	path := api.Host + api.HubAPIPath

	return listProjects(ctx, api.RateLimiter, path, hubKey, "", "", "", "", bearer.AccessToken)
}

func (api HubAPI) GetProjectDetails(ctx context.Context, hubKey, projectKey string) (result ForgeResponseObject, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.HubAPIPath

	return getProjectDetails(ctx, api.RateLimiter, path, hubKey, projectKey, bearer.AccessToken)
}

func (api HubAPI) GetTopFolders(ctx context.Context, hubKey, projectKey string) (result ForgeResponseArray, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.HubAPIPath

	return getTopFolders(ctx, api.RateLimiter, path, hubKey, projectKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */
func listProjects(ctx context.Context, limiter HttpRequestLimiter, path, hubKey, id, extension, page, limit string, token string) (result ForgeResponseArray, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET",
		path+"/"+hubKey+"/projects",
		nil,
	)

	if err != nil {
		return
	}

	params := req.URL.Query()
	if len(id) != 0 {
		params.Add("filter[id]", id)
	}
	if len(extension) != 0 {
		params.Add("filter[extension.type]", extension)
	}
	if len(page) != 0 {
		params.Add("page[number]", page)
	}
	if len(limit) != 0 {
		params.Add("page[limit]", limit)
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
		_ = decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}

func getProjectDetails(ctx context.Context, limiter HttpRequestLimiter, path, hubKey, projectKey, token string) (result ForgeResponseObject, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET",
		path+"/"+hubKey+"/projects/"+projectKey,
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
		_ = decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}

func getTopFolders(ctx context.Context, limiter HttpRequestLimiter, path, hubKey, projectKey, token string) (result ForgeResponseArray, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET",
		path+"/"+hubKey+"/projects/"+projectKey+"/topFolders",
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
		_ = decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}
