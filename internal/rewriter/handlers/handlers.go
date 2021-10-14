package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mqtt-rewriter/internal/app/config"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func DelayRewriteHanler(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	var delayRegexp = regexp.MustCompile(`^\$delay/(\d+)/(.*)$`)
	result := delayRegexp.FindAllStringSubmatch(topic, -1)
	if len(result) > 0 && len(result[0]) == 3 {
		interval, err := strconv.Atoi(result[0][1])
		actualTopic := result[0][2]
		if err == nil && interval >= 0 && actualTopic != "" {
			log.Printf("mqtt received delay message\nActual Topic: %s\nDelay: %dms\nPayload: %s\n", actualTopic, interval, string(message.Payload()))
			go func() {
				time.Sleep(time.Duration(interval) * time.Millisecond)
				client.Publish(actualTopic, message.Qos(), message.Retained(), message.Payload())
			}()
			return
		}
	}
}

func TemplateRewriteHandler(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	payload := message.Payload()
	payloadBuilder := strings.Builder{}
	for _, rule := range config.AppConfig.Rewriter.Template.Rules {
		if topic == rule.Topic {
			var obj interface{}
			switch rule.Type {
			case config.Raw:
				obj = string(payload)
			case config.Json:
				err := json.Unmarshal(payload, &obj)
				if err != nil {
					log.Println(err)
					return
				}
			default:
				obj = string(payload)
			}
			for _, target := range rule.Targets {
				if target.Template == "" {
					_ = client.Publish(target.Topic, message.Qos(), message.Retained(), payload)
					continue
				}
				t := template.New(fmt.Sprintf("%s_%s", rule.Topic, target.Topic))
				t.Funcs(template.FuncMap{"json": func(obj interface{}) string {
					jstring, err := json.Marshal(obj)
					if err != nil {
						return ""
					}
					return string(jstring)
				}})
				t, err := t.Parse(target.Template)
				if err != nil {
					log.Println(err)
					continue
				}
				err = t.Execute(&payloadBuilder, obj)
				if err != nil {
					log.Println(err)
					continue
				}
				client.Publish(target.Topic, message.Qos(), message.Retained(), payloadBuilder.String())
				payloadBuilder.Reset()
			}
		}
	}
}
