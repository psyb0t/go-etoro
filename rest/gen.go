// Package etororest is a Go client for the eToro public REST API, generated
// with oapi-codegen from the OpenAPI document eToro publishes.
package etororest

//go:generate bash ../scripts/fetch-spec.sh openapi.json
//go:generate go tool oapi-codegen --config=config.yaml ../api/openapi.json
