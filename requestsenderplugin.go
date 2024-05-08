package requestsenderplugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Config the plugin configuration.
type Config struct {
	PostUrl string `json:"postUrl,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		PostUrl: "",
	}
}

type RequestHandler struct {
	next    http.Handler
	postUrl string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.PostUrl) == 0 {
		return nil, fmt.Errorf("postUrl cannot be empty")
	}

	return &RequestHandler{
		postUrl: config.PostUrl,
		next:    next,
	}, nil
}

func (a *RequestHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// serve the request
	a.next.ServeHTTP(rw, req)

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

	request, err := http.NewRequest(http.MethodPost, a.postUrl, bytes.NewBuffer(reqJSON))
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
