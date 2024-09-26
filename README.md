# Go REST API

This project is a RESTful API built with Go, using the Echo framework and MySQL database. It provides CRUD operations for users and products, with authentication and role-based access control.

## Prerequisites

- Go 1.16 or later
- MySQL 5.7 or later

## Setup

1. Clone the repository:

   ```
   git clone
   cd go_rest_api
   ```

2. Install dependencies:

   ```
   go mod tidy
   ```

3. Set up the MySQL database:

   - Create a new database called `go_rest_api`
   - Update the database connection details in `internal/storage/mysql.go`

4. Run the database migrations:
   - Use a MySQL client to run the SQL commands in `db.sql`

## Running the Application

1. Start the server:

   ```
   go run cmd/api/main.go
   ```

2. The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

- POST /login

### Users

- GET /users (Admin only)
- POST /users (Admin only)
- GET /users/:id
- PUT /users/:id
- DELETE /users/:id (Admin only)

### Products

- GET /products
- POST /products (Admin only)
- GET /products/:id
- PUT /products/:id (Admin only)
- DELETE /products/:id (Admin only)

## Testing

You can use tools like cURL or Postman to test the API endpoints. Remember to include the JWT token in the Authorization header for authenticated requests.

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the MIT License.
