package bot

import (
	"fmt"
	"github.com/JILeXanDR/skypebot/skypeapi"
	"github.com/pkg/errors"
	"strings"
)

// https://docs.microsoft.com/en-us/azure/bot-service/rest-api/bot-framework-rest-connector-api-reference?view=azure-bot-service-4.0#conversation-operations
const serviceURL = "https://smba.trafficmanager.net/apis"

type API struct {
	appID     string
	appSecret string
	token     *skypeapi.TokenResponse
}

func (api *API) Authenticate() error {
	token, err := skypeapi.RequestAccessToken(api.appID, api.appSecret)
	if err != nil {
		return errors.Wrap(err, "can't request access token")
	}
	api.token = &token
	return nil
}

// Reply sends reply in incoming user's message
func (api *API) Reply(activity *skypeapi.Activity, text string) error {
	return skypeapi.SendReplyMessage(activity, text, api.token.AccessToken)
}

// SendToConversation sends message to a specific conversation
func (api *API) SendToConversation(conversationID, text string) error {
	activity := &skypeapi.Activity{
		Type: "message",
		Text: text,
	}
	url := fmt.Sprintf("%s/v3/conversations/%v/activities", serviceURL, conversationID)

	send := func() error {
		return skypeapi.SendActivityRequest(activity, url, api.token.AccessToken)
	}

	if err := send(); err != nil && strings.Contains(err.Error(), "401") {
		if err := api.Authenticate(); err != nil {
			return err
		}
		return send()
	}
	return nil
}

// TODO: is not used
func (api *API) ReplyToActivity(conversationID, activityID, text string) error {
	activity := &skypeapi.Activity{
		Type: "message",
		Text: text,
	}
	url := fmt.Sprintf(serviceURL+"/v3/conversations/%v/activities/%v", conversationID, activity)
	return skypeapi.SendActivityRequest(activity, url, api.token.AccessToken)
}

// TODO: move to "api" package
func newAPI(appID, appSecret string) *API {
	return &API{
		appID:     appID,
		appSecret: appSecret,
	}
}
