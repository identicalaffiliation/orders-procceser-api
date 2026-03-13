# ORDERS SYSTEM 

### MAIN
The project implements an event-driven architecture for order delivery with direct simulation of payment, delivery, and inventory. The project is fully containerized using Docker to simplify deployment and ensure a reproducible environment.

### TECH STACK
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)](https://www.rabbitmq.com/)

### WORKFLOW
``` txt
Order Service – creates and stores orders, updates their status, and responds to the client via API.

Inventory Service – checks and reserves goods.

Payment Service – simulates payment for orders.

Delivery Service – creates a delivery and updates the order status.
```

### GET START WITH API
``` zsh
1. git clone https://github.com/identicalaffiliation/orders-procceser-api.git
2. cd orders-procceser-api
3. make up
```

### DELETE API
``` zsh
1. make clean
```
