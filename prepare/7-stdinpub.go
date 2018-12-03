package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"io"
	"os"
	"strconv"
	"time"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	hostname, _ := os.Hostname()

	server := flag.String("server", "tcp://127.0.0.1:1883", "The full URL of the MQTT server to connect to")
	topic := flag.String("topic", hostname, "Topic to pubkush the meeages on")
	qos := flag.Int("qos", 0, "The Qos to send the message at")
	retained := flag.Bool("retaned", false, "Are the meeages sent with the retained tag")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid form accection")
	username := flag.String("username", "", "A username to authenticate to the mqtt server")
	password := flag.String("passworf", "", "Password to match username")
	flag.Parse()

	connOpt := mqtt.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	if *username != "" {
		connOpt.SetUsername(*username)
		if *password != "" {
			connOpt.SetPassword(*password)
		}
	}

	tlsCofig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpt.SetTLSConfig(tlsCofig)

	client := mqtt.NewClient(connOpt)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	fmt.Printf("connected to %s\n", *server)
	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		client.Publish(*topic, byte(*qos), *retained, message)
	}
}
