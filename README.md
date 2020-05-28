# Micro's POC
Micro is a set of tools that helps developers build and manage microservices. It contains two parts:

- [go-micro](https://github.com/micro/go-micro): A Golang microservices development framework. It is 
the core. By leveraging it, developers could build microservices quickly. The typical type of these 
microservices is gRPC.
- [micro](https://github.com/micro/micro): a command-line tool. Although not mandatory, it provides 
a lot of convenience for Micro development and management. For example, template project generating, 
runtime status inspecting, and services invoking. This tool is also based on go-micro.

In addition, [go-plugins](https://github.com/micro/go-plugins) are necessary in most cases, which is 
a series of plugins. It provides many different choices involving service discovery, 
asynchronous messaging, and transport protocols. Since go-micro is designed as plug-in architecture, 
with these plug-ins, a very flexible combination can be achieved to meet various needs. At any time, 
users can develop their own plug-ins for further extension.

![Micro's runtime](https://micro.mu/images/runtime10.svg?load)

– https://itnext.io/micro-in-action-getting-started-a79916ae3cac

## go-micro architecture
The goal of go-micro is to simplify microservices developing and distributed systems building. 
In practice, some work is always needed in every distributed system.
Therefore go-micro abstracts these common tasks into interfaces. This frees developers from 
underlying implementation details, reduces learning and development costs. make it possible 
to build a flexible and robust system very fast.

![Micro's architecture](https://miro.medium.com/max/1400/1*VdeGqjujc-pfL73JGLI3-w@2x.png)

– https://itnext.io/micro-in-action-getting-started-a79916ae3cac

## Server
The `server` module contains a microservice implementation using go-micro, it uses our implementation 
of a message broker connected with Azure.

## Client
The `client` module contains a client implementation using go-micro, a ping/pong client. It connects 
to `server` microservice via gRPC.

## Stream
The `stream` module contains a client implementation using go-micro, a stream client. It connects 
to `server` microservice via gRPC.

## Pub
The `pub` module contains a client implementation using go-micro, a publisher client. It connects 
to Azure through our `azqueue broker` sharing a queue between `server` module and it.

## Broker/AZQueue
The `azqueue` module contains the `pub/sub` implementation through Azure queues, ready to be hooked into 
the runtime of any go-micro server or client.  

## References
- https://micro.mu
- https://micro.mu/docs/index.html
- https://itnext.io/micro-in-action-getting-started-a79916ae3cac