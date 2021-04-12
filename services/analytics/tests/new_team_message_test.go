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

type teamMemberMessageMockedStorage struct {
	es.Interface
	Client *elastic.Client `json:"client"`
}

func (m teamMemberMessageMockedStorage) IsTeamMember(_ string) (teamMember bool, err error) {
	return true, nil
}

func (m teamMemberMessageMockedStorage) AddMessage(_ globals.Message, _ ...string) (err error) {
	return nil
}

func (m teamMemberMessageMockedStorage) QueryLastUserMessages(userID string) ([]globals.Message, error) {
	return []globals.Message{}, nil
}

func (m teamMemberMessageMockedStorage) QueryLabels(_ string) (hits []string, err error) {
	return []string{"rights"}, nil
}

func (m teamMemberMessageMockedStorage) QueryAnswers(_ []string, _ []string) (answers []globals.Answer, err error) {
	return answers, nil
}

func (m teamMemberMessageMockedStorage) QueryTools(_ string) (hits []string, err error) {
	return []string{"mock0", "mock1", "mock2"}, nil
}

type teamMemberMessageMockedEngine struct {
	engine.IEngine
	Client pb.EngineClient `json:"client"`
}

func (m teamMemberMessageMockedEngine) AnalyseMessageTools(_ *pb.Text) ([]pb.Category, error) {
	return []pb.Category{}, nil
}

func (m teamMemberMessageMockedEngine) AnalyseMessageLabels(_ *pb.Text) ([]pb.Category, error) {
	return []pb.Category{}, nil
}

var _ = Describe("In", func() {

	Describe("Test handler for new team messages", func() {
		It("Should not reply to team members posting messages", func() {
			client := teamMemberMessageMockedStorage{
				Client: nil,
			}
			engine := teamMemberMessageMockedEngine{
				Client: nil,
			}
			message := globals.Message{
				Type:           "",
				Status:         "",
				Labels:         nil,
				Tools:          nil,
				Text:           "",
				UserID:         "",
				UserName:       "",
				UserInfo:       globals.User{},
				Timestamp:      "",
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
				Action:      globals.Nothing,
				Text:        "Merci pour ton message.",
			}}

			a := analytics.Analyser{ESClient: client, Engine: engine}
			response, err := a.HandleMessage(message)

			Expect(err).To(Not(HaveOccurred()))
			Expect(response).To(Equal(expectedResponses))
		})
	})
})
