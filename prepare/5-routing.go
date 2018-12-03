package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"os"
	"time"
)

/*
*   参考: https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/routing/main.go
 */

/*----------------------------------------------------------------------
This sample is designed to demonstrate the ability to set individual
callbacks on a per-subscription basis. There are three handlers in use:
 brokerLoadHandler -        $SYS/broker/load/#
 brokerConnectionHandler -  $SYS/broker/connection/#
 brokerClientHandler -      $SYS/broker/clients/#
The client will receive 100 messages total from those subscriptions,
and then print the total number of messages received from each.
It may take a few moments for the sample to complete running, as it
must wait for messages to be published.
-----------------------------------------------------------------------*/

var (
	brokerLoad       = make(chan bool)
	brokerConnection = make(chan bool)
	brokerClients    = make(chan bool)
)

func brokerLoadHandler(client mqtt.Client, message mqtt.Message) {
	brokerLoad <- true
	// fmt.Printf("BrokerLoadHandler+")
	fmt.Printf("[%s]+", message.Topic())
	fmt.Printf("%s\n", message.Payload())
}

func brokerConnectionHandler(client mqtt.Client, message mqtt.Message) {
	brokerConnection <- true
	// fmt.Printf("BrokerConnectionHandler+")
	fmt.Printf("[%s]+", message.Topic())
	fmt.Printf("%s\n", message.Payload())
}

func brokerClientsHandler(client mqtt.Client, message mqtt.Message) {
	brokerClients <- true
	// fmt.Printf("BrokerClientsHandler+")
	fmt.Printf("[%s]+", message.Topic())
	fmt.Printf("%s\n", message.Payload())
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetClientID("router-sample")
	opts.SetCleanSession(true)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	topics := []string{
		"$SYS/broker/load/1",
		"$SYS/broker/connection/1",
		"$SYS/broker/clients/1",
		"$SYS/broker/load/2",
		"$SYS/broker/connection/2",
		"$SYS/broker/clients/2",
		"$SYS/broker/load/3",
		"$SYS/broker/connection/3",
		"$SYS/broker/clients/3",
	}
	fmt.Println("Sample Publisher Started")

	if token := c.Subscribe("$SYS/broker/load/#", 0, brokerLoadHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := c.Subscribe("$SYS/broker/connection/#", 0, brokerConnectionHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := c.Subscribe("$SYS/broker/clients/#", 0, brokerClientsHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	go func() {
		loadCount := 0
		connectionCount := 0
		clientsCount := 0
		for i := 0; i < 100; i++ {
			select {
			case <-brokerLoad:
				loadCount++
			case <-brokerConnection:
				connectionCount++
			case <-brokerClients:
				clientsCount++
			}
		}
		fmt.Printf("Received %3d Broker Load messages\n", loadCount)
		fmt.Printf("Received %3d Broker Connection messages\n", connectionCount)
		fmt.Printf("Received %3d Broker Clients messages\n", clientsCount)
	}()

	// 随机发布100条数据
	for i := 0; i < 10; i++ {
		if token := c.Publish(topics[rand.NewSource(time.Now().UnixNano()).Int63()%9], 1, false, "This is a message"); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
	}

	time.Sleep(6 * time.Second)
	fmt.Println("Over!")
	c.Disconnect(250)
	time.Sleep(1 * time.Second)
}
