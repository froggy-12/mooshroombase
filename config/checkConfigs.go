package config

import "errors"

func CheckIfFieldsAreEmpty(c Config) error {
	if c.Back_End_URL == "" {
		return errors.New("Back_End_URL is empty")
	}
	if len(c.FrontEndUrl) == 0 {
		return errors.New("FrontEndUrl is empty")
	}
	if c.EmailVerificationRoute == "" {
		return errors.New("EmailVerificationRoute is empty")
	}
	if c.PrimaryDB == "" {
		return errors.New("PrimaryDB is empty")
	}
	if c.JWTSecret == "" {
		return errors.New("JWTSecret is empty")
	}
	if len(c.ClientAPICredintials) == 0 {
		return errors.New("ClientAPICredintials is empty")
	}
	if c.JWTCookieAge == 0 {
		return errors.New("JWTCookieAge is empty")
	}
	if len(c.AllowedCorsOrigin) == 0 {
		return errors.New("AllowedCorsOrigin is empty")
	}
	if c.CorsHeadersMaxAge == 0 {
		return errors.New("CorsHeadersMaxAge is empty")
	}
	if !c.EmailVerificationAllowed {
		return errors.New("EmailVerificationAllowed is false")
	}
	if c.JWTTokenExpiration == 0 {
		return errors.New("JWTTokenExpiration is empty")
	}
	if len(c.RunningDatabaseContainers) == 0 {
		return errors.New("RunningDatabaseContainers is empty")
	}
	if c.MongoDBUsername == "" {
		return errors.New("MongoDBUsername is empty")
	}
	if c.MongoDBPassword == "" {
		return errors.New("MongoDBPassword is empty")
	}
	if c.VerifyEmailRouteClient == "" {
		return errors.New("VerifyEmailRouteClient is empty")
	}
	if c.SMTPServerAdress == "" {
		return errors.New("SMTPServerAdress is empty")
	}
	if c.SMTPServerPort == "" {
		return errors.New("SMTPServerPort is empty")
	}
	if c.SMTPEmailFrom == "" {
		return errors.New("SMTPEmailFrom is empty")
	}
	if c.SMTPPassword == "" {
		return errors.New("SMTPPassword is empty")
	}
	if c.MariaDBRootPassword == "" {
		return errors.New("MariaDBRootPassword is empty")
	}
	if c.BodySizeLimit == 0 {
		return errors.New("BodySizeLimit is empty")
	}
	return nil
}
