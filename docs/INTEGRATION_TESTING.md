# Integration / End-to-End Testing

This SDK ships an end-to-end test suite that runs against a **live Directus
instance**. It exercises every resource and every client method over real HTTP —
create/read/update/delete round-trips, auth, file upload, schema, utils, and the
action endpoints — verifying the SDK against actual Directus behavior rather than
mocks.

- Unit tests (mock HTTP, no server) live in `package directus` (`*_test.go`).
- Integration/e2e tests live in the **`directus/integration`** package and only
  build under the `integration` tag. See [API.md](API.md) and
  [DEVELOPMENT.md](DEVELOPMENT.md) for the SDK surface and conventions.

## TL;DR

```sh
task test:integration
```

That boots the Docker stack, waits until healthy, runs the tagged tests, and
tears the stack down again (even on failure).

## Prerequisites

- **Docker** + **Docker Compose v2** (`docker compose`). Used to run Directus,
  PostgreSQL/PostGIS, and Redis locally.
- **Go** (module toolchain per `go.mod`).
- **[Task](https://taskfile.dev)** (optional but recommended) for the wrapper
  tasks below.

## The stack

`docker-compose.yml` (repo root) defines three services:

| Service    | Image                       | Purpose                          |
|------------|-----------------------------|----------------------------------|
| `directus` | `directus/directus:latest`  | API under test, published on `:8055` |
| `database` | `postgis/postgis:13-master` | PostgreSQL data store            |
| `cache`    | `redis:6`                   | Cache / rate-limit store         |

All services run on the default compose bridge network; `directus` is reachable
from the host at `http://localhost:8055`. The Directus container's healthcheck
uses Node's `fetch` (not `wget`, which isn't in the image), so `--wait` reliably
blocks until the API is actually serving.

Admin credentials are bootstrapped from the compose environment on first boot:

| Setting  | Value                              |
|----------|------------------------------------|
| Email    | `test@example.com`                 |
| Password | `testAtExampleDotCom+1`            |
| Token    | `eiriezashohnai1xohjuC2aem7duuDie` (static admin token) |

## Running the tests

### With Task (recommended)

```sh
task test:integration   # up --wait -> go test -tags=integration -> compose down
```

Stack lifecycle helpers, if you want to keep it running while iterating:

```sh
task compose:up     # docker compose up -d --wait
task compose:logs   # docker compose logs -f
task compose:down   # docker compose down
```

### Manually

```sh
docker compose up -d --wait
go test -tags=integration -count=1 ./directus/integration/...
docker compose down
```

`-count=1` disables the test cache so the suite always hits the live server.

### Configuration (environment)

The client target is resolved from the environment, falling back to the
compose defaults, so you can point the suite at any instance:

| Variable         | Default                              | Meaning                    |
|------------------|--------------------------------------|----------------------------|
| `DIRECTUS_URL`   | `http://localhost:8055`              | Base URL of the instance   |
| `DIRECTUS_TOKEN` | `eiriezashohnai1xohjuC2aem7duuDie`   | Static admin bearer token  |

The admin email/password used by the auth tests are compile-time constants in
`setup_test.go` (they match `docker-compose.yml`); change both together if you
retarget to a different instance.

```sh
DIRECTUS_URL=https://my.directus.app DIRECTUS_TOKEN=xxxxx \
  go test -tags=integration -count=1 ./directus/integration/...
```

## Reports & CI output

Plain `go test` output is terse. Three options surface per-test results:

- **`task test:verbose`** (`go test -v`) — built-in, zero dependencies; prints a
  `--- PASS`/`--- FAIL` line per test.
- **`task test:report`** — runs the unit suite through
  [gotestsum](https://github.com/gotestyourself/gotestsum) (fetched on demand via
  `go run`, never added to `go.mod`). Prints a per-test `PASS`/`FAIL` line plus a
  `DONE N tests` summary, and writes machine-readable reports:
  - `test-reports/unit.xml` — JUnit XML (consumable by CI test reporters —
    GitHub/GitLab, Jenkins, etc.).
  - `test-reports/unit.json` — the raw `go test -json` event stream.
- **`task test:integration:report`** — the same for the live e2e suite (boots the
  stack, runs with `-tags=integration`, tears it down), writing
  `test-reports/integration.xml` and `test-reports/integration.json`.

On failure gotestsum prints a consolidated list of the failed tests and their
output after the run, so you can see exactly what broke. `test-reports/` is
gitignored and removed by `task clean`.

Without Task:

```sh
go test -v ./directus/...                     # verbose per-test PASS/FAIL
go test -json ./directus/... > report.json    # machine-readable event stream
go run gotest.tools/gotestsum@v1.13.0 --format testname \
  --junitfile junit.xml -- -count=1 ./directus/...
```

## How it is wired

- **Own package.** Tests are the external package `integration`
  (`directus/integration/`), importing the SDK. This keeps the live suite fully
  separate from the unit tests and from shipped code.
- **Build tag.** Every `*_test.go` starts with `//go:build integration`, so a
  normal `go test ./...` reports `directus/integration [no test files]` and
  never tries to reach a server. `doc.go` (untagged) keeps the package present
  for `go build ./...` / `go vet ./...`.
- **Dot import.** The SDK is dot-imported
  (`. "github.com/chop-sticks/directus-client-go/directus"`) so test bodies read
  naturally (`&User{...}`, `NewClient(...)`) instead of `directus.`-qualifying
  every type. This is a recognized pattern for external test packages and is not
  flagged by lint (which runs without the `integration` tag).

### File layout

```
directus/integration/
  doc.go            # package doc (no build tag)
  setup_test.go     # itestClient, env(), credential consts, server/auth smoke tests
  helpers_test.go   # uniqueName + fixtures (see below) + TestE2EFixtures
  access_test.go    # users, roles, policies, permissions
  files_test.go     # files, folders, assets, utils
  automation_test.go# flows, operations, panels, dashboards
  content_test.go   # presets, translations, shares, comments
  tracking_test.go  # notifications, activity, revisions, versions
  schema_test.go    # relations, settings, extensions, schema
  core_test.go      # collections, fields, items, singleton, aggregate
  server_test.go    # server, auth
```

### Shared harness

Defined in `setup_test.go` / `helpers_test.go` and reused across files:

- `itestClient(t) *Client` — client built from `DIRECTUS_URL`/`DIRECTUS_TOKEN`.
- `uniqueName(prefix) string` — collision-free names (atomic counter + time).
- Fixtures, each of which creates the object and registers a `t.Cleanup` to
  delete it:
  - `newTestCollection(t, c) string` — collection with int PK `id` + string `title`.
  - `newTestItem(t, c, coll, title) string` — inserts an item, returns its key.
  - `newTestRole(t, c)`, `newTestPolicy(t, c)`, `newTestDashboard(t, c)`,
    `newTestFolder(t, c)`, `newTestUser(t, c)` — return the new id.
  - `uploadTestFile(t, c) string` — uploads a small file, returns its id.

## Test design & assertion policy

Tests are **self-contained and order-independent**: each creates its own
fixtures with `uniqueName` and cleans them up via `t.Cleanup`. Because Directus
is shared mutable state, never assume a clean database or rely on another test's
data.

Every method is exercised at least once, under one of three policies:

1. **Achievable** (normal CRUD, reads, admin-token actions) — assert `err == nil`
   plus a minimal result check (non-nil, expected field, correct length).
2. **Must-fail** (invalid token/id/otp, unlicensed feature) — assert `err != nil`.
   e.g. `PasswordReset(badToken)`, `AcceptUserInvite(bad)`, `EnableTwoFactor(secret,"000000")`,
   `DeleteExtension(bad)`, `InstallRegistryExtension(bad,bad)`.
3. **Environment-dependent** (email transport, external network, extension
   registry) — call the method but do not fail on error
   (`if err != nil { t.Logf(...) }`). e.g. `PasswordRequest`, `InviteUser`,
   `InviteShare`, `ImportFile(url)`, `GetRegistryExtensions`.

## Adding a new e2e test

1. Put it in the matching `directus/integration/<group>_test.go` (or add a new
   file with the `//go:build integration` header + `package integration`).
2. Dot-import the SDK if you name any SDK type:
   `. "github.com/chop-sticks/directus-client-go/directus"`. Omit the import if
   you only call methods on a client and use local helpers (as `server_test.go`
   does).
3. Get a client with `itestClient(t)`, build fixtures with the shared helpers (or
   inline `Create*` calls with `uniqueName`), and register `t.Cleanup` for
   anything you create.
4. Choose the assertion policy above that fits the method.
5. Run `task test:integration` and iterate against the real server.

## Directus behaviors worth knowing (learned from the live suite)

These are quirks the e2e suite depends on; keep them in mind when extending:

- **Content versions** require the collection to have versioning enabled
  (`PatchCollection(coll, &CollectionRequest{Meta:&CollectionMeta{Versioning:true}}, nil)`)
  before `CreateContentVersion`.
- **Relations** need two collections plus a field on the "many" side before
  `CreateRelation`.
- **Permissions**: creating custom (non-full-access) permission rules is a
  licensed feature — `Patch*` permission calls return `403
  custom_permission_rules_enabled` on an unlicensed instance (treated as
  log-only). `POST /permissions` (batch) returns *all* permissions, not just the
  created rows, so derive created ids by querying the policy, not from the
  response length.
- **Notifications** use integer ids (`Notification.ID int`).
- **`UtilsImport`** sets the multipart part's `Content-Type` from the filename
  extension; Directus rejects `application/octet-stream`.
- **Sortable collections**: a non-nullable custom field makes item inserts
  require that field; make helper fields nullable when items are created without
  them.
- **Singletons** are created by enabling `Meta.Singleton`; `GetSingleton` may
  return an empty object before anything is written.

## Troubleshooting

- **`directus` container reports unhealthy but the API works.** The healthcheck
  runs inside the container via Node `fetch`; if you change it, ensure the tool
  exists in the image (`wget`/`curl` are not guaranteed; Node is).
- **Stale state / bootstrap mismatch.** The database is a bind mount
  (`./data/database`, gitignored). For a guaranteed-fresh instance:
  `docker compose down` then remove `data/database` before `compose up`. The
  static admin token is only applied on first bootstrap.
- **Port already in use.** Something else is on `:8055` — stop it or retarget
  with `DIRECTUS_URL`.
- **Left-over containers.** `task test:integration` always runs `docker compose
  down` (via a Task `defer`), but if a run was interrupted, `task compose:down`
  cleans up.
- **Image pulls are slow / arch note.** `postgis/postgis:13-master` runs under
  emulation on arm64 (Apple Silicon); it works but the first `up` pulls large
  images.
