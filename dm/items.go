package dm

import (
	"context"
	"encoding/json"
	"net/http"
)

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api FolderAPI) GetItemDetails(ctx context.Context, projectKey, itemKey string) (result ForgeResponseObject, err error) {

	// TO DO: take in optional header argument
	// https://forge.autodesk.com/en/docs/data/v2/reference/http/projects-project_id-items-item_id-GET/
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}

	path := api.Host + api.FolderAPIPath

	return getItemDetails(ctx, api.RateLimiter, path, projectKey, itemKey, bearer.AccessToken)
}

func (api FolderAPI) GetItemTip(ctx context.Context, projectKey, itemKey string) (result ForgeResponseObject, err error) {

	// TO DO: take in optional header argument
	// https://forge.autodesk.com/en/docs/data/v2/reference/http/projects-project_id-items-item_id-GET/
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}

	path := api.Host + api.FolderAPIPath

	return getItemDetails(ctx, api.RateLimiter, path, projectKey, itemKey, bearer.AccessToken)
}

func (api FolderAPI) GetItemVersions(ctx context.Context, projectKey, itemKey string) (result ForgeResponseArray, err error) {

	// TO DO: take in optional header argument
	// https://forge.autodesk.com/en/docs/data/v2/reference/http/projects-project_id-items-item_id-GET/
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}

	path := api.Host + api.FolderAPIPath

	return getItemVersions(ctx, api.RateLimiter, path, projectKey, itemKey, "", "", "", "", "", "", bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */
func getItemDetails(ctx context.Context, limiter HttpRequestLimiter, path, projectKey, itemKey, token string) (result ForgeResponseObject, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET", path+"/"+projectKey+"/items/"+itemKey, nil)
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

func getItemTip(ctx context.Context, limiter HttpRequestLimiter, path, projectKey, itemKey, token string) (result ForgeResponseObject, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET", path+"/"+projectKey+"/items/"+itemKey+"/tip", nil)
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

func getItemVersions(ctx context.Context, limiter HttpRequestLimiter, path, projectKey, itemKey, refType, id, extension, versionNumber, page, limit, token string) (result ForgeResponseArray, err error) {
	task := http.Client{}

	req, err := limiter.HttpRequest(ctx, "GET", path+"/"+projectKey+"/items/"+itemKey+"/versions", nil)
	if err != nil {
		return
	}

	params := req.URL.Query()
	if len(refType) != 0 {
		params.Add("filter[type]", refType)
	}
	if len(id) != 0 {
		params.Add("filter[id]", id)
	}
	if len(extension) != 0 {
		params.Add("filter[extension.type]", extension)
	}
	if len(versionNumber) != 0 {
		params.Add("filter[versionNumber]", versionNumber)
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
		decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}
