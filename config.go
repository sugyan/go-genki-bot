package main

import (
	"code.google.com/p/goconf/conf"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const (
	defaultConfigFileName = "config.ini"
)

// Config type
type Config struct {
	FileName   string
	configFile *conf.ConfigFile
}

// NewConfig creates new Config
// TODO: permission check?
func NewConfig(appName string) (config *Config, err error) {
	var u *user.User
	u, err = user.Current()
	if err != nil {
		return
	}
	// create directory
	dirPath := path.Join(u.HomeDir, ".config", appName)
	if err = os.MkdirAll(dirPath, 0755); err != nil {
		return
	}
	// config file path
	filePath := path.Join(dirPath, defaultConfigFileName)
	// create if it doesn't exist
	if _, err = os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			if err = ioutil.WriteFile(filePath, []byte{}, 0644); err != nil {
				return
			}
		} else {
			return
		}
	}

	var c *conf.ConfigFile
	if c, err = conf.ReadConfigFile(filePath); err != nil {
		return
	}
	return &Config{
		FileName:   filePath,
		configFile: c,
	}, nil
}

// GetAccessToken returns access token from config file
// return nil if config doesn't exist
func (c *Config) GetAccessToken(user *string) (token *oauth.AccessToken, err error) {
	var accessToken, accessTokenSecret string
	accessToken, err = c.configFile.GetString(*user, "access_token")
	if err != nil {
		return
	}
	accessTokenSecret, err = c.configFile.GetString(*user, "access_token_secret")
	if err != nil {
		return
	}

	return &oauth.AccessToken{
		Token:          accessToken,
		Secret:         accessTokenSecret,
		AdditionalData: map[string]string{},
	}, nil
}

// SetAccessToken stores access token to config file
func (c *Config) SetAccessToken(token *oauth.AccessToken) (err error) {
	if name, ok := token.AdditionalData["screen_name"]; ok {
		c.configFile.AddOption(name, "access_token", token.Token)
		c.configFile.AddOption(name, "access_token_secret", token.Secret)
	}
	c.configFile.AddOption("default", "access_token", token.Token)
	c.configFile.AddOption("default", "access_token_secret", token.Secret)
	err = c.configFile.WriteConfigFile(c.FileName, 0644, "")
	return err
}
