package test

import (
	"github.com/assi010/gotransip/v6"
	"github.com/assi010/gotransip/v6/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getMockServer(t *testing.T, url string, method string, statusCode int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, url, req.URL.String()) // check if right url is called
		assert.Equal(t, method, req.Method)    // check if the right request method is used
		rw.WriteHeader(statusCode)             // respond with given status code
		_, err := rw.Write([]byte(response))
		require.NoError(t, err, "error when writing mock response")
	}))
}

func getRepository(t *testing.T, url string, responseStatusCode int, response string) (Repository, func()) {
	server := getMockServer(t, url, "GET", responseStatusCode, response)
	config := gotransip.DemoClientConfiguration
	config.URL = server.URL
	client, err := gotransip.NewClient(config)
	require.NoError(t, err)

	// return tearDown method with which will close the test server after the test
	tearDown := func() {
		server.Close()
	}

	return Repository{Client: client}, tearDown
}

func TestRepository_Test(t *testing.T) {
	const apiResponse = `{ "ping":"pong" }`

	repo, tearDown := getRepository(t, "/api-test", 200, apiResponse)
	defer tearDown()

	require.NoError(t, repo.Test())
}

func TestRepository_TestError(t *testing.T) {
	const apiResponse = `{ "error":"blablabla" }`

	repo, tearDown := getRepository(t, "/api-test", 409, apiResponse)
	defer tearDown()

	err := repo.Test()

	if assert.Errorf(t, err, "server response error not returned") {
		assert.Equal(t, &rest.Error{Message: "blablabla", StatusCode: 409}, err)
	}
}
