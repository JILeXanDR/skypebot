package bot

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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

		b.On(EventMessage, func(activity *IncomingActivity) {
			assert.True(t, activity.SomeoneWroteToMe())
			assert.False(t, activity.IsGroup())
			assert.Equal(t, "Alexandr Shtovba", activity.FromUser().Name)
			assert.Equal(t, "test", activity.Text())
		})

		rr := emulateHookRequest(t, b.WebHookHandler(), privateMessageActivity)

		require.EqualValues(t, http.StatusNoContent, rr.Code)
	})

	t.Run("valid activity with massage handler (group message)", func(t *testing.T) {
		b := New(Config{})
		b.On(EventMessage, func(activity *IncomingActivity) {
			assert.True(t, activity.SomeoneWroteToMe())
			assert.True(t, activity.IsGroup())
			assert.Equal(t, "Alexandr Shtovba", activity.FromUser().Name)
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
	})

	require.NoError(t, b.Run())

	b.Send(ConversationID("8:jilexandr"), NewTextMessage("test"))
}
