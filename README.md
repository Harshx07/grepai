# grepai

A privacy-first semantic code search CLI tool. Unlike `grep` (exact text matching), `grepai` indexes the **meaning** of your code using vector embeddings, enabling natural language searches.

## Features

- **Semantic Search**: Search code by intent, not just keywords
- **Real-time Indexing**: Background daemon keeps your index always up-to-date
- **Privacy-first**: Supports local embeddings with Ollama (no data leaves your machine)
- **AI Agent Ready**: Built-in integration for Cursor and Claude Code
- **Cross-platform**: Single binary for macOS, Linux, and Windows

## Installation

```bash
go install github.com/yoanbernabeu/grepai/cmd/grepai@latest
```

Or download from [Releases](https://github.com/yoanbernabeu/grepai/releases).

## Quick Start

```bash
# Initialize grepai in your project
cd your-project
grepai init

# Start the indexing daemon
grepai watch

# Search your codebase semantically
grepai search "function that handles authentication errors"
```

## Commands

| Command | Description |
|---------|-------------|
| `grepai init` | Initialize grepai in current directory |
| `grepai watch` | Start real-time file watcher daemon |
| `grepai search <query>` | Search codebase with natural language |
| `grepai agent-setup` | Configure AI agents (Cursor, Claude Code) |

## Configuration

Configuration is stored in `.grepai/config.yaml`:

```yaml
version: 1
embedder:
  provider: ollama          # ollama | openai
  model: nomic-embed-text
  endpoint: http://localhost:11434
store:
  backend: gob              # gob | postgres
chunking:
  size: 512
  overlap: 50
```

## Embedding Providers

### Ollama (Default - Local)
```bash
# Install Ollama: https://ollama.ai
ollama pull nomic-embed-text
```

### OpenAI
Set your API key:
```bash
export OPENAI_API_KEY=sk-...
```

## Storage Backends

### GOB (Default)
In-memory index with file persistence. Best for individual projects.

### PostgreSQL with pgvector
For large monorepos or shared indexes:
```bash
docker compose up -d
```

## Requirements

- Go 1.22+ (for building from source)
- Ollama (for local embeddings) or OpenAI API key

## License

MIT License - Yoan Bernabeu 2026
