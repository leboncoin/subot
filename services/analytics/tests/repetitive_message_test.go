package analytics_test

import (
	elastic "github.com/elastic/go-elasticsearch/v6"
	es "github.com/leboncoin/subot/pkg/elastic"
	engine "github.com/leboncoin/subot/pkg/engine_grpc_client"
	pb "github.com/leboncoin/subot/pkg/engine_grpc_client/engine"
	"github.com/leboncoin/subot/pkg/globals"
	"github.com/leboncoin/subot/services/analytics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type repetitiveMockedStorage struct {
	es.Interface
	Client *elastic.Client `json:"client"`
}

func (m repetitiveMockedStorage) IsTeamMember(_ string) (teamMember bool, err error) {
	return false, nil
}

func (m repetitiveMockedStorage) QueryLastUserMessages(userID string) ([]globals.Message, error) {
	return []globals.Message{
		{
			ID:    "chuz&fzofzo23R92I",
			Type:     "topic",
			Status:   "",
			Labels:   []string{},
			Tools:    []string{},
			Text:     "Bonjour, pouvez-vous me donner les droits en lecture sur vault au path apps/team-engprod/support-analytics/prod/slack svp ?",
			UserID:   userID,
			UserName: "clement.mondion",
			UserInfo: globals.User{
			ID:         userID,
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

func (m repetitiveMockedStorage) QueryLabels(_ string) (hits []string, err error) {
	return []string{"rights"}, nil
}

func (m repetitiveMockedStorage) QueryAnswers(_ []string, _ []string) (answers []globals.Answer, err error) {
	return answers, nil
}

func (m repetitiveMockedStorage) QueryTools(_ string) (hits []string, err error) {
	return []string{"mock0", "mock1", "mock2"}, nil
}

func (m repetitiveMockedStorage) AddMessage(_ globals.Message, _ ...string) (err error) {
	return nil
}

type repetitiveMockedEngine struct {
	engine.IEngine
	Client pb.EngineClient `json:"client"`
}

func (m repetitiveMockedEngine) AnalyseMessageTools(_ *pb.Text) ([]pb.Category, error) {
	return []pb.Category{}, nil
}

func (m repetitiveMockedEngine) AnalyseMessageLabels(_ *pb.Text) ([]pb.Category, error) {
	return []pb.Category{}, nil
}

var _ = Describe("In", func() {

	Describe("Test handler for new repetitive messages", func() {
		It("Should remind the user to respect threads", func() {
			client := repetitiveMockedStorage{
				Client: nil,
			}
			engine := repetitiveMockedEngine{
				Client: nil,
			}
			message := globals.Message{
				Type:     "",
				Status:   "",
				Labels:   nil,
				Tools:    nil,
				Text:     "",
				UserID:   "UB210NGRK",
				UserName: "clement.mondion",
				UserInfo: globals.User{
					ID:         "UB210NGRK",
					SlackID:    0,
					Avatar:     "",
					TeamMember: false,
					Name:       "clement.mondion",
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
				Timestamp:      "123456789.000000",
				Reactions:      nil,
				Replies:        nil,
				EditedTs:       "",
				DeletedTs:      "",
				RemindAt:       "",
				ResponseTime:   0,
				ResolutionTime: 0,
				FeedbackStatus: "",
				FeedbackTs:     "",
			}

			expectedResponses := []*globals.SlackResponse{{
				Action:      globals.ReplyMessage,
				Text:        "Merci de respecter les threads.",
				Blocks:      nil,
				Ts:          "123456789.000000",
				ChanID:      "",
				UserID:      "",
				ResponseURL: "",
			}}

			a := analytics.Analyser{ESClient: client, Engine: engine}
			response, err := a.HandleMessage(message)
			Expect(err).To(Not(HaveOccurred()))
			Expect(response).To(Equal(expectedResponses))
		})
	})
})
