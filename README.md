# InnoTaxiUser API 

This is a microservice for creating and managing taxi orders. The service is written in Go using the Gin web framework.

## Installation

To install and run the service, follow these steps:

Install Go on your system if you haven't already done so.
Clone the repository to your local machine.
Run the following command in the root directory of the project:

    go run ./cmd/main.go

Also you can run project using docker-compose.
The service should now be running on localhost:8080.


## Run the tests

    go test ./internal/service 

## Code Description

# Project structure

The code for the microservice is organized into several packages:

- cmd/main.go contains the main function for the service.
- internal/app/app.go contains functions which sets up the API routes and starts the server.
- models/ contains the data models for the application. In this case, there is only one model - User.
- repositories/ contains the repository implementation for working with the databases. In service there are such databases as postgresql for store data, mongodb for logs and mongodb for cache.
- services/ contains the business logic services for the application.
- handlers/ contains the API request handlers for the application. Service provides handlers for registartion and auth user, also handlers for working with user's profile.


The file user.go contains the User structure, which represents the user data model. In this case, the user contains the following fields:

    type User struct {
        ID          uint64 
        Name        string 
        PhoneNumber string  
        Email       string 
        Raiting     float64
        Status      string 
    }

The Status field is an enumeration of the UserStatus type, which determines the status of the user:

    const (
	    StatusCreated string = "created"
	    StatusDeleted string = "deleted"
    )

Business logic services are located in the internal/services/ package. Service uses repository layer to get data, layer data implement such functions as: 

    type UserRepo interface {
        GetUserById(ctx context.Context, id string) (*model.User, error)
        UpdateUserById(ctx context.Context, id string, user *model.User) error
        DeleteUserById(ctx context.Context, id string) error
    } 

    type AuthRepo interface {
	    CreateUser(ctx context.Context, user UserSingUp) error
	    CheckUserByPhoneNumber(ctx context.Context, phone string) (*UserSingIn, error)
    }

    type TokenRepo interface {
	    AddToken(token string, expired time.Duration) error
	    GetToken(token string) bool
    }

This microservice demonstrates a simple way to create and manage taxi orders using Go and Gin. The code is organized into packages, making it easy to maintain and extend. The `UserService` and `UserHandler` objects provide the business logic and API endpoints, respectively.