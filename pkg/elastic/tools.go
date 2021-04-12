package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

//AddTool Add a new tool to the percolator index "tools"
func (es ES) AddTool(tool globals.Perco) error {
	if tool.Name == "" || tool.Query.Regexp.Input == "" {
		return errors.New("cannot create empty tool")
	}

	b, err := json.Marshal(tool)
	if err != nil {
		return err
	}

	_, err = es.Client.Index().
		Index("tools").
		Id(tool.ID).
		Refresh("true").
		Type("_doc").
		BodyString(string(b)).
		Do(es.Context)

	if err != nil {
		return fmt.Errorf("error creating document : %s", err.Error())
	}
	return nil
}


//DeleteTool Deletes a tool from the percolator index "tools"
func (es ES) DeleteTool(documentID string) error {
	if documentID == "" {
		return errors.New("cannot delete empty documentID")
	}

	tool, err := es.QueryToolByID(documentID)
	if err != nil {
		return err
	}

	toolInAnswers, err := es.QueryAnswersByTool(tool.Name)
	if err != nil {
		return err
	}

	if len(toolInAnswers) > 0 {
		return errors.New("cannot delete tool : already in use")
	}

	_, err = es.Client.Delete().
		Index("tools").
		Type("_doc").
		Id(documentID).
		Refresh("true").
		Do(es.Context)

	if err != nil {
		return err
	}

	return nil
}

// EditTool Modifies a tool matching the given documentID
func (es ES) EditTool(documentID string, tool globals.Perco) error {
	if documentID == "" {
		return errors.New("cannot edit tool without documentID")
	}

	b, err := json.Marshal(tool)
	if err != nil {
		return err
	}

	_, err = es.Client.Index().
		Index("tools").
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

//GetTools List all of the tools store in elasticsearch
func (es ES) GetTools() ([]globals.Perco, error) {
	query := elastic.NewMatchAllQuery()
	searchResult, err := es.Client.Search().
		Index("tools").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	tools := make([]globals.Perco, 0)

	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tools\n", searchResult.Hits.TotalHits)

		for _, hit := range searchResult.Hits.Hits {
			var p globals.Perco
			err := json.Unmarshal(*hit.Source, &p)
			p.ID = hit.Id
			if err != nil {
				log.Errorf("Unable to deserialize source into perco : %s", err)
			}
			tools = append(tools, p)
		}
	}
	log.Debug("Tools : ?", tools)
	return tools, nil
}

//QueryTools Query elasticsearch "tools" percolator to match tools to the given text
func (es ES) QueryTools(text string) ([]string, error) {
	pq := elastic.NewPercolatorQuery().
		Field("query").
		Document(map[string]interface{}{"input": text})

	searchResult, err := es.Client.Search().
		Index("tools").
		Query(pq).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var tools []string
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tools\n", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			var tool globals.Perco
			err := json.Unmarshal(*hit.Source, &tool)
			if err != nil {
				return nil, err
			}
			tools = append(tools, tool.Name)
		}
	}

	return tools, nil
}

// QueryToolByID returns the tool that matches the id
func (es ES) QueryToolByID(id string) (globals.Perco, error) {
	var tool globals.Perco
	searchResult, err := es.Client.Get().
		Index("tools").
		Id(id).
		Do(es.Context)

	if err != nil {
		return tool, err
	}

	if searchResult.Found {
		log.Debug("Found a tool matching this ID")
		err := json.Unmarshal(*searchResult.Source, &tool)
		tool.ID = searchResult.Id
		if err != nil {
			return tool, err
		}
	}

	return tool, nil
}

// QueryToolByName returns all the tools that match the input text based on percolate search
func (es ES) QueryToolByName(name string) ([]globals.Perco, error) {
	query := elastic.NewTermQuery("name.keyword", name)

	searchResult, err := es.Client.Search().
		Index("tools").
		Query(query).
		From(0).Size(1000).
		Pretty(true).
		Do(es.Context)

	if err != nil {
		return nil, err
	}

	var tools []globals.Perco
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tools\n", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			var tool globals.Perco
			err := json.Unmarshal(*hit.Source, &tool)
			tool.ID = hit.Id
			if err != nil {
				log.Errorf("Unable to deserialize source into answer : %s", err)
			}
			tools = append(tools, tool)
		}
	}

	return tools, nil
}
