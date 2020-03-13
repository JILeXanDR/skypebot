package skypeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

const (
	unexpectedHttpStatusCodeTemplate = "The microsoft servers returned an unexpected http status code: %v"
	requestTokenUrl                  = "https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token"
	replyMessageTemplate             = "%vv3/conversations/%v/activities/%v"
)

func RequestAccessToken(microsoftAppId string, microsoftAppPassword string) (TokenResponse, error) {
	var tokenResponse TokenResponse
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", microsoftAppId)
	values.Set("client_secret", microsoftAppPassword)
	values.Set("scope", "https://api.botframework.com/.default")
	if response, err := http.PostForm(requestTokenUrl, values); err != nil {
		return tokenResponse, err
	} else if response.StatusCode == http.StatusOK {
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&tokenResponse)
		return tokenResponse, err
	} else {
		return tokenResponse, fmt.Errorf(unexpectedHttpStatusCodeTemplate, response.StatusCode)
	}
}

func SendReplyMessage(activity *Activity, message, authorizationToken string) error {
	responseActivity := &Activity{
		Type:         activity.Type,
		From:         activity.Recipient,
		Conversation: activity.Conversation,
		Recipient:    activity.From,
		Text:         message,
		ReplyToID:    activity.ID,
	}
	replyUrl := fmt.Sprintf(replyMessageTemplate, activity.ServiceURL, activity.Conversation.ID, activity.ID)
	return SendActivityRequest(responseActivity, replyUrl, authorizationToken)
}

func SendActivityRequest(activity *Activity, replyUrl, authorizationToken string) error {
	jsonEncoded, err := json.Marshal(*activity)
	if err != nil {
		return err
	}
	_, err = PlainRequest(http.MethodPost, replyUrl, bytes.NewBuffer(*&jsonEncoded), authorizationToken)
	return err
}

func PlainRequest(method string, path string, body io.Reader, authorizationToken string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(authorizationHeaderKey, authorizationHeaderValuePrefix+authorizationToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(*&req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		return resp, nil
	}

	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return nil, err
	}

	errResp.Err.InnerHttpError.StatusCode = resp.StatusCode

	return resp, &errResp
}
