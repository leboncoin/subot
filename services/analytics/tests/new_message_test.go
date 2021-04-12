package analytics_test

import (
	elastic "github.com/elastic/go-elasticsearch/v6"
	new_elastic "github.com/olivere/elastic"
	es "github.com/leboncoin/subot/pkg/elastic"
	engine "github.com/leboncoin/subot/pkg/engine_grpc_client"
	pb "github.com/leboncoin/subot/pkg/engine_grpc_client/engine"
	"github.com/leboncoin/subot/pkg/globals"
	"github.com/leboncoin/subot/services/analytics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type newMessageMockedStorage struct {
	es.Interface
	Client *elastic.Client `json:"client"`
	NewClient *new_elastic.Client `json:"new_client"`
}

func (m newMessageMockedStorage) IsTeamMember(_ string) (teamMember bool, err error) {
	return false, nil
}

func (m newMessageMockedStorage) QueryLastUserMessages(userID string) ([]globals.Message, error) {
	return []globals.Message{}, nil
}

func (m newMessageMockedStorage) QueryLabels(_ string) (hits []string, err error) {
	return []string{"rights"}, nil
}

func (m newMessageMockedStorage) QueryAnswers(_ []string, _ []string) (answers []globals.Answer, err error) {
	return answers, nil
}

func (m newMessageMockedStorage) QueryTools(_ string) (hits []string, err error) {
	return []string{"mock0", "mock1", "mock2"}, nil
}

func (m newMessageMockedStorage) AddMessage(_ globals.Message, _ ...string) (err error) {
	return nil
}

type newMessageMockedEngine struct {
	engine.IEngine
	Client pb.EngineClient `json:"client"`
}

func (m newMessageMockedEngine) AnalyseMessageTools(_ *pb.Text) ([]pb.Category, error) {
	return []pb.Category{}, nil
}

func (m newMessageMockedEngine) AnalyseMessageLabels(_ *pb.Text) ([]pb.Category, error) {
	return []pb.Category{}, nil
}

var _ = Describe("In", func() {

	Describe("Test handler for new user messages", func() {
		It("Should save the correct message and reply default answer", func() {
			client := newMessageMockedStorage{
				Client: nil,
			}
			engine := newMessageMockedEngine{
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
				Text:        "Merci pour ton message.",
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
