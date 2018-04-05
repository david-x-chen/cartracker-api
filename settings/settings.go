package settings

import (
	"html/template"
	"log"
	"os"

	"cartracker.api/common"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var configs = map[string]string{
	"db": "settings/db_config.json",
}

// Init initialization
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	initServerConfig()
	initOAuthConfig()
	initDbConfig()
	initTemplates()
}

// InitServerConfig site configuration
func initServerConfig() {
	common.ServerCfg = new(common.ServerConfig)
	common.ServerCfg.Host = os.Getenv("HOST")
	common.ServerCfg.SubLocation = os.Getenv("SUB_LOCACTION")
	common.ServerCfg.ReadTimeout = 60
	common.ServerCfg.WriteTimeout = 60

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
	redirectURL := os.Getenv("AUTH_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = os.Getenv("REDIRECT_URL")
	}

	clientSecret := os.Getenv("AUTH_CLIENTSECRET")
	if clientSecret == "" {
		clientSecret = os.Getenv("CLIENT_SECRET")
	}

	clientID := os.Getenv("AUTH_CLIENTID")
	if clientID == "" {
		clientID = os.Getenv("CLIENTID")
	}

	cookieSecret := os.Getenv("AUTH_COOKIESECRET")
	if cookieSecret == "" {
		cookieSecret = os.Getenv("COOKIE_SECRET")
	}

	common.OAuthStore = sessions.NewCookieStore([]byte(cookieSecret))

	common.OAuth2Cfg = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  common.AuthorizeURL,
			TokenURL: common.TokenURL,
		},
		RedirectURL: redirectURL,
		Scopes:      []string{"openid", "email", "profile"},
	}
}

// InitDbConfig initialises the db config
func initDbConfig() {
	common.MongoConfig = new(common.DbConfig)

	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		dbhost = os.Getenv("MDB_HOST")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = os.Getenv("MDB_AUTHDB")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = os.Getenv("MDB_AUTHUSER")
	}

	dbPwd := os.Getenv("DB_PWD")
	if dbPwd == "" {
		dbPwd = os.Getenv("MDB_AUTHPWD")
	}

	common.MongoConfig.MongoDBHosts = dbhost
	common.MongoConfig.AuthDatabase = dbName
	common.MongoConfig.TestDatabase = dbName
	common.MongoConfig.AuthUserName = dbUser
	common.MongoConfig.AuthPassword = dbPwd
}

// InitTemplates load templates
func initTemplates() {
	common.Tmpls["home.html"] = template.Must(template.ParseFiles(common.TemplateDir+"home.html", common.DefaultLayout))
}
