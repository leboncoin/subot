package elastic

import (
	"encoding/json"
	"errors"
	"github.com/olivere/elastic"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/leboncoin/subot/pkg/globals"
)

// AddMessage stores a message in ES index
func (es ES) AddMessage(message globals.Message, id ...string) (err error) {
	var documentID string
	body := message

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if len(id) > 0 {
		documentID = id[0]
	}

	_, err = es.Client.Index().
		Index("messages").
		Id(documentID).
		Type("_doc").
		BodyString(string(b)).
		Do(es.Context)

	if err != nil {
		return err
	}
	return nil
}

// DeleteMessage stores a message in ES index
func (es ES) DeleteMessage(messageTs string) (err error) {
	query := elastic.NewTermQuery("ts", messageTs)

	_, err = es.Client.DeleteByQuery().
		Index("messages").
		Type("_doc").
		Query(query).
		Do(es.Context)

	if err != nil {
		return err
	}

	return nil
}

// EditMessage updates a message in the elasticsearch database
func (es ES) EditMessage(documentID string, message globals.Message) error {
	if documentID == "" {
		return errors.New("cannot edit tool without documentID")
	}

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = es.Client.Index().
		Index("messages").
		Type("_doc").
		Refresh("true").
		Id(documentID).
		BodyString(string(b)).
		Do(es.Context)

	if err != nil {
		return err
	}
	return nil
}

// QueryLastUserMessages returns the last messages of the given user in the last 2 minutes
func (es ES) QueryLastUserMessages(userID string) ([]globals.Message, error) {
	end := time.Now()
	start := end.Add(-2 * time.Minute)

	termQuery := elastic.NewTermQuery("userID", userID)
	rangeQuery := elastic.NewRangeQuery("ts").
		Gte(strconv.FormatInt(start.Unix(), 10)).
		Lte(strconv.FormatInt(end.Unix(), 10))

	query := elastic.NewBoolQuery()
	query.Filter(termQuery)
	query.Filter(rangeQuery)

	searchResult, err := es.Client.Search().
		Index("messages").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var messages []globals.Message
	log.Debugf("Found a total of %d messages\n", searchResult.Hits.TotalHits)
	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {
			var m globals.Message
			err := json.Unmarshal(*hit.Source, &m)
			m.ID = hit.Id
			if err != nil {
				log.Errorf("unable to deserialize source into answer : %s", err)
			}
			messages = append(messages, m)
		}
	}

	return messages, nil
}

// QueryReminderMessages returns a list of messages in a timestamp range
func (es ES) QueryReminderMessages() ([]globals.Message, error) {
	start := strconv.FormatInt(time.Now().Add(-1*time.Minute).Unix(), 10)
	end := strconv.FormatInt(time.Now().Unix(), 10)

	typeUserQuery := elastic.NewTermQuery("type", "user")
	statusFixedQuery := elastic.NewTermQuery("status", "fixed")
	statusDeletedQuery := elastic.NewTermQuery("status", "deleted")
	rangeQuery := elastic.NewRangeQuery("remind_at").
		Gte(start).
		Lte(end)

	query := elastic.NewBoolQuery()
	query.Filter(typeUserQuery)
	query.MustNot(statusFixedQuery)
	query.MustNot(statusDeletedQuery)
	query.Filter(rangeQuery)


	searchResult, err := es.Client.Search().
		Index("messages").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var messages []globals.Message
	log.Debugf("Found a total of %d messages\n", searchResult.Hits.TotalHits)
	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {
			var m globals.Message
			err := json.Unmarshal(*hit.Source, &m)
			m.ID = hit.Id
			if err != nil {
				log.Errorf("unable to deserialize source into answer : %s", err)
			}
			messages = append(messages, m)
		}
	}

	return messages, nil
}

// QueryRangeMessages returns a list of messages in a timestamp range
func (es ES) QueryRangeMessages(start string, end string) ([]globals.Message, error) {
	termQuery := elastic.NewTermQuery("type", "user")
	rangeQuery := elastic.NewRangeQuery("ts").
		Gte(start).
		Lte(end)

	query := elastic.NewBoolQuery()
	query.Filter(termQuery)
	query.Filter(rangeQuery)

	searchResult, err := es.Client.Search().
		Index("messages").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var messages []globals.Message
	log.Debugf("Found a total of %d messages\n", searchResult.Hits.TotalHits)
	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {
			var m globals.Message
			err := json.Unmarshal(*hit.Source, &m)
			m.ID = hit.Id
			if err != nil {
				log.Errorf("unable to deserialize source into answer : %s", err)
			}
			messages = append(messages, m)
		}
	}

	return messages, nil
}
