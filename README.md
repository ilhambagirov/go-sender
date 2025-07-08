# Go sender

## Table of Contents

- [About](#about)  
- [Features](#features)  
- [Requirements](#requirements)  
- [Installation](#installation)
- [API](#api)  

---

## About

Message Sender is a Go-based microservice that automatically processes and dispatches unsent messages from a PostgreSQL database. It exposes HTTP endpoints to start and stop the dispatch loop, and to retrieve message records (content, recipient phone number, sent status). Optionally, it can cache each sent message’s ID and timestamp in Redis for fast deduplication and monitoring.

---

## Features

- ✅ Automatic message dispatch  
- ✅ Start/Stop control endpoints  
- ✅ Query unsent messages
- ✅ Automatic new message insertion to database per 45 seconds
- 🚀 Pluggable sender strategies (Webhook, Twilio, etc.)  
- 🗄️ Redis caching of `messageId` + timestamp  

---

## Requirements

- Go 1.23+  
- Docker & Docker Compose
- PostgreSQL  
- Redis (optional)  

---

## Installation

1. Clone the repo  
   ```bash
   git clone https://github.com/ilhambagirov/go-sender.git
   cd go-sender

2. Copy example env file
   ```bash
   cp .env.example .env

3. Edit .env and set your values:
   ```dotenv
   PORT=9000
   
   DB_HOST=postgres
   DB_PORT=5432
   POSTGRES_USER=ilham
   POSTGRES_DB=go_send_db
   POSTGRES_PASSWORD=your_pass
   WEBHOOK_URL=your_webhook_url
   SENDER_TYPE=webhook

   REDIS_HOST=redis
   REDIS_PORT=6379
   REDIS_PASSWORD=your_redis_pass
   REDIS_DB=0

4. Run locally
   ```bash
   docker-compose up --build -d

5. Logs (Optional)
   ```bash
   docker-compose logs -f msg-sender

---

## API

## Swagger

`http://localhost:9000/swagger/index.html#/`

### List Messages

`GET /message`

- **Description**: Retrieve sent message content, recipient phone number, and sent status for each record.  
- **Responses**:
  - `200 OK`  
    ```json
    [
      {
        "content": "Hello, world!",
        "phone": "+1234567890",
        "is_sent": false
      },
      {
        "content": "Another message",
        "phone": "+1987654321",
        "is_sent": true
      }
    ]
    ```
  - `500 Internal Server Error`
 
---

### Start Dispatcher

`POST /start`

- **Description**: Begin the sender loop. If already running, returns an error.  
- **Responses**:
  - `200 OK`  
  - `400 Bad Request`  
    ```json
    { "error": "sender is already running" }
    ```
---

`POST /stop`

- **Description**: Cancel the in-flight dispatch loop. If not running, returns an error.  
- **Responses**:
  - `200 OK`  
  - `400 Bad Request`  
    ```json
    { "error": "sender is not running" }
    ```
