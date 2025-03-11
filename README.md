# Movies-Go API

A RESTful API for managing movie information built with Go, Gin, and PostgreSQL.

## Features

- CRUD operations for movies
- User management with registration and authentication
- JWT-based authentication and authorization
- Role-based access control
- Search functionality with pagination
- Containerized with Docker

## Tech Stack

- **Go**: Programming language
- **Gin**: HTTP web framework
- **Bun**: ORM for PostgreSQL
- **JWT**: Authentication
- **Docker**: Containerization

## Project Structure

```
Movies-Go/
├── cmd/                  # Application entry points
│   └── main.go           # Main application file
├── internal/             # Private application code
│   ├── controller/       # HTTP controllers
│   ├── entity/           # Domain models
│   ├── pkg/              # Shared packages
│   │   ├── auth/         # Authentication utilities
│   │   ├── config/       # Configuration management
│   │   ├── middleware/   # HTTP middleware
│   │   ├── repository/   # Database connection
│   │   └── script/       # Database migrations
│   ├── repository/       # Data access layer
│   ├── router/           # HTTP routes
│   └── util/             # Utility functions
├── conf.yaml             # Configuration file
├── Dockerfile            # Docker build instructions
├── docker-compose.yml    # Docker Compose configuration
└── README.md             # Project documentation
```

## API Endpoints

### Authentication

- `POST /api/movies/v1/auth/register`: Register a new user
- `POST /api/movies/v1/auth/login`: Authenticate user and get JWT token

### Users

- `GET /api/movies/v1/users`: Get all users (requires authentication)
- `GET /api/movies/v1/users/:id`: Get user by ID (requires authentication)
- `PUT /api/movies/v1/users/:id`: Update a user (requires admin role)
- `DELETE /api/movies/v1/users/:id`: Delete a user (requires admin role)

### Movies

- `GET /api/movies/v1/movies`: Get all movies
- `GET /api/movies/v1/movies/:id`: Get movie by ID
- `GET /api/movies/v1/movies/search?q=query&page=1&limit=10`: Search movies
- `POST /api/movies/v1/movies`: Create a new movie (requires authentication)
- `PUT /api/movies/v1/movies/:id`: Update a movie (requires authentication)
- `DELETE /api/movies/v1/movies/:id`: Delete a movie (requires authentication)

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL (if running locally)

### Running Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/Movies-Go.git
   cd Movies-Go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Configure the database in `conf.yaml`

4. Run the application:
   ```bash
   go run cmd/main.go
   ```

### Running with Docker

1. Build and start the containers:
   ```bash
   docker-compose up -d
   ```

2. The API will be available at `http://localhost:3000`

## Authentication

To access protected endpoints, you need to include a JWT token in the Authorization header:

```
Authorization: Bearer <your_token>
```

### Registration

You can register a new user by making a POST request to `/api/movies/v1/auth/register` with the following format:

```json
{
  "name": "John Doe",
  "email": "user@example.com",
  "password": "password123"
}
```

### Login

You can obtain a token by making a POST request to `/api/movies/v1/auth/login` with the following credentials:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

## License

This project is licensed under the MIT License. 