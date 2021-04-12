package analytics_test

import (
	elastic "github.com/elastic/go-elasticsearch/v6"
	es "github.com/leboncoin/subot/pkg/elastic"
	"github.com/leboncoin/subot/pkg/globals"
	"github.com/leboncoin/subot/services/analytics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type reactionMockedStorage struct {
	es.Interface
	Client *elastic.Client `json:"client"`
}

func (m reactionMockedStorage) IsTeamMember(_ string) (teamMember bool, err error) {
	return false, nil
}

func (m reactionMockedStorage) AddMessage(_ globals.Message, _ ...string) (err error) {
	return nil
}

func (m reactionMockedStorage) QueryRangeMessages(_ string, _ string) (hits []globals.Message, err error) {
	return []globals.Message{
		{
			ID:       "chuz&fzofzo23R92I",
			Type:     "topic",
			Status:   "",
			Labels:   nil,
			Tools:    nil,
			Text:     "Bonjour, pouvez-vous me donner les droits en lecture sur vault au path apps/team-engprod/support-analytics/prod/slack svp ?",
			UserID:   "UB210NGRK",
			UserName: "clement.mondion",
			UserInfo: globals.User{
				ID:         "UB210NGRK",
				SlackID:    0,
				Avatar:     "",
				TeamMember: false,
				Name:       "",
				Profile: globals.UserProfile{
					Email:                 "clement.mondion@adevinta.com",
					Avatar:                "https://avatars.slack-edge.com/2019-01-17/525636925408_32c9a24b43103a77c711_32.jpg",
					Avatar512:             "https://avatars.slack-edge.com/2019-01-17/525636925408_32c9a24b43103a77c711_512.jpg",
					RealName:              "Clément Mondion",
					LastName:              "Mondion",
					FirstName:             "Clément",
					DisplayName:           "clem",
					RealNameNormalized:    "Clement Mondion",
					DisplayNameNormalized: "clem",
				},
			},
			Timestamp:      "1592208201.000100",
			Reactions:      nil,
			Replies:        nil,
			EditedTs:       "",
			DeletedTs:      "",
			RemindAt:       "",
			ResponseTime:   0,
			ResolutionTime: 0,
			FeedbackStatus: "",
			FeedbackTs:     "",
		},
	}, nil
}

var _ = Describe("In", func() {
	Describe("Test handler for new reactions", func() {
		It("Should save the new status when reaction is :heavy_check_mark:", func() {
			client := reactionMockedStorage{
				Client: nil,
			}
			reaction := globals.Reaction{
				Name:      "clement.mondion",
				Timestamp: "123456789.000000",
			}

			expectedResponses := []globals.SlackResponse{{
				Action: "",
			}}
			a := analytics.Analyser{ESClient: client}
			response, err := a.HandleReaction(reaction)
			Expect(err).To(Not(HaveOccurred()))
			Expect(response).To(Equal(expectedResponses))
		})
	})
})
