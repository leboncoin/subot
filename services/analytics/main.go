package analytics

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote" // blank import for remote

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/leboncoin/subot/pkg/auth"
	"github.com/leboncoin/subot/pkg/config"
	"github.com/leboncoin/subot/pkg/elastic"
	engine "github.com/leboncoin/subot/pkg/engine_grpc_client"
	"github.com/leboncoin/subot/pkg/vault"
)

// Run starts the analytics service from outside of the package
func Run() {
	config.Initialize()

	es, err := elastic.Configure(false)
	if err != nil {
		log.Fatalf("Could not initialize elasticsearch connection %s", err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	engineClient, err := engine.Client(viper.GetString("engine_url"), opts)
	if err != nil {
		log.Fatal("Could not connect to analyser engine")
	}

	// Configure vault client
	if viper.GetBool("vault_enabled") {
		_, err := vault.Configure()
		if err != nil {
			log.Fatal("Could not initialize vault client: ", err)
		}
	}

	// Init analytics
	analyser := &Analyser{
		Engine:   engineClient,
		ESClient: es,
	}

	authHandler, authServer := auth.NewServer()

	runAPI(analyser, &authHandler, authServer)
}
