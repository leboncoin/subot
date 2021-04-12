package dex

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dexidp/dex/connector/ldap"
	dex "github.com/dexidp/dex/server"
	storage "github.com/dexidp/dex/storage"
	sql "github.com/dexidp/dex/storage/sql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	reflect "reflect"
	"strings"
)

type dexStore struct {
	Store storage.Storage
}

// NewDexServer returns a new instance of the dex server
func NewDexServer(dexURL string, authCallback string) (*dex.Server, error) {
	var s dexStore
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetReportCaller(true)

	storeConfig := sql.SQLite3{File: "/tmp/dex.db"}

	store, err := storeConfig.Open(logger)
	if err != nil {
		log.Error("store error", err)
		return nil, err
	}

	s.Store = store

	client := storage.Client{
		ID:           "support-analytics",
		Secret:       viper.GetString("dex_secret"),
		RedirectURIs: []string{authCallback},
	}

	if err := store.CreateClient(client); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Unable to create client")
		if err != storage.ErrAlreadyExists {
			log.Error("Err is not AlreadyExists")
			return nil, err
		}
	}

	for _, connector := range viper.GetStringSlice("dex_connectors") {
		log.Debugf("adding new dex connector : %s", connector)
		fn := fmt.Sprintf("New%sConnector", strings.Title(connector))
		res := reflect.ValueOf(s).MethodByName(fn).Call([]reflect.Value{})
		log.Debugf("connector returned : %s", res)
	}

	dexConfig := dex.Config{
		Issuer:             dexURL,
		Storage:            store,
		SkipApprovalScreen: true,
		Web: dex.WebConfig{
			Dir:     "pkg/auth/dex/web",
			LogoURL: "theme/logo.svg",
			Issuer:  "support-analytics",
			Theme:   "leboncoin",
		},
		Logger: logger,
	}

	server, err := dex.NewServer(context.Background(), dexConfig)
	if err != nil {
		return nil, err
	}

	return server, nil
}

// NewLdapConnector creates and attaches an ldap connector to the store
func (s dexStore) NewLdapConnector() error {
	var c = ldap.Config{
		Host:          viper.GetString("dex_ldap_host"),
		InsecureNoSSL: true,
		BindDN:        viper.GetString("dex_ldap_username"),
		BindPW:        viper.GetString("dex_ldap_password"),
	}
	c.UserSearch.BaseDN = viper.GetString("dex_ldap_usersearch_basedn")
	c.UserSearch.Filter = viper.GetString("dex_ldap_usersearch_filter")
	c.UserSearch.Username = viper.GetString("dex_ldap_usersearch_username")
	c.UserSearch.IDAttr = viper.GetString("dex_ldap_usersearch_idattr")
	c.UserSearch.EmailAttr = viper.GetString("dex_ldap_usersearch_emailattr")
	c.UserSearch.NameAttr = viper.GetString("dex_ldap_usersearch_nameattr")

	c.GroupSearch.BaseDN = viper.GetString("dex_ldap_groupsearch_basedn")
	c.GroupSearch.Filter = viper.GetString("dex_ldap_groupsearch_filter")
	c.GroupSearch.UserAttr = viper.GetString("dex_ldap_groupsearch_username")
	c.GroupSearch.GroupAttr = viper.GetString("dex_ldap_groupsearch_emailattr")
	c.GroupSearch.NameAttr = viper.GetString("dex_ldap_groupsearch_nameattr")

	jsonConfig, err := json.Marshal(c)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Error marshaling json config")
		return err
	}

	ldapConnector := storage.Connector{
		ID:     "ldap",
		Type:   "ldap",
		Name:   "ldap",
		Config: jsonConfig,
	}

	if err := s.Store.CreateConnector(ldapConnector); err != nil {
		if err != storage.ErrAlreadyExists {
			return err
		}
	}

	return nil
}
