# Go-Blog

Go-Blog is a RESTful blog application built with Go, Gin, and GORM, backed by a PostgreSQL database. The application supports user authentication with JWT, CRUD operations for blog posts, and search functionality with TF and DF statistics tracking.

## Features

- **User Authentication:** Secure user registration and login using JWT.
- **Blog Posts:** Create, read, update, and delete posts.
- **Search Functionality:** Search for words in posts with Term Frequency (TF) and Document Frequency (DF) statistics.
- **Authorization:** Middleware-based protection for sensitive routes.

## Prerequisites

- **Docker:** Ensure you have Docker and Docker Compose installed.
- **Go:** Required for local development.

## Getting Started

### Running with Docker

1. **Clone the repository:**

   ```bash
   git clone https://github.com/oel21sakka/go-blog.git
   cd go-blog
   ```

2. **Build and run the containers:**

   ```bash
   docker-compose up --build
   ```

3. **Access the application:**

   - The application will be available at `http://localhost:8080`.
   - The PostgreSQL database will be available at `localhost:5432`.

### Environment Variables

The application relies on the following environment variables, which are configured in the `docker-compose.yml`:

- `DB_HOST`: Database host (set to `db` in Docker).
- `DB_PORT`: Database port (`5432` by default).
- `DB_USER`: Database username (`postgres`).
- `DB_PASSWORD`: Database password (`postgres`).
- `DB_NAME`: Database name (`myblogdb`).

### API Endpoints

- **Authentication:**
  - `POST /register` - Register a new user.
  - `POST /login` - Log in and receive a JWT.
  
- **Posts:**
  - `GET /api/posts` - Get all posts.
  - `GET /api/posts/:id` - Get a specific post by ID.
  - `POST /api/posts` - Create a new post (protected).
  - `PUT /api/posts/:id` - Update an existing post (protected).
  - `DELETE /api/posts/:id` - Delete a post (protected).

- **Search:**
  - `POST /search_words` - Search for words in posts and get TF/DF statistics.

### Local Development

1. **Install dependencies:**

   ```bash
   go mod download
   ```

2. **Run the application:**

   ```bash
   go run main.go
   ```

3. **Run with Docker (without Compose):**

   ```bash
   docker build -t myblog .
   docker run --name myblog-app -p 8080:8080 myblog
   ```

## Project Structure

- **config/**: Database connection and configuration.
- **controllers/**: Route handlers for authentication, posts, and search.
- **middlewares/**: Authentication middleware.
- **models/**: Database models for users, posts, and search statistics.
- **routes/**: API routes setup.
- **utils/**: Utility functions, including JWT generation and validation.
- **tests/**: Config for initialization needed for testing and controllers tests for login, register and post create and edit functinoonality testing 

## Test Coverage
The current tests cover:

-  User Registration (POST /register):

   -  Verifies that a user can register successfully.
   -  Ensures that required fields are provided and properly validated.

-  User Login (POST /login):

   -  Tests successful login with valid credentials.
   -  Ensures login failure with incorrect credentials.
   -  Verifies the generation of a valid JWT token upon successful login.

-  Post Creation and Updates (POST /api/posts and PUT /api/posts/:id):

   -  Verifies that a user can create a new post.
   -  Tests that the user can update their post.
   -  Ensures that only authenticated users can create or update posts.

## Docker Setup

The project uses a multi-stage Dockerfile:

- **Build Stage:** Compiles the Go application.
- **Run Stage:** Runs the compiled binary in a minimal Alpine container.

The `docker-compose.yml` orchestrates the application and PostgreSQL container.

---
