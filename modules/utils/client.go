package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)



type User struct {
  Username string
  Api_key string
 }

type Context struct {
	username string
	api_key string
}

type Client struct {
	HTTPClient     *http.Client
	context        *Context
	Authly         *Authly
	progname       string
	currentVersion string
	versionHeader  string
}

func NewClient(client *http.Client, data *User) *Client {
context := &Context{data.Username, data.Api_key}
	return &Client{
		HTTPClient:     client,
		context:        context,
		Authly:         &Authly{},
		progname:       "sample",
		currentVersion: "0.3.0",
		versionHeader:  "Supported-Gulp",
	}
}

func (c *Client) detectClientError(err error) error {
	return fmt.Errorf("Failed to connect to api server.")
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	for headerKey, headerVal := range c.Authly.AuthMap {
		request.Header.Add(headerKey, headerVal)
	}

	request.Close = true
	response, err := c.HTTPClient.Do(request)

	if err != nil {
		return nil, err
	}
	supported := response.Header.Get(c.versionHeader)
	format := `################################################################

WARNING: You're using an unsupported version of %s.

You must have at least version %s, your current
version is %s.

################################################################

`
	if !validateVersion(supported, c.currentVersion) {
	    fmt.Println(format)
	    fmt.Println(supported)
	}
	if response.StatusCode > 399 {
		defer response.Body.Close()
		result, _ := ioutil.ReadAll(response.Body)
		return response, errors.New(string(result))
	}
	return response, nil


}

// validateVersion checks whether current version is greater or equal to
// supported version.
func validateVersion(supported, current string) bool {
	var (
		bigger bool
		limit  int
	)
	if supported == "" {
		return true
	}
	partsSupported := strings.Split(supported, ".")
	partsCurrent := strings.Split(current, ".")
	if len(partsSupported) > len(partsCurrent) {
		limit = len(partsCurrent)
		bigger = true
	} else {
		limit = len(partsSupported)
	}
	for i := 0; i < limit; i++ {
		current, err := strconv.Atoi(partsCurrent[i])
		if err != nil {
			return false
		}
		supported, err := strconv.Atoi(partsSupported[i])
		if err != nil {
			return false
		}
		if current < supported {
			return false
		}
		if current > supported {
			return true
		}
	}
	if bigger {
		return false
	}
	return true
}
