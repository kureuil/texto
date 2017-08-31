# texto - Simple Messaging Service

Texto is a simple directed chat server written in Go. Clients communicate with the server using the
[WebSocket Protocol](http://tools.ietf.org/html/rfc6455).

## Running

Prerequisites:

* docker & docker-compose 1.13+

```bash
$ git clone https://github.com/kureuil/texto.git
$ cd texto
$ docker-compose build
$ docker-compose up
```

The texto server should now be running. The server is bundled with a sample client written in JS, which if you used
docker, should be accessible here: http://localhost:8398/

## Architecture

This messaging server is built to be as flexible as possible. For this reason, there is an nginx proxy in front of the
application, allowing for the load balancing of the requests. Also, the messaging server doesn't store the messages (if
you send a message to an inexistant user, it won't ever be sent to anyone), and relies on Redis for the message
dispatching.

When a message is sent to the messaging server, it is transformed into a simpler message and published on the
recipient's redis channel. All messaging servers listen on every channel, and if a message to meant for a user known on
the current node, it is relayed.

This makes the system resilient to failure, if a messaging server is malfunctioning or stops you just have to start a
new one and register it into your load balancer (probably via your service discovery daemon). On the database side,
Redis provides a *Sentinel* mode which allow for easy replication and master-reelection in case of failure.

Finally, splitting the messaging server from the database allow for easier experimentation and gradual deployment. You
can easily migrate the instance one-by-one and see the effects of your upgrade in real-time, and easily rollback in case
of malfunction.

```
                               +------------------+          +------------------------+
                               |                  |          |                        |
                               |                  |          |                        |
                        +------> Messaging Server <-----+    | Redis Server           |
                        |      |                  |     |    | Replica 1              |
                        |      |                  |     |    |                        |
                        |      +------------------+     |    +-----------^------------+
                        |                               |                |
+-------------------+   |      +------------------+     |    +-----------v------------+
|                   <---+      |                  |     +---->                        |
|                   |          |                  |          |                        |
| NGINX             <----------> Messaging Server <----------> Redis Server           |
| Load Balancer     |          |                  |          | Master                 |
|                   <---+      |                  |     +---->                        |
+-------------------+   |      +------------------+     |    +-----------^------------+
                        |                               |                |
                        |      +------------------+     |    +-----------v------------+
                        |      |                  |     |    |                        |
                        |      |                  |     |    |                        |
                        +------> Messaging Server <-----+    | Redis Server           |
                               |                  |          | Replica 2              |
                               |                  |          |                        |
                               +------------------+          +------------------------+
```
*Example deployment*

## API Endpoints

If you wish to create your own client, you first need to be able to establish a WebSocket connection, as this is the
main channel of communication between a client and a server.

*All the exemples are written in Javascript, except when stated otherwise.*

### `/v1/texto`

This is the first version of the Texto Messaging Protocol. It relies on JSON messages sent through a WebSocket
connection.

#### Message Schema

Every JSON message follows the same schema:
```json
{
    # The client_id field stores the UUID of the current client/session.
    "client_id": "754cd3a0-27b3-4c51-a66e-466fed82b667",
    # The id field stores the UUID of the current message. When the server sends a response, it will use the id of the
    # request message.
    "id": "8a15b000-02d7-4823-8336-0cd0b0b13ae9",
    # The kind field indicates the type of the current message.
    "kind": "ack",
    # The data field is a dynamically shaped field which content depends on the kind field.
    "data": null
}
```

#### Message Kinds

##### `error`

The `error` message kind indicates that an error was encountered when processing a request.

**Payload**
```json
{
    # The code field stores the internal code of the error, which is intended for programmatic use.
    "code": "ENOMEM",
    # The description field stores a human readable of the error and can be displayed safely to a user or logged.
    "description": "Out-of-memory"
}
```

##### `registration`

The `registration` message kind is sent by the client when it wants to fetch information about its current session.

In the future, it could be used for authentication.

**Payload**
```json
null
```

##### `connection`

The `connection` message kind is sent in response to a `registration` request.

**Payload**
```json
{
    # The client_id field stores the UUID of the current client/session.
    "client_id": "754cd3a0-27b3-4c51-a66e-466fed82b667"
}
```

##### `send`

The `send` message kind is sent when a client wants to send a message to another client.

**Payload**
```json
{
    # The receiver_id field stores the UUID of the message's recipient.
    "receiver_id": "754cd3a0-27b3-4c51-a66e-466fed82b667",
    # The text of the message
    "text": "Lorem ipsum dolor sit amet..."
}
```

##### `receive`

The `receive` message kind is sent by the server when a client is receiving a message.

**Payload**
```json
{
    # The sender_id field stores the UUID of the message's sender.
    "sender_id": "754cd3a0-27b3-4c51-a66e-466fed82b667",
    # The text of the message
    "text": "Lorem ipsum dolor sit amet..."
}
```

##### `ack`

The `ack` message kind is sent to acknowledge of the reception of a `send` or a `receive` message.

**Payload**
```json
null
```

#### Examples

```javascript
// Establishing a WebSocket connection with
let ws = new WebSocket("ws://localhost:8398/v1/texto");
let sessionId = null;
ws.onmessage = function(event) {
    let message = JSON.parse(event.data)
    // Treat first received message independently
    if (sessionId === null && message.kind === 'connection') {
        sessionId = message.client_id;
        return;
    }
    // Process incoming message...
};
```

Once a connection is established with a client, the server sends a message containing the ID of the current session.
