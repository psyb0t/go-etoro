# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.5.0, 2026-10-29

- Regenerated from eToro spec version v1.402.0.
- REST: the social and feed endpoints are removed. Migration: stop calling these, the spec no longer has them.
- REST: removed the posts methods on `Client` and `ClientWithResponses`, including `PostPosts`, `GetPostsByPostId`, `PutPostsByPostId`, `DeletePostsByPostId`, `GetPostsByPostIdComments`, `PostPostsByPostIdComments`, and the like, follow, pin, save, share, poll vote and reply methods, each with its `WithBody` and `WithResponse` variants.
- REST: removed the feed methods `GetFeedsFollowing`, `GetFeedsForYou`, `GetFeedsNews`, `GetFeedsSaved`, `GetFeedsMarketsByMarketId`, `GetFeedsUsersByUserId` and `GetFeedsUsersByUserIdPinned`, and `PostAttachmentsWithBody`.
- REST: removed the same methods from `ClientInterface` and `ClientWithResponsesInterface`, so code that implements or mocks these interfaces must drop them.
- REST: removed the types `Attachment`, `Attachments`, `Discussion`, `DiscussionCreateRequest`, `DiscussionUpdateRequest`, `DiscussionsResponse`, `Emotion`, `Market`, `SharesResponse`, `User`, `TradingDemoApiPosition` and `TradingRealAdminApiOrder`, with their enum types and constants such as `MarketAssetType*`, `UserRoles*` and `EmotionType*`.
- REST: added `Order` and `TradingRealApiPosition`.
- WebSocket: `ws/etorows.gen.go` changes by one line, the spec version string.

## v0.4.2, 2026-10-28

- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.4.1, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.

## v0.4.0, 2026-10-23

- Regenerated from eToro spec version v1.400.0.
- WebSocket: `ConnectionAddressTemplate` is removed. Code that uses it must switch to the new per-operation constants.
- WebSocket: added `ConnectionAuthenticateAddressTemplate`, `ConnectionSubscribeAddressTemplate` and `ConnectionUnsubscribeAddressTemplate`. Each has the value "/ws".
- The WebSocket source header names v1.400.0. `OperationResponse` and its `OperationResponseOperation` values moved within the generated file, and the message name constants `MessageAuthenticateRequest` and `MessageOperationResponse` changed order. Their values are the same.

## v0.3.0, 2026-10-22

- Regenerated from eToro spec version v1.399.0.
- REST: `TradingRealApiPosition` is removed. The `positions` fields now use the new `TradingDemoApiPosition` type, which has the same fields, so code that names the old type must switch.
- Generated doc comments in the REST and WebSocket clients now spell the brand "etoro" instead of "eToro". The WebSocket source header names v1.399.0.

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
