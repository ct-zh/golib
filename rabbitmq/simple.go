/*
rabbitmq简单模式
*/
package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

// 获取简单模式下的rabbitmq实例
func NewSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

// 简单模式下发送消息
func (r *RabbitMQ) PublishSimple(message string) {
	// 申请队列， 如果队列不存在则创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		false,
		// 是否具有排他性(其他用户不能访问)
		false,
		// 是否阻塞
		false,
		// 额外参数
		nil,
	)
	if err != nil {
		panic(err)
	}

	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true，根据exchange类型和routekey规则，如果无法找到符合条件的队列，会将发送消息返回给发送者
		false,
		// 如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，会把消息发还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

}

func (r *RabbitMQ) ConsumeSimple() {
	// 申请队列， 如果队列不存在则创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		false,
		// 是否具有排他性(其他用户不能访问)
		false,
		// 是否阻塞
		false,
		// 额外参数
		nil,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := r.channel.Consume(
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 队列消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages, To exit  press CTRL + C")
	<-forever
}
