# addr_sh

A JSON API service for HTTP and networking utilities. Returns your IP address, request headers, reverse DNS lookups, and IPv4 CIDR calculations.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Returns your remote IP and a summary of available tools |
| GET | `/about` | Service info |
| GET | `/ip` | Returns your IPv4 address |
| GET/POST | `/headers` | Returns all request headers |
| GET | `/hostnames` | Reverse DNS lookup for your remote address |
| GET | `/hostnames/{address}` | Reverse DNS lookup for a specific IP |
| GET | `/cidr/v4/{cidr}` | Number of total and usable addresses for a CIDR prefix |
| GET | `/cidr/v4/{address}/{cidr}` | Full CIDR details: network, broadcast, first/last usable address |
| GET | `/cidr/v4/{network}/{cidr}/split/{count}` | Split a CIDR into N equal subnets |

All responses are JSON.

### Examples

```sh
# Your IP address
curl https://addr.sh/ip

# Request headers
curl https://addr.sh/headers

# Reverse DNS for an IP
curl https://addr.sh/hostnames/8.8.8.8

# Addresses in a /24
curl https://addr.sh/cidr/v4/24

# CIDR details for 192.168.1.50/24
curl https://addr.sh/cidr/v4/192.168.1.50/24

# Split 10.0.0.0/24 into 4 subnets
curl https://addr.sh/cidr/v4/10.0.0.0/24/split/4
```

## Configuration

Configuration is loaded from a file (YAML/TOML/JSON) and/or environment variables. Pass a config file path via the `CONFIG_FILE` environment variable or `--config-file` flag.

| Key | Default | Description |
|-----|---------|-------------|
| `LogLevel` | `info` | Log level |
| `ListenPort` | `:2000` | HTTP listen address |
| `TLSListenPort` | `:4443` | HTTPS listen address |
| `EnableTLS` | `false` | Enable TLS |
| `TLSCertFile` | `fullchain.pem` | Path to TLS certificate |
| `TLSKeyFile` | `privkey.pem` | Path to TLS private key |

Environment variables use the same key names (e.g., `LISTENPORT=:8080`).

## Building

```sh
# Build for current platform
make build

# Build for OpenBSD amd64
make build-openbsd

# Run tests
make test

# Run all static analysis (staticcheck, govulncheck, gosec, osv-scanner)
make check

# Install analysis tools
make tools
```

## Running

```sh
./addr_sh
# or with a config file
CONFIG_FILE=config.yaml ./addr_sh
```

The server starts on `:2000` by default. Prometheus metrics are available at the standard `/metrics` endpoint (via the default `DefaultServeMux` — note: the app uses a custom mux, so metrics require a separate scrape setup if needed).

## Notes

- Behind a reverse proxy, `X-Forwarded-For` is respected when the direct connection is from localhost.
- Dependencies are vendored under `vendor/`.
