package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx := context.Background()

	// Retrieve environment variables
	queueName := os.Getenv("SQS_QUEUE")
	queueUrl := os.Getenv("SQS_QUEUE_URL")
	localstackEndpoint := os.Getenv("LOCALSTACK_ENDPOINT")
	region := os.Getenv("AWS_DEFAULT_REGION")

	// Debug environment variables
	fmt.Printf("SQS_QUEUE: %s\n", queueName)
	fmt.Printf("SQS_QUEUE_URL: %s\n", queueUrl)
	fmt.Printf("LOCALSTACK_ENDPOINT: %s\n", localstackEndpoint)
	fmt.Printf("AWS_DEFAULT_REGION: %s\n", region)

	// Validate environment variables
	if queueName == "" || queueUrl == "" || localstackEndpoint == "" || region == "" {
		log.Fatal("Missing required environment variables: SQS_QUEUE, SQS_QUEUE_URL, LOCALSTACK_ENDPOINT, or AWS_DEFAULT_REGION")
	}

	// Load AWS configuration with anonymous credentials
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
	)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração AWS: %v", err)
	}

	// Create SQS client with custom endpoint
	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		fmt.Printf("Setting BaseEndpoint: %s\n", localstackEndpoint)
		o.BaseEndpoint = aws.String(localstackEndpoint)
	})

	// Send a message
	sendOutput, err := client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueUrl),
		MessageBody: aws.String("Hello Luis Rosada we are using SQS in LocalStack!"),
	})
	if err != nil {
		log.Fatalf("Erro ao enviar mensagem: %v", err)
	}
	fmt.Println("Mensagem enviada com sucesso. ID:", *sendOutput.MessageId)

	// Retry receiving the message up to 3 times
	for attempt := 1; attempt <= 3; attempt++ {
		fmt.Printf("Attempt %d: Receiving message...\n", attempt)
		receiveOutput, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueUrl),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20, // Long polling
			VisibilityTimeout:   0,  // Immediate visibility
		})
		if err != nil {
			log.Fatalf("Erro ao receber mensagem: %v", err)
		}

		if len(receiveOutput.Messages) > 0 {
			msg := receiveOutput.Messages[0]
			fmt.Println("Mensagem recebida:", *msg.Body)

			// Delete the message
			_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueUrl),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Fatalf("Erro ao deletar mensagem: %v", err)
			}
			fmt.Println("Mensagem deletada com sucesso.")
			return
		}
		fmt.Println("Nenhuma mensagem recebida. Retrying...")
		time.Sleep(2 * time.Second) // Wait before retrying
	}
	fmt.Println("Nenhuma mensagem recebida após 3 tentativas.")
}
