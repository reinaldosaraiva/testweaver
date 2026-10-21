# Current Kickoff

This file is the stable single-line kickoff entry point for a fresh stateless session in `reentry`.

## Use This Exact Kickoff

```text
Continue in stateless mode from plans/CURRENT.md
```

## Current Plan Pointer

- Current master plan: `plans/P001.md`
- Pointer session: `P001-S001`
- Pointer session status: `READY`
- Current executable session: `P001-S001`

## Journal View Rules

- `plans/INDEX.md` and `plans/CURRENT.md` are generated from `plans/journal.jsonl`.
- Do not edit `plans/INDEX.md`, `plans/CURRENT.md`, or `plans/journal.jsonl` directly.
- Route lifecycle changes through `plans/.reentry-codex-closeout.json`.

## Status Precedence

- `plans/journal.jsonl` is the only lifecycle authority.
- `plans/INDEX.md` and `plans/CURRENT.md` are generated views of the journal.
- Plan and session files provide execution scope and evidence only.
- Results, handoff, and state artifacts preserve evidence and observed state.
