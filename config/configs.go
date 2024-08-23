package config

type Config struct {
	Back_End_URL              string
	FrontEndUrl               []string
	EmailVerificationRoute    string
	PrimaryDB                 string
	JWTSecret                 string
	ClientAPICredintials      map[string]string
	JWTCookieAge              int
	AllowedCorsOrigin         []string
	CorsHeadersMaxAge         int
	EmailVerificationAllowed  bool
	JWTTokenExpiration        int
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
	GoogleKey                 string
	GoogleSecret              string
	DefaultProfilePicURL      string
}

var Configs Config
