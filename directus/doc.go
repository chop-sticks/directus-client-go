// Package directus is a typed Go client for the [Directus] REST API.
//
// It mirrors the official JavaScript/TypeScript SDK command surface, exposing
// CRUD and action methods for items, files, folders, users, roles, policies,
// permissions, flows, operations, panels, dashboards, presets, translations,
// shares, comments, notifications, activity, revisions, content versions,
// relations, settings, extensions, collections, and fields, plus the server,
// schema, utils, and auth endpoints.
//
// # Getting started
//
// Create a [Client] with an instance URL and a static token, then call a
// resource method:
//
//	host := "https://example.directus.app"
//	token := "your-static-token"
//
//	client, err := directus.NewClient(&host, &token)
//	if err != nil {
//		return err
//	}
//
//	users, err := client.GetUsers(&directus.Query{
//		Fields: []string{"id", "email"},
//		Filter: map[string]any{"status": map[string]any{"_eq": "active"}},
//	})
//
// To authenticate with credentials instead of a static token, call [Client.Login]
// and assign the returned access token to Client.Token.
//
// # Query parameters
//
// Most read and write methods accept a *[Query] that serializes the Directus
// global query parameters (fields, filter, sort, limit, offset, page, search,
// deep, aggregate, group by, alias, version, export, and meta). Pass nil when
// no query is needed.
//
// # Models and return types
//
// Core (system) collections have typed model structs named after the Directus
// type without the "Directus" prefix ([User], [Role], [File], and so on).
// Relational fields that Directus may return either as an ID string or as an
// expanded object are typed as any. Arbitrary user collections are handled
// generically as map[string]any via [Client.GetItems] and the related item
// methods. List and batch methods return slices; single-item reads, creates,
// and updates return pointers (nil when the API returns no data).
//
// # Conventions
//
// Method names follow Get/Create/Patch/Delete prefixes. "Update many by keys"
// methods (PatchXs) send {"keys":[...],"data":{...}}; batch variants
// (PatchXsBatch) send the raw item array. See the docs directory (DEVELOPMENT.md
// and API.md) for the full endpoint reference and extension guide.
//
// [Directus]: https://directus.io
package directus
