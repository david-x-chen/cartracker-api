package settings

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"cartracker.api/common"
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

	srvHost := os.Getenv("SRV_HOST")
	if srvHost != "" && len(srvHost) > 0 {
		common.ServerCfg.Host = srvHost
	}

	subLoc := os.Getenv("SRV_SUBDOMAIN")
	if subLoc != "" && len(subLoc) > 0 {
		common.ServerCfg.SubLocation = subLoc
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

	redirectURL := os.Getenv("AUTH_REDIRECT_URL")
	if redirectURL != "" && len(redirectURL) > 0 {
		common.OAuthCfgInfo.RedirectURL = redirectURL
	}

	clientSecret := os.Getenv("AUTH_CLIENTSECRET")
	if clientSecret != "" && len(clientSecret) > 0 {
		common.OAuthCfgInfo.ClientSecret = clientSecret
	}

	clientID := os.Getenv("AUTH_CLIENTID")
	if clientID != "" && len(clientID) > 0 {
		common.OAuthCfgInfo.ClientID = clientID
	}

	cookieSecret := os.Getenv("AUTH_COOKIESECRET")
	if cookieSecret != "" && len(cookieSecret) > 0 {
		common.OAuthCfgInfo.Secret = cookieSecret
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

	dbhost := os.Getenv("DB_HOST")
	if dbhost != "" && len(dbhost) > 0 {
		common.MongoConfig.MongoDBHosts = dbhost
	}

	dbName := os.Getenv("DB_NAME")
	if dbName != "" && len(dbName) > 0 {
		common.MongoConfig.AuthDatabase = dbName
		common.MongoConfig.TestDatabase = dbName
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser != "" && len(dbUser) > 0 {
		common.MongoConfig.AuthUserName = dbUser
	}

	dbPwd := os.Getenv("DB_PWD")
	if dbPwd != "" && len(dbPwd) > 0 {
		common.MongoConfig.AuthPassword = dbPwd
	}
}

// InitTemplates load templates
func initTemplates() {
	common.Tmpls["home.html"] = template.Must(template.ParseFiles(common.TemplateDir+"home.html", common.DefaultLayout))
}
