package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/url"
	"time"
)

/*
* 参考: https://www.cloudmqtt.com/docs-go.html
* https://github.com/CloudMQTT/go-mqtt-example
*/
func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func listen(uri *url.URL, topic string) {
	client := connect("sub", uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("* [%s] %s\n", message.Topic(), string(message.Payload()))
	})
}

func main() {
	// CLOUDMQTT_URL=mqtt://<user>:<pass>@<server>.cloudmqtt.com:<port>/<topic>
	// CLOUDMQTTURL := os.Getenv("CLOUDMQTT_URL")
	CLOUDMQTTURL := "mqtt://user:password@localhost:1883/topic"
	uri, err := url.Parse(CLOUDMQTTURL)
	if err != nil {
		log.Fatal(err)
	}

	topic := uri.Path[1:len(uri.Path)]
	if topic == "" {
		topic = "test"
	}

	go listen(uri, topic)

	client := connect("pub", uri)

	// This example sends a messages every second and
	// the same process receive the message and prints it to the console.
	// 设置定时器，1秒publish一条数据
	timer := time.NewTicker(1 * time.Second)
	for t := range timer.C {
		client.Publish(topic, 0, false, t.String())
	}
}
