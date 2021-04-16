package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var addrs = []string{"127.0.0.1:9092"}

func main() {
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}

	partitionList, err := consumer.Partitions("web_log")
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}

	fmt.Println(partitionList)
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}

		pc.AsyncClose()

		// 异步从每个分区消费信息
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
}
