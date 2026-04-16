# Go REST API Boilerplate

This project is a lightweight, high-performance Go REST API boilerplate featuring a custom tree-based router and a flexible middleware system. It serves as a foundation for building scalable web services with clean architecture and organized routing.

## Features

-   **Custom Tree-Based Router**: Efficient route matching using a tree structure (`RouteSchema`).
-   **Chainable Middleware System**: Support for `Use` (pre-handler), `After` (post-handler), and `Error` hooks.
-   **Clean Architecture**: Organized structure separating library code, application logic, and route definitions.
-   **Standard HTTP Methods**: Full support for `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `HEAD`, `CONNECT`, `TRACE`.
-   **Request/Response Wrappers**: Convenient abstractions for handling HTTP requests and responses.
-   **Docker Ready**: Includes `Dockerfile` and `docker-compose.yml` for containerized environments.
-   **Comprehensive Testing**: Pre-configured test suite with coverage reporting.

## Project Structure

```text
├── app/
│   ├── main.go               # Entry point of the application
│   ├── application/
│   │   └── handler/          # Business logic handlers
│   ├── lib/                  # Core library (router, middleware, wrappers)
│   └── routes/               # Centralized route definitions
├── docker/                   # Docker configuration files
├── test/                     # Comprehensive test suite
├── go.mod                    # Go module definition
└── Makefile                  # Common development tasks
```

## Getting Started

### Prerequisites

-   **Go**: 1.26 or higher
-   **Docker** (optional): For containerized execution

### Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/jmticonap/rest-api-boilerplate
    cd rest-api-boilerplate
    ```

2.  Install dependencies:
    ```sh
    go mod tidy
    ```

### Running Locally

To run the application directly:

```sh
make run
```

To build and run the binary:

```sh
make build
./rest-api
```

The server starts on `http://localhost:3000`.

## Usage

### Defining Routes

Routes are defined in `app/routes/routes.go` using a builder-style interface:

```go
package routes

import (
	"rest-api/app/application/handler"
	"rest-api/app/lib"
)

func InitRoutes() *lib.Routes {
	return lib.NewRoutes().
		Get(lib.Route{
			Path: "/example",
			Handler: lib.NewMiddleware().
				Use(MyMiddleware).
				Build(handler.MyHandler),
		})
}
```

### Creating Handlers

Handlers follow a specific signature that integrates with the middleware system:

```go
func MyHandler(r *http.Request) (*lib.MidResponse, error) {
    // Business logic here
    return &lib.MidResponse{
        Status: http.StatusOK,
        Data: map[string]string{"message": "Hello World"},
    }, nil
}
```

### Using Middleware

Middleware can be chained using `.Use()`, `.After()`, or `.Error()`:

```go
lib.NewMiddleware().
    Use(func(r *http.Request) (*lib.MidResponse, error) {
        // Pre-handler logic
        return nil, nil
    }).
    After(func(r *http.Request) (*lib.MidResponse, error) {
        // Post-handler logic
        return nil, nil
    }).
    Build(myHandler)
```

## Testing

The project uses a `Makefile` to simplify testing and coverage reporting:

| Command | Description |
| :--- | :--- |
| `make test` | Runs all tests with coverage. |
| `make test-v` | Runs tests in verbose mode. |
| `make test-coverage` | Opens the HTML coverage report. |

## Continuous Integration

The repository includes a GitHub Actions workflow (`.github/workflows/go-tests.yml`) that automatically runs unit tests on every pull request to the `main`, `develop`, and `release` branches, ensuring code quality and preventing regressions.

## Docker

To build and run the application using Docker:

```sh
# Build the image
make docker-build

# Start the application
make docker-up
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
