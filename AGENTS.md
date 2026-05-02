# Repository Guidelines

## Project Structure & Module Organization
This repository is a Go library (`github.com/akfaiz/go-mailgen`) for building HTML/plain-text emails.
- Core package files live at repo root: `builder.go`, `component.go`, `message.go`.
- Unit tests are colocated with source as `*_test.go` (for example `builder_test.go`).
- Built-in themes are in `templates/default` and `templates/plain`; template registration is in `templates/templates.go`.
- Runnable examples are under `examples/`, including integration usage in `examples/integration/main.go` and generated outputs in `examples/default` and `examples/plain`.
- CI is defined in `.github/workflows/ci.yml`.

## Build, Test, and Development Commands
Use `make` targets for consistent local workflows:
- `make install`: download and tidy modules (`go mod download && go mod tidy`).
- `make fmt`: format Go code with `go fmt ./...`.
- `make lint`: run `golangci-lint` with autofix.
- `make test`: run verbose tests with race detector and coverage profile.
- `make coverage`: generate `coverage.html` and `coverage.txt` from `coverage.out`.
- `make coverage-check`: enforce minimum total coverage (currently `90%`).
- `make update-examples`: regenerate example artifacts from `examples/main.go`.

## Coding Style & Naming Conventions
- Follow `.editorconfig`: tabs for Go/Makefile, 2 spaces for YAML/JSON/HTML/CSS, LF endings.
- Keep public APIs idiomatic Go: exported names in `PascalCase`, internal helpers in `camelCase`.
- Run `make fmt` and `make lint` before pushing; CI runs lint and tests across Go `1.23`, `1.24`, and `1.25`.

## Testing Guidelines
- Write table-driven tests where practical and keep tests beside implementation files.
- Name test functions with `TestXxx` and benchmarks with `BenchmarkXxx`.
- Run `make test` locally; use `make coverage-check` before opening PRs to avoid coverage regressions.

## Commit & Pull Request Guidelines
- Commit messages must follow Conventional Commits (enforced by `lefthook`):
  - Example: `feat: add custom logo fallback`
  - Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.
- Install hooks (`lefthook install`) so pre-commit runs format/tidy/lint checks.
- PRs should include: concise description, linked issue (if any), test updates for behavior changes, and regenerated examples/screenshots when template output changes.
