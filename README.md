# gator

gator is a small CLI RSS reader/aggregator written in Go.

## Prerequisites
- Go (1.20+ recommended)
- PostgreSQL running locally or reachable via URL

## Install
```sh
go install github.com/Heliotropoff/gator@latest
```
This installs the `gator` binary into your `GOBIN` (default: `~/go/bin`).

## Database setup
Create a database (for example `gator`) and run the migrations in `sql/schema` in order. Example with `psql`:
```sh
psql -d gator -f sql/schema/01_users.sql
psql -d gator -f sql/schema/02_feeds.sql
psql -d gator -f sql/schema/03_feed_follows.sql
psql -d gator -f sql/schema/04_feeds.sql
psql -d gator -f sql/schema/05_posts.sql
# include any later migrations if present
```

## Configuration
Create `~/.gatorconfig.json` with your database URL. The file is read on startup and also stores the current user after `login`/`register`.

Example:
```json
{
  "db_url": "postgres://user:pass@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

## Usage
Run the CLI with a subcommand:
```sh
gator <command> [args]
```

Common commands:
- `register <username>`: create a new user and set as current.
- `login <username>`: switch current user to an existing one.
- `addfeed <name> <url>`: add a feed owned by the current user and auto-follow it.
- `follow <feed_url>`: follow an existing feed.
- `feeds`: list all feeds with owners.
- `following`: list feeds the current user follows.
- `browse [limit]`: show recent posts from feeds you follow (default limit 2).
- `agg <duration>`: start the fetch loop (e.g., `agg 1m`); fetches feeds periodically.
- `unfollow <feed_url>`: unfollow a feed.
- `users`: list users (marks current).
- `reset`: delete all users (and cascaded data).

Stop the aggregator with Ctrl+C.

## Development
- Run `go test ./...` to execute tests (if present).
- Update SQL via `sql/schema` migrations and regenerate query code with `sqlc` if you change queries.

