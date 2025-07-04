# gator cli

This is a RSS blog aggregator written in Go, featuring the ability to follow and pull posts from RSS feeds. It's built with SQLC, Postgres, Goose, and Go.

## Installation

- Requires: Go, Postgres
- Run: `go install github.com/lordbaldwin1/gator`
- Create a config file, `.gatorconfig.json` with the following format:
- `{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable","current_user_name":"username"}`
- Create a Postgres database and in the `gator/sql/schema` directory run: `goose postgres [db_url] up`
- From root directory run: `go build`
- Command: `gator help` to see a list of commands

## Goals

This project was created to reinforce the fundamentals of HTTP clients, JSON parsing, and database integration in Go. In the process, I learned how to implement middleware in Go and, in general, build and architect a CLI tool in Go.