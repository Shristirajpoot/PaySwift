# BlackPay - A Wrapper for Paystack Payment

## Overview

BlackPay is a Go-based application designed to manage bills, payments, and notifications using the Paystack payment gateway. It integrates with various services such as email and SMS notifications, and it uses a PostgreSQL database for data storage.

## Features

- **User Management**: Create, read, update, and delete user records.
- **Bill Management**: Create, read, update, delete, and send bills to users.
- **Payment Processing**: Initiate payments using Paystack, verify transactions, and update payment status.
- **Notification System**: Send email and SMS notifications for new bills and payment updates.
- **Webhook Handling**: Handle webhooks from Paystack to update payment status in real-time.
- **API Endpoints**: Expose RESTful API endpoints for user, bill, and payment operations.

## Directory Structure

The project follows a structured directory layout to ensure maintainability and scalability:

```plaintext
blackprince001-blackpay/
├── main.go
├── services/
│   ├── sms.go
│   ├── email.go
│   └── paystack.go
├── go.sum
├── handlers/
│   ├── user.go
│   ├── webhook.go
│   ├── bill.go
│   └── payment.go
├── models/
│   ├── user.go
│   ├── bill.go
│   └── payment.go
├── migrations/
│   ├── create_payments_table.sql
│   ├── create_bills_table.sql
│   └── create_users_table.sql
├── config/
│   └── config.go
├── Dockerfile
├── go.mod
├── README.md
├── docker-compose.yml
├── middleware/
│   └── auth.go
└── utils/
    ├── db.sql
    ├── logger.go
    ├── notis.go
    └── database.go
```

## Installation and Setup

### Prerequisites

- Go 1.16 or later
- Docker and Docker Compose
- PostgreSQL database

### Steps to Set Up

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-repo/blackprince001-blackpay.git
   ```

2. **Navigate to the Project Directory**:
   ```bash
   cd blackprince001-blackpay
   ```

3. **Build and Run the Application Using Docker Compose**:
   ```bash
   docker-compose up --build
   ```

4. **Environment Variables**:
   Ensure the following environment variables are set in your `docker-compose.yml` file or in your environment:
   - `DATABASE_URL`
   - `PAYSTACK_SECRET_KEY`
   - `COURIER_API_KEY`
   - `ARKESEL_API_KEY`
   - `SERVER_PORT`
   - `DB_HOST`
   - `DB_USER`
   - `DB_PASSWORD`
   - `DB_NAME`
   - `DB_PORT`

## API Endpoints

### User Endpoints

- **Create User**:
  - Method: `POST`
  - URL: `/users`
  - Request Body: `{"name": "John Doe", "email": "john@example.com", "phone": "1234567890"}`
  - Response: `201 Created` with the created user object

- **Get User**:
  - Method: `GET`
  - URL: `/users/:id`
  - Response: `200 OK` with the user object

- **Update User**:
  - Method: `PUT`
  - URL: `/users/:id`
  - Request Body: `{"name": "Jane Doe", "email": "jane@example.com", "phone": "0987654321"}`
  - Response: `200 OK` with the updated user object

### Bill Endpoints

- **Create Bill**:
  - Method: `POST`
  - URL: `/bills`
  - Request Body: `{"user_id": 1, "amount": 100.00, "due_date": "2024-12-31T00:00:00Z", "status": "unpaid", "description": "Sample Bill"}`
  - Response: `201 Created` with the created bill object

- **Get Bill**:
  - Method: `GET`
  - URL: `/bills/:id`
  - Response: `200 OK` with the bill object

- **List Bills**:
  - Method: `GET`
  - URL: `/bills`
  - Response: `200 OK` with a list of bill objects

- **Update Bill**:
  - Method: `PUT`
  - URL: `/bills/:id`
  - Request Body: `{"amount": 150.00, "due_date": "2025-01-01T00:00:00Z", "status": "unpaid", "description": "Updated Bill"}`
  - Response: `200 OK` with the updated bill object

- **Delete Bill**:
  - Method: `DELETE`
  - URL: `/bills/:id`
  - Response: `200 OK` with a success message

- **Send Bill**:
  - Method: `POST`
  - URL: `/bills/:id/send`
  - Response: `200 OK` with a payment link and a success message

### Payment Endpoints

- **Initiate Payment**:
  - Method: `POST`
  - URL: `/payments`
  - Request Body: `{"bill_id": 1}`
  - Response: `200 OK` with the authorization URL and reference

- **Get Payment Info**:
  - Method: `GET`
  - URL: `/payments/:id`
  - Response: `200 OK` with the payment object

- **List Payments**:
  - Method: `GET`
  - URL: `/payments`
  - Response: `200 OK` with a list of payment objects

- **Verify Transaction**:
  - Method: `GET`
  - URL: `/payments/verify?reference=REFERENCE_ID`
  - Response: `200 OK` with the updated payment status

### Webhook Endpoint

- **Handle Webhook**:
  - Method: `POST`
  - URL: `/webhook`
  - Request Body: Paystack webhook payload
  - Response: `200 OK` with a success message

## Configuration

The application configuration is loaded from environment variables. Here are the key configurations:

- `DATABASE_URL`: The URL for the PostgreSQL database.
- `PAYSTACK_SECRET_KEY`: The secret key for the Paystack API.
- `COURIER_API_KEY`: The API key for the courier service.
- `ARKESEL_API_KEY`: The API key for the Arkesel SMS service.
- `SERVER_PORT`: The port on which the server will run.

## Development

### Building the Application

To build the application, you can use the provided `Dockerfile` and `docker-compose.yml`.

### Running Tests

Currently, this project does not include automated tests. However, you can add tests using Go's built-in `testing` package and frameworks like `testify`.

### Contributing

Contributions are welcome. Here are some steps to contribute:

1. **Fork the Repository**: Fork the repository to your own GitHub account.
2. **Create a Branch**: Create a new branch for your feature or bug fix.
3. **Make Changes**: Make the necessary changes and ensure they are well-documented.
4. **Commit Changes**: Commit your changes with a meaningful commit message.
5. **Open a Pull Request**: Open a pull request to the main branch.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Paystack API Documentation](https://paystack.com/docs/api/)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)

If you have any questions or need further assistance, feel free to open an issue or contact the maintainers.
