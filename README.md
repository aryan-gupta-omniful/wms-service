# Warehouse Management Service

## Directory Structure 
```
├── configs/
│   └── config.yaml
├── controllers/
│   ├── controller.go
│   └── inventory_controller.go
├── init/
│   └── init.go
├── models/
│   └── allModels.go
├── pkg/
│   ├── db/postgres
│   │   └── postgres.go
│   └── redis/
│       └── redis.go
├── routes/
│   └── routes.go
├── .env
├── .gitignore
├── docker-compose.yml
├── go.mod
├── go.sum
└── main.go
```

## Project Setup

1. Setup Dependencies
```go
  go mod tidy
```
2. Run Kafka Docker Image
```go
  docker compose up
```
3. Run Main Program
```go
  go run main.go
```
4. Make Database: ``` wms_service_db ``` in PostgreSQL

## API Endpoints

### Get All Hubs
- Endpoint: http://localhost:8081/api/v1/hubs
- Method: GET
- Description: Creates a new product.
  
![image](https://github.com/user-attachments/assets/76c67b79-a8a0-44bf-aed7-1294ae5a257a)

### Get Hub by ID
- Endpoint: http://localhost:8081/api/v1/hubs/:id
- Method: GET
- Description: Fetch Hub by ID.

![image](https://github.com/user-attachments/assets/50160b70-6e27-49eb-a70d-b859f2d55a46)

### Create New Hub
- Endpoint: http://localhost:8081/api/v1/hubs
- Method: POST
- Description: Creates a new Hub.
```sh
    {
      "id": 14,
      "tenant_id": 2,
      "manager_name": "Manager Z",
      "manager_contact": "9272862882",
      "manager_email": "manager.z@cloudtail.com"
    }
 ```

![image](https://github.com/user-attachments/assets/496d929f-1292-4181-932e-dc94ddf7cf82)

### Get All SKUs
- Endpoint: http://localhost:8081/api/v1/skus
- Method: GET
- Description: Get all SKUs.
  
![image](https://github.com/user-attachments/assets/56f4cb9a-5202-470b-90d4-696b6c579308)

### Get SKU by ID
- Endpoint: http://localhost:8081/api/v1/skus/:id
- Method: GET
- Description: Fetch Customer by ID.

![image](https://github.com/user-attachments/assets/c2674615-296f-4245-ad09-bb1efb9853d6)

### Create New SKU
- Endpoint: http://localhost:8081/api/v1/skus
- Method: POST
- Description: Creates a new customer.
```sh
    {
      "id": 31,
      "hub_id": 7,
      "seller_id": 8,
      "product_id": 15,
      "images": "earphone.jpg",
      "description": "Boat Earphones",
      "unit_price": 0,
      "fragile": false,
      "dimensions": "12x1x1"
    }
 ```

## Previews
- Hub ID and SKU ID Validation
![image](https://github.com/user-attachments/assets/c0e6fc4d-11ef-4524-bb80-5020396a175f)

- Inventory Check and Updation
![image](https://github.com/user-attachments/assets/30269fec-1af3-4fc0-9e26-35f3ed0103d4)


