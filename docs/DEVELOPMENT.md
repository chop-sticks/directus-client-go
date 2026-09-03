# Directus Go Client — Development Guide

This document is the single source of truth for extending this SDK. It captures
the architecture, conventions, the full Directus REST endpoint surface, and the
Go model mapping. Endpoint definitions were derived from the official JS/TS SDK
(`directus/directus` → `sdk/src/rest/commands` and `sdk/src/schema`).

> **Status: implemented.** Every resource below exists. This file is the guide
> for *how to extend* the SDK (architecture, helpers, conventions, endpoint→URL
> mapping, model derivation). For the authoritative inventory of *what exists*
> today — actual Go type names and method signatures — see
> [`API.md`](API.md).

## Architecture

The package is a thin, typed wrapper over the Directus REST API.

- `client.go` — `Client`, `NewClient`, `doRequest`, and the shared request
  helpers (`request[T]`, `requestRaw`, `execute`).
- `query_params.go` — the `Query` struct and URL builders.
- `doc.go` — package-level godoc overview.
- `<resource>.go` — one file per resource, holding its model struct(s) **and**
  its `Client` methods (this is the convention — types live with their
  behavior; there is no central `models.go`). `Collection*` types live in
  `collections.go`, `Field*` types in `fields.go`.
- `<resource>_test.go` — table/mock-server tests for that resource.

### Request helpers (use these; do not hand-roll `http.NewRequest`)

```go
// Decodes the Directus {"data": ...} envelope into T.
//   list:   request[[]User](c, "GET", "/users", q, nil)
//   single: request[*User](c, "GET", "/users/"+id, q, nil)   // nil if data is null
//   scalar: request[string](c, "GET", "/utils/random/string", q, nil)
func request[T any](c *Client, method, path string, q *Query, body any) (T, error)

// Returns the raw response body (no envelope). For /server/health, /server/ping,
// spec endpoints, plain-text and non-enveloped payloads.
func (c *Client) requestRaw(method, path string, q *Query, body any) ([]byte, error)

// Fire-and-forget: DELETE and action endpoints where only the error matters.
func (c *Client) execute(method, path string, q *Query, body any) error
```

`body` handling in `requestRaw`: `nil` → no body; `io.Reader` or `[]byte` sent
verbatim (used for multipart uploads); anything else is JSON-encoded.
`doRequest` accepts 200/201/204 and sets `Content-Type: application/json` unless
already set. An empty response body decodes to the zero value of `T`.

### The `Query` struct

`Query` serializes Directus global query parameters. Pass `nil` when no query is
needed. Fields: `Fields []string`, `Filter map[string]any`, `Sort []string`,
`Limit/Offset/Page *int`, `Search string`, `Deep/Aggregate map[string]any`,
`GroupBy []string`, `Alias map[string]string`, `Version/Export/Meta string`.
`Filter`, `Deep`, `Aggregate`, `Alias` are JSON-encoded; `Fields`/`Sort`/`GroupBy`
are comma-joined.

## Conventions

### Method naming (mirror the SDK, Go-idiomatic verbs)

| Operation                    | Prefix   | Example                                   |
|------------------------------|----------|-------------------------------------------|
| Read list                    | `Get`    | `GetUsers(q *Query) ([]User, error)`      |
| Read one                     | `Get`    | `GetUser(id string, q *Query) (*User, error)` |
| Read current/self            | `Get…Me` | `GetUsersMe(q *Query) (*User, error)`     |
| Create one                   | `Create` | `CreateUser(item *User, q *Query) (*User, error)` |
| Create many                  | `Create` | `CreateUsers(items []User, q *Query) ([]User, error)` |
| Update one                   | `Patch`  | `PatchUser(id string, item *User, q *Query) (*User, error)` |
| Update many (keys + data)    | `Patch`  | `PatchUsers(keys []string, item *User, q *Query) ([]User, error)` |
| Update many (batch of items) | `Patch`  | `PatchUsersBatch(items []User, q *Query) ([]User, error)` |
| Delete one                   | `Delete` | `DeleteUser(id string) error`             |
| Delete many                  | `Delete` | `DeleteUsers(keys []string) error`        |

Rules:
- Read one / create one / patch one return `*T`; lists and batches return `[]T`.
- A method takes `q *Query` **only if** the SDK command accepts a `query`.
  (Collection/field reads and most delete/action endpoints do not.)
- "Update many by keys" sends `{"keys": [...], "data": {...}}`.
  "Batch update" sends the array of items directly. Both PATCH the base path.
- "Delete many" sends the raw keys array as the body (except comments/items,
  which send `{"keys":[...]}` or `{"query":{...}}` — noted per resource below).
- Validate required path params up front and return an error before any HTTP
  call, e.g. `if id == "" { return nil, fmt.Errorf("id must be provided") }`.

### Models

- Go type name = the Directus type without the `Directus` prefix
  (`DirectusUser` → `User`, `DirectusFlow` → `Flow`). The resource headings and
  model reference below use the SDK's `DirectusX` names; the Go structs drop the
  prefix. Result/auxiliary structs live in the owning resource file.
- One flat struct per Directus collection (core collections have no meta/schema
  split — only `Field`, `Relation`, and `Collection` do).
- Relational fields that the API may return as either an ID string or an
  expanded object are typed `any` (e.g. `User any` / `UserCreated any`).
- Timestamps typed `'datetime'` in the SDK → Go `string` (ISO-8601 as returned).
- `Record<string, any>` → `map[string]any`; `T[] | null` → `[]T`.
- Use `omitempty` on optional write fields; keep required booleans/ints without
  `omitempty` only when the API needs the explicit value.
- Numeric IDs (activity, revisions, permissions, presets, fields) → `int`;
  UUID/string IDs → `string`.

### Testing

Helpers live in `client_test.go` / `collections_test.go` (same package):
- `newMockClient(t, handler http.HandlerFunc) *Client` — httptest server.
- `badHostClient(t) *Client` — client whose host fails `http.NewRequest` (for
  request-build-error tests).

Each resource test SHOULD cover, per method: path + HTTP method assertion,
happy-path decode, HTTP-error propagation, and (for reads/creates) a bad-JSON
case. For methods with input validation, assert the error fires before any HTTP
call. Do not run the whole suite in a subtask; the integrator runs it.

### Integration tests

The `directus/integration` package (its own external test package, build tag
`integration`, SDK dot-imported) exercises a live
Directus instance from `docker-compose.yml`. It is excluded from the normal
unit run. Add integration coverage here when a change needs end-to-end
verification against a real server. Run it with `task test:integration` (boots
the stack with `docker compose up -d --wait`, runs `go test -tags=integration`,
then `docker compose down`). Target/credentials come from `DIRECTUS_URL` /
`DIRECTUS_TOKEN`, defaulting to the compose values. Full guide:
[INTEGRATION_TESTING.md](INTEGRATION_TESTING.md).

## Endpoint reference

Format: `METHOD path` — notes. `{q}` marks endpoints that accept a `Query`.
The `data` envelope is implicit unless stated. IDs are string (UUID) unless the
model reference marks them `int`.

### items (generic, non-core collections) — `items.go`
- `GET /items/{collection}` {q} → `[]map[string]any` (or generic `[]T` — expose as `GetItems(collection string, q) ([]map[string]any, error)`)
- `GET /items/{collection}/{key}` {q} → single
- `POST /items/{collection}` {q} body=item → single; body=[]items → list (`CreateItem`/`CreateItems`)
- `PATCH /items/{collection}/{key}` {q} body=item → single
- `PATCH /items/{collection}` {q} body=`{keys,data}` → list; batch body=[]items → list
- `DELETE /items/{collection}/{key}`; `DELETE /items/{collection}` body=`{"keys":[...]}` or `{"query":{...}}`
- Item values are schema-defined; use `map[string]any` for payloads/returns.

### singleton — `singleton.go` (part of items.go acceptable)
- `GET /items/{collection}` {q} → single object (`GetSingleton`)
- `PATCH /items/{collection}` {q} body=item → single (`PatchSingleton`)

### aggregate — (in items.go)
- `GET /items/{collection}` with params `aggregate`, optional `groupBy`, plus query → `[]map[string]any`

### users — `users.go` (model `DirectusUser`)
- `GET /users` {q}; `GET /users/{id}` {q}; `GET /users/me` {q} (`GetUsersMe`)
- `POST /users` {q} (one/many); `PATCH /users/{id}` {q}; `PATCH /users` {q} keys+data / batch
- `PATCH /users/me` {q} body=item → single (`PatchUsersMe`)
- `DELETE /users/{id}`; `DELETE /users` body=[]keys

### roles — `roles.go` (`DirectusRole`)
- `GET /roles` {q}; `GET /roles/{id}` {q}; `GET /roles/me` {q} → list (`GetRolesMe`)
- `POST /roles` {q} one/many; `PATCH /roles/{id}` {q}; `PATCH /roles` {q} keys+data / batch
- `DELETE /roles/{id}`; `DELETE /roles` body=[]keys

### policies — `policies.go` (`DirectusPolicy`)
- `GET /policies` {q}; `GET /policies/{id}` {q}; `GET /policies/me/globals` → `{app_access,admin_access,enforce_tfa}` (`GetPolicyGlobals`)
- `POST /policies` {q} one/many; `PATCH /policies/{id}` {q}; `PATCH /policies` {q} keys+data / batch
- `DELETE /policies/{id}`; `DELETE /policies` body=[]keys

### permissions — `permissions.go` (`DirectusPermission`, id `int`)
- `GET /permissions` {q}; `GET /permissions/{id}` {q}
- `GET /permissions/me` → map (`GetUserPermissions`); `GET /permissions/me/{collection}[/{key}]` → item permissions (`GetItemPermissions`)
- `POST /permissions` {q} one/many; `PATCH /permissions/{id}` {q}; `PATCH /permissions` {q} keys+data / batch
- `DELETE /permissions/{id}`; `DELETE /permissions` body=[]keys

### files — `files.go` (`DirectusFile`)
- `GET /files` {q}; `GET /files/{id}` {q}
- `POST /files` multipart body (FormData) {q} → single (`UploadFile`, build with `mime/multipart`, set writer Content-Type)
- `POST /files/import` {q} body=`{url,data}` → single (`ImportFile`)
- `PATCH /files/{id}` {q} body=item (JSON) → single; `PATCH /files` {q} keys+data / batch
- `DELETE /files/{id}`; `DELETE /files` body=[]keys

### folders — `folders.go` (`DirectusFolder`)
- `GET /folders` {q}; `GET /folders/{id}` {q}
- `POST /folders` {q} one/many; `DELETE /folders/{id}`; `DELETE /folders` body=[]keys
- NOTE: the SDK has no folder update command — implement Get/Create/Delete only (no `PatchFolder`).

### assets — `assets.go`
- `GET /assets/{id}` {q} → raw bytes (`GetAsset(id, q) ([]byte, error)` via `requestRaw`)
- `POST /assets/files/` body=`{ids:[...]}` → raw ZIP bytes (`DownloadFilesZip`)
- `POST /assets/folder/{id}` → raw ZIP bytes (`DownloadFolderZip`)

### folders/files sub-note
Files upload uses `multipart/form-data`; the file bytes go in a `file` form
field; metadata fields are additional form fields.

### flows — `flows.go` (`DirectusFlow`)
- `GET /flows` {q}; `GET /flows/{id}` {q}
- `POST /flows` {q} one/many; `PATCH /flows/{id}` {q}; `PATCH /flows` {q} keys+data / batch (`PatchFlow`, `PatchFlows`, `PatchFlowsBatch`)
- `DELETE /flows/{id}`; `DELETE /flows` body=[]keys
- `GET|POST /flows/trigger/{id}` (in utils) → arbitrary (`TriggerFlow`, use `requestRaw`)

### operations — `operations.go` (`DirectusOperation`)
- `GET /operations` {q}; `GET /operations/{id}` {q}
- `POST /operations` {q} one/many; `PATCH /operations/{id}` {q}; `PATCH /operations` {q} keys+data / batch
- `DELETE /operations/{id}`; `DELETE /operations` body=[]keys

### panels — `panels.go` (`DirectusPanel`)
- Same CRUD shape as operations on `/panels`.

### dashboards — `dashboards.go` (`DirectusDashboard`)
- `GET /dashboards` {q}; `GET /dashboards/{id}` {q}
- `POST /dashboards` {q} one/many; `PATCH /dashboards/{id}` {q}; `PATCH /dashboards` {q} keys+data / batch
- `DELETE /dashboards/{id}`; `DELETE /dashboards` body=[]keys

### presets — `presets.go` (`DirectusPreset`, id `int`)
- `GET /presets` {q}; `GET /presets/{id}` {q}
- `POST /presets` {q} one/many; `PATCH /presets/{id}` {q}; `PATCH /presets` {q} keys+data / batch
- `DELETE /presets/{id}`; `DELETE /presets` body=[]keys

### translations — `translations.go` (`DirectusTranslation`)
- `GET /translations` {q}; `GET /translations/{id}` {q}
- `POST /translations` {q} one/many; `PATCH /translations/{id}` {q}; `PATCH /translations` {q} keys+data / batch
- `DELETE /translations/{id}`; `DELETE /translations` body=[]keys

### shares — `shares.go` (`DirectusShare`)
- `GET /shares` {q}; `GET /shares/{id}` {q}
- `POST /shares` {q} one/many; `PATCH /shares/{id}` {q}; `PATCH /shares` {q} keys+data / batch
- `DELETE /shares/{id}`; `DELETE /shares` body=[]keys
- utils: `POST /shares/auth` body=`{share,password,mode}` → auth data (`AuthenticateShare`); `POST /shares/invite` body=`{share,emails}` (`InviteShare`); `GET /shares/info/{id}` → share info (`ReadShareInfo`)

### comments — `comments.go` (`DirectusComment`)
- `GET /comments` {q}; `GET /comments/{id}` {q}
- `POST /comments` {q} one/many; `PATCH /comments/{id}` {q}
- `DELETE /comments/{id}`; `DELETE /comments` body=`{"keys":[...]}` (array) or `{"query":{...}}`

### notifications — `notifications.go` (`DirectusNotification`)
- `GET /notifications` {q}; `GET /notifications/{id}` {q}
- `POST /notifications` {q} one/many; `PATCH /notifications/{id}` {q}; `PATCH /notifications` {q} keys+data / batch
- `DELETE /notifications/{id}`; `DELETE /notifications` body=[]keys

### activity — `activity.go` (`DirectusActivity`, id `int`) — read-only
- `GET /activity` {q}; `GET /activity/{id}` {q}

### revisions — `revisions.go` (`DirectusRevision`, id `int`) — read-only
- `GET /revisions` {q}; `GET /revisions/{id}` {q}

### versions (content versions) — `versions.go` (`DirectusVersion`)
- `GET /versions` {q}; `GET /versions/{id}` {q}
- `POST /versions` {q} one/many; `PATCH /versions/{id}` {q}; `PATCH /versions` {q} keys+data / batch
- `DELETE /versions/{id}`; `DELETE /versions` body=[]keys
- `POST /versions/{id}/save` body=item → item state (`SaveToContentVersion`)
- `GET /versions/{id}/compare` → `{outdated,mainHash,current,main}` (`CompareContentVersion`)
- `POST /versions/{id}/promote` body=`{mainHash,fields?}` → item pk (`PromoteContentVersion`)

### relations — `relations.go` (`DirectusRelation`, meta id `int`)
- `GET /relations`; `GET /relations/{collection}`; `GET /relations/{collection}/{field}` (no query)
- `POST /relations` body=item → single (`CreateRelation`, no query)
- `PATCH /relations/{collection}/{field}` {q} body=item → single (`PatchRelation`)
- `DELETE /relations/{collection}/{field}` (`DeleteRelation`)

### settings — `settings.go` (`DirectusSettings`, singleton) 
- `GET /settings` {q} → single (`GetSettings`)
- `PATCH /settings` {q} body=item → single (`PatchSettings`)

### extensions — `extensions.go` (`DirectusExtension`)
- `GET /extensions/` → `[]Extension` (`GetExtensions`)
- `GET /extensions/registry` {q: search/limit/offset/sort/filter} → registry list
- `GET /extensions/registry/account/{pk}`; `GET /extensions/registry/extension/{pk}`
- `POST /extensions/registry/install` body=`{extension,version}` (`InstallRegistryExtension`)
- `DELETE /extensions/registry/uninstall/{id}`; `DELETE /extensions/{id}` (`DeleteExtension`)
- `PATCH /extensions/{name}` also exists in SDK update/extensions — body=`{meta:{enabled}}`; implement `PatchExtension(bundle?, name, item)` → `PATCH /extensions/{name}` or `/extensions/{bundle}/{name}`. Model `Extension` has `bundle`, `schema`, `meta{enabled}`.

### server — `server.go` (use `requestRaw`; NOT data-enveloped except info/specs/oas)
- `GET /server/health` → `{status,...}` NOT enveloped (`ServerHealth`)
- `GET /server/ping` → text `pong` (`ServerPing`)
- `GET /server/info` → `{data:{...}}` enveloped (`ServerInfo`)
- `GET /server/specs/oas` → OpenAPI JSON object, NOT enveloped (`ReadOpenAPISpec`)
- `GET /server/specs/graphql` | `/server/specs/graphql/system` → GraphQL SDL text (`ReadGraphqlSDL(scope)`)

### schema — `schema.go` (admin only; enveloped)
- `GET /schema/snapshot` optional params `includeCollections`/`excludeCollections` (comma) → `SchemaSnapshot`
- `POST /schema/diff` params `{force?,mode?}` body=snapshot → diff (or empty on no diff) (`SchemaDiff`)
- `POST /schema/apply` params `{force?}` body=diff → 204 (`SchemaApply`)
- Types: `SchemaSnapshotOutput{version int, directus, vendor string, collections/fields/systemFields/relations []map[string]any}`, `SchemaDiffOutput{hash string, diff map[string]any}`.

### utils — `utils.go`
- `POST /utils/cache/clear` (append `?system` when clearing schema cache) → 204 (`ClearCache(system bool)`)
- `POST /utils/export/{collection}` body=`{format,query,file}` → 204 (`UtilsExport`)
- `POST /utils/import/{collection}` multipart FormData → 204 (`UtilsImport`)
- `GET /utils/random/string` params `{length?}` → string (`RandomString`)
- `POST /utils/sort/{collection}` body=`{item,to}` → 204 (`UtilitySort`)

### auth — `auth.go`
- `POST /auth/login` (or provider endpoint) body=`{email,password,mode,otp?}` → `AuthenticationData{access_token,expires,refresh_token}` (`Login`)
- `POST /auth/refresh` body=`{mode,refresh_token?}` → auth data (`Refresh`)
- `POST /auth/logout` body=`{mode,refresh_token?}` → 204 (`Logout`)
- `POST /auth/password/request` body=`{email,reset_url?}` → 204 (`PasswordRequest`)
- `POST /auth/password/reset` body=`{token,password}` → 204 (`PasswordReset`)
- `GET /auth` (or `/auth?sessionOnly`) → `[]{name,driver,label,icon}` (`ReadProviders`)

### user utilities (in users.go or utils.go)
- `POST /users/invite` body=`{email,role,invite_url?}` (`InviteUser`)
- `POST /users/invite/accept` body=`{token,password}` (`AcceptUserInvite`)
- `POST /users/register` body=`{email,password,...}` (`RegisterUser`)
- `GET /users/register/verify-email?token=` (`RegisterUserVerify`)
- `POST /users/me/tfa/generate` body=`{password}` → `{secret,otpauth_url}` (`GenerateTwoFactorSecret`)
- `POST /users/me/tfa/enable` body=`{secret,otp}` (`EnableTwoFactor`)
- `POST /users/me/tfa/disable` body=`{otp}` (`DisableTwoFactor`)

## Model field reference (Directus schema → Go)

Core collection field lists (from `sdk/src/schema/*.ts`). Relational fields typed
`any`; `'datetime'` → `string`; `Record<string,any>` → `map[string]any`.

- **DirectusUser** (`id string`): status, first_name, last_name, email, password,
  token, last_access, last_page, external_identifier, tfa_secret, auth_data(map),
  provider, appearance, theme_light, theme_dark, theme_light_overrides(map),
  theme_dark_overrides(map), role(any), policies([]any), language, text_direction,
  avatar(any), title, description, location, tags([]string), email_notifications(bool).
- **DirectusRole** (`id string`): name, icon, description, parent(any),
  children([]any), policies([]any), users([]any).
- **DirectusPolicy** (`id string`): name, icon, description, ip_access, enforce_tfa,
  admin_access, app_access, permissions([]any), users([]any), roles([]any).
- **DirectusPermission** (`id int`): policy(any), collection, action, permissions(map),
  validation(map), presets(map), fields([]string).
- **DirectusFile** (`id string`): storage, filename_disk, filename_download, title,
  type, folder(any), created_on, uploaded_by(any), uploaded_on, modified_by(any),
  modified_on, charset, filesize, width(int), height(int), duration(int), embed,
  description, location, tags([]string), metadata(map), focal_point_x(int),
  focal_point_y(int), tus_id, tus_data(map).
- **DirectusFolder** (`id string`): name, parent(any).
- **DirectusFlow** (`id string`): name, icon, color, description, status, trigger,
  accountability, options(map), operation(any), date_created, user_created(any).
- **DirectusOperation** (`id string`): name, key, type, position_x(int), position_y(int),
  options(map), resolve(any), reject(any), flow(any), date_created, user_created(any).
- **DirectusPanel** (`id string`): dashboard(any), name, icon, color, show_header(bool),
  note, type, position_x(int), position_y(int), width(int), height(int), options(map),
  date_created, user_created(any).
- **DirectusDashboard** (`id string`): name, icon, note, date_created, user_created(any), color.
- **DirectusPreset** (`id int`): bookmark, user(any), role(any), collection, search,
  layout, layout_query(map), layout_options(map), refresh_interval(int), filter(map), icon, color.
- **DirectusTranslation** (`id string`): language, key, value.
- **DirectusShare** (`id string`): name, collection, item, role(any), password,
  user_created(any), date_created, date_start, date_end, times_used(int), max_uses(int).
- **DirectusComment** (`id string`): collection(any), item, comment, date_created,
  date_updated, user_created(any), user_updated(any).
- **DirectusNotification** (`id string`): timestamp, status, recipient(any), sender(any),
  subject, message, collection, item.
- **DirectusActivity** (`id int`): action, user(any), timestamp, ip, user_agent,
  collection, item, origin, revisions([]any).
- **DirectusRevision** (`id int`): activity(any), collection, item, data(map),
  delta(map), parent(any), version(any).
- **DirectusVersion** (`id string`): key, name, collection(any), item, hash,
  date_created, date_updated, user_created(any), user_updated(any), delta(map).
- **DirectusRelation**: collection, field, related_collection, meta{ id(int),
  junction_field, many_collection, many_field, one_allowed_collections, one_collection,
  one_collection_field, one_deselect_action, one_field, sort_field, system(bool) },
  schema{ column, constraint_name, foreign_key_column, foreign_key_schema,
  foreign_key_table, on_delete, on_update, table }.
- **DirectusSettings** (`id int` = 1): see `sdk/src/schema/settings.ts` — large; model the
  common fields (project_name, project_descriptor, project_url, default_language,
  default_appearance, project_color, project_logo(any), public_registration(bool),
  storage_asset_transform, custom_css, module_bar, mapbox_key, and the rest as needed).
  Add `Extra map[string]any` is NOT used — prefer explicit fields.
- **DirectusExtension**: `id string`, `bundle string|null`, `schema{ type, local, name, version, ... }|null`, `meta{ enabled bool }`.
