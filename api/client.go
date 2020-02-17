package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type apiClient struct {
	appID       string
	appPassword string
	baseURI     string
	token       *tokenResponse
}

func NewAPIClient(appID, appPassword string) *apiClient {
	return &apiClient{
		appID:       appID,
		appPassword: appPassword,
		baseURI:     "https://smba.trafficmanager.net/apis/v3",
	}
}

func (client *apiClient) authenticate() (*tokenResponse, error) {
	if client.token != nil {
		return client.token, nil
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", client.appID)
	form.Set("client_secret", client.appPassword)
	form.Set("scope", "https://apiClient.botframework.com/.default")

	resp, err := http.PostForm("https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token", form)
	if err != nil {
		return nil, err
	}

	var res tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	client.token = &res

	return &res, nil
}

func (client *apiClient) SendRequest() error {
	c := http.Client{}

	req, err := http.NewRequest(http.MethodPost, client.baseURI+"/ping", nil)

	token, err := client.authenticate()
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	log.Printf("%+v", res)

	return nil
}
