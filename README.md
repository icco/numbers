# Numbers

[![Go Report Card](https://goreportcard.com/badge/github.com/icco/numbers)](https://goreportcard.com/report/github.com/icco/numbers)

Tiny HTTP service that returns a single character from `book.txt`, indexed by
the current position within the UTC week.

## API

| Method | Path       | Description                                                          |
|--------|------------|----------------------------------------------------------------------|
| `GET`  | `/`        | Plain-text decimal codepoint of the current character.               |
| `GET`  | `/json`    | JSON: `character`, `seconds_passed`, `seconds`, `length`, `lookup`.  |
| `GET`  | `/healthz` | Liveness probe.                                                      |
| `GET`  | `/metrics` | OTel HTTP semconv metrics in Prometheus exposition format.           |

```ruby
require 'open-uri'
loop { puts open('https://numbersstation.blue/') {|f| f.read.to_i.chr }; sleep 1 }
```

## Running

```bash
go run .
```

```bash
docker build -t numbers .
docker run --rm -p 8080:8080 numbers
```
