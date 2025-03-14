# Bank API

A RESTful API service for managing bank information and SWIFT codes.

## Prerequisites

- Go 1.24
- Docker and Docker Compose
- Make (optional, but recommended)

## Quick Start

### Using Make (Recommended)

1. Clone the repository:

```bash
git clone https://github.com/MarcinZ20/bankAPI.git
```

2. Go to the project directory

```bash
cd bankAPI
```

3. Setup .env file (default copy from .env.example)

```bash
make setup
```

4. Start the development environment:

```bash
make docker-start
```

That's it! The API will be available at `http://localhost:8080`.

### Manual Setup

If you don't have Make installed, you can run the commands manually:

1. Clone the repository:

```bash
git clone https://github.com/MarcinZ20/bankAPI.git
```

2. Go to the project directory

```bash
cd bankAPI
```

3. Setup .env

```bash
cp .env.example .env
```

4. Start the Docker containers:

```bash
docker-compose up -d
```

## Development

### Available Make Commands

- `make setup` - Run setup
- `make docker-run` - Start all containers
- `make docker-stop` - Stop all containers
- `make test` - Run tests

### API Endpoints

- `GET /v1/swift-codes/:swiftCode` - Get bank details by SWIFT code
- `GET /v1/swift-codes/country/:ISO2Code` - Get bank data by ISO2 country code
- `POST /v1/swift-codes` - Add a new bank entry
- `DELETE /v1/swift-codes/:swiftCode` - Delete a bank entry

### Example Request

```bash
# Get bank details
curl http://localhost:8080/v1/swift-codes/BKSACLRMXXX

# Add a new bank
curl -X POST http://localhost:8080/api/v1/swift-codes \
  -H "Content-Type: application/json" \
  -d '{
    "swiftCode": "DEUTDEFFXXX",
    "bankName": "Deutsche Bank",
    "countryISO2": "DE",
    "countryName": "Germany",
    "address": "Taunusanlage 12",
    "isHeadquarter": true
  }'

# Get banks by ISO2 country code
curl http://localhost:8080/v1/swift-codes/country/CL

# Delete bank by SWIFT code
curl -X DELETE http://localhost:8080/v1/swift-codes/DEUTDEFFXXX
```

## Testing

Run the test suite:

```bash
make test
```

## Project Structure

```
bankAPI/
├── api/
|   |── handlers/       # API handlers for different routes
|   |── middleware/     # API context handlers
|   |── responses/      # API response templates
|   └── routes/         # API endpoint routes
├── cmd/
|   └── main/           # Application entry point
├── internal/
│   ├── app/            # Application specific operations
│   ├── database/       # Database operations
│   ├── parser/         # Data parsing
│   ├── repository/     # Database operations
│   ├── services/       # Business logic
│   ├── spreadsheet/    # Spreadsheet logic
│   ├── transform/      # Data transformation
│   └── validation/     # Validation logic
├── pkg/
│   ├──  models/        # Data models
|   └──  utils/         # API response templates
└── docker/             # Docker-related files
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
