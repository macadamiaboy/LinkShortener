package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testLinkResp struct {
	ID        int64  `json:"id"`
	ShortCode string `json:"short_code"`
	LongUrl   string `json:"long_url"`
	Clicks    int32  `json:"clicks"`
}

type codeURL struct {
	Code string `json:"code,omitempty"`
	URL  string `json:"url,omitempty"`
}

func Test_E2E_CreateAndGetLink(t *testing.T) {
	linkPayload := codeURL{Code: "test", URL: "https://test.com"}
	body, err := json.Marshal(linkPayload)
	require.NoError(t, err, "failed to marshall link payload")

	createRes, err := http.Post(appURL+"/shorten", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err, "failed to send POST request")
	defer createRes.Body.Close()

	require.Equal(t, http.StatusCreated, createRes.StatusCode)

	var createdLink testLinkResp
	err = json.NewDecoder(createRes.Body).Decode(&createdLink)
	require.NoError(t, err, "failed to decode POST response")
	linkCode := createdLink.ShortCode

	getRes, err := http.Get(appURL + "/" + linkCode)
	require.NoError(t, err, "failed to send GET request")
	defer getRes.Body.Close()

	require.Equal(t, http.StatusOK, getRes.StatusCode)

	var result codeURL
	err = json.NewDecoder(getRes.Body).Decode(&result)
	require.NoError(t, err, "failed to decode GET response")
	assert.Equal(t, linkPayload.URL, result.URL)
}

func Test_E2E_GetLinkAndClicks(t *testing.T) {
	linkCode := "test"
	linkURL := "https://test.com"

	getRes, err := http.Get(appURL + "/" + linkCode)
	require.NoError(t, err, "failed to send GET request")
	defer getRes.Body.Close()

	require.Equal(t, http.StatusOK, getRes.StatusCode)

	var result codeURL
	err = json.NewDecoder(getRes.Body).Decode(&result)
	require.NoError(t, err, "failed to decode GET response")
	require.Equal(t, linkURL, result.URL)

	getClicksRes, err := http.Get(appURL + "/stats/" + linkCode)
	require.NoError(t, err, "failed to send GET request")
	defer getClicksRes.Body.Close()

	require.Equal(t, http.StatusOK, getRes.StatusCode)

	var resultClicks struct {
		Clicks int32 `json:"clicks"`
	}
	err = json.NewDecoder(getClicksRes.Body).Decode(&resultClicks)
	require.NoError(t, err, "failed to decode GET response")
	assert.Equal(t, int32(2), resultClicks.Clicks)
}
