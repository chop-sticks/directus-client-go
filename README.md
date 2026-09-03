# directus-client-go

A typed Go client for the [Directus](https://directus.io) REST API.

It mirrors the official [JS/TS SDK](https://github.com/directus/directus/tree/main/sdk/src/rest)
command surface: items, files, folders, users, roles, policies, permissions,
flows, operations, panels, dashboards, presets, translations, shares, comments,
notifications, activity, revisions, content versions, relations, settings,
extensions, collections, fields, plus the server, schema, utils, and auth
endpoints.

## Install

```sh
go get github.com/chop-sticks/directus-client-go/directus
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/chop-sticks/directus-client-go/directus"
)

func main() {
	host := "http://localhost:8055"
	token := "your-static-token"

	client, err := directus.NewClient(&host, &token)
	if err != nil {
		panic(err)
	}

	// Read items from a user collection with a query.
	limit := 10
	articles, err := client.GetItems("articles", &directus.Query{
		Fields: []string{"id", "title", "author.name"},
		Filter: map[string]any{"status": map[string]any{"_eq": "published"}},
		Sort:   []string{"-date_created"},
		Limit:  &limit,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(articles)

	// Create an item.
	created, err := client.CreateItem("articles", map[string]any{
		"title":  "Hello",
		"status": "draft",
	}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(created["id"])

	// Core collections have typed models.
	me, err := client.GetUsersMe(&directus.Query{Fields: []string{"id", "email"}})
	if err != nil {
		panic(err)
	}
	fmt.Println(me.Email)
}
```

## Authentication

`NewClient` uses a static bearer token. To log in with credentials and use the
returned access token:

```go
auth, err := client.Login("admin@example.com", "password", "json", "")
if err != nil {
	panic(err)
}
client.Token = auth.AccessToken
```

## Query parameters

Most read and write methods accept `*Query` (pass `nil` for none). Fields are
serialized per the Directus [global query](https://docs.directus.io/reference/query.html):
`Fields`, `Filter`, `Sort`, `Limit`, `Offset`, `Page`, `Search`, `Deep`,
`Aggregate`, `GroupBy`, `Alias`, `Version`, `Export`, `Meta`.

## Method conventions

| Operation           | Method shape                                             |
|---------------------|---------------------------------------------------------|
| Read list           | `GetXs(q *Query) ([]X, error)`                          |
| Read one            | `GetX(id, q *Query) (*X, error)`                        |
| Create one / many   | `CreateX` / `CreateXs`                                  |
| Update one          | `PatchX(id, item, q)`                                   |
| Update many by keys | `PatchXs(keys, item, q)` → `{keys, data}`               |
| Update batch        | `PatchXsBatch(items, q)`                                |
| Delete one / many   | `DeleteX(id)` / `DeleteXs(keys)`                        |

## Contributing / extending

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for setup, workflow, coding
conventions, testing, and the PR process. For deeper reference:
[`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md) (architecture, request helpers,
naming, endpoint reference, model mapping) and [`docs/API.md`](docs/API.md)
(full inventory of implemented models and method signatures).

## Development

Common workflows are automated with [Task](https://taskfile.dev). Run `task`
to list them:

| Task | Description |
|------|-------------|
| `task build` | Compile all packages |
| `task test` | Run the test suite |
| `task test:cover` | Tests + coverage profile & HTML report |
| `task test:race` | Tests with the race detector |
| `task test:verbose` | Unit tests with per-test PASS/FAIL output |
| `task test:report` | Unit tests → JUnit XML + JSON in `test-reports/` |
| `task lint` | Run `golangci-lint` |
| `task fmt` | Format with `gofmt` |
| `task vet` | Run `go vet` |
| `task tidy` | Tidy & verify `go.mod`/`go.sum` |
| `task docs` | Serve browsable package docs (pkgsite) |
| `task ci` | Full gate: fmt check, vet, lint, test |
| `task test:integration` | Boot Directus (docker compose) and run integration tests |
| `task test:integration:report` | Integration tests → JUnit XML + JSON in `test-reports/` |
| `task compose:up` / `task compose:down` | Start / stop the Directus stack |

Without Task installed, the underlying commands still work directly, e.g.
`go test ./...`.

## Integration tests

Functional tests run against a live Directus instance defined in
[`docker-compose.yml`](docker-compose.yml) (Directus + PostgreSQL/PostGIS +
Redis). They live in the `directus/integration` package (`//go:build
integration`), so they are excluded from the normal `task test` / `go test
./...` run.

```sh
task test:integration   # up --wait, run tagged tests, then compose down
```

Or manually:

```sh
docker compose up -d --wait
go test -tags=integration -count=1 ./directus/integration/...
docker compose down
```

The target instance and credentials default to the values in
`docker-compose.yml` and can be overridden with `DIRECTUS_URL` and
`DIRECTUS_TOKEN`. See [`docs/INTEGRATION_TESTING.md`](docs/INTEGRATION_TESTING.md)
for the full guide — stack details, harness/fixtures, assertion policy, how to
add tests, and troubleshooting.
