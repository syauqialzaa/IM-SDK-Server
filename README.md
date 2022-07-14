<h1 align="center"><b>IM-SDK-Server</b></h1>

The backend of Instant Messaging Software Development Kit (IM-SDK-Server). Using Go Websocket-Based communication for chat protocols.

## Contents
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Message Protocol](#message-protocol)
- [Consume API]()
- [Run](#run)
- [Get Test With Frontend](#test-with-frontend)
- [Code Structure](#code-structure)
- [ElephantSQL](#elephantsql)
- [References](#references)

## Features
The following are the features in this project.

- Simple login and register
- Modify user avatar
- Search user and interact it
- Search group and interact it
- User-to-user chat
- User-in-group chat
- Create group
- Storing when user online and offline
- Text messages
- Picture messages
- Voice messages
- Video messages
- Clipboard picture
- Sending file
- Screens sharing
- Video calls (p2p video calls based on WebRTC)
- Deleting message for self/all
- leave chat/group chat
- Save data from another server (consume REST API)

## Tech Stack
The following are the technologies used in this project.

| Technolgy                       | Used                          |
| ------------------------------- | ----------------------------- |
| Programming Language            | Go v1.17                      |
| Framework                       | Gin v1.7.7                    |
| Database                        | ElephantSQL                   |
| Websocket                       | Gorilla Websocket v1.5.0      |
| Logs                            | Uber's Zap v1.21.0            |
| Configuration Management        | Viper v1.10.1                 |
| Object Relational Mapping (ORM) | Gorm v1.23.5                  |
| Communication Protocol          | Google's Proto Buffer v1.28.0 |
| Live Reload Apps                | Cosmtrek Air v1.29.0          |
| Distributed Platform            | Apache Kafka v1.32.0          |

## Message Protocol
The following is `message.proto` file

```proto3
syntax = "proto3";
package protocol;

message Message {
    string  avatar          = 1;    // avatar
    string  fromUsername    = 2;    // username who sent message
    string  from            = 3;    // user uuid who sent message
    string  to              = 4;    // target uuid
    string  content         = 5;    // text message content
    int32   contentType     = 6;    // message content type
    string  type            = 7;    // if it is a heartbeat message, the content is heartbeat
    int32   messageType     = 8;    // message type
    string  url             = 9;    // paths for pictures, videos, voices
    string  fileSuffix      = 10;   // file suffix
    bytes   file            = 11;   // in the case of binary images, files, videos, etc.
}
```
**Note**
- Message content type: [1] Text [2] Ordinary file [3] Picture [4] Audio [5] Video [6] Voice chat [7] Video chat
- Message type: [1] Single chat [2] Group chat
- Suffix, if the file suffix cannot be parsed through the binary header, use this suffix

### Generate Proto File
- If you modify message.proto, you need to recompile to generate the corresponding go file.

- If the proto file is not installed locally, you need to install first, otherwise the protoc command connot be found.

- Use gogoprotobuf, install protobuf library
    ```bash
    go get -u google.golang.org/protobuf
    ```

- Install protoc-gen-gogo
    ```bash
    go get -u github.com/golang/protobuf/protoc-gen-go
    ```

- Execute in the root directory
    ```bash
    protoc --gogo_out=.pkg/protocol/*.proto
    ```

- The protobuf frontend follows the protobuf backend 
- Technically the implementation depends on the frontend platform used

## Consume API
IM-SDK-Server can get user data from another REST API server. This feature is made so that users who already have data on other servers do not need to register. The data is not always in the form of user data, you can customize it according to your needs. For configure that thing, customize `pkg/common/response/http_req_res.go`
```go
// the response user data following json data from REST API target
type ResponseUserData struct {
	Username 	string     	`json:"username"`
	Password 	string     	`json:"password"`
	Nickname 	string     	`json:"nickname"`
	Avatar   	string     	`json:"avatar"`
	Email    	string     	`json:"email"`
}
```

Also configure paths url of REST API target in constant file `pkg/common/constant/constant.go`
```go
const (
    ...
    // paths
	BASE_URL		= "http://localhost:8080"
	GET_ALL			= "/students"
    ...
)
```

## Run
- Clone this repository
- Change directory to IM-SDK-Server directory
- Dependencies required

    ```bash
    go mod download
    ```
- Adjust to your needs in `config/config.go`
    ```go
    type Config struct {
        AppName        string		`mapstructure:"APP_NAME"`
        LogPath        string		`mapstructure:"LOG_PATH"`
        LogLevel       string		`mapstructure:"LOG_LEVEL"`
        StaticFile     string		`mapstructure:"STATIC_FILE"`
        ChannelType    string		`mapstructure:"CHANNEL_TYPE"`
        KafkaHost      string		`mapstructure:"KAFKA_HOST"`
        KafkaTopic     string		`mapstructure:"KAFKA_TOPIC"`
    }
    ```
- Then create `.env` file
    ```env
    APP_NAME        = "gin-chat-svc"
    LOG_PATH        = "logs/chat.log"
    LOG_LEVEL       = "debug"
    STATIC_FILE     = "static/file/"
    CHANNEL_TYPE    = "gochannel"
    KAFKA_HOST      = "kafka:9092"
    KAFKA_TOPIC     = "go-chat-message"
    ```
- For configure to connect to database, check right [here](#elephantsql)
- Then run the program
    ```bash
    go run .
    ```
- If you want to run the program with Cosmtrek Air, check right [here](https://github.com/cosmtrek/air)

## Test With Frontend
This IM-SDK-Server has been tested with a frontend. The frontend repository can be found right [here](https://github.com/kone-net/go-chat-web). Before you using the frontend repo, maybe you must do some configure that following IM-SDK-Server endpoints and params.

## Code Structure
```
├── app
│ ├── access                => prevent cors error and recovery
│ ├── api                   => api functions
│ ├── model                 => database model, one-to-one correspondence with tables
│ ├── service               => class called by service
│ ├── ws                    => websocket, 
│ ├── app.go                => the server core
│ └── router.go             => gin and controller class to bind
├── config                  => system global configuration, file configuration class
├── pkg             
│ ├── common                => constant variables, client-side req, server-side res, tools
│ │ ├── constant
│ │ ├── request
│ │ ├── response
│ │ └── suffix
│ ├── logger                => the log class encapsulated by global
│ ├── misprint              => encapsulated exception class
│ └── protocol              => protoc buffer automatically generated file, defined protoc buffer field
├── static                  => uploaded file, etc.
│ └── file
├── utility
│ ├── db                    => connect to database
│ └── kafka                 => kafka consumers and producers
├── config.env
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## ElephantSQL
ElephantSQL is a PostgreSQL database hosting service, for more details about ElephantSQL, you can check the official website right [here](https://www.elephantsql.com/docs/index.html).
- ElephantSQL configuration, first of all add a mapstructure varible in `config/config.go`
    ```go
    type Config struct {
        ...
        ElephantSQL		string		`mapstructure:"ELEPHANT_SQL"`
        ...
    }
    ```
- Then add your url ElephantSQL instance in `.env` file
    ```env
    ...
    ELEPHANT_SQL    = "<put your url right here>"
    ...
    ```
- For use using PostgreSQL, then change using the configuration as in `config/psql_config.go` and `utility/db/psql_db.go`

## References
- [https://github.com/kone-net/go-chat](https://github.com/kone-net/go-chat)
- [https://github.com/kone-net/go-chat-web](https://github.com/kone-net/go-chat-web)
- [https://github.com/nekonako/moechat](https://github.com/nekonako/moechat)
