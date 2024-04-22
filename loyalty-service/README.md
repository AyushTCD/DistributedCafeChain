
# Loyalty Service System

## Overview

The Loyalty Service System is designed to provide a globally accessible distributed service for a café chain, introducing a 'friends and family' loyalty-card scheme. This system rewards customers with discounts on in-store purchases based on their transaction/purchase history, allowing up to four people to share a loyalty-card account.

## Features

- **User Management**: Register, update, and manage user accounts.
- **Transaction Logging**: Record every purchase transaction to track user spending.
- **Loyalty Accounts**: Create and manage shared loyalty accounts for families or friends.
- **Points System**: Earn points for every purchase, with points convertible into discounts.
- **Discount Calculation**: Automatically calculate available discounts based on the transaction history.

## Technology Stack

- **Backend**: Go (Gin framework)
- **Database**: MySQL
- **Deployment**: Docker for containerization and Docker Compose for multi-container setups

## Project Structure

~~~
/loyalty-service
    /internal
		/account
               service.go     // Account management logic
          /api
               handler.go     // HTTP handlers for the web server
               router.go      // Router setup
          /invitation
               service.go
          /model              // Model definitions for each of the services
               account.go
               invitation.go
               transaction.go
               user.go
          /user
               service.go     // User management logic
          /transaction
               service.go     // Transaction processing logic
    /pkg
          /db
               database.go    // Database connection and initialization
     Dockerfile
     docker-compose.yml
     go.mod
     main.go                  // Entry point for the API server
     README.md
~~~

## Getting Started

### Prerequisites

- Go (version 1.15 or later)
- MySQL (version 8 or later)
- Docker and Docker Compose (for deployment)

### Installation

1. **Clone the Repository**:
~~~
git clone https://github.com/yourusername/loyalty-service.git
cd loyalty-service
~~~

2. **Set Up Environment Variables**:
~~~
MYSQL_URI="isabelle:password@tcp(127.0.0.1:3304)/loyalty_program?charset=utf8mb4&parseTime=True"
~~~

3. **Start MySQL**

See mysql-cluster-init/README.md

4. **Configure Database**

The configuration for the database must be specified in the `loyalty-service.toml` file, for example:
~~~
default = [
	"isabelle:password@tcp(10.100.2.2:3306)/loyalty_program?charset=utf8mb4&parseTime=True",
	"isabelle:password@tcp(10.100.2.3:3306)/loyalty_program?charset=utf8mb4&parseTime=True",
	"isabelle:password@tcp(10.100.2.4:3306)/loyalty_program?charset=utf8mb4&parseTime=True",
]
~~~

This configuration will be mounted into the Docker container automatically

5. **Run the Application**:
~~~
go run main.go
~~~

The service will start running on `http://localhost:8080`.

**OR**

To simplify the setup, use Docker Compose to run the service along with MySQL in containers:
~~~
docker compose up --build
~~~

The Go API servers (3 by default) will be available behind a Traefik load balancer on `http://localhost:8080`

To change the number of servers, do
~~~
docker compose scale <number>
~~~

## API Documentation

Endpoints include:

- POST `/users` - Register a new user
- GET `/users/:id` - Retrieve user details
- POST `/loyalty-accounts` - Create a new loyalty account
- GET `/loyalty-accounts/:id` - Get details of a loyalty account
- POST `/transactions` - Log a new transaction
- POST `/invitations/create` - Create an invitation token
- POST `/invitations/accept` - Accept an invitation to an account
- POST `/invitations/decline` - Decline an invitation to an account

## Test Scenario

### Create user1
~~~
curl -X POST http://localhost:8080/users \
     -H 'Content-Type: application/json' \
     -d '{"name": "John Doe", "email": "john.doe@example.com", "password": "password123"}'
~~~

### Create user2
~~~
curl -X POST http://localhost:8080/users \
     -H 'Content-Type: application/json' \
     -d '{"name": "Jane Doe", "email": "jane.doe@example.com", "password": "password123"}'
~~~

### Create an account and add user1 and user2
- Give them 100 points as a welcome gift
- For every 1 euro spent the user gets 20 points (equivalent to 20 cent)
~~~
curl -X POST http://localhost:8080/loyalty-accounts \
     -H "Content-Type: application/json" \
     -d '{"userIds": ["{user1ID}", "{user2ID}"], "points": 100}'
~~~

### Add a transaction
- User1  buys a coffee for 3.70e
- The account for user1 and user2 will receive 1 point per euro spent (rounded down)
- To spend points instead, append `?usePoints=true` to the URL
~~~
curl -X POST http://localhost:8080/transactions \
     -H 'Content-Type: application/json' \
     -d '{"AccountID": "{accountID}", "UserID": "{userID}", "amount": 3.70}'
~~~

### Create two more users to test invitation functionality
~~~
curl -X POST http://localhost:8080/users \
     -H 'Content-Type: application/json' \
     -d '{"name": "James Joyce", "email": "james.joyce@example.com", "password": "ulysses"}'
~~~
and
~~~
curl -X POST http://localhost:8080/users \
     -H 'Content-Type: application/json' \
     -d '{"name": "Homer Simpson", "email": "homer.simpson@example.com", "password": "donuts"}'
~~~

### User1 creates an invitation for user3
~~~
curl -X POST "http://localhost:8080/invitations/create" \
     -H "Content-Type: application/json" \
     -d '{"email": "james.joyce@example.com", "inviterID": "{user1ID}", "accountID": "{user1AccountID}"}'
~~~

### User1 creates an invitation for user4
~~~
curl -X POST "http://localhost:8080/invitations/create" \
     -H "Content-Type: application/json" \
     -d '{"email": "homer.simpson@example.com", "inviterID": "{user1ID}", "accountID": "{user1AccountID}"}'
~~~

### User3 accepts the invitation
- Invitation token can viewed in the 'invitations' collection
- User3 is now added to user1s account
- The invitation status is updated to 'accepted'
~~~
curl -X POST "http://localhost:8080/invitations/accept" \
     -H "Content-Type: application/json" \
     -d '{"token": "unique_invitation_token", "email": "james.joyce@example.com"}'
~~~

### User4 declines the invitation
- User3 is now added to user1s account
- User4 is not added to user1s account
- The invitation status is updated to 'declined'
~~~
curl -X POST "http://localhost:8080/invitations/decline" \
     -H "Content-Type: application/json" \
     -d '{"token": "unique_invitation_token", "email": "homer.simpson@example.com"}'
~~~


