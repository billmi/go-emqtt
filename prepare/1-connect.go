package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
	"strings"
)

var f mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", message.Topic())
	fmt.Printf("MSG: %s\n", message.Payload())
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// 连接mqtt-server， 设置客户端ID
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("mqtt_client")

	
	// 设置用户名
	opts.SetUsername("mqtt_test")
	// 设置心跳
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	// 创建object
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("connect Fail:",token.Error())
	}
	
	// 订阅topic
	if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	
	
	// 发布topic
	for i := 0; i < 10000000; i++ {
		fmt.Println(strings.Repeat("*",40))
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("go-mqtt/sample", 0, false, text)
		token.Wait()
		fmt.Println(strings.Repeat("*",40))
		time.Sleep(5*time.Second)
	}

	// 取消订阅topic
	time.Sleep(6 * time.Second)
	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
	time.Sleep(1 * time.Second)
}


