use actix_web::{get, App, HttpServer, Responder};
use futures::StreamExt;
use rdkafka::consumer::{Consumer, StreamConsumer};
use rdkafka::message::Message;
use rdkafka::{ClientConfig};
use serde::Deserialize;
use std::time::Duration;
use serde::de::{self, Deserializer, Visitor};

#[derive(Debug)]
enum MessageState {
    MessageCompleted,
    MessageProcessing,
    MessageFailed,
}

impl<'de> serde::Deserialize<'de> for MessageState {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let val = i32::deserialize(deserializer)?;
        match val {
            0 => Ok(MessageState::MessageCompleted),
            1 => Ok(MessageState::MessageProcessing),
            2 => Ok(MessageState::MessageFailed),
            _ => Err(de::Error::custom(format!("invalid state: {}", val))),
        }
    }
}


#[derive(Debug, Deserialize)]
struct MessagePayload {
    state: MessageState,
}

#[get("/")]
async fn index() -> impl Responder {
    "Consumer running!"
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    // Start Actix server
    let server = HttpServer::new(|| {
        App::new()
            .service(index)
    })
    .bind(("127.0.0.1", 8080))?
    .run();

    // Start Kafka consumer in a background task
    tokio::spawn(async move {
        consume_kafka().await;
    });

    server.await
}

async fn consume_kafka() {

    let brokers = "localhost:29092";
    let group_id = "rust-consumer-group";
    let topics = ["gods"];

    let consumer: StreamConsumer = ClientConfig::new()
        .set("group.id", group_id)
        .set("bootstrap.servers", brokers)
        .set("enable.partition.eof", "false")
        .set("session.timeout.ms", "6000")
        .set("enable.auto.commit", "true")
        .create()
        .expect("Consumer creation failed");

    consumer
        .subscribe(&topics)
        .expect("Can't subscribe to topics");

    println!("🚀 Kafka consumer started. Listening on topic 'gods'...");

    let mut stream = consumer.stream();
    while let Some(message) = stream.next().await {
        match message {
            Ok(m) => {
                if let Some(payload) = m.payload() {
                    match serde_json::from_slice::<MessagePayload>(payload) {
                        Ok(parsed) => {
                            println!("✅ Received message: {:?}", parsed);
                        }
                        Err(e) => {
                            eprintln!("⚠️ Failed to parse message: {}", e);
                        }
                    }
                }
            }
            Err(e) => {
                eprintln!("Kafka error: {}", e);
                tokio::time::sleep(Duration::from_secs(1)).await;
            }
        }
    }
}