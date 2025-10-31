# ya-runtime-salad

A demo runtime implementation for the [Golem Factory's Yagna platform](https://github.com/golemfactory/yagna) that simulates Salad nodes. This runtime exposes GPU-enabled nodes on the Golem Network.

## Prerequisites

- [Go](https://golang.org/dl/) 1.25.3 or later
- [Yagna](https://github.com/golemfactory/yagna) setup as a provider

It is recommended to follow the official [Golem provider installation guide](https://docs.golem.network/docs/providers/provider-installation) to set up your provider.

## Building

Build the runtime binary for Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o ya-runtime-salad ./cmd/ya-runtime-salad
```

Build the runtime binary for Windows:

```bash
GOOS=windows GOARCH=amd64 go build -o ya-runtime-salad.exe ./cmd/ya-runtime-salad
```

## Testing

Run the test suite:

```bash
go test ./...
```

## Installing

Build the runtime binary as described in the [Building](#building) section.

Find your Yagna provider's plugins directory. This is typically:

- **Windows**: `%APPDATA%\GolemFactory\yagna\plugins\`
- **Linux**: `~/.local/lib/yagna/plugins/`

Copy the built binary and configuration file to the plugins directory:

**Linux:**

```bash
# Create the plugins directory if it doesn't exist
mkdir -p ~/.local/lib/yagna/plugins

# Copy the binary
cp ya-runtime-salad ~/.local/lib/yagna/plugins/

# Copy the configuration
cp conf/ya-runtime-salad.json ~/.local/lib/yagna/plugins/
```

**Windows (PowerShell):**

```powershell
# Create the plugins directory if it doesn't exist
New-Item -ItemType Directory -Force -Path "$Env:APPDATA\GolemFactory\yagna\plugins"

# Copy the binary
Copy-Item "ya-runtime-salad.exe" "$Env:APPDATA\GolemFactory\yagna\plugins\"

# Copy the configuration
Copy-Item "conf\ya-runtime-salad.json" "$Env:APPDATA\GolemFactory\yagna\plugins\"
```

Restart your Yagna provider and verify the runtime is detected by checking the service logs:

```bash
yagna service run
```

## Configuring

Create an offer template configuration file:

- **Windows**: `%APPDATA%\ya-runtime-salad\template.json`
- **Linux**: `~/.config/ya-runtime-salad/template.json`

See the `examples` directory for examples.

Create a provider preset to make the runtime available for requestors:

```json
{
  "ver": "V1",
  "active": [
    "salad"
  ],
  "presets": [
    {
      "name": "salad",
      "exeunit-name": "salad",
      "pricing-model": "linear",
      "initial-price": 0.0,
      "usage-coeffs": {
        "golem.usage.duration_sec": 1.0e-6
      }
    }
  ]
}
```

## Links

- [Golem Network](https://golem.network/)
- [Yagna docs](https://docs.golem.network/)
- [Salad](https://salad.com/)
