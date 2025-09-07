# Simple Kafka Producer-Consumer Example

A basic Kafka implementation in Go using Fiber web framework for HTTP endpoints and Confluent's Kafka Go client.

## Features

- **Producer Service**: HTTP endpoint to send messages to Kafka topics
- **Consumer Service**: Background consumer that reads messages from Kafka topics
- **Web Interface**: Simple Fiber web server for health checks

## Project Structure

```
.
├── producer/
│   └── main.go          # HTTP server with Kafka producer
├── consumer/
│   └── main.go          # HTTP server with background Kafka consumer
├── docker-compose.yml   # Kafka + Zookeeper setup
└── README.md
```

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Port 3000, 4000, and 29092 available

## Quick Start

### 1. Start Kafka Infrastructure

```bash
docker-compose up -d
```

Wait for services to be healthy:
```bash
docker-compose ps
```

### 2. Start Consumer Service

```bash
cd consumer
go mod init consumer
go get github.com/confluentinc/confluent-kafka-go/kafka
go get github.com/gofiber/fiber/v2
go run main.go
```

Consumer runs on: http://localhost:3000

### 3. Start Producer Service

```bash
cd producer
go mod init producer  
go get github.com/confluentinc/confluent-kafka-go/v2/kafka
go get github.com/gofiber/fiber/v2
go run main.go
```

Producer runs on: http://localhost:4000

### 4. Send Messages

```bash
# Send batch of messages (every 1.5 seconds)
curl -X POST http://localhost:4000/k-producer
```

Watch the consumer terminal for incoming messages.

## API Endpoints

### Producer Service (Port 4000)

- `GET /` - Health check
- `POST /k-producer` - Send 10,000 messages to "gods" topic (1.5s intervals)

### Consumer Service (Port 3000)

- `GET /` - Health check
- Background process consumes from topics: `gods` and regex pattern `^aRegex.*[Tt]opic`

## Message Format

```json
{
  "State": 0  // 0=Completed, 1=Processing, 2=Failed
}
```

## Configuration

### Kafka Settings

- **Bootstrap Servers**: localhost:29092
- **Consumer Group**: mars
- **Auto Offset Reset**: earliest

### Topics

- `gods` - Main message topic
- Regex pattern: `^aRegex.*[Tt]opic` - Matches topics like "aRegexTestTopic"

## Docker Compose Services

```yaml
# Zookeeper - Kafka coordination
zookeeper:
  - Port: 22181
  
# Kafka Broker  
kafka:
  - Port: 29092 (external)
  - Port: 9092 (internal)
```

## Development

### Create Custom Topics

```bash
docker exec -it kafka kafka-topics \
  --create --topic your-topic \
  --bootstrap-server localhost:29092 \
  --partitions 1 --replication-factor 1
```

### Manual Testing

```bash
# Send test message
docker exec -it kafka kafka-console-producer \
  --topic gods --bootstrap-server localhost:29092

# Read messages manually  
docker exec -it kafka kafka-console-consumer \
  --topic gods --bootstrap-server localhost:29092 \
  --from-beginning
```

## Dependencies

```go
// Producer
github.com/confluentinc/confluent-kafka-go/v2/kafka
github.com/gofiber/fiber/v2

// Consumer  
github.com/confluentinc/confluent-kafka-go/kafka
github.com/gofiber/fiber/v2
```

## Troubleshooting

### Common Issues

1. **Consumer not receiving messages**
   - Check if producer sent messages first
   - Verify Kafka is running: `docker-compose ps`
   - Check topic exists: `docker exec -it kafka kafka-topics --list --bootstrap-server localhost:29092`

2. **Connection refused**
   - Ensure Kafka container is healthy
   - Check port 29092 is accessible
   - Verify bootstrap.servers configuration

3. **Slow message sending**
   - Producer sends 10,000 messages with 1.5s delay (4+ hours total)
   - Reduce message count for testing
   - Consider running producer in background goroutine

### Logs

```bash
# View Kafka logs
docker-compose logs kafka

# View Zookeeper logs  
docker-compose logs zookeeper
```

## Performance Notes

- Current producer sends messages sequentially with 1.5s delay
- For production, consider batch sending and async processing
- Consumer group "mars" will balance load if multiple consumer instances run

## License

MIT