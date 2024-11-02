# Slip - A Simple Note-Taking Application

## Overview

Slip is a simple note-taking application built using Go and Gin framework. It allows users to create, view, and manage notes through a series of APIs. The application stores notes in Markdown format and generates an index page for easy navigation.

## Development Status

Please note that the current code is still under development, and some features are not yet fully implemented. We welcome community participation and feedback, and hope to collaboratively maintain and improve this project. If you have any suggestions, issues, or would like to contribute code, please feel free to submit a pull request or open an issue.


## Features

- Create notes with a title and body.
- View notes in a web interface.
- Generate an index page for all notes.
- Written in Go with Gin framework for efficient web handling.
- Uses Markdown for note formatting.

## Getting Started

### Prerequisites

- Go 1.19 or later
- Docker (optional, for containerized deployment)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/slip.git
   cd slip
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Prepare required files:

   a. Configuration file:
   In the `configs` directory, copy `config.example.yaml` and rename it to `config.yaml`, then modify it according to your actual settings:

   ```yaml
   keys:
       client_id: "your_client_id" # Client ID
       secret_key: "your_secret_key1" # Secret key (16 bytes)
   ```

   b. Template file:
   Create a `templates` directory in the project root and copy the template file from the example:
   ```bash
   mkdir -p templates
   cp templates.example/index.html.tmpl templates/index.html.tmpl
   ```

4. Build the application:

   ```bash
   go build -o slip main.go
   ```

5. Run the application:

   ```bash
   ./slip
   ```

   The application will start on `http://localhost:8084`.

### Docker Deployment

To run the application using Docker, you can build and run the Docker container:

1. Build the Docker image:

   ```bash
   docker build -t slip .
   ```

2. Run the Docker container:

   ```bash
   docker run -p 8084:8084 slip
   ```

   Access the application at `http://localhost:8084`.

## Usage

- To log in and obtain a token, send a GET request to `/login?encrypted_string=&client_id=` to receive the generated `tokenstr`.
  - The two parameters represent the "encrypted string" and "client ID".
  - Each client provides an "encrypted string" during registration, and the server will provide a `clientId`, which is the "client ID".
  - The client must use the "encryption function" from `examples/auth.js` to generate the "encrypted string".
- To create a note, send a POST request to `/send-notes` with a JSON body containing the title and body of the note, and include `Authorization: token tokenstr` in the request header.

  Example:

  ```curl
  curl -X "POST" "http://127.0.0.1:8084/send-notes" \
       -H 'Authorization: token eyJhbGciOiJIUzI1NiIsInR5cC' \
       -H 'Content-Type: application/json; charset=utf-8' \
       -d $'{
    "title": "What is Slip",
    "body": "Slip is a simple and easy-to-use note-taking application that allows users to quickly create, view, and manage notes. It supports Markdown format and provides an intuitive web interface, enabling users to easily access and organize their notes. Whether for work records or life reflections, Slip helps users efficiently organize their thoughts."
  }'
  ```

- To view the index of notes, navigate to `/index`.
- To view a specific note, access `/notes/:title`, replacing `:title` with the actual title of the note.

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

## Contact

For any inquiries, please contact:

- Name: xzsj.wang
- Email: melody.wang1984@gmail.com
