# Shopping Cart Backend

This is the backend for the Shopping Cart application built with Golang, Gin, GORM, and PostgreSQL.

## Prerequisites
- Go 1.18+
- PostgreSQL (running on default port 5432)

## Environment Variables
Copy `.env.example` to `.env` and update values as needed.

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shopping_cart
DB_PORT=5432
```

## Setup
1. Install dependencies:
   ```
   go mod tidy
   ```
2. Start PostgreSQL and create the database `shopping_cart`.
3. Run the server:
   ```
   go run main.go migrate.go models.go user_handlers.go item_handlers.go cart_handlers.go order_handlers.go auth_middleware.go
   ```
   The server will start on `http://localhost:8080`.

## API Endpoints
- `POST   /users`         - Register new user
- `GET    /users`         - List all users
- `POST   /users/login`   - User login (returns token)
- `POST   /items`         - Create new item
- `GET    /items`         - List all items
- `POST   /carts`         - Add item to cart (auth required)
- `GET    /carts`         - List cart items (auth required)
- `POST   /orders`        - Convert cart to order (auth required)
- `GET    /orders`        - List all orders (auth required)

## Testing
Ginkgo tests will be added in the `tests/` directory.

## Notes
- Use the `Authorization: Bearer <token>` header for all cart and order related endpoints.
- Each user can only be logged in from one device at a time (single token per user).
