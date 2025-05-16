# 📨 AWS SQS & S3 with LocalStack, Terraform, and Go

> **TL;DR**: Simulate AWS S3 and SQS locally using LocalStack, manage infrastructure with Terraform, and interact with services using Go programs—all without an AWS account.

This project provides a complete setup to test **AWS S3** and **SQS** locally using **LocalStack**, with **Terraform** for infrastructure provisioning and **Go** programs for service interaction. Tailored for **WSL2** on Windows with Docker Desktop, this README offers step-by-step instructions, commands, and troubleshooting tips for a seamless experience.

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
   - [Docker Desktop](#docker-desktop)
   - [Terraform](#terraform)
   - [Go](#go)
   - [AWS CLI](#aws-cli)
3. [Project Structure](#project-structure)
4. [Setup Instructions](#setup-instructions)
   - [Start LocalStack](#start-localstack)
   - [Set Environment Variables](#set-environment-variables)
5. [Terraform Configuration](#terraform-configuration)
   - [Main Configuration](#main-configuration)
   - [Apply Terraform](#apply-terraform)
6. [Go Programs](#go-programs)
   - [S3 Program](#s3-program)
   - [SQS Program](#sqs-program)
7. [Test the Setup](#test-the-setup)
   - [Test S3](#test-s3)
   - [Test SQS](#test-sqs)
8. [Troubleshooting](#troubleshooting)
   - [S3 Issues](#s3-issues)
   - [SQS Issues](#sqs-issues)
9. [Go SDK Issues](#go-sdk-issues)
10. [Known Issues & Workarounds](#known-issues--workarounds)

---

## 🌟 Overview

This project demonstrates how to emulate AWS services locally using **LocalStack**, provision infrastructure with **Terraform**, and interact with services via **Go**. It’s designed for developers who want to test AWS workflows without incurring cloud costs. Key components include:

- **LocalStack**: Simulates AWS S3 and SQS on `localhost:4566`.
- **Terraform**: Creates an S3 bucket (`my-test-bucket`) locally.
- **Go Programs**:
  - `cmd/s3/main.go`: Uploads and retrieves a file (`go.mod`) from S3.
  - `cmd/sqs/main.go`: Sends, receives, and deletes messages from an SQS queue (`my-custom-sqs-queue`).
- **Environment**: Optimized for WSL2 on Windows with Docker Desktop, tested with LocalStack `4.4.1.dev15` (with a recommendation for `3.8.1` for stability).

The goal is to provide a lightweight, cost-free environment for testing AWS integrations.

---

## 🛠️ Prerequisites

Before starting, ensure the following tools are installed and configured:

- **Docker Desktop**:
  - Install from [Docker Desktop](https://www.docker.com/products/docker-desktop/).
  - Enable WSL2 integration in Docker Desktop settings under **Resources > WSL Integration**.
- **Terraform** (>= 1.5.0):

```bash
  sudo apt-get update && sudo apt-get install -y gnupg software-properties-common
  wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg
  echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
  sudo apt-get update && sudo apt-get install terraform
  terraform -version
```
---

## Go

```bash
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

## AWS CLI

```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
aws --version
```

---

## Project Structure

```
awslocalstack/
├── cmd/
│   ├── s3/
│   │   └── main.go          # Go program to interact with S3
│   └── sqs/
│       └── main.go          # Go program to interact with SQS
├── terraform/
│   └── localstack/
│       └── main.tf          # Terraform config for S3 bucket
├── go.mod                   # Go module dependencies
└── go.sum
```

---

## Setup Instructions

### 1. Start LocalStack

Run LocalStack in a Docker container:

```bash
docker run -d --name localstack-main -p 4566:4566 localstack/localstack:4.4.1.dev15
```

Verify LocalStack is running:

```bash
curl http://localhost:4566/_localstack/health
```

Expected output includes:

```json
{
  "services": {
    "s3": "running",
    "sqs": "available",
    ...
  },
  "edition": "community",
  "version": "4.4.1.dev15"
}
```

### 2. Set Environment Variables

```bash
export AWS_ACCESS_KEY_ID=dummy
export AWS_SECRET_ACCESS_KEY=dummy
export AWS_DEFAULT_REGION=us-west-2
export LOCALSTACK_ENDPOINT=http://localhost:4566
export S3_LOCALSTACK_ENDPOINT=http://s3.localhost.localstack.cloud:4566
export S3_BUCKET=my-test-bucket
export SQS_QUEUE=my-custom-sqs-queue
export SQS_QUEUE_URL=http://sqs.us-west-2.localhost.localstack.cloud:4566/000000000000/my-custom-sqs-queue

export TF_VAR_access_key=${AWS_ACCESS_KEY_ID}
export TF_VAR_secret_key=${AWS_SECRET_ACCESS_KEY}
export TF_VAR_region=${AWS_DEFAULT_REGION}
export TF_VAR_s3_localstack_endpoint=${S3_LOCALSTACK_ENDPOINT}
export TF_VAR_localstack_endpoint=${LOCALSTACK_ENDPOINT}
export TF_VAR_bucket_name=${S3_BUCKET}
export TF_VAR_sqs_queue_name=${SQS_QUEUE}
```

Verify:

```bash
printenv | grep -E 'AWS|SQS|S3|LOCALSTACK|TF_VAR'
```

To persist variables, add them to `~/.bashrc`:

```bash
echo 'export AWS_ACCESS_KEY_ID=dummy' >> ~/.bashrc
# Add all other variables similarly
source ~/.bashrc
```

---

## Terraform Configuration

The Terraform configuration (`terraform/localstack/main.tf`) creates an S3 bucket (`my-test-bucket`).

### `main.tf`

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

variable "access_key" {
  type    = string
  default = "dummy"
}
variable "secret_key" {
  type    = string
  default = "dummy"
}
variable "region" {
  type    = string
  default = "us-west-2"
}
variable "localstack_endpoint" {
  type    = string
  default = "http://localhost:4566"
}
variable "bucket_name" {
  type    = string
  default = "my-test-bucket"
}

provider "aws" {
  access_key                  = var.access_key
  secret_key                  = var.secret_key
  region                      = var.region
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  endpoints {
    s3 = var.localstack_endpoint
  }
}

resource "aws_s3_bucket" "test-bucket" {
  bucket = var.bucket_name
}
```

---

## Apply Terraform

Navigate to the Terraform directory:

```bash
cd terraform/localstack
```

Initialize Terraform:

```bash
terraform init
```

Apply the configuration:

```bash
terraform apply
```

If the bucket already exists:

```bash
terraform import aws_s3_bucket.test-bucket my-test-bucket
terraform apply
```

Verify the bucket:

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} s3api list-buckets
```

Expected output includes:

```json
{
  "Buckets": [
    {
      "Name": "my-test-bucket",
      "CreationDate": "2025-05-16T..."
    }
  ]
}
```

---

## Go Programs

### 1. S3 Program (`cmd/s3/main.go`)

This program uploads and retrieves a file (`go.mod`) from the S3 bucket.

```go
package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
    ctx := context.Background()

    bucketName := os.Getenv("S3_BUCKET")
    localstackEndpoint := os.Getenv("LOCALSTACK_ENDPOINT")
    region := os.Getenv("AWS_DEFAULT_REGION")

    fmt.Printf("S3_BUCKET: %s\n", bucketName)
    fmt.Printf("LOCALSTACK_ENDPOINT: %s\n", localstackEndpoint)
    fmt.Printf("AWS_DEFAULT_REGION: %s\n", region)

    if bucketName == "" || localstackEndpoint == "" || region == "" {
        log.Fatal("Missing required environment variables.")
    }

    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion(region),
        config.WithCredentialsProvider(aws.AnonymousCredentials{}),
    )
    if err != nil {
        log.Fatalf("Failed to load AWS config: %v", err)
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.BaseEndpoint = aws.String(localstackEndpoint)
        o.UsePathStyle = true
    })

    output, err := client.GetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String("go.mod"),
    })
    if err != nil {
        log.Fatalf("Failed to get object: %v", err)
    }
    defer output.Body.Close()

    b, err := io.ReadAll(output.Body)
    if err != nil {
        log.Fatalf("Failed to read object: %v", err)
    }

    fmt.Println("Object content:")
    fmt.Println(string(b))
}
```

#### Run the S3 Program

```bash
echo "module test" > go.mod
aws --endpoint-url=${S3_LOCALSTACK_ENDPOINT} s3 cp go.mod s3://${S3_BUCKET}/go.mod
go run cmd/s3/main.go
```

Expected output:

```
S3_BUCKET: my-test-bucket
LOCALSTACK_ENDPOINT: http://localhost:4566
AWS_DEFAULT_REGION: us-west-2
Object content:
module test
```

---

### 2. SQS Program (`cmd/sqs/main.go`)

This program sends, receives, and deletes a message from the SQS queue.

```go
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

    queueUrl := os.Getenv("SQS_QUEUE_URL")
    region := os.Getenv("AWS_DEFAULT_REGION")
    localstackEndpoint := os.Getenv("LOCALSTACK_ENDPOINT")

    if queueUrl == "" || region == "" || localstackEndpoint == "" {
        log.Fatal("Missing required environment variables.")
    }

    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion(region),
        config.WithCredentialsProvider(aws.AnonymousCredentials{}),
    )
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
        o.BaseEndpoint = aws.String(localstackEndpoint)
    })

    msg := "Hello Luis Rosada we are using SQS in LocalStack!"
    sendOutput, err := client.SendMessage(ctx, &sqs.SendMessageInput{
        QueueUrl:    aws.String(queueUrl),
        MessageBody: aws.String(msg),
    })
    if err != nil {
        log.Fatalf("Failed to send message: %v", err)
    }
    fmt.Println("Message sent. ID:", *sendOutput.MessageId)

    for attempt := 1; attempt <= 3; attempt++ {
        fmt.Printf("Attempt %d: Receiving message...\n", attempt)
        receiveOutput, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
            QueueUrl:            aws.String(queueUrl),
            MaxNumberOfMessages: 1,
            WaitTimeSeconds:     20,
            VisibilityTimeout:   0,
        })
        if err != nil {
            log.Fatalf("Failed to receive message: %v", err)
        }

        if len(receiveOutput.Messages) == 0 {
            fmt.Println("No message received.")
            time.Sleep(2 * time.Second)
            continue
        }

        msg := receiveOutput.Messages[0]
        fmt.Println("Received message:", *msg.Body)

        _, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
            QueueUrl:      aws.String(queueUrl),
            ReceiptHandle: msg.ReceiptHandle,
        })
        if err != nil {
            log.Fatalf("Failed to delete message: %v", err)
        }

        fmt.Println("Message deleted successfully.")
        break
    }
}
```

* Environment variables set:

  ```bash
  export LOCALSTACK_ENDPOINT=http://localhost:4566
  export AWS_DEFAULT_REGION=us-west-2
  export SQS_QUEUE=my-custom-sqs-queue
  export S3_BUCKET=my-test-bucket
  export S3_LOCALSTACK_ENDPOINT=http://localhost:4566
  ```

---

## ▶️ Step 1: Start LocalStack

```bash
docker run -d --name localstack-main -p 4566:4566 localstack/localstack:3.8.1
```

> ✅ Tip: Using version `3.8.1` avoids known SQS issues in newer dev releases.

---

## ▶️ Step 2: Create SQS Queue

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs create-queue --queue-name ${SQS_QUEUE}
```

---

## ▶️ Step 3: Run the Go Program

```bash
go run cmd/sqs/main.go
```

**Expected Output:**

```text
SQS_QUEUE: my-custom-sqs-queue
SQS_QUEUE_URL: http://sqs.us-west-2.localhost.localstack.cloud:4566/000000000000/my-custom-sqs-queue
LOCALSTACK_ENDPOINT: http://localhost:4566
AWS_DEFAULT_REGION: us-west-2
Setting BaseEndpoint: http://localhost:4566
Mensagem enviada com sucesso. ID: <message-id>
Attempt 1: Receiving message...
Mensagem recebida: Hello Luis Rosada we are using SQS in LocalStack!
Mensagem deletada com sucesso.
```

---

## ✅ Test the Setup

### 🪣 Test S3

#### List Buckets

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} s3api list-buckets
```

#### Upload a File

```bash
echo "module test" > go.mod
aws --endpoint-url=${S3_LOCALSTACK_ENDPOINT} s3 cp go.mod s3://${S3_BUCKET}/go.mod
```

#### List Bucket Contents

```bash
aws --endpoint-url=${S3_LOCALSTACK_ENDPOINT} s3 ls s3://${S3_BUCKET}
```

**Expected Output:**

```text
2025-05-16 12:45:32         12 go.mod
```

#### Retrieve File

```bash
aws --endpoint-url=${S3_LOCALSTACK_ENDPOINT} s3api get-object --bucket ${S3_BUCKET} --key go.mod output.txt
cat output.txt
```

---

### 📬 Test SQS

#### List Queues

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs list-queues
```

**Expected Output:**

```json
{
  "QueueUrls": [
    "http://sqs.us-west-2.localhost.localstack.cloud:4566/000000000000/my-custom-sqs-queue"
  ]
}
```

#### Send a Message

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs send-message \
  --queue-url ${SQS_QUEUE_URL} \
  --message-body "Test message"
```

#### Receive a Message

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs receive-message \
  --queue-url ${SQS_QUEUE_URL} \
  --max-number-of-messages 1 \
  --wait-time-seconds 20
```

#### Get Queue Attributes

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs get-queue-attributes \
  --queue-url ${SQS_QUEUE_URL} \
  --attribute-names All
```

---

## 🧪 Troubleshooting

### S3 Issues

#### ❗ `BucketAlreadyOwnedByYou`

```bash
terraform import aws_s3_bucket.test-bucket my-test-bucket
terraform apply
```

#### ❗ Cannot connect to localhost:4566

```bash
docker ps
curl http://localhost:4566/_localstack/health
```

#### ❗ File Not Found in S3

Ensure `go.mod` is uploaded:

```bash
aws --endpoint-url=${S3_LOCALSTACK_ENDPOINT} s3 ls s3://${S3_BUCKET}
```

---

### SQS Issues

#### ❗ `NonExistentQueue` Error

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs list-queues
```

If missing:

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs create-queue --queue-name ${SQS_QUEUE}
```

#### ❗ No Messages Received

* Check queue attributes:

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs get-queue-attributes \
  --queue-url ${SQS_QUEUE_URL} \
  --attribute-names All
```

* Set visibility timeout:

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs set-queue-attributes \
  --queue-url ${SQS_QUEUE_URL} \
  --attributes VisibilityTimeout=0
```

* Purge and retry:

```bash
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs purge-queue --queue-url ${SQS_QUEUE_URL}
aws --endpoint-url=${LOCALSTACK_ENDPOINT} sqs send-message --queue-url ${SQS_QUEUE_URL} --message-body "Test"
```

#### ❗ Connection Issues

```bash
docker logs localstack-main | grep sqs
```

---

### 🧰 Go SDK Issues

#### ❗ Missing Env Variables

```bash
printenv | grep -E 'AWS|SQS|S3|LOCALSTACK'
```

> ✅ Hardcode for testing:

```go
queueUrl := "http://sqs.us-west-2.localhost.localstack.cloud:4566/000000000000/my-custom-sqs-queue"
```

#### ❗ SDK Errors

Make sure you're using the correct versions:

```bash
go get github.com/aws/aws-sdk-go-v2@v1.30.3
go get github.com/aws/aws-sdk-go-v2/config@v1.27.27
go get github.com/aws/aws-sdk-go-v2/service/sqs@v1.38.5
```

---

## 🧪 Known Issues & Workarounds

| Issue                                        | Solution                                                 |
| -------------------------------------------- | -------------------------------------------------------- |
| SQS `ReceiveMessage` returns nothing         | Increase retry attempts or set `VisibilityTimeout=0`.    |
| WSL2 can't resolve `host.docker.internal`    | Use `localhost:4566` instead.                            |
| SQS queue behaves inconsistently             | Downgrade LocalStack to a stable version (`3.8.1`).      |
| No internet in Go container using LocalStack | Add `network_mode: host` to Docker Compose (Linux only). |

---




