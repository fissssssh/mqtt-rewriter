package rewriter

import (
	"fmt"
	"log"
	"mqtt-rewriter/config"
	"mqtt-rewriter/rewriter/handlers"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var once sync.Once
var instance *mqtt.Client

func Instance() *mqtt.Client {
	once.Do(func() {
		opts := mqtt.NewClientOptions()
		options := config.AppConfig.Mqtt
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", options.Broker, options.Port))
		opts.SetClientID(options.ClientId)
		opts.SetUsername(options.Username)
		opts.SetPassword(options.Password)
		opts.SetCleanSession(true)
		opts.SetKeepAlive(300)
		// 连接成功回调
		opts.OnConnect = func(c mqtt.Client) {
			log.Println("mqtt connected!")
			// 延时重写
			if config.AppConfig.Delay.Enable {
				log.Println("loading delay rewriter")
				c.Subscribe("$delay/+/+/#", 2, func(c mqtt.Client, m mqtt.Message) {
					go handlers.DelayRewriteHanler(c, m)
				})
			}
			// 模板重写
			if config.AppConfig.Template.Enable {
				log.Println("loading template rewriter")
				for _, rule := range config.AppConfig.Template.Rules {
					c.Subscribe(rule.Topic, 2, func(c mqtt.Client, m mqtt.Message) {
						go handlers.TemplateRewriteHandler(c, m)
					})
				}
			}
		}
		// 连接丢失回调
		opts.OnConnectionLost = func(c mqtt.Client, e error) {
			log.Printf("mqtt disconnected!\nreason:%s\n", e)
		}
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		instance = &client
	})
	return instance
}
