# Go Friend Management REST API Example
A RESTful API example for simple Friend Management application with Go, using **gin-gonic/gin** (A nice HTTP framework) and **docker** for deployment

## Structure
```
├── src
│   ├── main.go
│   ├── common           
│   │   └── util.go                         // Utility functions such as: check valid email, remove item in slice, etc...
│   │
│   ├── data                                // Data Access layer
│   │   ├── repository_*_mock.go            // Fake instance of repository files, use for service layer testing
│   │   ├── repository_*.go                 // Provide CRUD functions, work with database directly
│   │   └── db_config.go                    // Database configuration
│   │
│   ├── services                            // Service layer
│   │	├── service_*.go                    // Provide CRUD functions, work as middle layer between repository layer and endp
│   │   ├── service_*_mock.go               // Fake instance of service files, serve for testing trategy
│   │   └── service_*_test.go               // Handle service test cases
│   │
│   ├── endpoints                           // API layer
│   │   ├── endpoint_config.go              // Database configuration
│   │   ├── base_endpoint.go                // Standard functions for API function: responseOK, responseError
│   │   ├── user_endpoint_test.go           // Handle User's API test cases
│   │   ├── relationship_endpoint_test.go   // Handle Relationship's API test cases
│   │   ├── user_endpoint.go                // User's API
│   │   └── relationship_endpoint.go        // Friend Activities's API
│   │
│   ├── models         
│   │   └── *.go // Models for our application
│   │
│   └── docs     // API Documentation by Swagger
│       └── *.go // Auto generate by Swagger library which takes responsible to document API
│
├── db_migration
│   └── db_create.go // Script file to create database for testing
│
├── docker_compose.yml // Mandatory docker file
│
└── Dockerfile // Mandatory docker file
```

## Installation & Run

#### Enviroment
This project uses phpMyadmin database inside with docker-compose. You can compose with dockerfile or create your own phpMyadmin database without it
You should set the database config with yours or leave it as default on [db_config.go](https://github.com/s3corp-github/SP_FriendManagementAPI_Golang_KyTruong/blob/master/src/data/db_config.go)
```go
dbDriver := "mysql"
dbUser := "root"
dbPass := "123456@x@X"
dbName := "friendMgmt"
dbPort := "3306"
dbHost := "fullstack-mysql"
```

For run docker-compose, run these following commands in project's root folder:

```bash
docker-compose build
docker-compose up
```

#### phpMyadmin Database
```bash
# http://localhost:9090/
```
There will be empty database named **friendMgmt** and it's not ready yet. You need to create tables and sample data, the migration script is provided [here](https://github.com/s3corp-github/SP_FriendManagementAPI_Golang_KyTruong/blob/master/db_migration/db_create.sql)
Once finished these steps, the app is ready to go.

#### API Endpoint
```bash
# http://localhost:8081/swagger/index.html
```

## API Documentation
This is API self documentation by using Swagger. You can test all of them by expand specific api then click on try it out button.

--photo here


## Test Coverage
All APIs have been tested carefully by mocking strategy. 

--photo

## Achievement

- [x] Write the tests for all APIs.
- [x] Organize the code with packages
- [x] Make docs with Swagger
- [x] Building a deployment process 


