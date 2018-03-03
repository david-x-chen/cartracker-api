package common

import (
	"html/template"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	mgo "gopkg.in/mgo.v2"
)

const (
	// DefaultLayout is default layout html template
	DefaultLayout = "templates/shared/layout.html"
	// TemplateDir is directory of templates
	TemplateDir = "templates/"

	defaultConfigFile = "oauth_config.json"

	// AuthorizeURL is default URL
	AuthorizeURL = "https://dyntech.solutions/connect/authorize"
	// TokenURL is the default URL
	TokenURL = "https://dyntech.solutions/connect/token"
)

var (
	// ServerCfg for site configuration
	ServerCfg *ServerConfig
	// MongoSession is the session
	MongoSession *mgo.Session
	// MongoConfig is the configuration file
	MongoConfig *DbConfig
	// OAuthCfgInfo is the configuration file
	OAuthCfgInfo *OAuthConfigInfo
	// OAuth2Cfg is the configuration file
	OAuth2Cfg *oauth2.Config
	// OAuthStore stores session info
	OAuthStore *sessions.CookieStore
	// Tmpls the templates
	Tmpls = map[string]*template.Template{}

	// RequiredInfoTypes shows the allowed types
	RequiredInfoTypes = []string{
		"RPM", "SPEED", "STATUS", "ENGINE_LOAD", "SHORT_FUEL_TRIM_1", "LONG_FUEL_TRIM_1",
		"THROTTLE_POS", "COMMANDED_EQUIV_RATIO", "MAF", "INTAKE_TEMP", "COOLANT_TEMP",
		"CONTROL_MODULE_VOLTAGE", "TIMING_ADVANCE", "RUN_TIME"}
)

type (
	// ServerConfig site configuration
	ServerConfig struct {
		Host         string        `json:"host"`
		ReadTimeout  time.Duration `json:"readTimeout"`
		WriteTimeout time.Duration `json:"writeTimeout"`
		SubLocation  string        `json:"subLocation"`
	}

	// DbConfig is db configuration
	DbConfig struct {
		// MongoDBHosts host server
		MongoDBHosts string `json:"mongodbhosts"` //"localhost"
		// AuthDatabase database
		AuthDatabase string `json:"authdatabase"` //"cartracker"
		// AuthUserName user name of mongodb
		AuthUserName string `json:"authusername"`
		// AuthPassword password of mongodb user
		AuthPassword string `json:"authpassword"`
		// TestDatabase test db to connect
		TestDatabase string `json:"testdatabase"` //"cartracker"
	}

	// OAuthConfigInfo is the json config
	OAuthConfigInfo struct {
		ClientSecret string   `json:"clientSecret"`
		ClientID     string   `json:"clientID"`
		RedirectURL  string   `json:"redirectUrl"`
		Scopes       []string `json:"scopes"`
		Secret       string   `json:"cookieSecret"`
	}

	// AuthUserInfo authorized user information
	AuthUserInfo struct {
		Sub           string `json:"sub,omitempty"`
		UserName      string `json:"preferred_username,omitempty"`
		Name          string `json:"name,omitempty"`
		Email         string `json:"email,omitempty"`
		EmailVerified bool   `json:"email_verified,omitempty"`
	}

	// CarTrackInfo Car track information
	CarTrackInfo struct {
		TrackDate    float64 `json:"trackdate"`
		InfoType     string  `json:"infotype"`
		StringValue  string  `json:"stringvalue"`
		NumericValue float32 `json:"numericvalue"`
		ActualValue  string  `json:"actualvalue"`
	}

	// CarTrackEntity Car track information
	CarTrackEntity struct {
		TrackDate    time.Time `json:"trackdate"`
		InfoType     string    `json:"infotype"`
		StringValue  string    `json:"stringvalue"`
		NumericValue float32   `json:"numericvalue"`
		ActualValue  string    `json:"actualvalue"`
	}
)
