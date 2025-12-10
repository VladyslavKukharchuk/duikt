# Laboratory work 4 – Token-Based Authentication

Project Structure
```text
laboratory-work-4/
├── index.html      # User login form
└── main.go         # Go server with JWT authentication
```

### Requirements

- Go 1.18+
- Web browser to open index.html

## How to Run

1. Clone or copy the project to a local folder.
2. Initialize Go module and install library.
    ```shell
    go mod tidy
    ```
3. Start the server:
    ```shell
    go run main.go
    ```
4. Open your browser and navigate to:
   http://localhost:8080/index.html

## How to Test

1. Enter the username and password:
   ```text
    Username: user
    Password: 12345
   ```
2. After a successful login, you will receive a JWT.
3. Now you can use `Access Protected Resource` button to test your token.
