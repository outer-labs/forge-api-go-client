package dm

import (
	"encoding/json"
	"net/http"
)

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api HubAPI) ListProjects(hubKey string) (result DataDetails, err error) {
	
	// TO DO: take in optional arguments for query params: id, ext, page, limit
	// https://forge.autodesk.com/en/docs/data/v2/reference/http/hubs-hub_id-projects-GET/
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}

	path := api.Host + api.HubAPIPath

	return listProjects(path, hubKey, "", "", "", "", bearer.AccessToken)
}

func (api HubAPI) GetProjectDetails(hubKey, projectKey string) (result DataDetails, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.HubAPIPath

	return getProjectDetails(path, hubKey, projectKey, bearer.AccessToken)
}

func (api HubAPI) GetTopFolders(hubKey, projectKey string) (result DataDetails, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.HubAPIPath

	return getTopFolders(path, hubKey, projectKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */
func listProjects(path, hubKey, id, extension, page, limit string, token string) (result DataDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
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
    	err = &ErrorResult{StatusCode:response.StatusCode}
    	decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}

func getProjectDetails(path, hubKey, projectKey, token string) (result DataDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
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
    	err = &ErrorResult{StatusCode:response.StatusCode}
    	decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}

func getTopFolders(path, hubKey, projectKey, token string) (result DataDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
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
    	err = &ErrorResult{StatusCode:response.StatusCode}
    	decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}