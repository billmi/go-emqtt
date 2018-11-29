package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
	"strings"
)

// 在接受到匹配订阅消息时的函数回调
// 函数回调必须有一个func() 结构
// 当回调函数为空（nil)的时候，在库接受到消息后会调用客户端的默认消息处理进程
// 可以在结构体ClientOptions的SetDefaultPublishHandler()中设置
var f mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", message.Topic())
	fmt.Printf("MSG: %s\n", message.Payload())
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// 连接mqtt-server， 设置客户端ID
	// 2个必要的参数: 代理的URL + 使用的客户端ID
	// 创建一个新的ClientOptions结构体实例，包含代理的URL+客户端ID
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("mqtt_client")
	
	
	// 使用用户名
	opts.SetUsername("mqtt_test")
	// 使用密码
	opts.SetPassword("password")
	// 设置心跳
	opts.SetKeepAlive(2 * time.Second)
	// 使用协议版本4，连接协议，4代表3.1.1, 3代表3.1
	opts.SetProtocolVersion(4)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	// 创建object
	// c是mqtt.NewClient()返回的mqtt.Client
	c := mqtt.NewClient(opts)
	// 连接返回一个token，token用来被用来指示操作是否完成
	// token.Wait() 是一个阻塞函数，只有在操作完成时才返回
	// token.WaitTimeout() 会在操作完成后等待几毫秒后返回
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("connect Fail:",token.Error())
	}
	
	// 订阅
	// Subscribe()使用3个参数，一个订阅的字符串形式的topic,订阅的qos质量和一个在接受到匹配订阅消息时的函数回调
	if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	
	
	// 发布10000000条消息
	for i := 0; i < 10000000; i++ {
		fmt.Println(strings.Repeat("*",40))
		text := fmt.Sprintf("this is msg #%d!", i)
		// c.Publish("go-mqtt/sample", 0, false, text) 不使用token进行消息发布
		
		// 发布保留连接信息 c.Publish("test/topic",1, true, "Example Payload")
		// token := c.Publish("go-mqtt/sample", 0, false, text)
		// token.Wait()
		// 发布使用4个参数
		// 发布消息的字符串的topic
		// 消息的qos质量
		// 是否保持消息连接的bool
		// 即可以是字符串形式也可以是byte数组的消息体(payload)
		if token := c.Publish("go-mqtt/sample",1,false,text); token.Wait() && token.Error() != nil{
			fmt.Println(token.Error())
		}
		fmt.Println(strings.Repeat("*",40))
		time.Sleep(5*time.Second)
	}
	
	time.Sleep(6 * time.Second)
	
	// 取消订阅
	// Unsubscribe()可以接受>=1个的取消订阅的topic参数
	// 每个topic使用单独的字符串型数组参数分开
	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// 断开连接
	// Disconnect() 使用一个参数，该参数为线程中结束任何工作的毫秒数
	c.Disconnect(250)
	time.Sleep(1 * time.Second)
}


