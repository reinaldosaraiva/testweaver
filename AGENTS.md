# AGENTS.md — Governança Operacional do TestWeaver (MAS / Reentry 3.3 / Jev)

## 1. Propósito do Repositório
O **TestWeaver** é a plataforma autônoma para testes universais de contrato de API (OpenAPI 3.0/3.1), execução em 3 trilhos (`Track 1 Gate <= 2m`, `Track 2 Debounced Shards`, `Track 3 Nightly MARL/PBT`), grafo semântico de dependências (`SPDG` — *AutoRestTest ICSE 2025*), roteamento multi-modelo coarse-to-fine (*RESTGPT* + *LlamaRestTest FSE 2025*) e cobertura caixa-preta **`TCL0–TCL7`** (*Martín-López et al.*) com exportação nativa para SonarQube e Codecov.

## 2. Reentry 3.3 (Journal Mode) & Proteção de Estado
1. `plans/journal.jsonl` é a **autoridade única** de ciclo de vida.
2. `plans/CURRENT.md` e `plans/INDEX.md` são **views geradas** a partir de `plans/journal.jsonl`. **NUNCA** edite `plans/CURRENT.md`, `plans/INDEX.md` ou `plans/journal.jsonl` manualmente.
3. Ponto de entrada determinístico para qualquer sessão stateless:
   ```text
   Continue in stateless mode from plans/CURRENT.md
   ```
4. **Cascading Authorship**: Apenas uma sessão (`P00X-S00Y.md`) por plano ativo permanece materializada em estado `READY`. As sessões subsequentes ficam como previsão (`PLANNED`) no Plano-Mestre (`P00X.md`) e só são materializadas no closeout da sessão atual com base nas evidências reais (`plans/P00X-S00Y-results.md`).
5. **Closeout Atômico**:
   - Grave a evidência em `plans/<session_id>-results.md` contendo explicitamente `**Verdict:** PASS` (ou `**Verdict:** GO` em sessões de decisão/arquitetura).
   - Submeta a transação de closeout via `.reentry-agy-closeout.json` (ou `plans/.reentry-codex-closeout.json` conforme o host ativo).

## 3. Roteamento Multi-Agente (MAS / Jev) & Lentes de Engenharia
- **Coordenador (`root` — `Sol` / `Astra-high` / AGY Root)**:
  - Lê `plans/CURRENT.md`, sincroniza o GitHub Project Board, despacha cápsulas sanitizadas, integra o código e executa o closeout.
- **Consultor Read-Only (`GLM-5.3` via Jev Capsule)**:
  - Consultado antes do writer em decisões de concorrência Go, modelagem SQL, algoritmos de grafo (`SPDG`) e contratos OpenAPI.
- **Writer Único por Tarefa Bounded**:
  - **Go (`cmd/`, `server/`, `agent/`, `proxy/`, `tools/`, `db/`)**: usar obrigatoriamente as skills `$golang-pro` e `$test-driven-development`.
  - **Grafo SPDG, Sharding, Cálculo `TCL0–TCL7` e Recompensa MARL**: aplicar adicionalmente a lente `$knuth-style`.
  - **Runners / Adaptadores Python (`skills/`, scripts)**: aplicar `$hynek-style` e `$test-driven-development`.
- **Gate Independente de Revisão**:
  - Antes do closeout de qualquer sessão de código, rodar `$mas-review` em contexto limpo verificando regressões, aderência ao `P00X-design-contract.md` e cobertura de testes.

## 4. Governança Dual-Track: GitHub Project Board × Reentry
Repositório GitHub: `https://github.com/reinaldosaraiva/testweaver`
Project Board: **"TestWeaver — Produto e Entrega"**

| Estado no Reentry (`plans/CURRENT.md`) | Coluna no GitHub Project Board | Regra Operacional |
| :--- | :--- | :--- |
| Épico/Item ainda não fatiado | `Triagem` | Planos futuros (`P002..P006`) aguardando conclusão do plano predecessor. |
| `BLOCKED` (decisão funcional/produto) | `Definição funcional` | Aguardando definição humana de escopo/contrato. |
| Sessão materializada (`READY`) | `Pronto` | Posicionada para execução imediata pelo próximo loop stateless. |
| Plano `ACTIVE` / Sessão em curso | `Em desenvolvimento` | Card movido ao iniciar o trabalho da sessão. |
| `BLOCKED` (impedimento técnico/infra) | `Bloqueado` | Impedimento de ambiente ou dependência externa. |
| Validação visual/humana | `UAT` | Revisão de UI/Dashboard (`P006`). |
| `COMPLETE` (`Verdict: PASS` / `GO`) | `Concluído` | Movido no closeout junto com o fechamento da Issue no GitHub. |
