# portwatch

> CLI tool to monitor and alert on open ports and service changes on a host

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start monitoring all open ports on the host and alert on any changes:

```bash
portwatch watch
```

Scan once and print current open ports:

```bash
portwatch scan
```

Watch specific port ranges and send an alert on change:

```bash
portwatch watch --ports 1-1024 --interval 30s --alert webhook --webhook-url https://hooks.example.com/notify
```

### Example Output

```
[2024-01-15 08:32:01] INFO  Watching ports 1-65535 (interval: 60s)
[2024-01-15 08:33:01] INFO  No changes detected
[2024-01-15 08:34:01] ALERT Port 8080 opened — process: python3 (pid 19842)
[2024-01-15 08:35:01] ALERT Port 3306 closed
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--ports` | `1-65535` | Port range to monitor |
| `--interval` | `60s` | Scan interval |
| `--alert` | `stdout` | Alert method (`stdout`, `webhook`, `email`) |
| `--config` | `~/.portwatch.yaml` | Path to config file |

## Configuration

```yaml
ports: "1-65535"
interval: 60s
alert: webhook
webhook_url: https://hooks.example.com/notify
```

## License

MIT © [yourusername](https://github.com/yourusername)