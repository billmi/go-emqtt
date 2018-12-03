package main

import (
	"flag"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"os"
)

/*
*   参考: https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/sample/main.go
*   go build 4-sample.go
*   ./4-sample -action pub -message "Hello Emqtt" -num 10 -qos 1 -id myclientid -user testuser -password testpassword -broker 127.0.0.1:1883 -topic sample/test
*   ./4-sample -action sub -num 10 -qos 1 -id myid -user myuser -password mypassword -broker 127.0.0.1:1883 -topic sample/test
*/


/*
Options:
 [-help]                      Display help                                  显示帮助
 [-a pub|sub]                 Action pub (publish) or sub (subscribe)       动作|pub 发布 sub 订阅
 [-m <message>]               Payload to send                               发布一条message
 [-n <number>]                Number of messages to send or receive         发送接收数据的条数
 [-q 0|1|2]                   Quality of Service                            服务的质量
 [-clean]                     CleanSession (true if -clean is present)      清除session
 [-id <clientid>]             CliendID                                      客户端ID
 [-user <user>]               User                                          用户
 [-password <password>]       Password                                      密码
 [-broker <uri>]              Broker URI                                    代理URI
 [-topic <topic>]             Topic                                         主题
 [-store <path>]              Store Directory                               存储目录
*/

func main() {
	topic := flag.String("topic", "", "The topic name to/from which to publish/subscribe")
	broker := flag.String("broker", "tcp://127.0.0.1:1883", "The borker URI. ex: tcp://127.0.0.1:1883")
	password := flag.String("password", "", "The password (optional")
	user := flag.String("user", "", "The User (optional)")
	id := flag.String("id", "testid", "The ClientID(optional)")
	cleansess := flag.Bool("clean", false, "Set Clean Session (default false)")
	qos := flag.Int("qos", 0, "The Quality of Service 0,1,2(default 0)")
	num := flag.Int("num", 1, "The number of messages to publish or subscribe (default 1)")
	payload := flag.String("message", "", "The message text to publish (default empty)")
	action := flag.String("action", "", "Action publish or subscribe (required)")
	store := flag.String("store", ":memory", "The Store Directory (default use memory store)")
	flag.Parse()

	if *action != "pub" && *action != "sub" {
		fmt.Println("Invalid setting for -action, must be pub or sub")
		return
	}

	if *topic == "" {
		fmt.Println("Invalid seetting for -topic, must not be empty")
		return
	}

	fmt.Printf("Sample Info:\n")
	fmt.Printf("\taction:       %s\n", *action)
	fmt.Printf("\tbroker:       %s\n", *broker)
	fmt.Printf("\tclientid:     %s\n", *id)
	fmt.Printf("\tuser:         %s\n", *user)
	fmt.Printf("\tpassword:     %s\n", *password)
	fmt.Printf("\ttopic:        %s\n", *topic)
	fmt.Printf("\tmessage:      %s\n", *payload)
	fmt.Printf("\tqos:          %d\n", *qos)
	fmt.Printf("\tcleansess:    %v\n", *cleansess)
	fmt.Printf("\tnum:          %d\n", *num)
	fmt.Printf("\tstore:        %s\n", *store)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(*broker)
	opts.SetClientID(*id)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	opts.SetCleanSession(*cleansess)
	if *store != ":memory:" {
		opts.SetStore(mqtt.NewFileStore(*store))
	}

	if *action == "pub" {
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println("Sample Publisher Started")
		for i := 0; i < *num; i++ {
			fmt.Println("------ doing publish ------")
			token := client.Publish(*topic, byte(*qos), false, *payload)
			token.Wait()
		}

		client.Disconnect(250)
		fmt.Println("Sample Publisher Disconnected")
	} else {
		receiveCount := 0
		choke := make(chan [2]string)

		opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
			choke <- [2]string{message.Topic(), string(message.Payload())}
		})

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		if token := client.Subscribe(*topic, byte(*qos), nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		for receiveCount < *num {
			incoming := <-choke
			fmt.Printf("Received topic: %s message: %s\n", incoming[0], incoming[1])
			receiveCount++
		}

		client.Disconnect(250)
		fmt.Println("Sample Subscriber Disconnected")
	}

}
