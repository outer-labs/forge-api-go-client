package dm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// ObjectDetails reflects the data presented when uploading an object to a bucket or requesting details on object.
type ObjectDetails struct {
	BucketKey   string            `json:"bucketKey"`
	ObjectID    string            `json:"objectID"`
	ObjectKey   string            `json:"objectKey"`
	SHA1        string            `json:"sha1"`
	Size        uint64            `json:"size"`
	ContentType string            `json:"contentType, omitempty"`
	Location    string            `json:"location"`
	BlockSizes  []int64           `json:"blockSizes, omitempty"`
	Deltas      map[string]string `json:"deltas, omitempty"`
}

// BucketContent reflects the response when query Data Management API for bucket content.
type BucketContent struct {
	Items []ObjectDetails `json:"items"`
	Next  string          `json:"next"`
}

// UploadObject adds to specified bucket the given data (can originate from a multipart-form or direct file read).
// Return details on uploaded object, including the object URN. Check ObjectDetails struct.
func (api BucketAPI) UploadObject(bucketKey string, objectName string, reader io.Reader) (result ObjectDetails, err error) {
	bearer, err := api.Authenticate("data:write")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath

	return uploadObject(path, bucketKey, objectName, reader, bearer.AccessToken)
}

// DownloadObject returns the reader stream of the response body
// Don't forget to close it!
func (api BucketAPI) DownloadObject(bucketKey string, objectName string) (reader io.ReadCloser, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath

	return downloadObject(path, bucketKey, objectName, bearer.AccessToken)
}

// ListObjects returns the bucket contains along with details on each item.
func (api BucketAPI) ListObjects(bucketKey, limit, beginsWith, startAt string) (result BucketContent, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath

	return listObjects(path, bucketKey, limit, beginsWith, startAt, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */

func listObjects(path, bucketKey, limit, beginsWith, startAt, token string) (result BucketContent, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+bucketKey+"/objects",
		nil,
	)

	if err != nil {
		return
	}

	params := req.URL.Query()
	if len(beginsWith) != 0 {
		params.Add("beginsWith", beginsWith)
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

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)

	return
}

const maxUploadThreshold = 100000000

func uploadObject(path, bucketKey, objectName string, dataContent io.Reader, token string) (result ObjectDetails, err error) {
	buf := &bytes.Buffer{}
	nRead, err := io.Copy(buf, dataContent)

	if nRead > maxUploadThreshold {
		return putObjectChunked(path, bucketKey, objectName, buf, token)
	}

	return putObject(path, bucketKey, objectName, buf, token)
}

func putObject(path, bucketKey, objectName string, dataContent io.Reader, token string) (result ObjectDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("PUT",
		path+"/"+bucketKey+"/objects/"+objectName,
		dataContent)

	if err != nil {
		return result, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)

	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)

	return
}

const chunkSize = 10000000

// putObjectChunked runs a series of HTTP PUT operations, each of which
// represents the upload of a chunk of the object.
//
// The Forge API doesn't require the chunks to be uploaded in order, so
// in theory these operations could be executed in parallel. In practice,
// there's an account-level rate limit which would cause problems with
// large files, especially when other API requests are being made at the
// same time.
//
// Because of that rate limiting, and in the absence of a managed request
// pool, it's not safe to upload the chunks in parallel.
//
func putObjectChunked(path, bucketKey, objectName string, data *bytes.Buffer, token string) (result ObjectDetails, err error) {
	total := int64(data.Len())
	remaining := total
	sessionId := fmt.Sprintf("%d", time.Now().Unix()) // FIXME: better hash including object name?
	task := http.Client{}

	for remaining > 0 {
		chunk := &bytes.Buffer{}
		chunkSize, err := io.CopyN(chunk, data, chunkSize)
		if err != nil {
			return result, fmt.Errorf("failed to copy: %w", err)
		}

		req, err := http.NewRequest("PUT",
			path+"/"+bucketKey+"/objects/"+objectName+"/resumable",
			chunk,
		)

		if err != nil {
			return result, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", total-remaining, total-remaining+chunkSize-1, total))
		req.Header.Set("Session-Id", sessionId)

		response, err := task.Do(req)
		if err != nil {
			return result, fmt.Errorf("failed to execute request: %w", err)
		}
		defer response.Body.Close()

		switch response.StatusCode {

		// A chunk has been uploaded
		case http.StatusAccepted:
			remaining -= chunkSize
			continue

		// The final chunk has been uploaded,
		// and the process is complete
		case http.StatusOK:
			decoder := json.NewDecoder(response.Body)
			err = decoder.Decode(&result)
			return result, err

		default:
			content, _ := ioutil.ReadAll(response.Body)
			return result, errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		}
	}
	return
}

func downloadObject(path, bucketKey, objectName string, token string) (result io.ReadCloser, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+bucketKey+"/objects/"+objectName,
		nil)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)

	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] ")
		return
	}
	return response.Body, nil
}
