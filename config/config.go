package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Mqtt     Mqtt             `yaml:"mqtt"`
	Delay    DelayRewriter    `yaml:"delay"`
	Template TemplateRewriter `yaml:"template"`
}

type Mqtt struct {
	Broker   string `yaml:"broker"`
	Port     int    `yaml:"port"`
	ClientId string `yaml:"clientId"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type DelayRewriter struct {
	Enable bool `yaml:"enable"`
}

type TemplateRewriter struct {
	Enable bool           `yaml:"enable"`
	Rules  []TemplateRule `yaml:"rules"`
}

type TemplateRulePayloadType string

const (
	Raw  TemplateRulePayloadType = "raw"
	Json TemplateRulePayloadType = "json"
)

type TemplateRule struct {
	Topic   string                  `yaml:"topic"`
	Type    TemplateRulePayloadType `yaml:"type"`
	Targets []TemplateRuleTarget    `yaml:"targets"`
}

type TemplateRuleTarget struct {
	Topic    string `yaml:"topic"`
	Template string `yaml:"template"`
}

var AppConfig Config

func Init() error {
	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("mqtt.host", "localhost")
	viper.SetDefault("mqtt.port", 1883)
	viper.SetDefault("mqtt.clientId", "mqtt-rewriter")
	viper.SetDefault("mqtt.username", "")
	viper.SetDefault("mqtt.password", "")
	viper.SetDefault("mqtt.password", "")
	viper.SetDefault("delay.enable", false)
	viper.SetDefault("template.enable", false)
	viper.SetEnvPrefix("MQTTRWT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	viper.Unmarshal(&AppConfig)
	return nil
}
