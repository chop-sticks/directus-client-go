// Package integration holds end-to-end tests that exercise the directus SDK
// against a live Directus instance (see docker-compose.yml at the repo root).
//
// Every test file carries the `integration` build tag, so these tests are
// excluded from the normal `go test ./...` run. Bring up the stack and run
// them with:
//
//	task test:integration
//
// or manually:
//
//	docker compose up -d --wait
//	go test -tags=integration -count=1 ./directus/integration/...
//	docker compose down
//
// The SDK is dot-imported by the (build-tagged) test files so the test bodies
// read naturally; the target instance and credentials come from DIRECTUS_URL /
// DIRECTUS_TOKEN, defaulting to the docker-compose values.
package integration
