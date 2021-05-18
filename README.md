# Mobstore
Mobstore is an initial validation of the microservice architecture, using sync and async strategies to create resilient flows.

## Features
- Infrastructure as code, using terraform for provision resources
- Applications internally use grpc and asynchronous messages for communication
- Distributed database, each application has its own user for the database
- Two services in golang(users, products) and one in python(discounts)

## Main tech
- [terraform](https://www.terraform.io/)
- [mysql](https://www.mysql.com/)
- [golang](http://golang.org/)
- [python](https://www.python.org/)
- [docker](https://www.docker.com/)

## Architecture design
![Architecture design](Mobstore-design.png?raw=true "Design")

User application responsible for storing users and producing events with create, update and delete operations.
Product application responsible for storing products and producing events with create, update and delete operations.
Application discounts consume events produced by users and products through SNS-> SQS consumers and offer discounts based on the user and the product previously stored.

## How to run local

### Requirements
```
terraform=v0.15.3
python=3.9.5
docker=20.10.6
docker-compose=1.29.1
alembic=1.6.2
go=1.16.4
```

### Start
```sh
make run
```
This single command will provision resources with terraform using localstack (SNS, SQS and database users), start a MySQL server and Zipkin server (for spans). After that, users, products and discounts will start.

### Test
[![Run in Insomnia}](https://insomnia.rest/images/run.svg)](https://insomnia.rest/run/?label=Mobstore&uri=https%3A%2F%2Fgithub.com%2Fdalmarcogd%2Fmobstore%2Fblob%2Fmain%2FInsomnia_2021-05-18.json)
