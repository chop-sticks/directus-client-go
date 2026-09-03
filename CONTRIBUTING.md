# Contributing

Thanks for contributing to **directus-client-go**, a typed Go client for the
Directus REST API. This guide covers setup, the development workflow, coding
conventions, testing, and how to get a change merged.

## Prerequisites

- **Go** — matching `go.mod` (CI builds on Go 1.26; a newer toolchain works).
- **[Task](https://taskfile.dev)** — task runner for the workflows below
  (optional; the underlying `go`/`gofmt` commands work directly too).
- **[golangci-lint](https://golangci-lint.run) v2** — linting (`task lint`).
- **Docker** + **Docker Compose v2** — only for the integration/e2e suite.

## Getting started

```sh
git clone https://github.com/chop-sticks/directus-client-go
cd directus-client-go
task build      # compile everything
task ci         # fmt check + vet + lint + unit tests (the local gate)
```

Run `task` (or `task --list`) to see every available task.

## Repository layout

- `directus/` — the SDK. One file per resource holding **both** its model
  struct(s) and its `Client` methods (e.g. `users.go`, `files.go`); tests sit
  beside them (`users_test.go`). Shared plumbing lives in `client.go`,
  `query_params.go`, and `doc.go`.
- `directus/integration/` — the live end-to-end test suite (own package, build
  tag `integration`).
- `docs/` — developer documentation (see below).
- `Taskfile.yml`, `docker-compose.yml`, `.golangci.yml` — tooling.

## Documentation map

Read these before making non-trivial changes; keep them in sync with your work:

- [`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md) — architecture, shared request
  helpers, naming conventions, the full endpoint reference, and how to add or
  extend a resource. **Start here.**
- [`docs/API.md`](docs/API.md) — the implemented surface: every model and method
  signature, plus non-obvious Directus behaviors ("gotchas"). Update it when you
  add or change public methods.
- [`docs/INTEGRATION_TESTING.md`](docs/INTEGRATION_TESTING.md) — running and
  writing the live e2e tests.

## Development workflow

| Step | Task |
|------|------|
| Format | `task fmt` (writes) / `task fmt:check` (verifies) |
| Vet | `task vet` |
| Lint | `task lint` (`task lint:fix` to auto-fix) |
| Unit tests | `task test` (`task test:verbose` for per-test output) |
| Coverage | `task test:cover` → `coverage.html` |
| Full local gate | `task ci` |

**Before opening a PR, `task ci` must pass** (this mirrors what reviewers and CI
expect). Keep the tree `gofmt`-clean and `golangci-lint`-clean (0 issues).

## Coding conventions

Full detail is in [`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md); the essentials:

- **One convention only.** Match the existing patterns; do not introduce a second
  way of doing something that already has a home.
- **Use the shared HTTP helpers** — `request[T]`, `requestRaw`, `execute`. Never
  hand-roll `http.NewRequest`/envelope decoding (the two multipart methods that
  must set a boundary are the only, documented, exceptions).
- **Method naming** mirrors the SDK: `Get`/`Create`/`Patch`/`Delete`; single
  reads/creates/updates return `*T`, lists/batches return `[]T`, deletes return
  `error`. A method takes `q *Query` only where the endpoint accepts a query.
- **Models** are flat structs co-located with the resource; relational fields that
  may be an id or an expanded object are `any`; optional/partial-write fields use
  `,omitempty`; ids are `int` or `string` per the resource.
- **Validate** required path params and return an error before any HTTP call.
- **Doc comments** on every exported symbol, starting with its name (godoc-ready).
- Run the project formatter/linter — do not hand-format.

## Testing

- **Unit tests** (no server) live next to the code. Each method should cover:
  path + HTTP method, happy-path decode, HTTP-error propagation, bad-JSON for
  reads/creates, and pre-HTTP validation. Reuse the shared helpers
  (`newMockClient`, `badHostClient`); don't add duplicate package-level helpers.
- **Integration/e2e tests** run against a live Directus via Docker — see
  [`docs/INTEGRATION_TESTING.md`](docs/INTEGRATION_TESTING.md):

  ```sh
  task test:integration          # up --wait -> tagged tests -> compose down
  task test:integration:report   # same, plus JUnit XML + JSON in test-reports/
  ```

- **Reports.** `task test:report` / `task test:integration:report` emit per-test
  PASS/FAIL and machine-readable JUnit XML + JSON under `test-reports/` (gitignored).

Add or update tests for any behavioral change. If you add or change a public
method, add unit coverage and, where it exercises real server behavior, extend
the e2e suite.

## Continuous integration

GitHub Actions (`.github/workflows/go-tests.yml`) runs `go test -v -cover
./directus/...` on Go 1.26 for pushes and PRs to `main`. Integration tests are
**not** run in CI (they need a live stack) — run them locally when your change
touches request building, payloads, or endpoints. Ensure `task ci` is green
before pushing.

## Commit messages

This repo uses **[Conventional Commits](https://www.conventionalcommits.org)**
(the JetBrains Conventional Commit plugin is configured). Format:

```
<type>(<optional scope>): <summary>
```

Common types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `ci`.
Examples:

```
feat(users): add two-factor endpoints
fix(notifications): treat id as int (Directus returns a number)
docs(api): note POST /permissions returns all rows
test(integration): cover schema snapshot/diff/apply
```

Keep commits focused; write the body to explain *why* when it isn't obvious.

## Pull requests

1. Branch from `main`.
2. Make the change with tests and updated docs (`docs/API.md` for new/changed
   methods; the relevant guide for behavior).
3. Run `task ci` (and `task test:integration` if your change affects HTTP
   payloads/endpoints).
4. Open the PR with a clear description of what changed and why; link any related
   issue.
5. Keep the diff clean — no stray formatting churn, no committed artifacts
   (`coverage.*`, `test-reports/`, `data/` are gitignored).

## Reporting issues

Open a GitHub issue with: the SDK version/commit, Directus version, a minimal
reproduction (code + the request/response if relevant), and what you expected
versus what happened.
