package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/leboncoin/subot/pkg/globals"
	"github.com/leboncoin/subot/pkg/slack"
	"github.com/leboncoin/subot/services/replier"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type replyMockedSender struct {
	slack.Interface
	Channel  slack.Chan `json:"channel"`
	Token    string     `json:"token"`
	BotToken string     `json:"bot_token"`
	BotID    string     `json:"bot_id"`
}

func (m replyMockedSender) GetMessageType(_ slack.Event) globals.MessageType {
	return globals.Thread
}

func (m replyMockedSender) GetEvent(e slack.Event) globals.Event {
	return m.GetReply(e)
}

func (m replyMockedSender) GetReply(_ slack.Event) globals.Reply {
	return globals.Reply{
		UserID:    "",
		UserName:  "",
		UserInfo:  globals.User{},
		Timestamp: "",
		ThreadTs:  "",
		Text:      "",
		FromBot:   false,
	}
}

func (m replyMockedSender) IsValidToken(_ slack.EventRequest) bool {
	return true
}

func (m replyMockedSender) IsWatchedChannel(_ slack.Event) bool {
	return true
}

var _ = Describe("In", func() {
	var mockAnalyticsServer *httptest.Server

	var analyticsReq struct {
		globals.Message
	}

	var analyticsRes replier.AnalyticsAPIResponse

	BeforeEach(func() {
		var analyticsJSONBody []byte
		var err error

		mockAnalyticsServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.RequestURI != fmt.Sprintf("/v1/analytics/%s", globals.Thread) {
				res.WriteHeader(500)
			}
			res.WriteHeader(200)
			_, err = res.Write(analyticsJSONBody)
		}))

		analyticsRes = replier.AnalyticsAPIResponse{
			Error: "",
			Responses: []globals.SlackResponse{
				{
					Action: globals.Nothing,
					Ts:     analyticsReq.Timestamp,
				},
			},
		}

		analyticsJSONBody, err = json.Marshal(analyticsRes)

		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		mockAnalyticsServer.Close()
	})

	Describe("Test handler for message replies", func() {
		It("Should call the correct analytics api endpoint", func() {
			s := replyMockedSender{
				Channel:  slack.Chan{},
				Token:    "",
				BotToken: "",
				BotID:    "",
			}

			h := replier.Handler{Slack: s, ApiUrl: mockAnalyticsServer.URL}
			h.HandleNewEvent(slack.EventRequest{})
			Expect("/v1").To(Equal("/v1"))
			Expect(analyticsRes.Error).To(BeEmpty())
			Expect(analyticsRes.Responses).To(HaveLen(1))
			Expect(analyticsRes.Responses[0].Action).To(Equal(globals.Nothing))
		})
	})
})
