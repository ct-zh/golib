package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var addrs = []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"}

func main() {
	config := sarama.NewConfig()

	// 发送完数据需要leader和follow都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的消息将在success channel返回
	config.Producer.Return.Successes = true

	fmt.Println("开始连接kafka")

	// 连接kafka
	client, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	defer client.Close()

	fmt.Println("连接成功： ", client)

	// 构造一个消息
	msg := &sarama.ProducerMessage{
		Topic: "web_log",
		Value: sarama.StringEncoder("this is a test log"),
	}

	fmt.Println("开始发送消息")

	// 发送消息
	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("发送成功： partition:%v offset:%v\n", partition, offset)
}
