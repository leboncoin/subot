package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// AddLabel stores the label in the elastic search index
func (es ES) AddLabel(label globals.Perco) error {
	if label.Name == "" || label.Query.Regexp.Input == "" {
		return errors.New("cannot create empty label")
	}

	b, err := json.Marshal(label)
	if err != nil {
		return err
	}

	_, err = es.Client.Index().
		Index("labels").
		Id(label.ID).
		Type("_doc").
		Refresh("true").
		BodyString(string(b)).
		Do(es.Context)

	if err != nil {
		return fmt.Errorf("error creating document : %s", err.Error())
	}
	return nil
}

// DeleteLabel removes the label from the elastic search index
func (es ES) DeleteLabel(documentID string) error {
	if documentID == "" {
		return errors.New("cannot delete empty documentID")
	}

	label, err := es.QueryLabelByID(documentID)
	if err != nil {
		return err
	}

	labelInAnswers, err := es.QueryAnswersByLabel(label.Name)
	if err != nil {
		return err
	}
	if len(labelInAnswers) > 0 {
		return errors.New("cannot delete label : already in use")
	}

	_, err = es.Client.Delete().
		Index("labels").
		Type("_doc").
		Refresh("true").
		Id(documentID).
		Do(es.Context)

	if err != nil {
		return err
	}

	return nil
}

// EditLabel Modifies a label matching the given documentID
func (es ES) EditLabel(documentID string, label globals.Perco) error {
	if documentID == "" {
		return errors.New("cannot edit label without documentID")
	}

	b, err := json.Marshal(label)
	if err != nil {
		return err
	}

	_, err = es.Client.Index().
		Index("labels").
		Type("_doc").
		Id(documentID).
		Refresh("true").
		BodyString(string(b)).
		Do(es.Context)

	if err != nil {
		return err
	}
	return nil
}

// GetLabels returns the list of all the labels (using the match_all query)
func (es ES) GetLabels() ([]globals.Perco, error) {
	query := elastic.NewMatchAllQuery()
	searchResult, err := es.Client.Search().
		Index("labels").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var labels []globals.Perco
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d labels\n", searchResult.Hits.TotalHits)

		for _, hit := range searchResult.Hits.Hits {
			var p globals.Perco
			err := json.Unmarshal(*hit.Source, &p)
			p.ID = hit.Id
			if err != nil {
				log.Errorf("Unable to deserialize source into perco : %s", err)
			}
			labels = append(labels, p)
		}
	}

	return labels, nil
}

// QueryLabels returns all the labels that match the input text based on percolate search
func (es ES) QueryLabels(text string) ([]string, error) {
	pq := elastic.NewPercolatorQuery().
		Field("query").
		Document(map[string]interface{}{"input": text})

	searchResult, err := es.Client.Search().
		Index("labels").
		Query(pq).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var labels []string
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d labels\n", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			var label globals.Perco
			err := json.Unmarshal(*hit.Source, &label)
			if err != nil {
				return nil, err
			}
			labels = append(labels, label.Name)
		}
	}

	return labels, nil
}

// QueryLabelByID returns the label that matches the id
func (es ES) QueryLabelByID(id string) (globals.Perco, error) {
	var label globals.Perco
	searchResult, err := es.Client.Get().
		Index("labels").
		Id(id).
		Do(es.Context)

	if err != nil {
		return label, err
	}

	if searchResult.Found {
		log.Debug("Found a label matching this ID")
		err := json.Unmarshal(*searchResult.Source, &label)
		label.ID = searchResult.Id
		if err != nil {
			return label, err
		}
	}

	return label, nil
}

// QueryLabelByName returns the label that match the given Name
func (es ES) QueryLabelByName(name string) ([]globals.Perco, error) {
	query := elastic.NewTermQuery("name.keyword", name)

	searchResult, err := es.Client.Search().
		Index("labels").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var labels []globals.Perco
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d labels\n", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			var label globals.Perco
			err := json.Unmarshal(*hit.Source, &label)
			label.ID = hit.Id
			if err != nil {
				log.Errorf("Unable to deserialize source into answer : %s", err)
			}
			labels = append(labels, label)
		}
	}

	return labels, nil
}
