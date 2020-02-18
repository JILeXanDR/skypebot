package bot

import (
	"bytes"
	"github.com/JILeXanDR/skypebot/bot/message"
	"github.com/JILeXanDR/skypebot/skypeapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func loadTestFile(t *testing.T, name string) []byte {
	f, err := os.Open(path.Join("./testdata", name))
	if err != nil {
		t.Fatal(err.Error())
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err.Error())
	}

	return data
}

func emulateHookRequest(t *testing.T, handler http.HandlerFunc, json []byte) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	return rr
}

func TestBot_WebHookHandler(t *testing.T) {
	privateMessageActivity := loadTestFile(t, "private_message_activity.json")
	groupMessageActivity := loadTestFile(t, "group_message_activity.json")

	t.Run("nil body", func(t *testing.T) {
		b := New(Config{})

		rr := emulateHookRequest(t, b.WebHookHandler(), nil)

		require.EqualValues(t, http.StatusBadRequest, rr.Code)
		require.EqualValues(t, `bad activity: EOF`, rr.Body.String())
	})

	t.Run("empty json without massage handler", func(t *testing.T) {
		b := New(Config{})

		rr := emulateHookRequest(t, b.WebHookHandler(), []byte(`{}`))

		require.EqualValues(t, http.StatusNoContent, rr.Code)
	})

	t.Run("valid activity without massage handler", func(t *testing.T) {
		b := New(Config{})

		rr := emulateHookRequest(t, b.WebHookHandler(), privateMessageActivity)

		require.EqualValues(t, http.StatusNoContent, rr.Code)
	})

	t.Run("valid activity with massage handler (private message)", func(t *testing.T) {
		b := New(Config{})

		b.On(EventMessage, func(activity *Activity) {
			assert.True(t, activity.SomeoneWroteToMe())
			assert.False(t, activity.IsGroup())
			assert.Equal(t, "Alexandr Shtovba", activity.Sender().account.Name)
			assert.Equal(t, "test", activity.Text())
		})

		rr := emulateHookRequest(t, b.WebHookHandler(), privateMessageActivity)

		require.EqualValues(t, http.StatusNoContent, rr.Code)
	})

	t.Run("valid activity with massage handler (group message)", func(t *testing.T) {
		b := New(Config{})
		b.On(EventMessage, func(activity *Activity) {
			assert.True(t, activity.SomeoneWroteToMe())
			assert.True(t, activity.IsGroup())
			assert.Equal(t, "Alexandr Shtovba", activity.Sender().account.Name)
			assert.Equal(t, "help", activity.Text())
		})

		rr := emulateHookRequest(t, b.WebHookHandler(), groupMessageActivity)

		require.EqualValues(t, http.StatusNoContent, rr.Code)
	})
}

func TestBot_Send(t *testing.T) {
	b := New(Config{
		AppID:     os.Getenv("SKYPE_APP_ID"),
		AppSecret: os.Getenv("SKYPE_APP_SECRET"),
		Logger:    log.New(os.Stdout, "", 0),
	})

	require.NoError(t, b.Run())

	t.Run("to specific contact id", func(t *testing.T) {
		err := b.Send(ConversationID("8:jilexandr"), message.TextMessage("test 1"))
		require.NoError(t, err)
	})

	t.Run("to specific group id", func(t *testing.T) {
		err := b.Send(ConversationID("19:58b03afc025e48d3a34e12d370412971@thread.skype"), message.TextMessage("test 2"))
		require.NoError(t, err)
	})

	activity := &Activity{
		&skypeapi.Activity{
			From: skypeapi.ChannelAccount{
				ID: "8:jilexandr",
			},
			Conversation: skypeapi.ConversationAccount{
				ID: "19:58b03afc025e48d3a34e12d370412971@thread.skype",
			},
		},
	}
	t.Run("reply to activity (directly to group)", func(t *testing.T) {
		err := b.Send(activity, message.TextMessage("test 3"))
		require.NoError(t, err)
	})

	t.Run("reply to activity (personally to contact who sent message)", func(t *testing.T) {
		err := b.Send(activity.Sender(), message.TextMessage("test 4"))
		require.NoError(t, err)
	})
}
