package config

import (
	"path"

	"github.com/spf13/viper"
)

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"
var domain = "sites.kana.li"
var configFolderName = ".config/kana"

type AppConfig struct {
	Xdebug        bool
	Type          string
	Local         bool
	PHP           string
	AdminUsername string
	AdminPassword string
	AdminEmail    string
	Domain        string
	RootKey       string
	RootCert      string
	SiteCert      string
	SiteKey       string
	Viper         *viper.Viper
}

// LoadAppConfig gets config information that transcends sites such as app and default settings
func (c *Config) LoadAppConfig() error {

	dynamicConfig, err := c.loadAppViper()
	if err != nil {
		return err
	}

	c.App.Viper = dynamicConfig
	c.App.Xdebug = dynamicConfig.GetBool("xdebug")
	c.App.Local = dynamicConfig.GetBool("local")
	c.App.AdminEmail = dynamicConfig.GetString("admin.email")
	c.App.AdminPassword = dynamicConfig.GetString("admin.password")
	c.App.AdminUsername = dynamicConfig.GetString("admin.username")
	c.App.PHP = dynamicConfig.GetString("php")
	c.App.Type = dynamicConfig.GetString("type")

	return err

}

// loadAppViper loads the app config using viper and sets defaults
func (c *Config) loadAppViper() (*viper.Viper, error) {

	dynamicConfig := viper.New()

	dynamicConfig.SetDefault("xdebug", false)
	dynamicConfig.SetDefault("type", "site")
	dynamicConfig.SetDefault("local", false)
	dynamicConfig.SetDefault("php", "7.4")
	dynamicConfig.SetDefault("admin.username", "admin")
	dynamicConfig.SetDefault("admin.password", "password")
	dynamicConfig.SetDefault("admin.email", "admin@mykanasite.localhost")

	dynamicConfig.SetConfigName("kana")
	dynamicConfig.SetConfigType("json")
	dynamicConfig.AddConfigPath(path.Join(c.Directories.App, "config"))

	err := dynamicConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			err = dynamicConfig.SafeWriteConfig()
			if err != nil {
				return dynamicConfig, err
			}
		} else {
			return dynamicConfig, err
		}
	}

	changeConfig := false

	// Reset default "site" type if there's an invalid type in the config file
	if !CheckString(dynamicConfig.GetString("type"), validTypes) {
		changeConfig = true
		dynamicConfig.Set("type", "site")
	}

	// Reset default php version if there's an invalid version in the config file
	if !CheckString(dynamicConfig.GetString("php"), validPHPVersions) {
		changeConfig = true
		dynamicConfig.Set("php", "7.4")
	}

	if changeConfig {
		err = dynamicConfig.WriteConfig()
		if err != nil {
			return dynamicConfig, err
		}
	}

	return dynamicConfig, nil
}
