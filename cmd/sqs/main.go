package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx := context.Background()

	sqsQueueURL := os.Getenv("SQS_QUEUE_URL") // URL completa da fila
	localstackEndpoint := os.Getenv("LOCALSTACK_ENDPOINT")
	region := os.Getenv("AWS_DEFAULT_REGION")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatal("Erro ao carregar configuração AWS:", err)
	}

	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(localstackEndpoint)
	})

	// Enviando uma mensagem
	sendOutput, err := client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(sqsQueueURL),
		MessageBody: aws.String("Olá do Go usando SQS no LocalStack!"),
	})
	if err != nil {
		log.Fatal("Erro ao enviar mensagem:", err)
	}

	fmt.Println("Mensagem enviada com sucesso. ID:", *sendOutput.MessageId)

	// Recebendo a mensagem
	receiveOutput, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sqsQueueURL),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     5,
	})
	if err != nil {
		log.Fatal("Erro ao receber mensagem:", err)
	}

	if len(receiveOutput.Messages) > 0 {
		msg := receiveOutput.Messages[0]
		fmt.Println("Mensagem recebida:", *msg.Body)

		// Apagando a mensagem da fila
		_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(sqsQueueURL),
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			log.Fatal("Erro ao deletar mensagem:", err)
		}
		fmt.Println("Mensagem deletada com sucesso.")
	} else {
		fmt.Println("Nenhuma mensagem recebida.")
	}
}
