# AGENTS.md

## Monorepo layout

```
sentinel/                    # Git root
├── sentinel-api/            # Python/FastAPI — primary, fully functional
└── sentinel-worker/         # Go — skeleton, not functional yet
```

No root-level tooling. All commands run inside `sentinel-api/`.

## Commands (run from `sentinel-api/`)

```bash
uv sync --group dev          # install deps
docker compose up -d         # start PostgreSQL 17 (creates sentinel_db)
alembic upgrade head         # run migrations (requires .env with DATABASE_URL)
uv run uvicorn app.main:app --reload   # dev server (APScheduler starts in-process)

uv run ruff check .          # lint
uv run ruff format .         # format
pyright                      # type-check (must be installed globally; not via uv)

uv run pytest                # all tests
uv run pytest app/tests/api/test_monitors.py::test_create_monitor  # single test
```

## Environment

Copy `.env.example` or create `.env` inside `sentinel-api/`:

```
DATABASE_URL=postgresql+asyncpg://postgres:postgres@127.0.0.1:5432/sentinel_db
SECRET_KEY=...
ALGORITHM=HS256
ACCESS_TOKEN_EXPIRE_MINUTES=30
```

Alembic's `env.py` imports `app.core.config.Settings` which reads `.env` from CWD — so both `alembic` and `uv run` must be run from `sentinel-api/`.

## Test prerequisites

Tests require a real PostgreSQL database `sentinel_tests_db` on localhost:5432. The test DB URL is **hardcoded** in `conftest.py:13`:

```
postgresql+asyncpg://postgres:postgres@127.0.0.1:5432/sentinel_tests_db
```

Tables are created/dropped per function via `SQLModel.metadata.create_all`/`drop_all`. No migrations are run for tests.

## Architecture notes

- **Checker registry**: New check types implement `BaseChecker` and self-register with `@register`. The scheduler discovers them via `get_checker(check_type)` — no changes needed in API or scheduler code.
- **State machine**: Alerts fire **only on transitions** (healthy→unhealthy = "down", unhealthy→healthy = "recovery"). Every individual check produces a `CheckResult`; `alert` rows are emitted only on state changes.
- **Scheduler**: APScheduler `AsyncIOScheduler` runs `check_all_monitors` every 10 seconds inside the FastAPI lifespan. No separate worker process needed.
- **All DB access is async** (asyncpg, async SQLAlchemy sessions, async Alembic).

## sentinel-worker (Go)

Stub only — most function bodies are empty. Builds with `go build ./cmd/worker/` but does nothing meaningful. Shares the same PostgreSQL schema; designed to poll `monitor` and write `check_result` + `alert` independently.