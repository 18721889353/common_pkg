package main

import (
	"common_pkg/rabbitmq/common"
	"common_pkg/rabbitmq/deadletter/constant"
	"log"
)

func main() {
	// # ========== 1.创建连接 ==========
	mq := common.NewRabbitMQ()
	defer mq.Close()
	mqCh := mq.Channel

	// # ========== 2.消费消息 ==========
	msgsCh, err := mqCh.Consume(constant.NormalQueue, "", false, false, false, false, nil)
	common.FailOnError(err, "消费normal队列失败")

	forever := make(chan bool)
	go func() {
		for d := range msgsCh {
			// 要实现的逻辑
			log.Printf("接收的消息: %s", d.Body)

			// 手动应答
			d.Ack(false)
			//d.Reject(true)
		}
	}()
	log.Printf("[*] Waiting for message, To exit press CTRL+C")
	<-forever
}
