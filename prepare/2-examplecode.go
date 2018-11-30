package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"sync"
	"time"
)

/*
* 参考: https://gist.github.com/atotto/6406c0e579c6cd8c920ba53ba952f0f5
*/

const TOPIC = "mytopic/test"

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	var wg sync.WaitGroup
	wg.Add(1)

	if token := client.Subscribe(TOPIC, 0, func(client mqtt.Client, message mqtt.Message) {
		if string(message.Payload()) != "mymessage" {
			fmt.Printf("want mymessage, got %s,", message.Payload())
		} else {
			fmt.Println(string(message.Payload()))
		}
		wg.Done()
	}); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	time.Sleep(10*time.Second)
	if token := client.Publish(TOPIC, 0, false, "mymessage"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	time.Sleep(10*time.Second)
	wg.Wait()
}

/*
# 连接信息
节点            客户端ID                                         用户名     IP地址     端口   清除会话 版本协议 心跳 连接时间
emqx@127.0.0.1 Mjg0NzI0MTQ5ODgyMDY4ODE1NDczMzQ2MDkzNjQxMjM2NDI undefined 127.0.0.1 55830 true    4      30  2018-xx-xx xx:58:29
*/