package rabbitmqDemo

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// url格式 amqb://账号:密码@rabbitmq服务器地址:端口号/vhost名称
const MQURL = "amqp://user:user@127.0.0.1:5672/testdb"

type RabbitMQ struct {
	conn    *amqp.Connection // 连接句柄
	channel *amqp.Channel

	QueueName string // 队列名称
	Exchange  string // 交换机
	Key       string
	MqUrl     string // 连接信息
}

// 实例化RabbitMQ
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbit := RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MqUrl:     MQURL,
	}

	var err error
	rabbit.conn, err = amqp.Dial(rabbit.MqUrl)
	rabbit.failOnErr(err, "创建连接错误")
	rabbit.channel, err = rabbit.conn.Channel()
	rabbit.failOnErr(err, "获取Channel失败")

	return &rabbit
}

// 销毁连接
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
		panic(fmt.Sprintf("%s: %s", message, err))
	}
}
