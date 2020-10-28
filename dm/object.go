package dm

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
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
	bearer, err := api.Authenticate("data:write data:read")
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
	if err != nil {
		return result, err
	}

	if nRead > maxUploadThreshold {
		if _, err := putObjectChunked(path, bucketKey, objectName, buf, token); err != nil {
			return ObjectDetails{}, err
		}

		return waitForObjectRecombination(path, bucketKey, objectName, token)
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

const chunkSize = 5000000

func putObjectChunked(path, bucketKey, objectName string, data *bytes.Buffer, token string) (result ObjectDetails, err error) {
	total := int64(data.Len())
	sessionId := fmt.Sprintf("%x-%d", md5.Sum([]byte(objectName)), time.Now().Unix())

	wg := sync.WaitGroup{}
	errChan := make(chan error)
	resultChan := make(chan ObjectDetails)

	go func() {
		remaining := total
		for remaining > 0 {
			chunk := &bytes.Buffer{}
			size := int64(chunkSize)
			if chunkSize > remaining {
				size = remaining
			}
			_, err := io.CopyN(chunk, data, size)
			if err != nil {
				errChan <- fmt.Errorf("failed to copy: %w", err)
				return
			}

			wg.Add(1)

			go func(remaining, size int64, chunk *bytes.Buffer) {
				defer wg.Done()

				task := http.Client{}
				req, err := http.NewRequest("PUT",
					path+"/"+bucketKey+"/objects/"+objectName+"/resumable",
					chunk,
				)

				if err != nil {
					errChan <- err
					return
				}

				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", total-remaining, total-remaining+size-1, total))
				req.Header.Set("Session-Id", sessionId)
				req.Header.Set("Content-Type", "application/stream")
				req.Header.Set("Content-Length", fmt.Sprintf("%d", size))

				response, err := task.Do(req)
				if err != nil {
					errChan <- fmt.Errorf("failed to execute request: %w", err)
				}
				defer response.Body.Close()

				switch response.StatusCode {

				// A chunk has been uploaded
				case http.StatusAccepted:
					return

				// The final chunk has been uploaded,
				// and the process is complete
				case http.StatusOK:
					output := ObjectDetails{}
					decoder := json.NewDecoder(response.Body)
					err = decoder.Decode(&output)
					resultChan <- output

				default:
					content, _ := ioutil.ReadAll(response.Body)
					errChan <- errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
				}
			}(remaining, size, chunk)

			remaining -= size
		}

		wg.Wait()
	}()

	select {
	case err := <-errChan:
		return ObjectDetails{}, err
	case res := <-resultChan:
		return res, nil
	}
}

// The Forge API doesn't give us many clues out when a chunked upload is recombined.
// The only way to be sure is to poll the object details API until the SHA1 hash
// is populated.
func waitForObjectRecombination(path, bucketKey, objectName, token string) (result ObjectDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+bucketKey+"/objects/"+objectName+"/details",
		nil,
	)
	if err != nil {
		return ObjectDetails{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeout := time.After(10 * time.Minute)

	for {
		select {
		case <-ticker.C:
			response, err := task.Do(req)
			if err != nil {
				return ObjectDetails{}, fmt.Errorf("failed to execute request: %w", err)
			}
			defer response.Body.Close()

			switch response.StatusCode {

			case http.StatusOK:
				output := ObjectDetails{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&output)

				if output.SHA1 != "" {
					return output, nil
				}

			default:
				content, _ := ioutil.ReadAll(response.Body)
				return ObjectDetails{}, errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
			}

		case <-timeout:
			return ObjectDetails{}, fmt.Errorf("timed out waiting for file recombination")
		}
	}
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
