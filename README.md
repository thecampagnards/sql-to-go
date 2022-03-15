# SQL to Go

Generates structures based on an sql script containing create tables.

## Install

```bash
go get github.com/thecampagnards/sql-to-go
sql-to-go -sql-file example.sql -model-type bun -output-folder out
```

## Dev

```bash
go run . -sql-file example.sql -model-type bun -output-folder out
```