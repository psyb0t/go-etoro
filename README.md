# go-etoro

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-etoro.svg)](https://pkg.go.dev/github.com/psyb0t/go-etoro)
[![CI](https://github.com/psyb0t/go-etoro/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-etoro/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-etoro/badges/version.svg)](https://github.com/psyb0t/go-etoro/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-etoro/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-etoro/badges/importers.svg)](https://github.com/psyb0t/go-etoro/blob/badges/importers.md)

Unofficial Go clients for the [eToro public API](https://api-portal.etoro.com), generated from the specs eToro publishes.

| Import path | Package | API | Generator |
|---|---|---|---|
| `github.com/psyb0t/go-etoro` (repo root) | `etoro` | REST, `https://public-api.etoro.com` | [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) |
| `github.com/psyb0t/go-etoro/ws` | `etorows` | WebSocket, `wss://ws.etoro.com/ws` | [gen-asyncapi](https://github.com/psyb0t/gen-asyncapi) |

`client.gen.go` in the root and everything in `ws/` except `gen.go` and `config.yaml` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-etoro
```

## REST

Every endpoint in the spec has a typed method on `ClientWithResponses`. eToro authenticates with two headers, `x-api-key` and `x-user-key`, taken from the API key section of your account settings.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/psyb0t/go-etoro"
)

func main() {
	auth := etoro.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("x-api-key", os.Getenv("ETORO_API_KEY"))
		req.Header.Set("x-user-key", os.Getenv("ETORO_USER_KEY"))

		return nil
	})

	client, err := etoro.NewClientWithResponses("https://public-api.etoro.com", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.GetBalancesWithResponse(
		context.Background(),
		&etoro.GetBalancesParams{XRequestId: uuid.New()},
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

## WebSocket

`etorows` is a typed client for the WebSocket API. It sends the authenticate, subscribe and unsubscribe requests, and decodes the instrument rate, candle and market status pushes, the order and position updates, and the operation responses. Every message is checked against eToro's schema both ways.

```go
package main

import (
	"context"
	"crypto/rand"
	"log/slog"
	"os"

	"github.com/psyb0t/gen-asyncapi/asyncapirt/gorillaws"
	"github.com/psyb0t/go-etoro/ws"
)

type handler struct {
	etorows.UnimplementedHandler
}

func (handler) HandleInstrumentRateMessage(
	ctx context.Context,
	envelope etorows.InstrumentRateEnvelope,
) error {
	for _, item := range envelope.Messages {
		rate, err := item.ContentValue()
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "rate", "topic", item.Topic, "ask", *rate.Ask)
	}

	return nil
}

func main() {
	ctx := context.Background()

	conn, err := gorillaws.Dial(ctx, etorows.ServerEtoroURL)
	if err != nil {
		panic(err)
	}

	client := etorows.NewClient(conn)
	defer client.Close()

	err = client.Authenticate(ctx, etorows.AuthenticateEnvelope{
		ID:        rand.Text(),
		Operation: "Authenticate",
		Data: etorows.AuthenticateData{
			APIKey:  os.Getenv("ETORO_API_KEY"),
			UserKey: os.Getenv("ETORO_USER_KEY"),
		},
	})
	if err != nil {
		panic(err)
	}

	err = client.Subscribe(ctx, etorows.SubscribeEnvelope{
		ID:        rand.Text(),
		Operation: "Subscribe",
		Data:      etorows.SubscribeData{Topics: []string{"instrument:1001"}},
	})
	if err != nil {
		panic(err)
	}

	if err := client.Run(ctx, handler{}); err != nil {
		panic(err)
	}
}
```

Push frames are envelopes whose items carry a `type` and a JSON-encoded `content` string; `ContentValue()` decodes and checks it. Rate values arrive as strings on the wire (`"1.25"`), so the `InstrumentRate` fields are strings too.

Pushes are told apart by their `type`; both order update types go to `HandleOrderUpdateMessage`. An envelope with a type the spec doesn't list goes to `HandleUnknownMessage`; the embedded `UnimplementedHandler` logs a warning and keeps reading.

## Regenerating

eToro's specs are not committed. `make generate` downloads the current OpenAPI and AsyncAPI documents into the gitignored `api/` directory and regenerates both packages:

```sh
make generate
```

Each package's `gen.go` declares its own steps, so `make generate` is just `go generate ./...`:

- `gen.go` in the root downloads `openapi.json` and runs oapi-codegen with `config.yaml`.
- `ws/gen.go` downloads `asyncapi.json`, runs `scripts/match-wire-types.sh`, and runs gen-asyncapi with `ws/config.yaml`. The match step turns the number fields of every schema whose description says its values arrive as JSON strings, such as `InstrumentRate`, into strings, so real frames decode.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Agent integrations

The [skill](.agents/skills/go-etoro) teaches a coding agent how to use go-etoro when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-etoro@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-etoro@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-etoro
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the specs and regenerate
```

## License

[MIT](LICENSE). The generated code comes from eToro's published API specifications. The specifications themselves are not redistributed here.
