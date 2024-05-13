package requestsenderplugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

// Config the plugin configuration.
type Config struct {
	PostUrl         string   `json:"postUrl,omitempty"`
	DenylistedPaths []string `json:"denylistedPaths,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		PostUrl:         "",
		DenylistedPaths: []string{},
	}
}

type RequestHandler struct {
	next                 http.Handler
	postUrl              string
	denylistedPathsRegex []*regexp.Regexp
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.PostUrl) == 0 {
		return nil, fmt.Errorf("postUrl cannot be empty")
	}

	var denylistedPathsRegex []*regexp.Regexp

	for _, regex := range config.DenylistedPaths {
		regex, err := regexp.Compile(regex)
		if err != nil {
			return nil, fmt.Errorf("invalid regex: %s", regex)
		}
		denylistedPathsRegex = append(denylistedPathsRegex, regex)
	}

	os.Stdout.WriteString(name + " plugin initialized\n")

	return &RequestHandler{
		postUrl:              config.PostUrl,
		next:                 next,
		denylistedPathsRegex: denylistedPathsRegex,
	}, nil
}

func (handler *RequestHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// serve the request
	handler.next.ServeHTTP(rw, req)
	// check if request path is deny-listed
	for _, regex := range handler.denylistedPathsRegex {
		match := regex.MatchString(req.URL.Path)
		if match {
			return
		}
	}
	// send post request to URL
	reqBody := struct {
		RequestHost string `json:"requestHost"`
	}{
		RequestHost: req.Host,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	request, err := http.NewRequest(http.MethodPost, handler.postUrl, bytes.NewBuffer(reqJSON))
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}
	defer response.Body.Close()
}
