package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

// 新建路由模式实例
func NewRouting(exchangeName string, routingKey string) *RabbitMQ {
	return NewRabbitMQ("", exchangeName, routingKey)
}

func (r *RabbitMQ) PublishRouting(message string) {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", // 路由模式类型
		true,
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare  an exchange")

	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

func (r *RabbitMQ) ReceiveRouting() {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", // 路由模式类型
		true,
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare  an exchange")

	// 创建队列
	q, err := r.channel.QueueDeclare(
		"", // 随机生成队列名称
		false,
		false,
		true, // 设置为排他
		false,
		nil)
	r.failOnErr(err, "failed to declare  an exchange")

	// 绑定队列
	err = r.channel.QueueBind(
		q.Name,
		// 在pub/sub 模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil)

	// 消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare  an exchange")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages, To exit  press CTRL + C")
	<-forever

}
