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
