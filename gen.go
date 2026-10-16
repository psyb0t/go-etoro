// Package etoro is a Go client for the eToro public REST API, generated with
// oapi-codegen from the OpenAPI document eToro publishes. The WebSocket client
// lives in the ws subpackage.
package etoro

//go:generate bash scripts/fetch-spec.sh openapi.json
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.json
