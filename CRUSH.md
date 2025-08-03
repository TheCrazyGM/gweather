# CRUSH.md – gweather Agent Usage & Code Guidelines

## Build & Run Commands
- Build: `go build` or `make`
- Install: `go install` or `make install`
- Clean: `make clean`
- Uninstall: `make uninstall`
- Run: `./gweather [city]` or `./gweather -m [city]`

## Test Commands
- Run all tests: `go test ./...`
- Verbose tests: `go test -v ./...`
- Run a single test: `go test -run TestName`
- Test coverage: `go test -cover ./...`
- Generate coverage report: `go test -coverprofile=coverage.out && go tool cover -html=coverage.out`

## Lint & Format
- Format: `gofmt` or `go fmt ./...`
- Lint: `golint ./...`
- Static analysis: `go vet ./...`

## Code Style Guidelines
- **Imports:** Stdlib, then blank line, then external; use goimports/go fmt ordering
- **Formatting:** Always `gofmt` before committing
- **Types & Naming:**
  - CamelCase for exported names
  - lowerCamel for unexported
  - Acronyms capitalized (e.g., `APIKey`)
- **Error handling:**
  - Always check errors, wrap with `fmt.Errorf("context: %w", err)`
  - Return errors, don't panic except for truly exceptional conditions
- **Documentation:**
  - Exported functions/types must use complete sentences
- **Tests:**
  - Use `httptest.NewServer` for HTTP mocking
  - Use variables (not constants) for values that are changed in tests
  - Restore changed global states with `defer`
- **Environment Variables:**
  - Use a `.env` file via godotenv (if present)
- **CLI Options:**
  - Use `-m/--metric` for Celsius, `-n/--no-emoji` to disable weather icons
  
## Project Notes
- Do **not** commit `.crush/` directory
- See README.md and CLAUDE.md for further detailed app usage and contribution patterns
