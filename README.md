# Chirpy

A backend web server for **Chirpy** - a Twitter-like social media platform built using **Go** as part of the [Boot.dev](https://www.boot.dev) backend engineering curriculum.

---

## Features

* **HTTP Server**: Built from scratch using Go's standard library (`net/http`).
* **File Server & Middleware**: Serves static HTML files and tracks server metrics via custom HTTP middleware.
* **JSON API**: Handles user management, authentication, and creating/fetching "chirps" (short messages).
* **Database Integration**: Persistence layer for storing users and posts securely.

---

## Prerequisites

Make sure you have the following installed on your machine:
* [Go](https://golang.org/) (v1.22 or higher)
* [Git](https://git-scm.com/)
* [Boot.dev CLI](https://boot.dev) (optional, if running course tests)

---

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone [https://github.com/i-am-ok-15/chirpy.git](https://github.com/i-am-ok-15/chirpy.git)

   cd chirpy


2. **Initialize / Download dependencies:**
    ```bash
    go mod tidy

3. **Build and run the server:**
    ```bash
    go build -o out && ./out

## API Endpoints

| Endpoint          | Method        | Description |
| -----------       | -----------   | -----------
| /app/             | GET           | Servers the frontend web client |
| /api/healthz      | GET           | Health check endpoint (returns OK) |
| /api/metrics      | GET           | Returns the total count of page/file server hits |
| /api/reset        | POST          | Resets the metrics counter back to 0 |
| /api/chirps       | GET / POST    | Retrieve or create chirps
