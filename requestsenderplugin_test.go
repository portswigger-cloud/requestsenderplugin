package requestsenderplugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendAndFilterRequest(t *testing.T) {
	testCases := []struct {
		name               string
		denylistedPaths    []string
		requestPath        string
		expectedRequestUrl string
	}{
		{
			name:               "Not denylisted path",
			denylistedPaths:    []string{"^denylisted/.*"},
			requestPath:        "/notdenylisted/abc",
			expectedRequestUrl: "/post",
		},
		{
			name:               "Empty denylisted paths",
			denylistedPaths:    []string{},
			requestPath:        "/abc",
			expectedRequestUrl: "/post",
		},
		{
			name:               "Denylisted path",
			denylistedPaths:    []string{"^/denylisted/.*"},
			requestPath:        "/denylisted/abc",
			expectedRequestUrl: "",
		},
		{
			name:               "Multiple denylisted paths",
			denylistedPaths:    []string{"^/denylisted/.*", "$denylisted2/.*"},
			requestPath:        "/denylisted/abc",
			expectedRequestUrl: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var requestUrl string

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				requestUrl = req.URL.String()
				rw.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			cfg := CreateConfig()
			cfg.DenylistedPaths = tc.denylistedPaths
			cfg.PostUrl = string(server.URL) + "/post"

			ctx := context.Background()
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			handler, err := New(ctx, next, cfg, "plugin")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+tc.requestPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			if recorder.Code != http.StatusOK {
				t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
			}

			if requestUrl != tc.expectedRequestUrl {
				t.Errorf("Expected request URL %s, but got %s", tc.expectedRequestUrl, requestUrl)
			}
		})
	}
}
