# TestWeaver

<p align="center">
  <img src="docs/assets/testweaver-architecture-whiteboard.png" alt="TestWeaver: Multi-Agent & Multi-Model API Testing Architecture (3,000+ Endpoints)" width="100%" />
</p>

> **Autonomous Polyglot API Contract Testing, Cooperative Multi-Agent Fuzzing (SPDG + MARL), 3-Tier Fast Pipelines (`<= 2 min` Gate), and OpenAPI `TCL0–TCL7` Maturity Coverage Platform.**
>
> Derived from [**ScopeWeaver**](https://github.com/cskwork/scopeweaver) and [**ARTEX**](https://github.com/Autumn-27/ARTEX) under **AGPL-3.0** (see [`docs/PROVENANCE.md`](docs/PROVENANCE.md)).

---

## Why TestWeaver?

Modern engineering organizations run polyglot portfolios (**Go**, **Python / FastAPI / Django**, **TypeScript / Next.js**, **Java / Kotlin**) scaling from dozens to **3,000+ REST API endpoints**. Traditional API test suites fail at scale in two ways:

1. **CI Bottleneck**: Running sequential black-box/integration tests across 3,000 endpoints takes 30–60+ minutes, blocking developers and autonomous coding agents on every push.
2. **Shallow Coverage & Blind Spots**: Hand-written tests only exercise happy paths (`200 OK`), missing inter-endpoint producer-consumer dependencies (`POST /users` -> `GET /accounts/{id}` -> `PATCH /limit`), undocumented `500 Internal Server Error` edge cases, RBAC/Auth scope leaks (`401`/`403`), and silent OpenAPI schema drift (`RUNTIME_ONLY` / `SPEC_ONLY`).

**TestWeaver** solves both problems in a **single Go binary + embedded Next.js Dashboard + PostgreSQL Graph + MITM Proxy (`:8788`)**:

- **1. Semantic Property Dependency Graph (`SPDG`) in PostgreSQL**: Parses OpenAPI 3.0/3.1 specifications and builds a live graph of endpoints (`assets`), parameter producer-consumer edges (`exploration_edges`), and AST-level semantic diffs (`oasdiff`).
- **2. Coarse-to-Fine Multi-Model Routing (`llm_profiles`)**:
  - **Tier A — Reasoning LLM (`GLM-5.3` / `Claude`)**: Runs *once* per `openapi_schema_hash` change to extract hidden natural-language constraints, domain invariants, and parameter dependencies from OpenAPI docstrings.
  - **Tier B — High-Throughput SLM Pool (`Qwen3.8-27B` / Fast Coding Pool)**: Drives parallel workers (`work#1..work#N`) for low-latency payload synthesis, boundary mutation, and stateful sequence generation.
- **3. Cooperative Multi-Agent Loop (`MARL` + Specialized Agents)**:
  - `Planner Agent`: Decomposes testing value across tracks and prioritizes frontier operations using reward feedback $R = \alpha \cdot \Delta\text{TCL} + \beta \cdot \mathbb{I}(\text{status} \ge 500) - \gamma \cdot \text{Flake}$.
  - `Dependency Agent`: Resolves stateful `POST -> GET -> PATCH -> DELETE` chains over the SPDG.
  - `Operation & Auth Agent`: Exercises RBAC/scope matrices (`401`/`403`) and security headers.
  - `Value & PBT Fuzzer Agent`: Combines Tier B SLMs with Property-Based Testing (Schemathesis / Hypothesis / RESTler strategies).
  - `Retester & Self-Healing Oracle Agent`: Replays content-addressed SHA-256 HTTP snapshots (`traffic_evidence_snapshots`) to confirm regressions and automatically quarantine flaky endpoints (`quarantined = true`).
- **4. 3-Tier Fast Execution Pipeline**:
  - **Track 1 — Sync Fast Gate (`<= 2 min`)**: Executes only the dirty subgraph (`dirty_in_diff = true`) on pre-push / PR with HTTP keep-alive connection pooling.
  - **Track 2 — Async Debounced Full Suite (`10 min` window)**: Coalesces commit bursts and runs deterministic domain shards across `work#1..work#N`.
  - **Track 3 — Nightly Deep MARL & Stateful PBT**: Runs deep property-based fuzzing, stateful sequences, latency SLO checks (`p95`/`p99`), and flake quarantine verification.
- **5. Live MITM Proxy (`:8788`) & `TCL0–TCL7` Maturity Oracle**: Passively inspects all test traffic against the OpenAPI contract, grades every endpoint from **`TCL0` (Untested)** to **`TCL7` (Stateful Flows)**, and exports `openapi-cobertura.xml` and `openapi-sonar-coverage.xml` directly to **SonarQube** and **Codecov** with zero blocking overhead.

---

## Scientific Foundations (2019–2025 Literature)

TestWeaver translates five peer-reviewed software engineering research breakthroughs directly into production infrastructure:

| Paper / Venue | Authors | Core Mechanism Integrated in TestWeaver | TestWeaver Subsystem |
| :--- | :--- | :--- | :--- |
| **AutoRestTest** (*ICSE 2025*) | M. Kim, T. Stennett, S. Sinha, A. Orso (*Georgia Tech & IBM Research*) | **Semantic Property Dependency Graph (`SPDG`) + Multi-Agent Reinforcement Learning (`MARL`)** with value decomposition across Operation, Parameter, Dependency, and Value agents. | `db/schema.sql` (`assets`, `exploration_nodes`, `exploration_edges`) + `agent/planner` |
| **RESTGPT** (*ICSE / ISSTA 2024*) | M. Kim et al. | **Natural-Language Constraint & Rule Extraction** from OpenAPI descriptions using reasoning LLMs prior to test generation. | **Tier A Model Pool** (`GLM-5.3` / `Claude`) cached by `openapi_schema_hash` |
| **LlamaRestTest** (*ESEC/FSE 2025*) | M. Kim, T. Stennett, S. Sinha, A. Orso | **Small/Fast Language Models (`SLMs`)** for high-throughput, low-cost REST payload generation and constraint satisfaction. | **Tier B Model Pool** (`Qwen3.8-27B` / Fast Pool) across `work#1..work#N` |
| **Test Coverage Levels (`TCL0–TCL7`)** (*ICWS / ACM TOSEM 2019–2022*) | A. Martín-López, S. Segura, A. Ruiz-Cortés | **Black-Box Specification Coverage Hierarchy** (`TCL0` Untested -> `TCL1` Paths -> `TCL2` Ops -> `TCL3` Content-Types -> `TCL4` Params/Auth -> `TCL5` Status Codes -> `TCL6` Response Schema -> `TCL7` Stateful Flows). | `proxy/` Live OpenAPI Validator + `server/` Cobertura/Sonar XML Exporter |
| **RESTler** (*ICSE 2019*) & **Schemathesis** (*ICSE 2022*) | V. Atlidakis, P. Godefroid, M. Polishchuk; Z. Hatfield-Dodds, D. Dygalo | **Stateful Producer-Consumer Fuzzing & Property-Based Shrinking** with deterministic HTTP replay. | `skills/stateful-pbt-fuzzer` + `traffic_evidence_snapshots` + `agent/retester` |

<p align="center">
  <img src="docs/assets/testweaver-3tier-pipeline-whiteboard.png" alt="TestWeaver: Universal Polyglot 3-Tier Fast Pipeline & OpenAPI TCL Coverage" width="100%" />
</p>

---

## Architecture Overview

```mermaid
flowchart LR
    subgraph SPDG_LAYER["1. OpenAPI 3.1 Ingestion & SPDG Graph"]
        direction TB
        POLY["Polyglot Services (Go / Python / Next.js)<br/>3,000+ Endpoints (/openapi.json)"]
        OASDIFF["Semantic AST Diff (oasdiff)<br/>openapi_schema_hash + dirty_in_diff"]
        TIER_A["Tier A Reasoning LLM (GLM-5.3 / Claude)<br/>Extracts Semantic Rules & Hidden Constraints"]
        SPDG["PostgreSQL SPDG (db/schema.sql)<br/>assets + exploration_edges (sim >= τ)"]

        POLY --> OASDIFF --> TIER_A --> SPDG
    end

    subgraph MARL_LAYER["2. Cooperative Multi-Agent & 3-Tier Engine"]
        direction TB
        PLANNER["Planner & Value Decomposer Agent<br/>R = α·ΔTCL + β·I(5xx) - γ·Flake"]
        DEP_AG["Dependency Agent (SPDG Chains)"]
        AUTH_AG["Operation & RBAC Agent (401/403)"]
        VAL_AG["Value & PBT Fuzzer Agent<br/>Tier B SLM Pool (Qwen3.8-27B) + Schemathesis"]
        RET_AG["Retester & Self-Healing Oracle Agent<br/>SHA-256 Replay + Flake Quarantine"]

        T1["Track 1: Sync Gate (<= 2 min)<br/>Dirty-in-Diff Subgraph Only"]
        T2["Track 2: Async Debounced (10m)<br/>Domain Shards work#1..work#N"]
        T3["Track 3: Nightly MARL & PBT<br/>Deep Stateful Fuzzing"]

        SPDG --> PLANNER
        PLANNER --> T1 & T2 & T3
        T1 & T2 & T3 --> DEP_AG & AUTH_AG & VAL_AG
        VAL_AG --> RET_AG
    end

    subgraph ORACLE_LAYER["3. MITM Proxy :8788 & TCL0-TCL7 Oracle"]
        direction TB
        PROXY["MITM Proxy :8788 + Traffic Evidence SHA-256<br/>Live Request/Response Schema Validation"]
        TCL["TCL0–TCL7 Maturity Calculator"]
        SONAR["Cobertura XML & SonarQube / Codecov Exporter"]

        DEP_AG & AUTH_AG & VAL_AG & RET_AG -->|"HTTP_PROXY=:8788"| PROXY
        PROXY --> TCL --> SONAR
        TCL -.->|"Reward Signal (ΔTCL)"| PLANNER
    end
```

---

## Quick Start

### 1. Prerequisites
- **Go** `1.24+`
- **Node.js** `20+` (for building the embedded Next.js WebUI and MCP/Pi agent adapter)
- **PostgreSQL** `15+` (schema in [`db/schema.sql`](db/schema.sql) auto-migrates idempotently on startup)

### 2. Build Single Binary (Backend + Embedded WebUI)

```bash
# 1. Build static Next.js frontend and copy to server/webui/dist
cd web && npm ci && npm run build && cd ..
rm -rf server/webui/dist && cp -r web/out server/webui/dist

# 2. Compile single binary with embedded UI
go build -tags embedui -o testweaver ./cmd/artex
```

### 3. Run TestWeaver

```bash
export TESTWEAVER_DB_DSN="postgres://postgres:postgres@127.0.0.1:5432/testweaver?sslmode=disable"
./testweaver
```

- **Web Dashboard & REST API**: `http://127.0.0.1:8787`
- **Live OpenAPI MITM Proxy & TCL Calculator**: `http://127.0.0.1:8788`

---

## Development & Autonomous Execution (MAS / Reentry 3.3 / Jev)

This repository is managed via **Reentry 3.3 (Journal Mode)** inDual-Track synchronization with the **GitHub Project Board ("TestWeaver — Produto e Entrega")**:

- **Kickoff Command**: `Continue in stateless mode from plans/CURRENT.md`
- **Governance Rules**: [`AGENTS.md`](AGENTS.md) and [`plans/STATELESS-CONTRACT.md`](plans/STATELESS-CONTRACT.md)
- **Active Master Plans**:
  - [`plans/P001.md`](plans/P001.md) — Bootstrap Standalone `TestWeaver`, Reentry 3.3 Loop, CLI/MCP Rebranding & Typed `EndpointTestMetadata`
  - [`plans/P002.md`](plans/P002.md) — OpenAPI 3.0/3.1 Asset Sync, `SPDG` Builder & Semantic Diff (`oasdiff`)
  - [`plans/P003.md`](plans/P003.md) — 3-Tier Execution Engine (`Track 1 <= 2m`, `Track 2 Debounced Shards`, `Track 3 Nightly MARL/PBT & Quarantine`)
  - [`plans/P004.md`](plans/P004.md) — Proxy `:8788` Live OpenAPI Schema Validator, `TCL0–TCL7` Calculator & SonarQube/Codecov XML Exporter
  - [`plans/P005.md`](plans/P005.md) — Polyglot QA Skills (`openapi-sync`, `contract-drift-guard`, `polyglot-api-runner`, `stateful-pbt-fuzzer`) & Multi-Model Routing (`RESTGPT` + `LlamaRestTest`)
  - [`plans/P006.md`](plans/P006.md) — Next.js Dashboard UI (`TCL0–TCL7` Coverage Matrix, 3-Tier Track Filter, Contract Breach Diff Viewer & Flake Quarantine)

---

## License & Provenance

Licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**. See [`LICENSE`](LICENSE) and [`docs/PROVENANCE.md`](docs/PROVENANCE.md) for the complete upstream lineage (`Autumn-27/ARTEX` -> `cskwork/scopeweaver` -> `reinaldosaraiva/testweaver`).
