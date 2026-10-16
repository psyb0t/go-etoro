# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-15

- Added `etororest`, a REST client generated with oapi-codegen from eToro's OpenAPI document (spec version v1.395.0).
- Added `etorows`, a WebSocket client generated with gen-asyncapi from eToro's AsyncAPI document: typed requests, a `Handler` per push and update type, strict decoding, and a gorilla/websocket transport through `gorillaws.Dial`. Rate values are strings, as eToro sends them.
- `make generate` downloads the current specs from the eToro API portal and regenerates both packages. The specs are not committed.
