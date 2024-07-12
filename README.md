# URL Shortener

This project is a simple URL shortener built with Go. It allows users to shorten URLs and redirect to the original URL using the generated short link. The application uses local file storage to persist the shortened URLs across restarts.

## Features

- Shorten URLs with a simple HTTP request.
- Redirect to the original URL using the generated short link.
- Persist shortened URLs in a file for future use.

## Getting Started

### Prerequisites

- Go (version 1.14 or higher recommended)

### Installation

1. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/MikeMatovu/url-shortener.git
    ```

2. Navigate to the project directory:

    ```bash
    cd url-shortener
    ```

3. Build the project (optional):

    ```bash
    go build
    ```

4. To start the URL shortener, run:

    ```bash
    go run main.go
    ```

The server will start on port 8080. You can access the application via [http://localhost:8080](http://localhost:8080).

## Usage

- **Shorten a URL:** Send a POST request to [http://localhost:8080/shorten](http://localhost:8080/shorten) with the original URL in the request body.
- **Access a Shortened URL:** Navigate to `http://localhost:8080/short/<shortened_part>` in your browser, where `<shortened_part>` is the unique identifier for the shortened URL.

## Stopping the Application

To stop the application, use the interrupt signal (Ctrl+C in most terminals). The application will automatically save the current state of shortened URLs to a file named `urls.gob`.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

FREE
