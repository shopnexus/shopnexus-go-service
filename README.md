# ShopNexus Go Service

A gRPC-based microservice for handling e-commerce operations including user accounts, product management, cart functionality, and payments.

![Flow](https://raw.githubusercontent.com/shopnexus/shopnexus-go-service/refs/heads/main/flow.png)

## Features

- **Account Management**
  - User registration and authentication
  - Admin account management
  - JWT-based authentication
  - Role-based access control

- **Product Management**
  - Product model and inventory tracking
  - Brand management
  - Product categorization with tags
  - Image management for products and brands

- **Shopping Cart**
  - Add/remove items
  - Cart management
  - Checkout process

- **Payment Processing**
  - Multiple payment methods (Cash, Momo, VNPay)
  - Order tracking
  - Payment status management

## Technology Stack

- Go
- gRPC
- PostgreSQL
- Prisma (for database schema management)
- Protocol Buffers
- SQLC (for type-safe SQL)

## Project Structure

```
.
├── cmd/
│ ├── main.go # Service entry point
│ └── client/ # gRPC client examples
├── config/ # Configuration management
├── gen/ # Generated code (protobuf, sqlc)
├── internal/
│ ├── model/ # Domain models
│ ├── repository/ # Database operations
│ ├── server/ # gRPC server implementation
│ ├── service/ # Business logic
│ └── util/ # Utility functions
├── prisma/ # Database schema and migrations
└── proto/ # Protocol buffer definitions

```

## Getting Started

### Prerequisites

- Go 1.21 or later
- PostgreSQL
- Node.js (for Prisma)
- Protocol Buffer Compiler

### Installation

1. Clone the repository:

```bash
git clone https://github.com/shopnexus/shopnexus-go-service.git
cd shopnexus-go-service
```

2. Install dependencies:

```bash
go mod download
npm install
```

3. Initialize the database:

```bash
make init-migrate
```

4. Generate required code:

```bash
make proto    # Generate protobuf code
make sqlc     # Generate SQL code
```

### Configuration

Create a `.env` file in the root directory with the following variables:

```env
DATABASE_URL="postgresql://user:password@localhost:5432/shopnexus?schema=public"
APP_STAGE="development"
JWT_SECRET="your-secret-key"
```

### Running the Service

Development mode:

```bash
make dev
```

Production mode:

```bash
APP_STAGE=production make run
```

## API Documentation

The service exposes gRPC endpoints for:

- Account Service: User/Admin authentication and management
- Product Service: Product and inventory management
- Cart Service: Shopping cart operations
- Payment Service: Payment processing

For detailed API documentation, refer to the proto files in the `proto/` directory.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
