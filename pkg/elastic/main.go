package elastic

import (
	"context"
	olivere "github.com/olivere/elastic"
	"github.com/spf13/viper"
	"log"
)

// ES is a struct representing the elasticsearch instance
type ES struct {
	Client  *olivere.Client `json:"Client"`
	Context context.Context `json:"context"`
}

// Configure returns an instance of the ES struct
func Configure(testing bool) (instance ES, err error) {
	host := viper.GetString("elastic_url")
	ctx := context.Background()

	client, err := olivere.NewClient(
		olivere.SetURL(host),
		olivere.SetSniff(false), olivere.SetHealthcheck(false),
	)
	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	if !testing {
		_, _, err = client.Ping(host).Do(ctx)
		if err != nil {
			// Handle error
			log.Fatalf("Unable to ping elasticsearch : %s", err)
		}

		for _, index := range []string{"answers", "firemen", "ignores", "labels", "messages", "overrides", "team", "tools"} {
			exists, err := client.IndexExists("answers").Do(ctx)
			if err != nil || !exists {
				// Handle error
				log.Fatalf("Missing index %s", index)
			}
		}
	}

	instance = ES{Client: client, Context: ctx}

	return
}
