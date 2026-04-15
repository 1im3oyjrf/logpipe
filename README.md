# logpipe

A lightweight CLI tool for tailing and filtering structured JSON logs from multiple sources with live grep support.

---

## Installation

```bash
go install github.com/yourname/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logpipe.git
cd logpipe
go build -o logpipe .
```

---

## Usage

Tail a log file and filter by a JSON field value:

```bash
logpipe tail app.log --filter "level=error"
```

Pipe logs from multiple sources and grep in real time:

```bash
logpipe tail app.log worker.log --grep "timeout"
```

Pretty-print raw JSON log output:

```bash
cat app.log | logpipe fmt
```

### Flags

| Flag | Description |
|------------|--------------------------------------|
| `--filter` | Filter by JSON key=value pair |
| `--grep` | Live grep across streamed log output |
| `--pretty` | Pretty-print JSON output |
| `--follow` | Continuously tail the file (default) |

---

## Requirements

- Go 1.21+

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)