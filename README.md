# SQL to Go

Generates golang structures based on an sql script file containing create tables.

## Install

```bash
go install github.com/thecampagnards/sql-to-go
sql-to-go -model-type bun -output-folder out examples/*.sql
```

### Using docker

```bash
alias sql-to-go='docker run --rm -ti --user $(id -u):$(id -g) -v $(pwd):/app -w /app ghcr.io/thecampagnards/sql-to-go:master'
sql-to-go -model-type bun -output-folder out examples/*.sql
```

### Output

```go
// Code generated by sql-to-go; DO NOT EDIT.

package models

import (
	"github.com/uptrace/bun"
)

type A struct {
	bun.BaseModel `bun:"table:a"`
	ID            int64   `bun:"id,pk,notnull"`
	AToBs         []*AToB `bun:"rel:has-many"`
}
// Code generated by sql-to-go; DO NOT EDIT.

package models

import (
	"github.com/uptrace/bun"
)

type AToB struct {
	bun.BaseModel `bun:"table:a_to_b"`
	AID           *int64  `bun:"a_id,pk"`
	BID           *int64  `bun:"b_id,pk"`
	Example       *string `bun:"example,"`
	A             *A      `bun:"rel:belongs-to"`
	B             *B      `bun:"rel:belongs-to"`
}
// Code generated by sql-to-go; DO NOT EDIT.

package models

import (
	"github.com/uptrace/bun"
)

type B struct {
	bun.BaseModel `bun:"table:b"`
	ID            int64   `bun:"id,pk,notnull"`
	AToBs         []*AToB `bun:"rel:has-many"`
}
// Code generated by sql-to-go; DO NOT EDIT.

package models

import (
	"github.com/uptrace/bun"
)

type Book struct {
	bun.BaseModel `bun:"table:book"`
	ID            int64   `bun:"id,pk,notnull"`
	Name          *string `bun:"name,"`
	UserID        *int64  `bun:"user_id,"`
	User          *User   `bun:"rel:belongs-to"`
}
// Code generated by sql-to-go; DO NOT EDIT.

package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user"`
	City          *string    `bun:"city,"`
	Country       *string    `bun:"country,"`
	DateOfBirth   *time.Time `bun:"date_of_birth,"`
	Email         *string    `bun:"email,"`
	ID            int64      `bun:"id,pk,notnull"`
	Name          *string    `bun:"name,"`
	PostalCode    *string    `bun:"postal_code,"`
	Surname       *string    `bun:"surname,"`
	Books         []*Book    `bun:"rel:has-many"`
}
```

## Dev

```bash
go run . -model-type bun -output-folder out examples/*.sql
```