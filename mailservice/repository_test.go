package mailservice

import (
	"github.com/assi010/gotransip/v6"
	"github.com/assi010/gotransip/v6/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockServer struct is used to test the how the client sends a request
// and responds to a servers response
type mockServer struct {
	t                   *testing.T
	expectedURL         string
	expectedMethod      string
	statusCode          int
	expectedRequestBody string
	response            string
	skipRequestBody     bool
}

func (m *mockServer) getHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(m.t, m.expectedURL, req.URL.String()) // check if right expectedURL is called

		if m.skipRequestBody == false && req.ContentLength != 0 {
			// get the request body
			// and check if the body matches the expected request body
			body, err := ioutil.ReadAll(req.Body)
			require.NoError(m.t, err)
			assert.Equal(m.t, m.expectedRequestBody, string(body))
		}

		assert.Equal(m.t, m.expectedMethod, req.Method) // check if the right expectedRequestBody expectedMethod is used
		rw.WriteHeader(m.statusCode)                    // respond with given status code

		if m.response != "" {
			_, err := rw.Write([]byte(m.response))
			require.NoError(m.t, err, "error when writing mock response")
		}
	}))
}

func (m *mockServer) getClient() (*repository.Client, func()) {
	httpServer := m.getHTTPServer()
	config := gotransip.DemoClientConfiguration
	config.URL = httpServer.URL
	client, err := gotransip.NewClient(config)
	require.NoError(m.t, err)

	// return tearDown method with which will close the test server after the test
	tearDown := func() {
		httpServer.Close()
	}

	return &client, tearDown
}

func TestRepository_AddDNSEntriesDomains(t *testing.T) {
	expectedRequestBody := `{"domainNames":["example.com","another.com"]}`
	server := mockServer{t: t, expectedMethod: "POST", expectedURL: "/mail-service", statusCode: 201, expectedRequestBody: expectedRequestBody}
	client, tearDown := server.getClient()
	defer tearDown()

	repo := Repository{Client: *client}
	exmapleDomains := []string{"example.com", "another.com"}
	err := repo.AddDNSEntriesDomains(exmapleDomains)

	require.NoError(t, err)
}

func TestRepository_GetInformation(t *testing.T) {
	responseBody := `{"mailServiceInformation":{ "username": "test@vps.transip.email", "password": "KgDseBsmWJNTiGww", "usage": 54, "quota": 1000, "dnsTxt": "782d28c2fa0b0bdeadf979e7155a83a15632fcddb0149d510c09fb78a470f7d3" } }`
	server := mockServer{t: t, expectedMethod: "GET", expectedURL: "/mail-service", statusCode: 200, response: responseBody}
	client, tearDown := server.getClient()
	defer tearDown()

	repo := Repository{Client: *client}
	mailServiceInfo, err := repo.GetInformation()

	require.NoError(t, err)
	assert.Equal(t, "test@vps.transip.email", mailServiceInfo.Username)
	assert.Equal(t, "KgDseBsmWJNTiGww", mailServiceInfo.Password)
	assert.Equal(t, float32(54), mailServiceInfo.Usage)
	assert.Equal(t, float32(1000), mailServiceInfo.Quota)
	assert.Equal(t, "782d28c2fa0b0bdeadf979e7155a83a15632fcddb0149d510c09fb78a470f7d3", mailServiceInfo.DNSTxt)
}

func TestRepository_RegeneratePassword(t *testing.T) {
	server := mockServer{t: t, expectedMethod: "PATCH", expectedURL: "/mail-service", statusCode: 204}
	client, tearDown := server.getClient()
	defer tearDown()

	repo := Repository{Client: *client}
	err := repo.RegeneratePassword()

	require.NoError(t, err)
}
