# branch

A CLI tool for creating git branches with consistent naming patterns. Never struggle with branch naming conventions again!

## Features

- 🎯 **Consistent naming**: Automatically formats branches as `<type>/<ticket>-<description>`
- 🎫 **Smart ticket detection**: Recognizes common ticket formats (GitHub issues, Jira, Linear, etc.)
- ⚙️ **Fully configurable**: Customize both branch commands and ticket patterns to match your workflow
- 🚀 **Fast**: Create and switch to branches in one command

## Installation

### From Source

```bash
git clone https://github.com/owenrumney/branch.git
cd branch
go build -o branch
sudo mv branch /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/owenrumney/branch@latest
```

Make sure `$GOPATH/bin` or `$GOBIN` is in your `PATH`.

## Usage

### Basic Commands

By default, the tool provides several subcommands for different branch types:

- `branch feat` - Create a feature branch
- `branch fix` - Create a bugfix branch
- `branch tests` - Create a tests branch
- `branch chore` - Create a chore branch
- `branch docs` - Create a documentation branch

> **Note**: These commands are configurable! You can customize which branch types are available by editing your config file (see [Configuration](#configuration) below).

### Examples

#### With a ticket number

```bash
# Linear/Jira style tickets
branch feat PIP-1234 implement new authentication
# Creates: feat/pip-1234-implement-new-authentication

branch fix INFRA-567 update database connection
# Creates: fix/infra-567-update-database-connection

# GitHub issues
branch docs #123 add api documentation
# Creates: docs/123-add-api-documentation
```

#### Without a ticket number

```bash
branch feat implement user dashboard
# Creates: feat/implement-user-dashboard

branch fix resolve memory leak
# Creates: fix/resolve-memory-leak

branch chore update dependencies
# Creates: chore/update-dependencies
```

#### Special characters

The tool automatically handles special characters, converting them appropriately:

```bash
branch feat add new feature!
# Creates: feat/add-new-feature

branch fix update user's profile
# Creates: fix/update-users-profile
```

## Branch Naming Format

Branches follow this pattern:
```
<type>/<ticket>-<description>
```

- **Type**: A branch type command (defaults: `feat`, `fix`, `tests`, `chore`, `docs` - configurable)
- **Ticket** (optional): Automatically detected if it matches a known pattern
- **Description**: The rest of your input, converted to a URL-friendly slug

## Configuration

### Default Ticket Patterns

The tool recognizes these ticket formats by default:

- GitHub issues: `#123`
- Jira/Linear style: `PIP-1234`, `INFRA-124`
- Underscore variant: `PIP_1234`

### Custom Configuration

You can customize both branch commands and ticket patterns by creating a config file at:
- `$XDG_CONFIG_HOME/branch/config.json` (if `XDG_CONFIG_HOME` is set)
- `~/.config/branch/config.json` (default)

#### Customizing Branch Commands

You can define your own branch type commands. For example, if you prefer `feature` instead of `feat`, or want to add custom types like `hotfix` or `release`:

```json
{
  "branch_commands": [
    "feature",
    "bugfix",
    "hotfix",
    "release",
    "chore"
  ],
  "ticket_patterns": [
    "^#\\d+$",
    "^[A-Z]+-\\d+$"
  ]
}
```

After updating your config, the new commands will be available:
```bash
branch feature PIP-1234 add new functionality
# Creates: feature/pip-1234-add-new-functionality

branch hotfix fix critical bug
# Creates: hotfix/fix-critical-bug
```

#### Customizing Ticket Patterns

You can also customize which ticket patterns are recognized:

```json
{
  "ticket_patterns": [
    "^#\\d+$",
    "^[A-Z]+-\\d+$",
    "^CUSTOM-\\d+$"
  ]
}
```

The patterns are regular expressions. If you don't specify any patterns, the defaults will be used.

#### Complete Example

Here's a complete configuration example:

```json
{
  "branch_commands": [
    "feat",
    "fix",
    "tests",
    "chore",
    "docs",
    "hotfix"
  ],
  "ticket_patterns": [
    "^#\\d+$",
    "^[A-Z]+-\\d+$",
    "^CUSTOM-\\d+$"
  ]
}
```

## Examples in Action

```bash
# Feature with Linear ticket
$ branch feat PIP-1234 add dark mode toggle
Created and switched to branch: feat/pip-1234-add-dark-mode-toggle

# Bugfix without ticket
$ branch fix resolve login timeout issue
Created and switched to branch: fix/resolve-login-timeout-issue

# Documentation with GitHub issue
$ branch docs #456 update installation guide
Created and switched to branch: docs/456-update-installation-guide

# Chore task
$ branch chore update npm packages
Created and switched to branch: chore/update-npm-packages
```

## Requirements

- Go 1.25.5 or later
- Git (must be in a git repository to create branches)

## Development

### Building

```bash
go build ./...
```

### Running Tests

```bash
go test ./...
```

### Running Specific Tests

```bash
go test -v ./... -run TestName
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
