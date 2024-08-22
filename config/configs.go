package config

import "time"

type Config struct {
	AppUrl                    string
	EmailVerificationRoute    string
	PrimaryDB                 string
	JWTSecret                 string
	ClientAPICredintials      map[string]string
	JWTCookieAge              int
	AllowedCorsOrigin         []string
	CorsHeadersMaxAge         int
	EmailVerificationAllowed  bool
	JWTTokenExpiration        time.Time
	RunningDatabaseContainers []string
	MongoDBUsername           string
	MongoDBPassword           string
	VerifyEmailRouteClient    string
	SMTPServerAdress          string
	SMTPServerPort            string
	SMTPEmailFrom             string
	SMTPPassword              string
	MariaDBRootPassword       string
	Authentication            bool
	ChatFunctions             bool
	BodySizeLimit             int
	GithubKey                 string
	GithubSecret              string
	DiscordKey                string
	DiscordSecret             string
	FacebookKey               string
	FaceBookSecret            string
	GoogleKey                 string
	GoogleSecret              string
	MicrosoftKey              string
	MicrosoftSecret           string
	LinkedInKey               string
	LinkedInSecret            string
	TwitterKey                string
	TwitterSecret             string
	AppleKey                  string
	AppleSecret               string
}

var Configs Config
