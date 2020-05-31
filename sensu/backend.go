package sensu

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/apex/log"
	v2 "github.com/sensu/sensu-go/api/core/v2"
)

// NbEventMaxPerIter : Number of event to retrieve per request
// it's going to set pagination value
// Refer to the following URL for detail:
// https://docs.sensu.io/sensu-go/latest/api/overview/#pagination
var NbEventMaxPerIter uint = 200

// ExtractEvents : Take json data as []byte.
// It will return an erray of sensu event, with error
func ExtractEvents(data []byte) ([]v2.Event, error) {
	ctx := log.WithFields(log.Fields{
		"file":     "sensu/backend.go",
		"function": "ExtractEvents",
	})

	var events []v2.Event

	err := json.Unmarshal(data, &events)
	if err != nil {
		return nil, err
	}

	ctx.Errorf("Extrating %d events", len(events))

	return events, nil
}

// ExtractJSONWithHeader :  function used to call the backend and to retrieve events.
// Auth token have to be provided in the header map
func ExtractJSONWithHeader(rawURL string, header map[string]string, filter map[string]string) ([]v2.Event, error) {
	ctx := log.WithFields(log.Fields{
		"file":     "sensu/backend.go",
		"function": "extractJSONWithHeader",
	})
	var method string = "GET"
	var body string = ""
	var eventResults []v2.Event = []v2.Event{}

	reqURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	reqURLQuery := reqURL.Query()
	reqURLQuery.Set("limit", strconv.Itoa(int(NbEventMaxPerIter)))

	if len(filter) > 0 {
		for key, value := range filter {
			reqURLQuery.Add(key, value)
		}
		// body = strings.Join(filter, "%20%26%26%20")
	}

	client := &http.Client{}

	var resp *http.Response
	for ok := true; ok; ok = len(resp.Header.Get("Sensu-Continue")) != 0 {
		// Strange behavior here.
		// when using encore, spaces are translated to + instead of %20.
		// Using the replace to force %20 instead
		reqURL.RawQuery = strings.Replace(reqURLQuery.Encode(), "+", "%20", -1)
		// Ask for the next batch of data
		req, err := http.NewRequest(method, reqURL.String(), bytes.NewBuffer([]byte(body)))
		if err != nil {
			return nil, err
		}

		for key, value := range header {
			req.Header.Add(key, value)
		}

		// Alway application/json format
		req.Header.Add("Content-Type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		ctx.Errorf("Request to backend performed. Code: %d", resp.StatusCode)
		if resp.StatusCode != 200 {
			return nil, err
		}

		// Successfull auth
		// extracting the token
		token, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// fmt.Println("Error with the API Request")
			return nil, err
		}
		ctx.Errorf("Reading %d bytes(s) in response body", len(token))
		ctx.Debugf("buf: %s", string(token))
		events, err := ExtractEvents(token)
		if err != nil {
			return nil, err
		}
		eventResults = append(eventResults, events...)

		reqURLQuery.Set("continue", resp.Header.Get("Sensu-Continue"))
	}

	ctx.Debugf("Total Reading %d events", len(eventResults))

	return eventResults, nil
}

// ExtractJSONWithKey : Extract events from API with an APIKey
func ExtractJSONWithKey(url string, apikey string, filter map[string]string) ([]v2.Event, error) {
	ctx := log.WithFields(log.Fields{
		"file":     "sensu/backend.go",
		"function": "ExtractJSONWithKey",
	})

	header := map[string]string{
		"Authorization": "Key " + apikey,
	}
	eventResults, err := ExtractJSONWithHeader(url, header, filter)
	if err != nil {
		return nil, err
	}

	ctx.Debugf("Total Reading %d events", len(eventResults))
	return eventResults, nil
}

// ExtractJSONWithUser : Extract event using a login/password
func ExtractJSONWithUser(url string, namespace string, user string, password string, filter map[string]string) ([]v2.Event, error) {
	ctx := log.WithFields(log.Fields{
		"file":     "sensu/backend.go",
		"function": "ExtractJSONWithUser",
	})

	bearerKey, err := LoginUserPassword(user, password, url)
	if err != nil {
		return nil, err
	}
	ctx.Infof("Backend auth successful")
	ctx.Debugf("BearerKey %s", bearerKey)

	// Request events
	eventURL := url + "/api/core/v2/namespaces/" + namespace + "/events"

	header := map[string]string{
		"Authorization": "Bearer " + bearerKey,
	}
	eventResults, err := ExtractJSONWithHeader(eventURL, header, filter)
	if err != nil {
		return nil, err
	}

	ctx.Debugf("Total Reading %d events", len(eventResults))
	return eventResults, nil
}

// LoginUserPassword : Function used to log on the backend.
// Will return a token to be used for auth in following API request
func LoginUserPassword(user string, password string, sensuURL string) (string, error) {
	ctx := log.WithFields(log.Fields{
		"file":     "sensu/backend.go",
		"function": "ExtractJSONWithKey",
	})

	client := &http.Client{}
	uriAuth := sensuURL + "/auth"
	ctx.Debugf("Auth URL: %s", uriAuth)
	req, err := http.NewRequest("GET", uriAuth, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	auth := b64.StdEncoding.EncodeToString([]byte(user + ":" + password))
	req.Header.Add("Authorization", "Basic "+auth)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	ctx.Debugf("Auth request to backend performed. Code: %d", resp.StatusCode)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error trying to auth with backend")
	}

	// Successfull auth
	// extracting the token
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ctx.Debugf("Reading %d bytes(s) in response body", len(token))

	type sensuResponse struct {
		Access     string `json:"access_token"`
		Reauth     string `json:"refresh_token"`
		Expiration int    `json:"expires_at,omitempty"`
	}

	var reply sensuResponse

	err = json.Unmarshal([]byte(token), &reply)
	if err != nil {
		return "", err
	}

	ctx.Debugf("Sensu Backend auth info: access_token -> %s // refresh_token -> %s // Expiration -> %d", reply.Access, reply.Reauth, reply.Expiration)

	return reply.Access, nil
}
