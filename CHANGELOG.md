# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.2.0, 2026-10-21

- Regenerated from eToro spec version v1.398.0.
- REST: `TradingDemoApiPosition` is removed. The `positions` fields now use the new `TradingRealApiPosition` type, so code that names the old type must switch.
- REST: `EquityDetailsData` has new fields `Invested`, `OpenOrdersValue`, `OpenPositionsCount`, `PositionsValue` and `RealizedPNL`.
- WebSocket: `InstrumentCandlesAddress` and `InstrumentRatesAddress` are removed.
- WebSocket: `InstrumentCandlesAddressTemplate`, `InstrumentRatesAddressTemplate` and `PrivatePortfolioAddressTemplate` now have the value "/ws" instead of "candle:{instrumentId}", "instrument:{instrumentId}" and "private".
- The pinned Go toolchain moves to 1.26.10 in `go.mod` and `Dockerfile.dev`, which fixes the govulncheck finding GO-2026-6809 in `crypto/x509`. The agent skill setup notes name the same version.

## v0.1.5, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.4, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.3, 2026-10-16

- Added the `go-etoro` agent skill. The README's new Agent integrations section has the install commands for Claude Code, Codex and ClawHub, and the README now says the clients are unofficial.

## v0.1.2, 2026-10-16

- The WebSocket client is generated with gen-asyncapi v0.2.3.

## v0.1.1, 2026-10-16

- The REST client is now package `etoro` at the module root (`github.com/psyb0t/go-etoro`). The WebSocket client is package `etorows` at `github.com/psyb0t/go-etoro/ws`.

## v0.1.0, 2026-10-15

- Added `etororest`, a REST client generated with oapi-codegen from eToro's OpenAPI document (spec version v1.395.0).
- Added `etorows`, a WebSocket client generated with gen-asyncapi from eToro's AsyncAPI document: typed requests, a `Handler` per push and update type, strict decoding, and a gorilla/websocket transport through `gorillaws.Dial`. Rate values are strings, as eToro sends them.
- `make generate` downloads the current specs from the eToro API portal and regenerates both packages. The specs are not committed.
