# Fog (getfog.dev)

Fog turns your local machine into a "personal cloud" for AI coding agents.

This repo is intentionally local-first:
- Fog does not host an LLM.
- Fog runs existing AI coding CLIs (Cursor Agent, Claude Code, Gemini CLI, Aider) on the user's machine.
- Isolation and auditability come from Git worktrees and Git history.

## Components

- `cmd/wtx`: Git worktree manager (no AI, no networking).
- `cmd/fog`: CLI for one-off tasks + onboarding + managed repo registry.
- `cmd/fogd`: local daemon (HTTP API + session execution).
- `cmd/fogapp`: Wails desktop UI that talks to `fogd` and starts an embedded daemon when needed.
- `cmd/fogcloud`: optional distribution/control-plane (present in repo, not the current focus).

## Current Product Direction

Desktop-first sessions are the primary UX.

- A **session** is a long-lived branch + worktree "conversation" that can be followed up.
- A **run** is one execution inside a session (one prompt, one lifecycle).
- Follow-ups and re-runs operate on the same session worktree.
- Forking is explicit: it creates a new branch/worktree from the current session head.
- Streaming output is persisted as run events and can be consumed via SSE.

Slack/cloud code exists, but if you are adding features, default to desktop + local API unless explicitly asked to work on Slack/cloud.

## Data And Security

Fog home:
- `FOG_HOME` (default `~/.fog`)
- SQLite DB: `FOG_HOME/fog.db` (repos, settings, secrets, sessions, runs, run_events, tasks)
- Master key: `FOG_HOME/master.key` (file-based AES-256-GCM key)

Rules:
- Never store GitHub PATs in plaintext.
- Secrets are encrypted at rest using AES-GCM with a local key file (no keychain dependency).
- Avoid logging tokens; sanitize git args that may contain headers.

Managed repos layout (default):
- `FOG_HOME/repos/<owner>/<repo>/repo.git` (bare clone)
- `FOG_HOME/repos/<owner>/<repo>/base` (base worktree)

## Key Invariants

- Safety: operate in worktrees; do not mutate user repos in-place.
- Cancellation: only the latest active run is cancelable; cancellation stops the entire process group.
- PR creation: uses `gh` CLI; PR is created once per session (draft) when enabled; follow-ups update the same branch/PR.

## Where To Change Things

- AI tool adapters and streaming: `internal/ai/*`
- Session engine (follow-ups, fork, cancellation, run events): `internal/runner/session.go`
- HTTP API endpoints: `internal/api/*`
- State (SQLite schema, encryption): `internal/state/*`
- Desktop UI: `cmd/fogapp/frontend/*` (plain HTML/CSS/JS embedded into Wails)

## Development Commands

- Unit tests: `go test ./...`
- Build binaries: `make all`
- Desktop dev: `make fogapp-dev` (requires `wails`)
- Desktop build: `make fogapp-build`

Build tags:
- Desktop uses `-tags desktop` (`cmd/fogapp` enforces this).

## Notes For Other Agents

- Prefer small, reviewable changes and keep code simple.
- Avoid third-party deps unless there is a clear payoff.
- Keep docs in sync with behavior (especially session semantics and storage layout).

