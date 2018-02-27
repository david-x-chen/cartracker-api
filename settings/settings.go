package settings

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/david-x-chen/cartracker.api/common"
	"golang.org/x/oauth2"
)

var configs = map[string]string{
	"server": "settings/server_config.json",
	"oauth":  "settings/oauth_config.json",
	"db":     "settings/db_config.json",
}

// Init initialization
func Init() {
	initServerConfig()
	initOAuthConfig()
	initDbConfig()
	initTemplates()
}

// InitServerConfig site configuration
func initServerConfig() {
	serverConfigFile, err := ioutil.ReadFile(configs["server"])
	if err != nil {
		log.Printf("serverConfigFile.Get err #%v ", err)
	}

	err = json.Unmarshal(serverConfigFile, &common.ServerCfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

// InitOAuthConfig initialises the oauth2
func initOAuthConfig() {
	jsonConfigFile, err := ioutil.ReadFile(configs["oauth"])
	if err != nil {
		log.Printf("jsonConfigFile.Get err #%v ", err)
	}

	err = json.Unmarshal(jsonConfigFile, &common.OAuthCfgInfo)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	common.OAuth2Cfg = &oauth2.Config{
		ClientID:     common.OAuthCfgInfo.ClientID,
		ClientSecret: common.OAuthCfgInfo.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  common.AuthorizeURL,
			TokenURL: common.TokenURL,
		},
		RedirectURL: common.OAuthCfgInfo.RedirectURL,
		Scopes:      common.OAuthCfgInfo.Scopes,
	}
}

// InitDbConfig initialises the db config
func initDbConfig() {
	dbConfigFile, err := ioutil.ReadFile(configs["db"])
	if err != nil {
		log.Printf("dbConfigFile.Get err #%v ", err)
	}

	err = json.Unmarshal(dbConfigFile, &common.MongoConfig)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

// InitTemplates load templates
func initTemplates() {
	common.Tmpls["home.html"] = template.Must(template.ParseFiles(common.TemplateDir+"home.html", common.DefaultLayout))
}
