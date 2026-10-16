# Testing

Refines `~/.claude/rule-details/go/writing-tests.md`.

- **TST1** Tests are black-box and use only the exported API, as `ws/messages_test.go` does from package `etorows_test`.
- **TST2** Tests never call eToro. They decode literal frames and write through a fake transport.
- **TST3** Every spec patch or generator config change that alters wire behavior gets a test with a frame shaped like what eToro sends, for example rate values as JSON strings.
- **TST4** Table tests use `testCases` and `tc`, call `t.Parallel()` in the test and in every subtest, and assert errors with `require.ErrorIs` against a sentinel.
- **TST5** `make test` runs `go test -race` across every package and must pass.
