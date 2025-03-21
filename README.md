# Task Management API

A RESTful API service for managing maintenance tasks, built with Go and Fiber framework. This project implements a task
management system with role-based access control, allowing technicians to manage their tasks and managers to oversee all
operations.

## Features

- **Role-Based Access Control**
    - Manager: Can view all tasks, delete tasks, and receive notifications
    - Technician: Can create, view, and update their own tasks

- **Task Management**
    - Create new tasks (max 2500 characters summary)
    - List tasks (filtered by role permissions)
    - Update task status
    - Delete tasks (manager only)
    - Task completion tracking

- **Security**
    - JWT-based authentication
    - Role-based authorization
    - Request rate limiting
    - Secure middleware configuration

## Tech Stack

- **Backend**: Go with Fiber framework
- **Database**: MySQL
- **Containerization**: Docker
- **Authentication**: JWT
- **API Documentation**: OpenAPI/Swagger (coming soon)

## Prerequisites

- Go 1.24+
- Docker and Docker Compose
- MySQL 8.0+

## Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/task-management-api
   cd task-management-api
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env file with your configuration
   ```

3. **Start the application with Docker**
   ```bash
   docker-compose up -d
   ```

   The API will be available at `http://localhost:3000`

## API Endpoints

### Authentication

- `POST /api/v1/users/login` - User login
- `POST /api/v1/users` - Create new user

### Tasks

- `POST /api/v1/tasks` - Create a new task
- `GET /api/v1/tasks` - List tasks (filtered by user role)
- `PUT /api/v1/tasks/:taskPid` - Update task status
- `DELETE /api/v1/tasks/:taskPid` - Delete task (manager only)

## Development

### Local Setup

1. Install dependencies:
   ```bash
   go mod tidy
   go mod vendor
   ```

2. Run the database:
   ```bash
   docker-compose up -d mysql
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

### Testing

Run the test suite:

```bash
go test ./...
```

## Security Features

- Request rate limiting
- CORS protection
- Helmet security headers
- Request compression
- Request ID tracking
- Panic recovery middleware

## Database Schema

The application uses two main tables:

- `users`: Stores user information and roles
- `tasks`: Stores task details with relationships to users

## K8S

```bash
kubectl apply -f k8s-manifest.yaml
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.