# Directus Go Client — API Reference (implemented surface)

Authoritative inventory of what currently exists in the `directus` package.
Generated from source; keep in sync when adding methods. For *how to extend*
(architecture, helpers, conventions, endpoint→URL mapping, model derivation) see
[`DEVELOPMENT.md`](DEVELOPMENT.md).

Status: **complete** — full parity with the JS/TS SDK REST command surface.
- 32 source files, 226 exported `*Client` methods, ~95% test coverage.
- Build/test: `go build ./...` and `go test ./...` (from repo root). Both green.
- Module: `github.com/chop-sticks/directus-client-go`, Go 1.26.

## Core (client + query)

- `client.go` — `Client{HostURL, HTTPClient, Token}`, `NewClient(host, token *string)`,
  `doRequest`, and the generics: `request[T](c, method, path, q, body) (T, error)`,
  `(c) requestRaw(...) ([]byte, error)`, `(c) execute(...) error`.
- `query_params.go` — `Query{...}`, `buildQueryString`, `buildURL`.
- `doc.go` — package-level godoc.
- Model structs are co-located with their resource: each `<resource>.go` holds
  its own type(s) plus methods. The `Collection*`/`CollectionRequest` types live
  in `collections.go`; the `Field*` types live in `fields.go`. There is no
  central `models.go`.

## Models (Go type → id type)

Go type names drop the `Directus` prefix used by the JS SDK.

| Go type | id | Go type | id | Go type | id |
|---|---|---|---|---|---|
| `Collection` | name (string) | `Field` | composite | `User` | string |
| `Role` | string | `Policy` | string | `Permission` | **int** |
| `File` | string | `Folder` | string | `Flow` | string |
| `Operation` | string | `Panel` | string | `Dashboard` | string |
| `Preset` | **int** | `Translation` | string | `Share` | string |
| `Comment` | string | `Notification` | **int** | `Activity` | **int** (read-only) |
| `Revision` | **int** (read-only) | `Version` | string | `Relation` | composite |
| `Settings` | int (singleton) | `Extension` | string | | |

Auxiliary result/struct types: `AuthenticationData`, `AuthProvider`,
`PolicyGlobals`, `ShareInfo`, `VersionCompare`, `ServerHealth`, `SchemaSnapshot`,
`SchemaDiff`, `ExtensionMeta`, `RelationMeta`, `RelationSchema`.

## Method inventory (by resource file)

### collections.go
- `GetCollections() ([]Collection, error)`
- `GetCollectionByName(name string) (*Collection, error)`
- `CreateCollection(req *CollectionRequest, q *Query) (*Collection, error)`
- `PatchCollection(name string, req *CollectionRequest, q *Query) (*Collection, error)`
- `PatchCollectionsBatch(items []CollectionRequest, q *Query) ([]Collection, error)`
- `DeleteCollection(name string) error`

### fields.go
- `GetFields() ([]Field, error)`
- `GetFieldsByCollection(collection string) ([]Field, error)`
- `GetFieldByCollectionAndName(collection, name string) (*Field, error)`
- `CreateField(collection string, field *Field, q *Query) (*Field, error)`
- `PatchField(collection, name string, field *Field, q *Query) (*Field, error)`
- `PatchFields(collection string, fields []Field, q *Query) ([]Field, error)`
- `DeleteField(collection, name string) error`

### items.go (generic + singleton + aggregate; values are `map[string]any`)
- `GetItems(collection string, q *Query) ([]map[string]any, error)`
- `GetItem(collection, key string, q *Query) (map[string]any, error)`
- `CreateItem(collection string, item map[string]any, q *Query) (map[string]any, error)`
- `CreateItems(collection string, items []map[string]any, q *Query) ([]map[string]any, error)`
- `PatchItem(collection, key string, item map[string]any, q *Query) (map[string]any, error)`
- `PatchItems(collection string, keys []any, item map[string]any, q *Query) ([]map[string]any, error)`
- `PatchItemsBatch(collection string, items []map[string]any, q *Query) ([]map[string]any, error)`
- `DeleteItem(collection, key string) error`
- `DeleteItems(collection string, keys []any) error` — body `{"keys":[...]}`
- `DeleteItemsByQuery(collection string, query map[string]any) error` — body `{"query":{...}}`
- `GetSingleton(collection string, q *Query) (map[string]any, error)`
- `PatchSingleton(collection string, item map[string]any, q *Query) (map[string]any, error)`
- `Aggregate(collection string, q *Query) ([]map[string]any, error)` — set `q.Aggregate` / `q.GroupBy`

### users.go
- `GetUsers(q) ([]User, error)`; `GetUser(id string, q) (*User, error)`; `GetUsersMe(q) (*User, error)`
- `CreateUser(item *User, q) (*User, error)`; `CreateUsers(items []User, q) ([]User, error)`
- `PatchUser(id string, item *User, q) (*User, error)`; `PatchUsers(keys []string, item *User, q) ([]User, error)`; `PatchUsersBatch(items []User, q) ([]User, error)`; `PatchUsersMe(item *User, q) (*User, error)`
- `DeleteUser(id string) error`; `DeleteUsers(keys []string) error`
- `InviteUser(email, role, inviteURL string) error`; `AcceptUserInvite(token, password string) error`
- `RegisterUser(email, password string, opts map[string]any) error`; `RegisterUserVerify(token string) error`
- `GenerateTwoFactorSecret(password string) (map[string]any, error)`; `EnableTwoFactor(secret, otp string) error`; `DisableTwoFactor(otp string) error`

### roles.go
- `GetRoles(q)`; `GetRole(id string, q)`; `GetRolesMe(q) ([]Role, error)`
- `CreateRole`/`CreateRoles`; `PatchRole(id)`/`PatchRoles(keys)`/`PatchRolesBatch`; `DeleteRole(id)`/`DeleteRoles(keys)`

### policies.go
- `GetPolicies(q)`; `GetPolicy(id string, q)`; `GetPolicyGlobals() (*PolicyGlobals, error)`
- `CreatePolicy`/`CreatePolicies`; `PatchPolicy`/`PatchPolicies`/`PatchPoliciesBatch`; `DeletePolicy`/`DeletePolicies`

### access.go
- `GetAccesses(q)`; `GetAccess(id string, q)`
- `CreateAccess`/`CreateAccesses`; `PatchAccess(id)`/`PatchAccesses(keys)`/`PatchAccessesBatch`; `DeleteAccess(id)`/`DeleteAccesses(keys)`

### permissions.go (id `int`)
- `GetPermissions(q)`; `GetPermission(id int, q)`
- `GetUserPermissions() (map[string]any, error)` — `/permissions/me`
- `GetItemPermissions(collection, key string) (map[string]any, error)` — `/permissions/me/{collection}[/{key}]`
- `CreatePermission`/`CreatePermissions`; `PatchPermission(id int)`/`PatchPermissions(keys []int)`/`PatchPermissionsBatch`; `DeletePermission(id int)`/`DeletePermissions(keys []int)`

### files.go / folders.go / assets.go
- `GetFiles(q)`; `GetFile(id, q)`; `ImportFile(url string, data *File, q) (*File, error)`
- `UploadFile(data io.Reader, filename string, fields map[string]string, q) (*File, error)`
- `PatchFile(id)`/`PatchFiles(keys)`/`PatchFilesBatch`; `DeleteFile(id)`/`DeleteFiles(keys)`
- `GetFolders(q)`; `GetFolder(id, q)`; `CreateFolder`/`CreateFolders`; `DeleteFolder(id)`/`DeleteFolders(keys)` — **no folder update endpoint exists**
- `GetAsset(id string, q) ([]byte, error)`; `DownloadFilesZip(ids []string) ([]byte, error)`; `DownloadFolderZip(id string) ([]byte, error)`

### flows.go / operations.go / panels.go / dashboards.go
Each has the standard CRUD set: `GetXs(q)`, `GetX(id, q)`, `CreateX`/`CreateXs`,
`PatchX(id)`/`PatchXs(keys)`/`PatchXsBatch`, `DeleteX(id)`/`DeleteXs(keys)`.
- flows.go additionally: `TriggerFlow(method, id string, data map[string]string) ([]byte, error)`

### presets.go (id `int`)
- `GetPresets(q)`; `GetPreset(id int, q)`; `CreatePreset`/`CreatePresets`; `PatchPreset(id int)`/`PatchPresets(keys []int)`/`PatchPresetsBatch`; `DeletePreset(id int)`/`DeletePresets(keys []int)`

### translations.go / dashboards.go
Standard string-id CRUD (Get/GetOne/Create/Creates/Patch/Patch-by-keys/Batch/Delete/Deletes).

### notifications.go (id `int`)
Standard CRUD with int ids: `GetNotifications(q)`, `GetNotification(id int, q)`,
`CreateNotification`/`CreateNotifications`, `PatchNotification(id int)`/`PatchNotifications(keys []int)`/`PatchNotificationsBatch`, `DeleteNotification(id int)`/`DeleteNotifications(keys []int)`.

### shares.go
- Standard CRUD, plus:
- `AuthenticateShare(share, password, mode string) (*AuthenticationData, error)`
- `InviteShare(share string, emails []string) error`
- `ReadShareInfo(id string) (*ShareInfo, error)`

### comments.go
- `GetComments(q)`; `GetComment(id, q)`; `CreateComment`/`CreateComments`; `PatchComment(id)`; `DeleteComment(id)`; `DeleteComments(keys []string) error` — body `{"keys":[...]}` (object, **not** a raw array)

### activity.go / revisions.go (read-only, id `int`)
- `GetActivities(q)`; `GetActivity(id int, q)`
- `GetRevisions(q)`; `GetRevision(id int, q)`

### versions.go
- Standard CRUD on `/versions`, plus:
- `SaveToContentVersion(id string, item map[string]any) (map[string]any, error)`
- `CompareContentVersion(id string) (*VersionCompare, error)`
- `PromoteContentVersion(id, mainHash string, fields []string) (any, error)`

### relations.go
- `GetRelations()`; `GetRelationsByCollection(collection string)`; `GetRelation(collection, field string) (*Relation, error)`
- `CreateRelation(item *Relation) (*Relation, error)` — no query
- `PatchRelation(collection, field string, item *Relation, q) (*Relation, error)`; `DeleteRelation(collection, field string) error`

### settings.go
- `GetSettings(q) (*Settings, error)`; `PatchSettings(item *Settings, q) (*Settings, error)`

### extensions.go
- `GetExtensions() ([]Extension, error)`; `GetRegistryExtensions(q) ([]map[string]any, error)`
- `PatchExtension(name string, item map[string]any) (*Extension, error)`; `PatchBundleExtension(bundle, name string, item map[string]any) (*Extension, error)`
- `DeleteExtension(id string) error`; `InstallRegistryExtension(extensionID, version string) error`; `UninstallRegistryExtension(id string) error`

### server.go
- `ServerHealth() (*ServerHealth, error)`; `ServerPing() (string, error)`; `ServerInfo() (map[string]any, error)`; `ReadOpenAPISpec() (map[string]any, error)`; `ReadGraphqlSDL(scope string) (string, error)`

### schema.go
- `GetSchemaSnapshot(includeCollections, excludeCollections []string) (*SchemaSnapshot, error)`
- `SchemaDiffSnapshot(snapshot *SchemaSnapshot, force bool, mode string) (*SchemaDiff, error)`
- `SchemaApply(diff *SchemaDiff, force bool) error`

### utils.go
- `ClearCache(system bool) error`; `RandomString(length int) (string, error)`; `UtilitySort(collection string, item, to any) error`; `UtilsExport(collection, format string, query, file map[string]any) error`; `UtilsImport(collection string, data io.Reader, filename string) error`

### auth.go
- `Login(email, password, mode, otp string) (*AuthenticationData, error)`; `Refresh(refreshToken, mode string) (*AuthenticationData, error)`; `Logout(refreshToken, mode string) error`
- `PasswordRequest(email, resetURL string) error`; `PasswordReset(token, password string) error`; `ReadProviders(sessionOnly bool) ([]AuthProvider, error)`

## Resolved decisions & gotchas (non-obvious, honor when extending)

- **Response envelope**: `request[T]` unwraps `{"data": …}`. Endpoints that are
  NOT enveloped use `requestRaw` + manual decode: `/server/health`,
  `/server/ping` (plain `pong`), `/server/specs/oas`, `/server/specs/graphql*`.
  `/server/info`, `/schema/*`, `/utils/random/string`, auth, and all CRUD ARE
  enveloped.
- **Empty body → zero value**: 204/empty decodes to the zero `T` (nil pointer,
  nil slice). `SchemaDiffSnapshot` returns `nil, nil` when there is no diff.
- **id validation**: string ids validated `== ""`; int ids validated `<= 0`;
  both return an error *before* any HTTP call (mirrors `fields.go`).
- **Delete-many body**: raw `[]keys` for most resources via `execute("DELETE", …, keys)`;
  EXCEPTIONS — `DeleteComments` and `DeleteItems` send `{"keys":[...]}`,
  `DeleteItemsByQuery` sends `{"query":{...}}`.
- **Update-many-by-keys** body is `{"keys":[...], "data":{...}}`; **batch update**
  sends the raw `[]items`. Both PATCH the collection base path.
- **`q *Query` presence** mirrors the SDK: reads of collections/fields, all
  deletes, relations create, and action endpoints take NO query.
- **Multipart exceptions**: `UploadFile` and `UtilsImport` hand-build
  `*http.Request` (via `buildURL` + `c.doRequest`) because the multipart boundary
  must be set in `Content-Type` — the only two methods that bypass the shared
  helpers. Marked with comments in-source.
- **Auth defaults**: `Login`/`Refresh`/`Logout` default `mode` to `"json"` when
  empty; optional `otp`/`refresh_token`/`reset_url`/`invite_url` are omitted from
  the body when empty.
- **Schema query params** (`includeCollections`, `excludeCollections`, `force`,
  `mode`) are NOT part of `Query`; they are appended to the path via `url.Values`
  inside `schema.go`. `TriggerFlow` similarly builds its own GET query string.
- **`GetItemPermissions`** omits the `/{key}` path segment when `key == ""`.
- **Relational fields** (`user`, `role`, `folder`, `user_created`, …) are typed
  `any` because Directus returns either an id string or an expanded object
  depending on `fields`.
- **`UtilsImport` media type**: the multipart file part sets `Content-Type` from
  the filename extension (`.json`/`.csv`/`.xml`/`.yaml`); Directus rejects the
  default `application/octet-stream`.
- **`Field`/`Relation` partial writes**: `Type`/`Meta`/`Schema` (Field) and
  `Collection`/`Field`/`RelatedCollection`/`Meta`/`Schema` (Relation) are
  `omitempty` so a partial `Patch*` never sends an empty required key. Field
  create/patch also inject `collection` (and field name) from the URL via
  `withFieldLocation`.
- **Notification ids are integers** despite the JS SDK typing them `string`;
  verified against a live instance.
- **`POST /permissions` (batch) returns ALL permissions**, not just the created
  rows — do not infer created ids from the response length; query by policy or
  create singly. Custom permission rules (non-full-access `Patch*`) are a
  licensed feature and return 403 `custom_permission_rules_enabled` on an
  unlicensed instance.

## Testing conventions

Shared test helpers (same package, do not redefine): `newMockClient(t, handler)`
in `client_test.go`, `badHostClient(t)` in `collections_test.go`, `RoundTripFunc`,
and `intPtr(int) *int` in `query_params_test.go`. Each method is covered for
path+method, happy-path decode, HTTP-error propagation, bad-JSON (reads/creates),
and pre-HTTP validation (via a `called` bool in the handler).

### End-to-end integration tests

The `directus/integration` package (build tag `integration`, SDK dot-imported)
exercises EVERY method against a live Directus (docker compose). Shared fixtures
live in `directus/integration/helpers_test.go`
(`itestClient`, `uniqueName`, `newTestCollection`/`newTestItem`/`newTestRole`/
`newTestPolicy`/`newTestDashboard`/`newTestFolder`/`newTestUser`/`uploadTestFile`),
each self-registering `t.Cleanup`. Achievable methods assert success; methods
needing external prerequisites (email transport, extension registry, invalid
tokens/OTPs) assert the expected error or are log-only. Run with
`task test:integration`. Full guide: [INTEGRATION_TESTING.md](INTEGRATION_TESTING.md).
