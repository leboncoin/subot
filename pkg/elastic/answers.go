package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// AddAnswer Adds an answer matching the given tool and label
func (es ES) AddAnswer(answer globals.Answer) error {
	if err := es.checkAnswer(answer); err != nil {
		return err
	}

	b, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	_, err = es.Client.Index().
		Index("answers").
		Type("_doc").
		BodyString(string(b)).
		Refresh("true").
		Do(es.Context)

	if err != nil {
		return fmt.Errorf("error creating document : %s", err.Error())
	}
	return nil
}

// checkAnswer verifies if the tool and label specified in the answer really exists and that the answer is not empty
func (es ES) checkAnswer(answer globals.Answer) error {
	if answer.Answer == "" {
		return errors.New("no answer provided")
	}

	if answer.Tool != "" {
		tools, err := es.QueryToolByName(answer.Tool)
		if err != nil {
			return err
		}
		if len(tools) != 1 {
			return errors.New("specified tool does not exist")
		}
	}

	if answer.Label != "" {
		labels, err := es.QueryLabelByName(answer.Label)
		if err != nil {
			return err
		}
		if len(labels) != 1 {
			return errors.New("specified label does not exist")
		}
	}
	return nil
}

// DeleteAnswer removes an answer from the ES index
func (es ES) DeleteAnswer(documentID string) error {
	_, err := es.Client.Delete().
		Index("answers").
		Type("_doc").
		Id(documentID).
		Refresh("true").
		Do(es.Context)

	if err != nil {
		return err
	}

	return nil
}

// EditAnswer Modifies an answer matching the given documentID
func (es ES) EditAnswer(documentID string, answer globals.Answer) error {
	if documentID == "" {
		return errors.New("cannot edit answer without documentID")
	}

	if err := es.checkAnswer(answer); err != nil {
		return err
	}

	b, err := json.Marshal(answer)
	if err != nil {
		// Handle error
		return fmt.Errorf("cannot marshall answer document : %s", err)
	}

	_, err = es.Client.Index().
		Index("answers").
		Type("_doc").
		Id(documentID).
		Refresh("true").
		BodyString(string(b)).
		Do(es.Context)

	if err != nil {
		return fmt.Errorf("error indexing document ID=%s : %s", documentID, err.Error())
	}
	return nil
}

// GetAnswers returns all of the answers stored in elasticsearch
func (es ES) GetAnswers() ([]globals.Answer, error) {
	query := elastic.NewMatchAllQuery()
	searchResult, err := es.Client.Search().
		Index("answers").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var answers []globals.Answer
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d answers\n", searchResult.Hits.TotalHits)

		for _, hit := range searchResult.Hits.Hits {
			var a globals.Answer
			err := json.Unmarshal(*hit.Source, &a)
			a.ID = hit.Id
			if err != nil {
				log.Errorf("Unable to deserialize source into answer : %s", err)
			}
			answers = append(answers, a)
		}
	}

	return answers, nil
}

// QueryAnswers returns all answers matching at least one tool or one label from given parameters
func (es ES) QueryAnswers(tools []string, labels []string) ([]globals.Answer, error) {
	l := stringToInterface(labels)
	t := stringToInterface(tools)
	query := elastic.NewBoolQuery()
	if len(tools) > 0 {
		query = query.Filter(elastic.NewTermsQuery("tool.keyword", t...))
	} else {
		query = query.Should().MustNot(elastic.NewExistsQuery("tool"))
	}
	if len(labels) > 0 {
		query = query.Filter(elastic.NewTermsQuery("label.keyword", l...))
	} else {
		query = query.Should().MustNot(elastic.NewExistsQuery("label"))
	}

	searchResult, err := es.Client.Search().
		Index("answers").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		log.Errorf("Error while querying elastic %s", err)
		return nil, err
	}

	var answers []globals.Answer
	log.Debugf("Found a total of %d answers\n", searchResult.Hits.TotalHits)
	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {
			var a globals.Answer
			err := json.Unmarshal(*hit.Source, &a)
			if err != nil {
				log.Errorf("unable to deserialize source into answer : %s", err)
			}
			answers = append(answers, a)
		}
	}

	return answers, nil
}

// QueryAnswersByLabel returns all answers matching the given label
func (es ES) QueryAnswersByLabel(label string) ([]globals.Answer, error) {
	query := elastic.NewTermQuery("label", label)

	searchResult, err := es.Client.Search().
		Index("answers").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		log.Errorf("Error while querying elastic %s", err)
		return nil, err
	}

	var answers []globals.Answer
	log.Debugf("Found a total of %d answers\n", searchResult.Hits.TotalHits)
	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {
			var a globals.Answer
			err := json.Unmarshal(*hit.Source, &a)
			if err != nil {
				log.Errorf("unable to deserialize source into answer : %s", err)
			}
			answers = append(answers, a)
		}
	}

	return answers, nil
}

// QueryAnswersByTool returns all answers matching at the given tool
func (es ES) QueryAnswersByTool(tool string) ([]globals.Answer, error) {
	query := elastic.NewTermQuery("tool", tool)

	searchResult, err := es.Client.Search().
		Index("answers").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		log.Errorf("Error while querying elastic %s", err)
		return nil, err
	}

	var answers []globals.Answer
	log.Debugf("Found a total of %d answers\n", searchResult.Hits.TotalHits)
	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {
			var a globals.Answer
			err := json.Unmarshal(*hit.Source, &a)
			if err != nil {
				log.Errorf("unable to deserialize source into answer : %s", err)
			}
			answers = append(answers, a)
		}
	}

	return answers, nil
}
