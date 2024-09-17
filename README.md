# Reminder App

This is a simple reminder application built with Go. It allows users to create, read, update, and delete events, and sends email reminders for upcoming events.

## Features

- RESTful API for managing events
- SQLite database for persistent storage
- Email notifications for upcoming events
- API key authentication
- Input validation using go-playground/validator

## Project Structure

```
reminder-app/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   └── events.go
│   │   ├── middleware/
│   │   │   └── auth.go
│   │   └── routes.go
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   └── event.go
│   ├── storage/
│   │   ├── database.go
│   │   └── sqlite.go
│   └── service/
│       ├── email.go
│       └── events.go
├── pkg/
│   └── utils/
│       └── json.go
├── scripts/
│   └── migration.sql
├── .env
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- Go 1.16 or higher
- SQLite

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/reminder-app.git
   cd reminder-app
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Create a `.env` file in the root directory with the following content:
   ```
   DB_PATH=reminder.db
   EMAIL_FROM=your_email@gmail.com
   EMAIL_TO=your_email@gmail.com
   EMAIL_PASSWORD=your_app_password
   API_KEY=your_secret_api_key
   SERVER_PORT=8080
   ```
   Replace the email and API key values with your own.

4. Run the application:
   ```
   go run cmd/api/main.go
   ```

## API Endpoints

All endpoints require the `X-API-Key` header for authentication.

- `GET /events`: Retrieve all events
- `POST /events`: Create a new event
- `PUT /events/{id}`: Update an existing event
- `DELETE /events/{id}`: Delete an event

### Example: Creating an Event

```bash
curl -X POST http://localhost:8080/events \
     -H "Content-Type: application/json" \
     -H "X-API-Key: your_secret_api_key" \
     -d '{"name":"Birthday Party","date":"2024-03-15T18:00:00Z","recurring":false}'
```

## Validation

The application uses `go-playground/validator` for input validation. Event creation and updates are validated with the following rules:

- `name`: Required, must not be empty
- `date`: Required, must be a future date
- `recurring`: Optional boolean field

## Email Notifications

The application sends email reminders for upcoming events. It uses Gmail's SMTP server, so make sure to use an "App Password" for the `EMAIL_PASSWORD` in your `.env` file.

## Development

To run the application in development mode:

```bash
go run cmd/api/main.go
```

## Testing

Soon...

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
