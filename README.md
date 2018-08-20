# Hacker News Daily

An app for scraping Hacker News and then building a front page for a given, scraped date.

## Getting Started

### Prereq

- `go1.10.2`
- `dep v0.4.1`

### Building

1. `dep ensure`

### Running

Run the scraping process with `go run /path/to/hn-daily/main.go &`.

When you're ready to create your daily digest, stop your scraping process, and run `go run digest/main.go`. This will create you JSON formated digest in the root directory. 

## Limitation

This project uses boltDB. Bolt keeps a lock on the file for the DB. The current implementation of the scraper keeps a connection open to the DB while the process is running. This means you cannot create a digest while also scraping.
