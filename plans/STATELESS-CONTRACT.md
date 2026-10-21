# Stateless Execution Contract — TestWeaver

This repository uses **Reentry 3.3 (Journal Mode)** for continuous, stateless multi-agent execution across **AGY**, **Codex**, and **MAS / Jev**.

## 1. Disposable Chat, Durable Files
- Do not rely on conversational memory from prior turns.
- All durable context lives in:
  1. `plans/journal.jsonl` (single lifecycle authority)
  2. `plans/CURRENT.md` and `plans/INDEX.md` (generated views — never hand-edit)
  3. `plans/BASELINE.md` (narrative landmark history)
  4. `plans/P00X.md` and `plans/P00X-design-contract.md` (master plan & frozen architectural contract)
  5. `plans/P00X-S00Y.md` and `plans/P00X-S00Y-results.md` (bounded session spec & delivery evidence)

## 2. Entry Point
Every fresh session starts from:
```text
Continue in stateless mode from plans/CURRENT.md
```

## 3. Execution Discipline
1. Read `plans/CURRENT.md`, then read the active master plan (`plans/P00X.md`), its frozen design contract (`plans/P00X-design-contract.md`), and the current executable session (`plans/P00X-S00Y.md`).
2. Respect the bounded file scope declared in `plans/P00X-S00Y.md`.
3. Enforce **Test-Driven Development (`$test-driven-development`)**: write failing tests (`RED`), implement minimal code to pass (`GREEN`), and refactor cleanly.
4. Run independent code review (`$mas-review`) before authoring `plans/P00X-S00Y-results.md`.
5. Materialize the next session file (`plans/P00X-S00(Y+1).md`) via **Cascading Authorship** when closing `S00Y`.
