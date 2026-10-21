# P001 Design Contract — TestWeaver Core & Multi-Agent Architecture Invariants

- **Plan:** `P001`
- **Status:** DRAFT (Target: `FROZEN` upon `P001-S001` `GO` verdict)

## 1. Provenance & Module Stability (`AGPL-3.0`)
1. **License Continuity**: `TestWeaver` is licensed under **AGPL-3.0**, preserving the complete upstream chain (`Autumn-27/ARTEX` -> `cskwork/scopeweaver` -> `cskwork/testweaver`) in `LICENSE` and `docs/PROVENANCE.md`.
2. **Go Module Path Stability**: Keep `module github.com/Autumn-27/artex` in `go.mod` and `./cmd/artex` as the main package directory so upstream bugfixes from `ScopeWeaver` apply cleanly without import-path rewrites. The compiled binary output MUST be `testweaver`.

## 2. Database Schema & `EndpointTestMetadata` Contract (`db/schema.sql`)
1. Preserve the existing idempotent PostgreSQL tables (`assets`, `explorations`, `exploration_nodes`, `exploration_edges`, `exploration_anchors`, `activity`, `llm_profiles`, `traffic_evidence_snapshots`, `finding_traffic_bindings`).
2. For `assets` with `type = 'endpoint'`, `assets.extra` (`JSONB`) MUST serialize/deserialize through a canonical Go struct `EndpointTestMetadata`:
   - `operation_id` (`string`): Unique OpenAPI operation identifier.
   - `openapi_schema_hash` (`string`): SHA-256 of normalized request/response schema + parameters + security scopes.
   - `tcl_level` (`int`, `0..7`): Current Martín-López Test Coverage Level (`TCL0` Untested through `TCL7` Stateful Flows).
   - `status_codes_covered` (`[]int`): Distinct HTTP status codes observed via `:8788` proxy traffic evidence.
   - `dirty_in_diff` (`bool`): Flagged by semantic OpenAPI diff (`oasdiff`) for Track 1 Sync Fast Gate (`<= 2 min`).
   - `drift_status` (`string`): `"OK" | "SPEC_ONLY" | "RUNTIME_ONLY" | "SCHEMA_MISMATCH"`.
   - `execution_track` (`string`): `"track1_gate" | "track2_debounced" | "track3_nightly"`.
   - `semantic_rules` (`[]string`): Natural-language constraints extracted once per `openapi_schema_hash` by Tier A Reasoning LLM (*RESTGPT*).
   - `quarantined` (`bool`): Automatic flake quarantine flag managed by the `retester` agent.

## 3. Multi-Model & Multi-Agent Routing Invariants
1. **Tier A (Reasoning LLM — `GLM-5.3` / `Claude`)**: Invoked only when `openapi_schema_hash` changes to populate `semantic_rules` and infer inter-parameter `SPDG` edges (`sim(e_i, e_j) >= tau`).
2. **Tier B (High-Throughput SLM Pool — `Qwen3.8-27B` / Fast Pool)**: Invoked by worker goroutines (`work#1..work#N`) for payload generation and stateful mutations (*LlamaRestTest*).
3. **MARL Reward Signal (*AutoRestTest*)**: Planner prioritizes frontier nodes using $R = \alpha \cdot \Delta\text{TCL} + \beta \cdot \mathbb{I}(\text{status} \ge 500) - \gamma \cdot \text{Flake}$.

## 4. Environment & MCP Adapter Backwards Compatibility
1. Configuration readers MUST check `TESTWEAVER_*` environment variables first and fall back to `SCOPEWEAVER_*` and `ARTEX_*`.
2. `adapters/agent` MUST expose `testweaver_*` MCP/Pi tool names while keeping `scopeweaver_*` aliases functional during transition.
