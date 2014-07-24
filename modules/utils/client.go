package utils

import (
//	"crypto/x509"
	"errors"
	"fmt"
//	"github.com/tsuru/config"
	"io/ioutil"
	"net/http"
//	"net/url"
	"strconv"
	"strings"
)

const (
  target = "https://api.megam.co"
  )

type User struct {
  Username string
  Api_key string
 }

type Context struct {
	username string
	api_key string
//	Stdout io.Writer
//	Stderr io.Writer
//	Stdin  io.Reader
}

type Client struct {
	HTTPClient     *http.Client
	context        *Context
	Authly         *Authly
	progname       string
	currentVersion string
	versionHeader  string
}

//func NewClient(client *http.Client, context *Context, manager *Manager) *Client {
func NewClient(client *http.Client, data *User) *Client {
//func NewClient(client *http.Client, manager *Manager) *Client {
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
//	urlErr, _ := err.(*url.Error)	
	//if !ok {
	//	fmt.Println("=================>7===============")
	//	fmt.Println(err)
	//	return err
	//}
	//switch urlErr.Err.(type) {
	//case x509.UnknownAuthorityError:
		//target, _ := config.GetString("api:server")
	//	return fmt.Errorf("Failed to connect to api server (%s): %s", target, urlErr.Err)
	//}
	//target, _ := config.GetString("api:server")
	return fmt.Errorf("Failed to connect to api server.")
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	for headerKey, headerVal := range c.Authly.AuthMap {
		request.Header.Add(headerKey, headerVal)
	}

	request.Close = true
	response, err := c.HTTPClient.Do(request)
	
	//err = c.detectClientError(err)
	if err != nil {
		return nil, err
	}
	supported := response.Header.Get(c.versionHeader)
	format := `################################################################

WARNING: You're using an unsupported version of %s.

You must have at least version %s, your current
version is %s.

Please go to http://docs.tsuru.io/en/latest/install/client.html
and download the last version.

################################################################

`
	if !validateVersion(supported, c.currentVersion) {
	    fmt.Println(format)
	    fmt.Println(supported)
		//fmt.Fprintf(format, c.progname, supported, c.currentVersion)
	}
	if response.StatusCode > 399 {
		defer response.Body.Close()
		result, _ := ioutil.ReadAll(response.Body)
		return response, errors.New(string(result))
	}
	return response, nil

	/*
			import (
		    "bytes"
		    "fmt"
		    "net/http"
		    "net/url"
		)

		func main() {
		    apiUrl := "https://api.com"
		    resource := "/user/"
		    data := url.Values{}
		    data.Set("name", "foo")
		    data.Add("surname", "bar")

		    u, _ := url.ParseRequestURI(apiUrl)
		    u.Path = resource
		    urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

		    client := &http.Client{}
		    r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
		    r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
		    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		    resp, _ := client.Do(r)
		    fmt.Println(resp.Status)
		}
	*/

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
