// Package etorows is a Go client for the eToro public WebSocket API, generated
// with gen-asyncapi from the AsyncAPI document eToro publishes.
package etorows

//go:generate bash ../scripts/fetch-spec.sh asyncapi.json
//go:generate bash ../scripts/match-wire-types.sh
//go:generate go tool gen-asyncapi --config=config.yaml ../api/asyncapi.wire.json
