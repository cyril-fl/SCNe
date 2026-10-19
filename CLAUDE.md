@AGENTS.md

## Claude Code notes

`AGENTS.md` and `.agents/skills/` are the single source of truth shared with
Codex. `.claude/skills` is a symlink to `.agents/skills`; edit skills there,
never in a copy.

Skill references in `AGENTS.md` and in skills use Codex `$name` syntax. In
Claude Code, invoke them with the Skill tool:

- `$go-test-create` → `go-test-create`
- `$go-test-run` → `go-test-run`
- `$go-tests` → `go-tests`
- `$go-documentation` → `go-documentation` (manual only: `/go-documentation`)
- `$primitive-first-review` → `primitive-first-review`
- `$plum-integration` → `plum-integration`

The `agents/openai.yaml` files inside skills are Codex-only metadata; ignore them.
