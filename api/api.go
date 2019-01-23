package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/gathering-gg/parser/config"
)

// Token set by cli
var Token string
var client = &http.Client{Timeout: 600 * time.Second}

// Upload creates a request to send JSON to the server. Only accepts data
// that can be JSON marshalled.
func Upload(path string, body interface{}) (*http.Request, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := createDefaultRequest("POST", path, bytes.NewBuffer(json), map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
	return req, err
}

// UploadFile uploads a file
func UploadFile(path, name string, file io.Reader) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(name, name)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := createDefaultRequest("POST", path, body, map[string]string{
		"Content-Type": writer.FormDataContentType(),
	})
	return req, err
}

// Do sends an API request and returns the API response.
func Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = url.String()
				return nil, e
			}
		}
		return nil, err
	}
	defer func() {
		// Drain up to 512 bytes and close the body to let
		// the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}
	return resp, err
}

func createDefaultRequest(method, path string, body io.Reader, headers map[string]string) (*http.Request, error) {
	uri := config.Root + path
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("gathering/%s", config.Version))
	req.Header.Set("Authorization", "token "+Token)
	return req, nil
}
