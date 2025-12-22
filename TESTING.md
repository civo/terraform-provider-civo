# Testing Guide

## Docker-Based Testing

This project includes a Docker-based testing environment for reproducible test runs in an isolated environment.

### Quick Start

```bash
# Build the test image
docker build -f Dockerfile.test -t terraform-provider-civo:test .

# Run unit tests
docker run --rm terraform-provider-civo:test

# Run build verification
docker run --rm terraform-provider-civo:test make build

# Run format check
docker run --rm terraform-provider-civo:test make fmtcheck
```

### Test Environment

- **Base Image**: golang:1.24-alpine
- **Dependencies**: git, make, bash
- **Go Modules**: Pre-downloaded for faster execution
- **Default Command**: `make test` (runs unit test suite)

### Benefits

- **Reproducible**: Same environment for all developers and CI/CD
- **Isolated**: No dependency on local Go installation or system packages
- **Fast**: Layer caching speeds up subsequent builds
- **Clean**: Tests run in a fresh environment every time

### Local Testing (Without Docker)

```bash
# Run unit tests
make test

# Run acceptance tests (requires CIVO_TOKEN)
TF_ACC=1 make testacc

# Build the provider
make build

# Format code
make fmt

# Check formatting
make fmtcheck
```

### CI/CD Integration

The Docker test environment can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
test:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v2
    - name: Build test image
      run: docker build -f Dockerfile.test -t terraform-provider-civo:test .
    - name: Run tests
      run: docker run --rm terraform-provider-civo:test
```

## Test Coverage

The test suite covers:

- Provider configuration and initialization
- All resource CRUD operations (Create, Read, Update, Delete)
- Data source queries
- Input validation and error handling
- SSH key generation and validation

## Requirements

- Docker (for containerized testing)
- Go 1.24+ (for local testing)
- CIVO API token (for acceptance tests only)

## Contributing

When adding new features:

1. Add unit tests for new functionality
2. Run the full test suite: `make test`
3. Verify formatting: `make fmtcheck`
4. Build verification: `make build`

For acceptance tests that interact with the Civo API, follow existing patterns in the `civo/acceptance` directory.
