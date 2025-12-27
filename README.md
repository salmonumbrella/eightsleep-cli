# 🛏️ Eight Sleep CLI — Smart bed control in your terminal.

Control your Eight Sleep Pod from the command line.

## Features

- **Adjustable base** - control base angle, presets, and vibration settings
- **Alarms** - create, update, delete, snooze, and dismiss alarms
- **Audio playback** - play tracks, manage favorites, control volume
- **Autopilot insights** - view autopilot details, history, and configuration
- **Household management** - view household summary, schedule, and users
- **Metrics & insights** - analyze sleep trends, intervals, and aggregated data
- **Power & temperature control** - turn pod on/off, set temperature levels
- **Schedule daemon** - run automated temperature schedules from config file
- **Secure authentication** - browser-based OAuth login with token caching
- **Sleep analytics** - view sleep metrics for individual days or date ranges
- **Temperature modes** - control nap mode, hot-flash mode, and view events
- **Temperature schedules** - manage cloud-based temperature schedules
- **Travel management** - manage trips, plans, and flight status

## Installation

### Homebrew

```bash
brew install salmonumbrella/tap/eightsleep
```

### Go Install

```bash
go install github.com/salmonumbrella/eightsleep-cli/cmd/eightsleep@latest
```

## Quick Start

### 1. Authenticate

```bash
eightsleep login
```

Opens your browser for secure OAuth authentication.

### 2. Check Pod Status

```bash
eightsleep status
```

### 3. Control Temperature

```bash
# Set temperature level (-100 to 100)
eightsleep temp 20

# Or use Fahrenheit/Celsius
eightsleep temp 68F
eightsleep temp 20C
```

## Configuration

### Priority

Configuration sources are applied in this order (highest priority first):

1. Command-line flags
2. Environment variables (prefixed with `EIGHTSLEEP_`)
3. Config file (`~/.config/eightsleep-cli/config.yaml`)

### Environment Variables

- `EIGHTSLEEP_EMAIL` - Eight Sleep account email
- `EIGHTSLEEP_PASSWORD` - Eight Sleep account password
- `EIGHTSLEEP_USER_ID` - Eight Sleep user ID (auto-resolved if not provided)
- `EIGHTSLEEP_CLIENT_ID` - OAuth client ID (optional, defaults to app client)
- `EIGHTSLEEP_CLIENT_SECRET` - OAuth client secret (optional)
- `EIGHTSLEEP_TIMEZONE` - IANA timezone or `local` (default: local)
- `EIGHTSLEEP_OUTPUT` - Output format: `table`, `json`, or `csv` (default: table)
- `EIGHTSLEEP_VERBOSE` - Enable verbose logging (true/false)

### Config File

Create `~/.config/eightsleep-cli/config.yaml`:

```yaml
email: "you@example.com"
password: "your-password"
# user_id: "optional"              # auto-resolved via /users/me
# timezone: "America/New_York"     # defaults to local
# output: "table"                  # table|json|csv
# verbose: false
```

Set restrictive permissions:

```bash
chmod 600 ~/.config/eightsleep-cli/config.yaml
```

## Security

### Credential Storage

- **Browser-based login** - uses OAuth flow via browser (recommended)
- **Token caching** - tokens are cached locally to reduce API calls
- **Config file permissions** - CLI warns if permissions are too permissive (>0600)
- **Auto-resolution** - user ID and device ID are automatically resolved

### Best Practices

- Use `eightsleep login` for secure browser-based authentication
- Set `chmod 600` on config files to protect credentials
- Avoid passing `--password` on command line (visible in shell history)
- Use environment variables or config file for credentials

### API Limitations

Eight Sleep does not publish a stable public API. This CLI uses undocumented endpoints:

- **Rate limiting** - the API enforces rate limits; repeated logins may return 429
- **Header mimicking** - client mimics Android app headers to reduce throttling
- **Token reuse** - tokens are cached and reused to minimize auth calls
- **No local control** - only HTTPS cloud API; no Bluetooth or local network

## Commands

### Authentication

```bash
eightsleep login                    # Authenticate via browser (recommended)
eightsleep logout                   # Clear cached authentication token
eightsleep whoami                   # Show configured user ID
```

### Power & Temperature

```bash
eightsleep on                       # Turn pod on
eightsleep off                      # Turn pod off
eightsleep status                   # Show device status
eightsleep temp 20                  # Set temperature level (-100 to 100)
eightsleep temp 68F                 # Set temperature in Fahrenheit
eightsleep temp 20C                 # Set temperature in Celsius
eightsleep presence                 # Check if user is in bed
```

### Sleep Analytics

```bash
eightsleep sleep day                          # Today's sleep metrics
eightsleep sleep day --date 2024-12-15        # Specific date
eightsleep sleep range --from 2024-12-01 --to 2024-12-15
```

### Alarms

```bash
eightsleep alarm list                         # List all alarms
eightsleep alarm create --time 07:00 --days mon,wed,fri
eightsleep alarm update <alarmId> --time 07:30
eightsleep alarm delete <alarmId>
eightsleep alarm snooze <alarmId>
eightsleep alarm dismiss <alarmId>
eightsleep alarm dismiss-all
eightsleep alarm vibration-test
```

### Temperature Schedules

```bash
eightsleep schedule list                      # List all schedules
eightsleep schedule next                      # Show next upcoming events
eightsleep schedule create --name "Evening" --time 21:00 --level 20
eightsleep schedule update <scheduleId> --level 30
eightsleep schedule delete <scheduleId>
```

### Temperature Modes

```bash
# Nap mode
eightsleep tempmode nap on
eightsleep tempmode nap off
eightsleep tempmode nap extend
eightsleep tempmode nap status

# Hot-flash mode
eightsleep tempmode hotflash on
eightsleep tempmode hotflash off
eightsleep tempmode hotflash status

# Temperature events
eightsleep tempmode events
```

### Audio

```bash
eightsleep audio tracks                       # List available tracks
eightsleep audio categories                   # List track categories
eightsleep audio state                        # Show current playback state
eightsleep audio play <trackId>
eightsleep audio pause
eightsleep audio next
eightsleep audio seek <position>
eightsleep audio volume <level>
eightsleep audio pair                         # Pair audio device
eightsleep audio favorites list
eightsleep audio favorites add <trackId>
eightsleep audio favorites remove <trackId>
```

### Adjustable Base

```bash
eightsleep base info                          # Base information
eightsleep base angle <degrees>               # Set base angle
eightsleep base presets                       # List available presets
eightsleep base preset-run <presetName>       # Run a preset
eightsleep base vibration-test                # Test vibration
```

### Device Information

```bash
eightsleep device info                        # Device information
eightsleep device peripherals                 # Connected peripherals
eightsleep device owner                       # Device owner
eightsleep device warranty                    # Warranty information
eightsleep device online                      # Online status
eightsleep device priming-tasks               # Priming tasks
eightsleep device priming-schedule            # Priming schedule
```

### Metrics & Insights

```bash
eightsleep metrics trends                     # Sleep trends
eightsleep metrics intervals                  # Sleep intervals
eightsleep metrics summary                    # Summary statistics
eightsleep metrics aggregate                  # Aggregated metrics
eightsleep metrics insights                   # Sleep insights
```

### Autopilot

```bash
eightsleep autopilot details                  # Autopilot details
eightsleep autopilot history                  # Autopilot history
eightsleep autopilot recap                    # Autopilot recap
eightsleep autopilot set-level-suggestions <level>
eightsleep autopilot set-snore-mitigation <enabled>
```

### Travel

```bash
eightsleep travel trips                       # List trips
eightsleep travel create-trip --destination "Tokyo" --start-date 2024-12-20
eightsleep travel delete-trip <tripId>
eightsleep travel plans                       # List travel plans
eightsleep travel create-plan --trip-id <id> --type jetlag
eightsleep travel update-plan <planId> --status active
eightsleep travel tasks                       # List travel tasks
eightsleep travel airport-search --query "JFK"
eightsleep travel flight-status --flight-number AA100
```

### Household

```bash
eightsleep household summary                  # Household summary
eightsleep household schedule                 # Household schedule
eightsleep household current-set              # Current settings
eightsleep household invitations              # Pending invitations
eightsleep household devices                  # Household devices
eightsleep household users                    # Household users
eightsleep household guests                   # Guest accounts
```

### Schedule Daemon

```bash
eightsleep daemon --dry-run                   # Preview without executing
eightsleep daemon --pid-file /tmp/eightsleep.pid
```

## Output Formats

### Table

Human-readable tables with aligned columns:

```bash
$ eightsleep alarm list
ID            STATUS    TIME     DAYS
alarm_123     ACTIVE    07:00    Mon, Wed, Fri
alarm_456     INACTIVE  08:30    Tue, Thu
```

### JSON

Machine-readable output for scripting:

```bash
$ eightsleep status --output json
{
  "device_id": "device_abc123",
  "online": true,
  "temperature": 20,
  "power": "on"
}
```

### CSV

Export data for spreadsheets:

```bash
$ eightsleep sleep range --from 2024-12-01 --to 2024-12-07 --output csv
date,sleep_score,duration_hours,deep_sleep_minutes,rem_minutes
2024-12-01,82,7.5,90,120
2024-12-02,85,8.0,95,130
```

Data goes to stdout, errors and progress to stderr for clean piping.

## Examples

### Morning Routine Automation

```bash
# Create morning alarm with gradual warming
eightsleep alarm create --time 06:30 --days mon,tue,wed,thu,fri
eightsleep schedule create --name "Morning Warmup" --time 06:00 --level 50
```

### Track Sleep Quality Over Time

```bash
# Export last 30 days to CSV
eightsleep sleep range \
  --from $(date -v-30d +%Y-%m-%d) \
  --to $(date +%Y-%m-%d) \
  --output csv > sleep_data.csv
```

### Pre-bed Cooling Schedule

```bash
# Create evening cooling schedule
eightsleep schedule create --name "Evening Cool" --time 21:00 --level -30
```

### Automation

Use `--output json` for scripting and pipeline integration:

```bash
# Get current temperature as JSON
eightsleep status --output json | jq '.temperature'

# List alarms and filter by status
eightsleep alarm list --output json | jq '.[] | select(.status == "ACTIVE")'
```

## Global Flags

All commands support these flags:

- `--config <path>` - Config file path (default: `~/.config/eightsleep-cli/config.yaml`)
- `--email <email>` - Eight Sleep account email
- `--password <password>` - Eight Sleep account password
- `--user-id <id>` - Eight Sleep user ID (auto-resolved if not provided)
- `--timezone <tz>` - IANA timezone or `local` (default: local)
- `--output <format>` - Output format: `table`, `json`, or `csv` (default: table)
- `--fields <fields>` - Comma-separated list of fields to display
- `--verbose`, `-v` - Enable verbose logging
- `--quiet` - Suppress config loading banner
- `--help`, `-h` - Show help for any command

## Shell Completions

Generate shell completions for your preferred shell:

### Bash

```bash
# macOS (Homebrew):
eightsleep completion bash > $(brew --prefix)/etc/bash_completion.d/eightsleep

# Linux:
eightsleep completion bash > /etc/bash_completion.d/eightsleep

# Or source directly:
source <(eightsleep completion bash)
```

### Zsh

```zsh
eightsleep completion zsh > "${fpath[1]}/_eightsleep"
```

### Fish

```fish
eightsleep completion fish > ~/.config/fish/completions/eightsleep.fish
```

### PowerShell

```powershell
eightsleep completion powershell | Out-String | Invoke-Expression
```

## Development

After cloning, install git hooks:

```bash
make setup
```

This installs [lefthook](https://github.com/evilmartians/lefthook) pre-commit and pre-push hooks for linting and testing.

## License

MIT

## Links

- [GitHub Repository](https://github.com/salmonumbrella/eightsleep-cli)

## Prior Work

This project builds on work from the Eight Sleep community:

- [clim8](https://github.com/blacktop/clim8) - Go CLI
- [8sleep-mcp](https://github.com/elizabethtrykin/8sleep-mcp) - MCP server (Node/TypeScript)
- [pyEight](https://github.com/mezz64/pyEight) - Python library
- [eight_sleep](https://github.com/lukas-clarke/eight_sleep) - Home Assistant integration
- [homebridge-eightsleep](https://github.com/nfarina/homebridge-eightsleep) - Homebridge plugin
