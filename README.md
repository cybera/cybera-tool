# Cybera Tool

## Requirements
Install Go through Homebrew if running from source instead of downloading the [latest release](https://github.com/cybera/racometer/releases):

```bash
brew install go
```

## Installation

```bash
git clone https://github.com/cybera/cybera-tool
cd cybera-tool
go install
```

## Usage

Login:

```bash
cybera-tool -c <username>:<password> | tail -1 | pbcopy
export CYBERA_KEY=$(pbpaste)
```

Log:

```bash
cybera-tool log DevOps "I've spent today on cybera-tool"
```

Also:

```bash
cybera-tool help
```
