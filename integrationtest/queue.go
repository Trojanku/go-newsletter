package integrationtest

import (
	"Goo/messaging"
	"Goo/utils"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// CreateQueue for testing.
// Usage:
// queue, cleanup := CreateQueue()
// defer cleanup()
// ...
func CreateQueue() (*messaging.Queue, func()) {
	utils.MustLoad("../.env-test")

	name := utils.GetStringOrDefault("QUEUE_NAME", "jobs")
	queue := messaging.NewQueue(messaging.NewQueueOptions{
		Config: getAWSConfig(),
		Name:   name,
	})

	createQueueOutput, err := queue.Client.CreateQueue(context.Background(), &sqs.CreateQueueInput{
		QueueName: &name,
	})
	if err != nil {
		panic(err)
	}

	return queue, func() {
		_, err := queue.Client.DeleteQueue(context.Background(), &sqs.DeleteQueueInput{
			QueueUrl: createQueueOutput.QueueUrl,
		})
		if err != nil {
			panic(err)
		}
	}
}
