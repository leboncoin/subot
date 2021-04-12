package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"

	"github.com/leboncoin/subot/pkg/globals"
)

// AddTeamMember stores a team member in ES index
func (es ES) AddTeamMember(user globals.TeamMember) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = es.Client.Index().
		Index("team").
		Type("_doc").
		BodyString(string(b)).
		Refresh("true").
		Do(es.Context)

	if err != nil {
		return fmt.Errorf("error creating document : %s", err.Error())
	}
	return nil
}

// DeleteTeamMember removes a team member from the ES index
func (es ES) DeleteTeamMember(documentID string) error {
	if documentID == "" {
		return errors.New("cannot delete without document ID")
	}

	_, err := es.Client.Delete().
		Index("team").
		Type("_doc").
		Id(documentID).
		Refresh("true").
		Do(es.Context)

	if err != nil {
		return err
	}

	return nil
}

// EditTeamMember modifies team member info in the ES index
func (es ES) EditTeamMember(documentID string, teamMember globals.TeamMember) error {
	if documentID == "" {
		return errors.New("cannot delete without document ID")
	}

	b, err := json.Marshal(teamMember)
	if err != nil {
		return err
	}

	_, err = es.Client.Index().
		Index("team").
		Type("_doc").
		BodyString(string(b)).
		Refresh("true").
		Id(documentID).
		Do(es.Context)

	if err != nil {
		return err
	}

	return nil
}

//GetTeamMembers Retrieve all of the document at the index "team"
func (es ES) GetTeamMembers() ([]globals.TeamMember, error) {
	query := elastic.NewMatchAllQuery()
	searchResult, err := es.Client.Search().
		Index("team").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var users []globals.TeamMember
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d team members\n", searchResult.Hits.TotalHits)

		for _, hit := range searchResult.Hits.Hits {
			var u globals.TeamMember
			err := json.Unmarshal(*hit.Source, &u)
			u.ID = hit.Id
			if err != nil {
				log.Errorf("Unable to deserialize source into message : %s", err)
			}
			users = append(users, u)
		}
	}

	return users, nil
}

// IsTeamMember looks for the userId in the team index
func (es ES) IsTeamMember(userID string) (bool, error) {
	query := elastic.NewTermQuery("slack_id", userID)

	searchResult, err := es.Client.Search().
		Index("team").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return false, err
	}

	return searchResult.Hits.TotalHits > 0, nil
}
