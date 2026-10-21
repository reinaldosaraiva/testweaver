# Provenance and Attribution

## Upstream Lineage

1. **ARTEX (`Autumn-27/ARTEX`)**:
   - Upstream repository: `https://github.com/Autumn-27/ARTEX`
   - Fork base commit: `160fe13` (`160fe133560b443ff472d44dc2a60b67d57a679b`)
   - License: `AGPL-3.0`
2. **ScopeWeaver (`cskwork/scopeweaver`)**:
   - Derivative repository: `https://github.com/cskwork/scopeweaver`
   - Added content-addressed HTTP traffic evidence (`traffic_evidence_snapshots`, `finding_traffic_bindings`), local MCP/Pi agent adapters (`adapters/agent`), and multi-model pool routing (`llm_profiles`).
3. **TestWeaver (`reinaldosaraiva/testweaver`)**:
   - Derivative repository: `https://github.com/reinaldosaraiva/testweaver`
   - Dedicated derivative of ScopeWeaver re-engineered for **Universal Polyglot API Contract Testing (OpenAPI 3.0/3.1)**, **Semantic Property Dependency Graphs (`SPDG`) + Cooperative Multi-Agent Testing (`MARL`)**, **Coarse-to-Fine Multi-Model Routing (`RESTGPT` + `LlamaRestTest`)**, **3-Tier Fast Pipelines (`<= 2 min` Gate, Debounced Shards, Nightly PBT)**, and **Black-Box `TCL0–TCL7` Maturity Coverage + SonarQube/Codecov Export**.

## Compatibility Notes

- The Go module path remains `github.com/Autumn-27/artex` and main entrypoint package remains `./cmd/artex` to avoid unnecessary internal import churn and allow clean cherry-picking of infrastructure fixes between `ScopeWeaver` and `TestWeaver`.
- The compiled binary name is `testweaver`, and configuration loaders read `TESTWEAVER_*` environment variables first with fallback to `SCOPEWEAVER_*` / `ARTEX_*`.
