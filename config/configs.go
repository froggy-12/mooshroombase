package config

type Config struct {
	PrimaryDB                 string
	RunningDatabaseContainers []string
	MongoDBUsername           string
	MongoDBPassword           string
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
