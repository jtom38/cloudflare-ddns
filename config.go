package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	ConfigEmail   = "EMAIL"
	ConfigToken   = "API_TOKEN"
	ConfigDomain  = "DOMAIN"
	ConfigHosts   = "HOSTS"
	ConfigIpStack = "IP_STACK" // ipv6/ipv4/dual
)

type ConfigModel struct {
	Email   string   `yaml:"Email"`
	Token   string   `yaml:"Token"`
	Domain  string   `yaml:"Domain"`
	Hosts   []string `yaml:"Hosts"`
	IpStack string   `yaml:"IpStack"`
}

type ConfigClient struct{}

func NewConfigClient() ConfigClient {
	c := ConfigClient{}
	c.RefreshEnv()

	return c
}

func (cc *ConfigClient) LoadConfig() ConfigModel {
	// load yaml first

	model, err := LoadYaml("config.yaml")
	if err != nil {
		log.Print(err)
	}

	// refresh env to make sure its current
	cc.RefreshEnv()

	// if no domains pulled from yaml, load from env
	if len(model.Hosts) == 0 {
		envHosts := cc.GetConfig(ConfigHosts)
		model.Hosts = append(model.Hosts, strings.Split(envHosts, ",")...)
	}
	if model.Domain == "" {
		model.Domain = cc.GetConfig(ConfigDomain)
	}
	if model.Email == "" {
		model.Email = cc.GetConfig(ConfigEmail)
	}
	if model.Token == "" {
		model.Token = cc.GetConfig(ConfigToken)
	}
	if model.IpStack == "" {
		model.IpStack = cc.GetConfig(ConfigIpStack)
	}

	return model
}

func (cc *ConfigClient) GetConfig(key string) string {
	res, filled := os.LookupEnv(key)
	if !filled {
		log.Printf("Missing the a value for '%v'.  Could generate errors.", key)
	}
	return res
}

func (cc *ConfigClient) GetFeature(flag string) (bool, error) {
	cc.RefreshEnv()

	res, filled := os.LookupEnv(flag)
	if !filled {
		errorMessage := fmt.Sprintf("'%v' was not found", flag)
		return false, errors.New(errorMessage)
	}

	b, err := strconv.ParseBool(res)
	if err != nil {
		return false, err
	}
	return b, nil
}

// Use this when your ConfigClient has been opened for awhile and you want to ensure you have the most recent env changes.
func (cc *ConfigClient) RefreshEnv() {
	// Check to see if we have the env file on the system
	_, err := os.Stat(".env")

	// We have the file, load it.
	if err == nil {
		_, err := os.Open(".env")
		if err == nil {
			loadEnvFile()
		}
	}
}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func (cc *ConfigClient) IsConfigInCurrentDirectory(fileName string) error {
	_, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	return nil
}

func LoadYaml(yamlFile string) (ConfigModel, error) {
	var results = ConfigModel{}

	// check for the config in the current directory
	_, err := os.Stat(yamlFile)
	if err != nil {
		return ConfigModel{}, err
	}

	content, err := os.ReadFile(yamlFile)
	if err != nil {
		return ConfigModel{}, err
	}

	err = yaml.Unmarshal(content, &results)
	if err != nil {
		return ConfigModel{}, err
	}

	return results, nil
}
