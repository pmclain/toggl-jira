# Toggl-Jira Time Sync

A Go tool to sync Toggl time entries to Jira work logs. This version is distributed as a single binary with no dependencies.

## Features

- Single binary distribution - no runtime dependencies
- Cross-platform support (Linux, macOS, Intel, ARM)
- Project filtering with `JIRA_PROJECTS` environment variable
- Automatic worklog updates for existing entries
- Detailed logging
- Configurable time rounding

## Installation

### Download Pre-built Binary

Download the appropriate binary for your system from the [releases page](../../releases):

- `toggl-jira-darwin-amd64` for macOS (Intel)
- `toggl-jira-darwin-arm64` for macOS (Apple Silicon)
- `toggl-jira-linux-amd64` for Linux (x86_64)
- `toggl-jira-linux-arm64` for Linux (ARM64)

Make it executable:
```bash
chmod +x toggl-jira-*
```

### Build from Source

Requires Go 1.21 or later:

```bash
git clone https://github.com/pmclain/toggl-jira
cd toggl-jira
go build
```

## Configuration

Create a `.env` file with your credentials:
```bash
cp .env.sample .env
```

Edit the `.env` file:
```
TOGGL_API_TOKEN=your_toggl_token
JIRA_EMAIL=your_email@company.com
JIRA_API_TOKEN=your_jira_token
JIRA_HOST=your.jira.host.com
JIRA_PROJECTS=PROJ1,PROJ2  # Optional: comma-separated list of project keys
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `TOGGL_API_TOKEN` | Yes | Your Toggl API token |
| `JIRA_EMAIL` | Yes | Your Jira account email |
| `JIRA_API_TOKEN` | Yes | Your Jira API token |
| `JIRA_HOST` | Yes | Your Jira instance hostname |
| `JIRA_PROJECTS` | No | Comma-separated list of Jira project keys to sync |

## Usage

Simply run the binary:
```bash
./toggl-jira
```

The tool will:
1. Fetch Toggl time entries from the last 48 hours
2. Parse Jira ticket numbers from the Toggl entry descriptions
3. Create or update work logs in Jira for each entry

### Toggl Entry Format

Your Toggl time entries should include the Jira issue key at the start of the description:

```
PROJ-123 Working on feature
PROJ-456 Bug fix
```

If `JIRA_PROJECTS` is set, only entries matching those project keys will be synced.

### Automation

For Linux/macOS users, you can set up a cron job to run the sync automatically:

```bash
# Run every hour
0 * * * * /path/to/toggl-jira
```

## Development

### Running Tests

```bash
go test -v ./...
```

### Building for All Platforms

```bash
# Build for all supported platforms
mkdir -p dist
GOOS=darwin GOARCH=amd64 go build -o dist/toggl-jira-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o dist/toggl-jira-darwin-arm64
GOOS=linux GOARCH=amd64 go build -o dist/toggl-jira-linux-amd64
GOOS=linux GOARCH=arm64 go build -o dist/toggl-jira-linux-arm64
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 