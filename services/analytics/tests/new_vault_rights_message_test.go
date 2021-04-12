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

type vaultRightsMockedStorage struct {
	es.Interface
	Client *elastic.Client `json:"client"`
}

func (m vaultRightsMockedStorage) IsTeamMember(_ string) (teamMember bool, err error) {
	return false, nil
}

func (m vaultRightsMockedStorage) QueryLastUserMessages(userID string) (hits []globals.Message, err error) {
	return []globals.Message{}, nil
}

func (m vaultRightsMockedStorage) QueryLabels(_ string) (hits []string, err error) {
	return []string{"rights"}, nil
}

func (m vaultRightsMockedStorage) QueryAnswers(_ []string, _ []string) (answers []globals.Answer, err error) {
	return []globals.Answer{
		{
			Tool:   "vault",
			Label:  "rights",
			Answer: "As-tu bien vérifié le path de ton secret ?\nFormat: `apps/team-<team name>/<app name>/<environment>/<secret name>`\n(sans `/` en début de path :wink:)\nPlus d information disponible dans cette <https://confluence.mpi-internal.com/display/LBCCORE/Vault|documentation>",
		},
	}, nil
}

func (m vaultRightsMockedStorage) QueryTools(_ string) (hits []string, err error) {
	return []string{"mock0", "mock1", "mock2"}, nil
}

func (m vaultRightsMockedStorage) AddMessage(_ globals.Message, _ ...string) (err error) {
	return nil
}

type vaultRightsMockedEngine struct {
	engine.IEngine
	Client pb.EngineClient `json:"client"`
}

func (m vaultRightsMockedEngine) AnalyseMessageTools(_ *pb.Text) ([]pb.Category, error) {
	return []pb.Category{}, nil
}

func (m vaultRightsMockedEngine) AnalyseMessageLabels(_ *pb.Text) ([]pb.Category, error) {
	return []pb.Category{}, nil
}


var _ = Describe("In", func() {

	Describe("Test handler for new user messages with known answers", func() {
		It("Should reply the known answer", func() {
			client := vaultRightsMockedStorage{
				Client: nil,
			}
			engine := vaultRightsMockedEngine{
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

			expectedResponses := []*globals.SlackResponse{
				{
					Action:      globals.ReplyMessage,
					Text:        "Merci pour ton message.\nAs-tu bien vérifié le path de ton secret ?\nFormat: `apps/team-<team name>/<app name>/<environment>/<secret name>`\n(sans `/` en début de path :wink:)\nPlus d information disponible dans cette <https://confluence.mpi-internal.com/display/LBCCORE/Vault|documentation>",
					Blocks:      nil,
					Ts:          "123456789.000000",
					ChanID:      "",
					UserID:      "",
					ResponseURL: "",
				}, {
					Action:      "reply",
					Text:        "",
					Blocks:      []interface{}{},
					Ts:          "123456789.000000",
					ChanID:      "",
					UserID:      "",
					ResponseURL: "",
				},
			}

			a := analytics.Analyser{ESClient: client, Engine: engine}
			response, err := a.HandleMessage(message)
			Expect(err).To(Not(HaveOccurred()))
			Expect(response[0]).To(Equal(expectedResponses[0]))
			//TODO : test the feedback message
		})
	})
})
