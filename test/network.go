package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"loan-api/app"
	"loan-api/config"
)

func MakeRequest(
	Method string, CONFIG config.Config, url string, requestBody interface{}, headers map[string]string,
) *httptest.ResponseRecorder {
	router := app.SetupRouter(CONFIG)

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(Method, url, bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")
	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func MakePostRequest(
	CONFIG config.Config, url string, body interface{}, headers map[string]string,
) *httptest.ResponseRecorder {
	return MakeRequest(
		"POST", CONFIG, url, body, headers,
	)
}

func MakeGetRequest(
	CONFIG config.Config, url string, urlParams map[string]interface{}, headers map[string]string,
) *httptest.ResponseRecorder {
	router := app.SetupRouter(CONFIG)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	if urlParams != nil {
		q := req.URL.Query()
		for param, value := range urlParams {
			q.Add(param, fmt.Sprintf("%v", value))
		}
		req.URL.RawQuery = q.Encode()
	}

	router.ServeHTTP(w, req)

	return w
}
