package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mqtt-rewriter/config"
	"mqtt-rewriter/rewriter/renderer"
	"regexp"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var delayRegexp = regexp.MustCompile(`^\$delay/(\d+)/(.*)$`)

func DelayRewriteHanler(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	result := delayRegexp.FindAllStringSubmatch(topic, -1)
	if len(result) > 0 && len(result[0]) == 3 {
		interval, err := strconv.Atoi(result[0][1])
		if err != nil {
			log.Printf("error interval format: %s", err)
			return
		}
		actualTopic := result[0][2]
		if interval >= 0 && actualTopic != "" {
			time.Sleep(time.Duration(interval) * time.Millisecond)
			client.Publish(actualTopic, message.Qos(), message.Retained(), message.Payload())
		}
	}
}

func TemplateRewriteHandler(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	payload := message.Payload()
	for _, rule := range config.AppConfig.Template.Rules {
		if topic == rule.Topic {
			var data interface{}
			switch rule.Type {
			case config.Raw:
				data = string(payload)
			case config.Json:
				err := json.Unmarshal(payload, &data)
				if err != nil {
					log.Println(err)
					return
				}
			default:
				data = string(payload)
			}
			for _, target := range rule.Targets {
				if target.Template == "" {
					_ = client.Publish(target.Topic, message.Qos(), message.Retained(), payload)
					continue
				}
				name := fmt.Sprintf("%s_%s", rule.Topic, target.Topic)
				rendered, err := renderer.Render(name, data, target.Template)
				if err != nil {
					log.Println(err)
					continue
				}
				client.Publish(target.Topic, message.Qos(), message.Retained(), rendered)
			}
		}
	}
}
