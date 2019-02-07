package main

import (
	"fmt"
	"time"
)

type ServerStruct struct {
	channel_subscriber chan Sub
	channel_publisher chan Message
}

type Sub struct {
	topic string
	news chan string
}

type Message struct {
	topic string
	body string
}

func Server(server ServerStruct) {
	subscribedTopics := make(map[string] []chan string)
	var sub Sub
	var msg Message
	for {
		select {
		case sub = <- server.channel_subscriber:
			manageTopic(sub, subscribedTopics)
			fmt.Printf("added new subscriber to topic %s\n", sub.topic)
		case msg = <- server.channel_publisher:
			sendNews(msg, subscribedTopics)
			fmt.Printf("send topic %s to all subscribers\n", sub.topic)
		}
	}
}

func manageTopic(sub Sub, subscribedTopics map[string][]chan string) {
	if val, ok := subscribedTopics[sub.topic]; ok {
		//fmt.Printf("adding new subscriber to topic %s, actualLength %d", sub.topic, len(val))
		val = append(val, sub.news)
		//fmt.Printf(" new Length %d\n", len(val))
		subscribedTopics[sub.topic] = val
	} else {
		//fmt.Printf("adding new topic %s\n", sub.topic)
		var array [] chan string
		array = append(array, sub.news)
		subscribedTopics[sub.topic] = array
	}
}

func sendNews(msg Message, subscribedTopics map[string][]chan string) {
	array, ok := subscribedTopics[msg.topic]
	var counter int = 1
	if ok {
		//fmt.Printf("\n---\nArray size %d for topic %s\n---\n", len(array), msg.topic)
		for _, channel  := range array {
			channel <- msg.body
			fmt.Printf("Message send to %d subscribers\n", counter)
			counter++
		}
	}
}

func Client_Publish(channel chan Message, msg []Message) {
	//fmt.Printf("Publisher got array with length: %d\n", len(msg))
	for _, messages := range msg {
	//	fmt.Printf("Sending topic %s message %s\n", messages.topic, messages.body)
		channel <- messages
		time.Sleep(5 * 1e9)
	}
}

func Client_Subscribe(channel chan Sub, topic string, id int) {
	var listen = make(chan string)
	var sub = Sub{topic, listen}
	channel <- sub
	fmt.Printf("ID %d Subscribed for topic %s\n", id, topic)
	for {
		var news = <- listen
		fmt.Printf("ID %d got news: %s\n\n",id, news)
	}
}

func main() {
	var messageChannel = make(chan Message)
	var subscriberChannel = make(chan Sub)
	var server = ServerStruct{subscriberChannel, messageChannel}
	var messages1 [] Message
	messages1 = initMessagesPub1(messages1)
	var messages2 [] Message
	messages2 = initMessagesPub2(messages2)
	go Server(server)
	go Client_Subscribe(server.channel_subscriber, "go",1)
	go Client_Subscribe(server.channel_subscriber, "gym",2)
	go Client_Subscribe(server.channel_subscriber, "car",3)
	go Client_Subscribe(server.channel_subscriber, "gym",4)
	go Client_Publish(server.channel_publisher, messages1)
	go Client_Publish(server.channel_publisher, messages2)
	for {
		time.Sleep(1*1e9)
	}
}

func initMessagesPub1(messages [] Message) [] Message {
	var msg1 = Message{"go", "Learning Go"}
	var msg2 = Message{"go", "Testing Go Code"}
	var msg3 = Message{"gym", "Take proteins"}
	var msg4 = Message{"gym", "Lift weights"}
	var msg5 = Message{"car", "Drive"}
	var msg6 = Message{"car", "Stop"}
	messages = append(messages, msg1, msg2, msg3, msg4, msg5, msg6)
	return messages
}

func initMessagesPub2(messages [] Message) [] Message {
	var msg1 = Message{"gym", "Do 8 reps"}
	var msg2 = Message{"gym", "Do 3 sets"}
	messages = append(messages, msg1, msg2)
	return messages
}