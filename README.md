# Daggerheart Adversary Tracker

A web application for tracking Daggerheart Adversaries during combat. This tool helps Game Masters manage adversaries and combat encounters for tabletop roleplaying games.

## Features

- Create and store adversary statblocks for quick reference
- Build and save encounters with multiple adversaries
- Track initiative, health, and conditions during combat
- Organize adversaries by type, challenge rating, and more

## Tech Stack

- Go (Golang) for backend
- Chi router for HTTP routing
- SQLite for data persistence
- HTML/CSS with Tailwind CSS for styling
- HTMX and Alpine.js for interactivity

## Development

### Prerequisites

- Go 1.21 or higher
- SQLite

### Setup

```bash
# Clone the repository
git clone https://github.com/juthrbog/adversarytracker.git
cd adversarytracker

# Run the application
go run cmd/app/main.go
```

The application will be available at http://localhost:8080

## Project Structure

```
/cmd/app/             # Main application entry point
/internal/            # Private app packages
/pkg/                 # Reusable code
/templates/           # HTML templates
/static/              # CSS/JS assets
/web/handlers/        # HTTP handlers
/web/middleware/      # Custom middleware
/db/                  # Database access
/data/                # SQLite database file
```

## License

MIT
