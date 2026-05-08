# Gauge → Ginkgo Migration — Plain English Plan

---

## What are we actually doing and why?

Right now our tests are written in a tool called **Gauge**. Think of Gauge like a recipe book — you write test steps in plain English in `.spec` files, and Go code in `steps/` executes each step.

We want to move to **Ginkgo**, which is the standard Go testing framework used across the entire Tekton / OpenShift ecosystem. The benefit: one language (Go), one toolchain, easier for new contributors to pick up, and it works natively with how Konflux CI runs tests today.

**Nothing about what we test changes. Only how we write and run the tests changes.**

---

## The 7 steps, explained simply

---

### Step 1 — What tests do we actually have?

We have **135 test cases** spread across **31 test files**, covering every part of OpenShift Pipelines:

| Area | # Tests | What it covers |
|------|---------|----------------|
| Ecosystem Tasks | 36 | buildah, S2I, skopeo, git-cli, helm, kn — the bundled tasks |
| Triggers | 23 | EventListeners, TLS, GitHub/GitLab/Bitbucket webhooks, CEL |
| Operator | 33 | Auto-pruning, RBAC, addon config, HPA, upgrade (pre + post) |
| Pipelines Core | 16 | Run/cancel/timeout pipelines, resolvers (git, cluster, bundles, http) |
| PAC | 7 | Pipelines as Code on GitHub and GitLab |
| Chains | 2 | Signing task runs and images |
| Results | 2 | Storing task/pipeline run records |
| Manual Approval Gate | 2 | Approve / reject a gate pipeline |
| Metrics | 1 | Prometheus metrics |
| Versions / Sanity | 2 | All component versions are correct |
| OLM / Install | 3 | Install, upgrade, uninstall via OLM |
| Console Icon | 1 | The icon shows up in the OpenShift web console |

---

### Step 2 — Which tests are ours to own, and which are already done elsewhere?

Not all 135 tests are unique to us. Some of them test generic Tekton behaviour that the upstream Tekton project already tests on vanilla Kubernetes.

**Tests we must keep (~110 tests, ~82%)** — these only make sense on OpenShift:
- Anything that uses the **OpenShift Operator** (install, upgrade, addon config, auto-prune, RBAC)
- Anything that uses **OpenShift-specific images** (S2I, ImageStreams, buildah with OpenShift registry)
- Anything that uses **OLM** (the OpenShift package manager)
- **PAC, Chains, Results, Manual Approval Gate** — all managed by the Operator
- **Ecosystem tasks** — these are shipped only in the downstream catalog

**Tests we can likely drop (~25 tests, ~18%)** — upstream already covers these on Kubernetes, and Konflux now runs upstream tests on OpenShift nodes too:
- Basic "run a pipeline and check it succeeds"
- Basic "pipeline timeout failure"
- Basic "cancel a pipeline"
- CronJob trigger, pipeline tutorial basics

> **Action:** Before deleting anything, confirm the equivalent upstream test runs in Konflux on an OpenShift cluster. Only then remove the duplicate.

---

### Step 3 — Where do the new tests live?

**Stay in the same repo** (`release-tests`). No new repo needed.

Here's what changes inside the repo:

```
TODAY                          AFTER MIGRATION
─────────────────────────────  ──────────────────────────────────
specs/                         test/                   ← NEW
  pipelines/run.spec      →      pipelines/pipelines_test.go
  triggers/eventlistener… →      triggers/triggers_test.go
  operator/auto-prune…    →      operator/operator_test.go
  … etc                          … etc

steps/                         (deleted — no longer needed)
  pipeline/pipeline.go
  triggers/triggers.go
  … etc

pkg/                           pkg/                    ← UNCHANGED
  pipelines/                     pipelines/
  k8s/                           k8s/
  operator/                      operator/
  … etc                          … etc
```

The key insight: **all the real test logic already lives in `pkg/`**. The `steps/` folder is just a thin connector layer between Gauge and `pkg/`. Ginkgo calls `pkg/` directly, so `steps/` becomes unnecessary and gets deleted.

---

### Step 4 — What changes in CI?

Almost nothing. The pipeline that runs our tests today (in `.tekton/acceptance-tests-pr.yaml`) does:

1. Provision an OpenShift cluster
2. Install the Operator via OLM
3. **Run the tests** ← only this line changes
4. Collect results, upload to ReportPortal / Polarion
5. Tear down the cluster

**Today (Gauge):**
```
gauge run --tags sanity specs/pipelines/
```

**After (Ginkgo):**
```
go test ./test/pipelines/... --ginkgo.label-filter=sanity
```

That's the only change in CI. Everything else — cluster setup, reporting, Slack notifications — stays exactly the same.

The Docker image used to run tests also just swaps `gauge` for the `ginkgo` CLI tool.

---

### Step 5 — How do Gauge specs map to Ginkgo tests?

Think of it as a direct translation. Every concept in Gauge has an exact equivalent in Ginkgo:

| Gauge (what we have) | Ginkgo (what we write) | In plain English |
|----------------------|------------------------|-----------------|
| `# Suite title` at top of spec file | `Describe("Suite title", ...)` | The name of the test group |
| `## Scenario name: TC-ID` | `It("Scenario name", ...)` | One individual test |
| `Tags: e2e, sanity` | `Label("e2e", "sanity")` | Labels to filter which tests run |
| `Pre condition: * step` | `BeforeEach(...)` | Setup that runs before every test |
| `Steps: * do this \n * do that` | lines of code inside `It(...)` | The actual test steps |
| `gauge.BeforeScenario` | `BeforeEach` | "Before each test, do this setup" |
| `gauge.AfterScenario` | `AfterEach` | "After each test, do this cleanup" |
| Data table `\|row1\|row2\|` | `DescribeTable` + `Entry(...)` | Running the same test with different inputs |
| Concepts (`.cpt` files) | A shared Go helper function | Reusable steps called by multiple tests |

**Concrete example** — this Gauge scenario:

```
## Run sample pipeline: PIPELINES-03-TC01
Tags: e2e, pipelines, sanity
Steps:
  * Verify that image stream "golang" exists
  * Create pvc.yaml and pipelinerun.yaml
  * Verify pipelinerun "output-pipeline-run-v1b1" is "successful"
```

Becomes this Ginkgo test:

```go
It("Run sample pipeline", Label("PIPELINES-03-TC01", "e2e", "pipelines", "sanity"), func() {
    openshift.VerifyImageStreamExists(cs, "golang")
    k8s.CreateFromFile(cs, ns, "testdata/pvc/pvc.yaml")
    k8s.CreateFromFile(cs, ns, "testdata/v1beta1/pipelinerun/pipelinerun.yaml")
    pipelines.ValidatePipelineRun(cs, "output-pipeline-run-v1b1", "successful", ns)
})
```

Same test, same steps, just written in Go instead of plain English.

---

### Step 6 — What order do we convert things in?

We do this in 4 phases over ~10 weeks, so we're never in a broken state:

#### Phase 1 — Lay the foundation (Week 1–2)
Get the plumbing right before converting any tests:
- Add Ginkgo to the project
- Create the test entry point (`test/suite_test.go`) with setup/teardown that mirrors what Gauge does today
- Remove Gauge-specific code from shared packages (`pkg/store`, `pkg/config`)
- Update the Docker image (swap `gauge` → `ginkgo`)

No tests converted yet — CI still runs Gauge.

#### Phase 2 — Convert the quick wins first: Sanity tests (Week 3–4)
Start with the **sanity-tagged tests** since they run in every CI job and are the most critical:
- Versions check (2 tests)
- Pipeline timeout + results (2 tests)
- Addon enable/disable (2 tests)
- Ecosystem buildah (1 test)
- Chains signature (1 test)
- Results (1 test)

**Gate:** Run the Ginkgo sanity tests on a real cluster alongside the Gauge ones. Both must pass with the same results.

#### Phase 3 — Convert the full suite (Week 5–8)
Convert everything else, in order of value:

1. **Auto-prune** (15 tests) — most tests, high business value
2. **EventListeners** (17 tests) — complex but well-understood
3. **Ecosystem tasks** (21 tests) — largest spec; use table-driven tests
4. **S2I ecosystem** (9 tests) — quick with tables
5. **Operator features** (addon, RBAC, roles, HPA)
6. **Chains, Results, Metrics, PAC**
7. **Upgrade tests** (pre + post) — run last, most sensitive
8. **OLM install/icon** — dedicated install job

#### Phase 4 — Clean up (Week 9–10)
- Delete `specs/` folder
- Delete `steps/` folder
- Remove Gauge from `go.mod`
- Update the CI command in `.tekton/`
- Update `README.md`

---

### Step 7 — How do we know the Ginkgo version is as good as the Gauge version?

Before we retire each Gauge spec, we run a side-by-side check:

| What to check | What "pass" looks like |
|--------------|------------------------|
| Same number of tests | Each Gauge `##` scenario has a matching Ginkgo `It(...)` |
| All sanity tests pass | Ginkgo sanity run = same result as Gauge sanity run on same cluster |
| Tag filtering works | `--ginkgo.label-filter=sanity` picks the same tests as `--tags sanity` in Gauge |
| Failure behaviour | When a test fails, the test namespace is **kept** (for debugging), not deleted |
| Reports work | JUnit XML output from Ginkgo is accepted by the Polarion uploader |

During Phases 2–3, both Gauge and Ginkgo run in CI at the same time. Once a suite passes the side-by-side check, we drop Gauge for that suite. By the end of Phase 4, only Ginkgo runs.

---

## Summary on one page

```
WHAT:  Move 135 tests from Gauge (.spec files) to Ginkgo (_test.go files)
WHY:   One language (Go), standard toolchain, easier for contributors
WHERE: Same repo (release-tests), new test/ folder
WHEN:  10 weeks, 4 phases

WHAT DOESN'T CHANGE:
  - The actual tests themselves (same scenarios, same assertions)
  - The pkg/ helper packages
  - The CI pipeline structure
  - Cluster setup / teardown
  - Reporting to Polarion / ReportPortal

WHAT CHANGES:
  - specs/ + steps/ folders → deleted
  - test/ folder → created
  - Run command: "gauge run" → "go test ./test/..."
  - Docker image: gauge binary → ginkgo binary
  - ~25 duplicate-of-upstream tests → removed after confirmation

RISK MITIGATION:
  - Gauge + Ginkgo run in parallel during transition
  - Convert sanity suite first (fastest feedback loop)
  - Each suite has a parity gate before Gauge is dropped
```
