package elastic

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// AddFireman stores the firemen info and period in ES index
func (es ES) AddFireman(message globals.Message) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = es.Client.Index().
		Index("firemen").
		Type("_doc").
		BodyString(string(b)).
		Do(es.Context)

	if err != nil {
		return err
	}

	return nil
}

// QueryRangeFireman returns the fireman of the chan for the range
func (es ES) QueryRangeFireman(start string, end string) ([]globals.Message, error) {
	startTs, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, err
	}
	endTs, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, err
	}

	rangeQuery := elastic.NewRangeQuery("ts").
		Gte(strconv.FormatInt(startTs.Unix(), 10)).
		Lte(strconv.FormatInt(endTs.Unix(), 10))

	searchResult, err := es.Client.Search().
		Index("firemen").
		Query(rangeQuery).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var messages []globals.Message
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var m globals.Message
			err := json.Unmarshal(*hit.Source, &m)
			if err != nil {
				log.Errorf("unable to deserialize source into message : %s", err)
			}
			messages = append(messages, m)
		}
	}

	return messages, nil
}