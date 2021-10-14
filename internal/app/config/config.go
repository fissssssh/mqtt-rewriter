package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mqtt     Mqtt     `yaml:"mqtt"`
	Rewriter Rewriter `yaml:"rewriter"`
}

type Mqtt struct {
	Broker   string `yaml:"broker"`
	Port     int    `yaml:"port"`
	ClientId string `yaml:"clientId"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Rewriter struct {
	Delay    DelayRewriter    `yaml:"delay"`
	Template TemplateRewriter `yaml:"template"`
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

func Init() {
	config, _ := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(config, &AppConfig)
	log.Println(AppConfig)
}
