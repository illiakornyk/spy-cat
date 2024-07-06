# Spy Cat Management System

## Overview

This project is a CRUD application for managing spy cats, their missions, and the targets they are assigned to. The application is built to demonstrate the creation of RESTful APIs, interaction with a SQL-like database, integration with third-party services, and optional user interface creation.

## Features

- **Spy Cats Management**:

  - Create, update, list, and delete spy cats.
  - Each cat has a name, years of experience, breed, and salary.
  - Breed validation using TheCatAPI.

- **Missions and Targets Management**:
  - Create missions along with targets.
  - Each mission can have multiple targets.
  - Update and delete missions and targets.
  - Assign cats to missions.
  - Mark targets and missions as complete.

## Technical Details

- **Framework**: Go-Chi for routing.
- **Database**: SQLite for simplicity and ease of use. The project structure allows easy integration with other databases like PostgreSQL or MySQL.
- **Validation**: Go-Playground Validator for request validation.
- **Logging**: Custom logger using Go's built-in `slog` package.

## Project Structure

The project is organized to separate concerns and allow easy extension:

- `cmd/spy-cat`: Entry point of the application.
- `internal/handlers`: Contains HTTP handlers for different entities.
- `internal/storage`: Database interactions.
- `internal/lib`: Common libraries and utilities.
- `config`: Configuration files.
- `migrations`: Database migration files.

## Getting Started

### Prerequisites

- Docker
- Go (if running without Docker)

### Running the Application with Docker

1. **Build the Docker image**:

   ```sh
   docker build -t spy-cat-app .
   ```

2. **Run the Docker container**:

   ```sh
   docker run -d -p 8082:8082 --name spy-cat-container spy-cat-app
   ```

3. **Make requests** to the application:
   - Base URL: `http://localhost:8082/api/v1`
   - Use the provided [Postman collection](./Spy%20Cats.postman_collection.json) to test the API.

### Configuration

The application configuration is managed using YAML files. For Docker, ensure you use the `dev.yml` configuration:

```yaml
env: "dev"
storage_path: "./storage/storage.db"
http_server:
  address: "0.0.0.0:8082"
  timeout: 4s
  idle_timeout: 30s
```

### Endpoints

Refer to the [Postman collection](./Spy%20Cats.postman_collection.json) in the repository for detailed information about available endpoints and their usage.

## Additional Improvements

- Implement tests.
- Add authentication.
- Improve error handling.
- Extend with more detailed logging and monitoring.
