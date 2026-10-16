# Go code

Refines `~/.claude/rule-details/go/wrapping-errors.md` and `~/.claude/rule-details/coding-style.md`.

- **GO1** `client.gen.go` and everything in `ws/` except `gen.go` and `config.yaml` is generated. Never edit it by hand. Change the generator config (`config.yaml`, `ws/config.yaml`) or the spec patch in `scripts/`, then run `make generate` and commit the result.
- **GO2** Spec fixes live in `scripts/`, applied to the downloaded spec before generation, as `scripts/match-wire-types.sh` does for fields eToro sends as JSON strings. Each patch says in its header which spec problem it fixes and matches on the spec's own content, not on a hardcoded list of names where a rule can be stated.
- **GO3** The downloaded specs in `api/` are gitignored and never committed.
- **GO4** Hand-written Go (the `gen.go` files and tests) wraps errors with `ctxerrors.Wrap` against `ctxerrors/commerr` sentinels. No `fmt.Errorf`, no `errors.New` outside a sentinel declaration.
- **GO5** The exported API is whatever the generators emit. A config change that renames or removes an exported name is a breaking release (see `releases.md`).
- **GO6** `make lint` passes with zero issues. Hand-written lines are at most 80 columns, the `.golangci.yml` limit.
