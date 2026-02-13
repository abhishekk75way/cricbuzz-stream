# Live Cricket Score Streaming (Cricbuzz-style)

This project is a simple real-time live score backend built using Go and Server-Sent Events (SSE).

If you’ve ever wondered how platforms like Cricbuzz push live scores without refreshing the page, this project demonstrates the same core idea in a clean and minimal way.

No WebSockets.
No polling.
Just efficient server-to-client streaming.

## What This Project Does

* Admin pushes score updates through an API
* The server broadcasts updates instantly
* Clients receive live score changes in real time
* One persistent connection per match

This is ideal for:

* Live cricket scores
* Ball-by-ball updates
* Commentary feeds
* Any read-only real-time stream

## High-Level Architecture

```
Admin (Postman / Curl)
        |
        | POST /admin/update-score
        v
Go Server
        |
        | Broadcast via Broker
        v
SSE Endpoint
        |
        | GET /events?matchId=123
        v
Clients (Browser / Mobile / Curl)
```

Clients subscribe once and stay connected.
The server pushes updates whenever scores change.


## Why SSE Instead of WebSockets?

This system only requires one-way communication: server to client.

Clients do not need to send messages back in real time.

| Feature                  | SSE | WebSocket |
| ------------------------ | --- | --------- |
| Server to Client         | Yes | Yes       |
| Client to Server         | No  | Yes       |
| Auto reconnect           | Yes | No        |
| Complexity               | Low | Higher    |
| Suitable for live scores | Yes | Overkill  |

SSE is simpler and perfectly suited for live score streaming.

## Getting Started

### 1. Run the Server

```bash
go run cmd/main.go
```

The server starts at:

```
http://localhost:8080
```

### 2. Open a Live Stream

In a terminal:

```bash
curl "http://localhost:8080/events?matchId=123"
```

You should see:

```
data: connected
```

This confirms that the SSE connection is working.

### 3. Send a Score Update

In another terminal:

```bash
curl -X POST http://localhost:8080/admin/update-score \
  -H "Content-Type: application/json" \
  -d '{"matchId":"123","score":"100/2","overs":"10.3"}'
```

The first terminal will instantly show:

```
data: {"id":"123","score":"100/2","overs":"10.3"}
```

That is real-time streaming.

## API Endpoints

### GET /events?matchId=123

* Opens an SSE connection
* Keeps the HTTP request open
* Streams updates for that match

Important:

* You must include matchId
* Use curl or a browser
* Postman does not properly display streaming responses

### POST /admin/update-score

Used by admin systems to push new scores.

Request body:

```json
{
  "matchId": "123",
  "score": "100/2",
  "overs": "10.3"
}
```

The server broadcasts this update to all clients subscribed to that match.

## Internal Design Overview

* Each connected client gets its own Go channel
* Clients are grouped by matchId
* When a score update arrives:

  * The broker finds clients for that match
  * Sends the update through their channels
  * Flushes the response immediately

The HTTP connection stays open the entire time.

SSE is essentially HTTP that never closes.

## Common Issues

If something is not working, check the following:

* Missing matchId in the request
* Using Postman to listen to SSE
* Forgetting to call Flush()
* Using unbuffered channels
* Running behind Nginx without disabling buffering

If you see `data: connected`, the stream is working.