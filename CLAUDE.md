# gweather - Go Weather App Guidelines

## Build Commands

- Build: `go build` or `make`
- Install: `make install` or `go install`
- Clean: `make clean`
- Uninstall: `make uninstall`
- Run: `./gweather [city]` or with metric units: `./gweather -m [city]`

## Testing

- Run all tests: `go test ./...`
- Run tests with verbose output: `go test -v ./...`
- Run single test: `go test -run TestName`
- Test coverage: `go test -cover ./...`
- Generate coverage report: `go test -coverprofile=coverage.out && go tool cover -html=coverage.out`

## Test Patterns

- Use httptest.NewServer for API mocking
- Use function variables instead of constants for easier mocking
- Capture command output with cobra.Command.SetOut and bytes.Buffer
- Restore original values with defer statements

## Code Style Guidelines

- Formatting: Use `gofmt` or `go fmt ./...`
- Linting: `golint ./...` and `go vet ./...`
- Imports: Standard grouping (stdlib first, then external packages)
- Naming: CamelCase for exported identifiers, lowercase for unexported
- Error handling: Explicit checks, wrap errors with fmt.Errorf and %w verb
- Comments: Document exported functions with full sentences
- Variables: Use variables for values that need to be modified in tests
- Environment variables: Load from .env file using godotenv
