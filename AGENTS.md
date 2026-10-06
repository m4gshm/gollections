# AGENTS.md — Gollections

Go collections library (slices, maps, iterators) for Go 1.27+.

## Quick Start

```bash
go test ./...          # all tests
go test ./slice/...    # slice package
```

## Structure

```
���── slice/      # slice functions
���── maps/       # map functions
���── seq/        # iter.Seq, Seq2, SeqE
���── internal/   # examples, test fixtures
```

Each package is a separate `go.mod` module. Public API is only the root `api.go`.

## Rules

- Add new functions only to the corresponding package.
- No external dependencies — only stdlib + `golang.org/exp/slices/maps`.
- Test coverage ≥ 80%.
- Naming: `CamelCase` for exports, `camelCase` for internal.
- Document every export (godoc comment).

## Anti-patterns

- Avoid chained `Filter(Convert(Filter(...)))` — allocates intermediate slices.
- Use `seq` for lazy pipelines instead of `slice`.