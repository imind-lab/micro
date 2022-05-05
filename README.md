## Imind
Imind is a framework based on gRPC for microservices development 

#### Overview
Imind provides the core requirements for microservices development. Microservices communicate with each other based on gRPC, Service Registry and Service Discovery using the service of kubernetes, decouple and dynamically update the configuration file using the configmap and secret of kubernetes, and service circuit breaker, degradation and traffic limit using Sentinel developed by Alibaba, service distributed tracing using OpenTelemetry and Jaeger. Imind scaffolding can automatically generate common microservice code.

#### Features
Go Micro abstracts away the details of micro services. Here are the main features.
- Base on gRPC

  Imind implements communication between services based on gRPC and simplifies gRPC operations for developers. Integrated gRPC-Gateway for restful API support.

- Distributed tracing

  Imind is based on OpenTelemetry and Jaeger to implements distributed tracing, which can automatically implements link tracing among services and easily implements the tracing among methods within services manually.

- Circuitbreaker,Degradation and Traffic limit

  Imind is based on Sentinel developed by Alibaba implements service circuit breaker, degradation and traffic limit. It supports dynamic data sources and can implement dynamic updating of rules.

- Distributed log

  Imind injects trace id into logs through a customized gRPC interceptor and automatically collects application information for easy program debugging and log analysis

- Kubernetes deployment.

  Imind automatically generates helm chart, implements one-click deployment of applications to kubernetes, and supports gitlab ci.

- Customize tag in Proto

  Imind supports custom tag in proto, which can verify request parameters based on tag rules and unify response format based on tag.

#### Getting Started

Install scaffolding and use scaffolding to generate sample microservice code. Before you do that, you need to complete the [environment configuration](docs/prerequisite.md)

##### 1. Install scaffolding

```shell
go install https://github.com/imind-lab/micro/microctl@latest

# View scaffolding help
microctl help

# View scaffolding version
microctl version
```

##### 2. Run scaffolding to generate the sample microservice

```shell
cd $GOPATH/src

# init subcommand to initialize the microservice
# -d Code repository domain name, default is github.com
# -p project name, default is imind-lab
# -s microservice name, default is greeter
# -l service type. Currently api(aggregation service) and srv(backend service) are supported, default is srv
microctl init -d gitlab.imind.tech -p daniel -s greeter -l srv

cd gitlab.imind.tech/daniel/greeter

# Deploy the Greeter service to Kubernetes(You need to have operational Kubernetes)
make deploy
```

#### Advanced usage

To be continued...