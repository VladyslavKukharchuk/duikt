# Laboratory work 4 – Token-Based Authentication

Project Structure
laboratory-work-4/
├── index.html      # User login form
└── main.go         # Go server with JWT authentication

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
5. Enter the username and password:
   ```text
    Username: user
    Password: 12345
   ```
6. After a successful login, you will receive a JWT, which can be used to access the protected route.

## Notes

- The token is valid for 2 minutes only.
- Failed login attempts and invalid token usage are logged in the server console.