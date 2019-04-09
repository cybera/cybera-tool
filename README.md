# Cybera Tool

## Requirements
Install Go through Homebrew.

```bash
brew install go
```

## Installation

```bash
go get github.com/cybera/cybera-tool
```

## Usage

Login:

```bash
export CYBERA_KEY=$(cybera-tool authenticate -c <username:password> | tail -1)
```

Alternatively you can leverage ENV variables:
``` bash
CYBERA_CREDS="<username>:<password>"
export CYBERA_KEY=$(cybera-tool authenticate | tail -1)
```

Log:

```bash
cybera-tool log DevOps "I've spent today working on cybera-tool"
```

This command will fail if it has been longer than 15 minutes since the
authentication command was run. (Session expiration)

Also:

```bash
cybera-tool help
```
