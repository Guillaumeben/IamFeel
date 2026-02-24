# IamFeel

AI-powered training assistant for personalized sport session planning. Built with Go and Claude.

## What is this?

IamFeel helps you plan training sessions based on:
- Your schedule and availability
- Training history (6 months tracking)
- Current training phase (strength, technique, conditioning, etc.)
- Personal goals (short, medium, long-term)
- Club memberships and available sessions
- Nutrition and supplement tracking

**Sport-agnostic design:** While the initial version focuses on boxing, the architecture allows anyone to adapt it for their sport (ironman, golf, climbing, etc.).

## Status

🚧 **Early Development** - Project structure established, implementation in progress.

See [docs/PLAN.md](docs/PLAN.md) for detailed roadmap and progress.

## Features (Planned)

### MVP
- ✅ Interactive onboarding wizard
- 📋 Weekly training plan generation
- 🖥️ Simple web dashboard
- 💬 Chat interface for adjustments
- 📊 Training history tracking (6 months)
- 🍎 Nutrition and supplement logging

### Future
- 📅 Google Calendar integration
- 📈 Advanced analytics and progress tracking
- 🏆 Goal achievement tracking
- 📱 Mobile-friendly PWA
- 🔌 Wearable device integration

See [docs/BACKLOG.md](docs/BACKLOG.md) for all future ideas.

## Tech Stack

- **Language:** Go
- **Database:** SQLite (pure Go, no CGo)
- **LLM:** Anthropic Claude API (configurable model)
- **Web:** Standard library + chi router
- **Templates:** html/template
- **Config:** YAML

## Project Structure

```
iamfeel/
├── cmd/
│   ├── server/          # Web server
│   └── cli/             # CLI commands (onboard, plan, etc.)
├── internal/
│   ├── agent/           # Claude AI integration
│   ├── db/              # Database models & queries
│   ├── config/          # Configuration loading
│   └── api/             # HTTP handlers
├── web/
│   ├── templates/       # HTML templates
│   └── static/          # CSS/JS assets
├── configs/
│   └── sports/          # Sport-specific configs (boxing.yaml, etc.)
├── data/                # Runtime data (gitignored)
│   ├── coach.db         # SQLite database
│   └── user_config.yaml # User profile
└── docs/
    ├── PLAN.md          # Development roadmap
    └── BACKLOG.md       # Future ideas
```

## Getting Started

### Prerequisites

- Go 1.21+
- Anthropic API key

### Installation

```bash
# Clone the repository
git clone https://github.com/tuxnam/iamfeel.git
cd iamfeel

# Install dependencies
go mod download

# Set up your API key
export ANTHROPIC_API_KEY="your-key-here"
```

### Quick Start

```bash
# Run onboarding (creates your profile)
go run cmd/cli/main.go onboard

# Generate weekly plan
go run cmd/cli/main.go plan

# Start web server
go run cmd/server/main.go

# Open browser to http://localhost:8080
```

## Configuration

### Sport Configs

Sport configurations define:
- Training phases (strength, technique, conditioning, etc.)
- Session types (bag work, sparring, roadwork, etc.)
- Periodization guidelines
- Sport-specific considerations

Example: `configs/sports/boxing.yaml`

### User Profile

Your profile is stored in `data/user_config.yaml` and includes:
- Sports you practice
- Available equipment/gym access
- Club sessions and memberships
- Fitness level and experience
- Goals and preferences
- Weekly availability

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build CLI
go build -o bin/coach cmd/cli/main.go

# Build server
go build -o bin/server cmd/server/main.go
```

### Roadmap

Follow progress in [docs/PLAN.md](docs/PLAN.md). Current phase: **Phase 1 - Foundation**

## Design Philosophy

- **Agent-native:** File-based state, transparent operation
- **Local-first:** Runs on your machine, your data stays local
- **Simple:** No over-engineering, clear code
- **Sport-agnostic:** Easy to adapt for different sports
- **Privacy-focused:** Configurable LLM, local database

## Contributing

This is a personal project, but contributions are welcome!

For others wanting to adapt this for their sport:
1. Copy `configs/sports/template.yaml`
2. Fill in your sport's phases and session types
3. Run onboarding with your sport config
4. Adjust prompts in `internal/agent/prompts.go` if needed

## License

MIT License - see LICENSE file

## Roadmap Highlights

See [docs/PLAN.md](docs/PLAN.md) for complete roadmap.

---

Built with for better training.
