# QualityFlow Configuration Guide

This directory contains the multi-project configuration system for QualityFlow.
Each project gets its own subdirectory under `projects/` with YAML files that
control every aspect of test plan generation, test code generation, and
pipeline behavior.

## How Config Loading Works

Every QualityFlow command (`/stp-builder`, `/std-builder`, `/generate-go-tests`,
`/generate-python-tests`) invokes the **project-resolver** skill as Step 0:

1. Parse the Jira ID to extract the prefix (e.g., `CNV` from `CNV-66855`)
2. Look up the prefix in `routing.yaml` to find the project (e.g., `cnv`)
3. Load `_defaults.yaml` (shared defaults for all projects)
4. Load `projects/<project>/project.yaml` (project-specific overrides)
5. Merge feature toggles (project values override defaults)
6. Return `project_context` with `config_dir`, `feature_toggles`, and identity

Agents then read only the config files they need from the resolved `config_dir`.

## Directory Structure

```text
config/
  _schema.yaml                          # Validation rules for project configs
  _defaults.yaml                        # Shared defaults (all projects inherit)
  routing.yaml                          # Jira prefix -> project routing
  projects/
    cnv/                                # OpenShift Virtualization
      project.yaml                      # Identity, toggles, scope boundaries
      repositories.yaml                 # Repos, orgs, build system
      components.yaml                   # Component -> package mappings
      jira.yaml                         # Jira instance config
      tier1.yaml                        # Go/Ginkgo config
      tier2.yaml                        # Python/pytest config
      environment.yaml                  # Platform requirements
      pii_exceptions.yaml               # PII allowlist
      patterns/                         # Pattern detection rules
        tier1_patterns.yaml             # Go/Ginkgo code patterns
        tier2_patterns.yaml             # Python/pytest code patterns
      reference/                        # Reference test files
        tier1/                          # Example Go test files
        tier2/                          # Example Python test files
      templates/                        # Code/document templates
        stp/                            # STP document templates
        std/                            # STD YAML templates
        tier1/                          # Go test file templates
        tier2/                          # Python test file templates
    mtv/                                # Migration Toolkit for Virtualization
      ...                               # Same structure (no tier1.yaml)
```

## Adding a New Project

### Step 1: Create the project directory

```bash
mkdir -p config/projects/<name>
```

### Step 2: Add routes in `routing.yaml`

Add one or more prefix entries that map Jira issue prefixes to your project:

```yaml
routes:
  - prefix: "MYPROJ"
    project: "<name>"
  - prefix: "MYBUGS"
    project: "<name>"
```

### Step 3: Create required YAML files

The schema (`_schema.yaml`) requires these files in every project directory:

| File | Purpose |
|------|---------|
| `project.yaml` | Core identity, feature toggles, scope boundaries |
| `repositories.yaml` | Repository locations and build configuration |
| `components.yaml` | Component-to-package mappings for code analysis |
| `jira.yaml` | Jira instance URL, prefixes, custom fields |
| `environment.yaml` | Platform, cluster, operator requirements |
| `pii_exceptions.yaml` | Allowed names and vendor replacements |

### Step 4: Create optional files based on feature toggles

| File | Required when |
|------|---------------|
| `tier1.yaml` | `feature_toggles.tier1_tests: true` |
| `tier2.yaml` | `feature_toggles.tier2_tests: true` |

### Step 5: Add optional directories

These directories are optional but recommended:

- `patterns/` -- Pattern YAML files for code generation
- `reference/` -- Example test files the generators learn from
- `templates/` -- STP/STD/test file templates

### Step 6: Deploy and test

```bash
uv run deploy.py --target both
# Then run a command against a Jira ID with your new prefix
```

## File Reference

### routing.yaml

Maps Jira issue prefixes to project directories.

```yaml
version: "1.0"

routes:
  - prefix: "CNV"           # Jira prefix to match
    project: "cnv"           # Directory name under projects/

default_project: null        # null = fail on unknown prefix
                             # Set to a project ID for fallback
```

Multiple prefixes can route to the same project. For example, CNV routes
`CNV`, `VIRTSTRAT`, and `OCPBUGS` all to the `cnv` project.

### project.yaml

Core identity and behavior configuration.

```yaml
project_id: "cnv"                                    # Must match directory name
display_name: "OpenShift Virtualization (CNV)"       # Human-readable name
description: "KubeVirt-based virtualization for OpenShift"
```

**feature_toggles** -- Override defaults from `_defaults.yaml`:

```yaml
feature_toggles:
  polarion: true             # Include Polarion markers in generated tests
  tier1_tests: true          # Enable Go/Ginkgo test generation
  tier2_tests: true          # Enable Python/pytest test generation
```

See [Feature Toggles Reference](#feature-toggles-reference) for all toggles.

**versioning** -- Product and platform version strings used in STP documents:

```yaml
versioning:
  product_name: "OpenShift Virtualization"
  platform_name: "OCP"
  current_version: "4.22"
```

**stp_document** -- STP document header text:

```yaml
stp_document:
  header: "Openshift-virtualization-tests Test plan"
```

**sig_mappings** -- Map SIG labels to feature areas (used for test organization):

```yaml
sig_mappings:
  sig-network: "network"
  sig-compute: "compute"
  sig-storage: "storage"
```

**decorator_mappings** -- Map SIG labels to Go decorator calls (tier 1 only):

```yaml
decorator_mappings:
  sig-network: "decorators.SigNetwork"
  sig-compute: "decorators.SigCompute"
```

**scope_boundaries** -- Define what is in/out of scope for this project:

```yaml
scope_boundaries:
  validation_gate: "Would removing KubeVirt make this test meaningless?"
  in_scope_resources:       # Kubernetes resources this project owns
    - "VirtualMachine"
    - "VirtualMachineInstance"
  out_of_scope_if_only:     # Generic resources (only in scope if combined with above)
    - "Pod"
    - "Deployment"
  cli_tools:                # CLI tools relevant to this project
    - "virtctl"
    - "oc"
  domain_vocabulary:        # Domain-specific terms and abbreviations
    - "VMI"
    - "NAD"
```

### repositories.yaml

Repository locations for code analysis and test generation.

**primary_repo** (required) -- The main source code repository:

```yaml
primary_repo:
  name: "kubevirt"                     # Repository name
  org: "kubevirt"                      # GitHub organization
  full_name: "kubevirt/kubevirt"       # org/name
  url: "https://github.com/kubevirt/kubevirt"
  local_path_env: "KUBEVIRT_REPO_PATH" # Env var pointing to local clone
  default_branch: "main"
  language: "go"                       # Primary language
  build_system: "bazel"               # Build system (bazel, make, etc.)
  build_command: "bazel test //tests/{package}:go_default_test"
```

**tier2_repo** (optional) -- Separate repository for tier 2 tests:

```yaml
tier2_repo:
  name: "openshift-virtualization-tests"
  org: "RedHatQE"
  full_name: "RedHatQE/openshift-virtualization-tests"
  url: "https://github.com/RedHatQE/openshift-virtualization-tests"
  local_path_env: "TIER2_REPO_PATH"
  default_branch: "main"
  language: "python"
```

**additional_repos** -- Other repositories for cross-repo analysis:

```yaml
additional_repos:
  - name: "containerized-data-importer"
    org: "kubevirt"
    full_name: "kubevirt/containerized-data-importer"
```

**pr_url_patterns** -- URL patterns for matching PR links:

```yaml
pr_url_patterns:
  - "https://github.com/{org}/{repo}/pull/{number}"
```

### components.yaml

Maps source code components to package paths and features. Used by the
regression-analyzer and code generation agents to understand project structure.

**component_package_map** -- Each key is a component name:

```yaml
component_package_map:
  virt-controller:
    package_path: "pkg/virt-controller/"     # Root path in the repo
    features:                                 # Features within this component
      - { name: "VM Lifecycle", path: "pkg/virt-controller/vm/" }
      - { name: "Live Migration", path: "pkg/virt-controller/migration/" }
```

**path_to_feature** -- Reverse mapping from file paths to feature names:

```yaml
path_to_feature:
  "pkg/virt-controller/vm/": "VM Lifecycle"
  "pkg/network/": "Networking"
```

### jira.yaml

Jira instance configuration for the jira-collector agent.

```yaml
instance:
  url: "https://issues.redhat.com"              # Jira instance URL
  browse_pattern: "https://issues.redhat.com/browse/{key}"

prefixes:                                        # Valid Jira prefixes
  - "CNV"
  - "VIRTSTRAT"

custom_fields:                                   # Custom field names
  feature_link: "Feature Link"
  git_pull_request: "Git Pull Request"

pr_url_scan_pattern: "https://github.com/.*/pull/\\d+"  # Regex for PR URLs
```

Note: MTV uses a different Jira instance (`https://redhat.atlassian.net`) and
includes a `git_pull_request_field_id` for direct field lookup.

### tier1.yaml

Go/Ginkgo test generation configuration. Only required when
`feature_toggles.tier1_tests` is `true`.

```yaml
enabled: true
language: "go"
framework: "ginkgo-v2"
default_package: "network"         # Default test package
import_base: "kubevirt.io/kubevirt" # Base import path
```

**imports** -- Organized by category for generated test files:

| Category | Purpose |
|----------|---------|
| `dot_imports` | Ginkgo/Gomega (imported with `.`) |
| `standard` | Go standard library (`context`, `time`) |
| `project_api` | Project API types with aliases |
| `k8s_core` | Kubernetes core types with aliases |
| `test_framework` | Test framework packages (decorators, testsuite) |
| `network` | Network-specific imports |

**helper_libraries** -- Test helper packages the generator can use:

```yaml
helper_libraries:
  libvmifact: "kubevirt.io/kubevirt/tests/libvmifact"
  libnet: "kubevirt.io/kubevirt/tests/libnet"
  libwait: "kubevirt.io/kubevirt/tests/libwait"
```

**os_patterns** -- VM creation and login patterns per guest OS:

```yaml
os_patterns:
  fedora:
    vm_creator: "libvmifact.NewFedora"
    login: "console.LoginToFedora"
```

**timeout_constants** -- Named timeout constants available in the test framework:

```yaml
timeout_constants:
  tiny: "StartupTimeoutSecondsTiny"
  medium: "StartupTimeoutSecondsMedium"
  migration: "MigrationWaitTime"
```

**nad_patterns** -- Network Attachment Definition creation helpers:

```yaml
nad_patterns:
  localnet: "libnet.NewPasstNetAttachDef"
  bridge: "libnet.NewBridgeNetAttachDef"
  sriov: "libnet.NewSriovNetAttachDef"
```

**context_init** -- Statements to initialize test context (added to `BeforeEach`):

```yaml
context_init:
  - statement: "ctx := context.Background()"
    variable: "ctx"
    type: "context.Context"
  - statement: "namespace := testsuite.GetTestNamespace(nil)"
    variable: "namespace"
    type: "string"
```

### tier2.yaml

Python/pytest test generation configuration. Only required when
`feature_toggles.tier2_tests` is `true`.

```yaml
enabled: true
language: "python"
framework: "pytest"
```

**python_packages** -- Python packages used by generated tests:

```yaml
python_packages:
  wrapper: "openshift-python-wrapper"
  ocp_resources: "ocp_resources"
```

**import_patterns** -- Organized by category:

| Category | Purpose |
|----------|---------|
| `standard` | `pytest`, `logging` |
| `utilities_network` | Network helper functions |
| `utilities_virt` | VM creation and lifecycle |
| `utilities_storage` | Storage helpers |
| `utilities_constants` | Timeout and type constants |

**polarion** -- Polarion test case marker configuration:

```yaml
polarion:
  marker_format: 'pytest.mark.polarion("{ID}")'
  id_prefix: "CNV-"
```

Set `enabled: false` for projects that don't use Polarion (e.g., MTV).

**global_fixtures** -- pytest fixtures available in all test files:

```yaml
global_fixtures:
  - "admin_client"
  - "unprivileged_client"
```

**pytest_markers** -- Custom pytest markers (MTV example):

```yaml
pytest_markers:
  - "tier0"
  - "warm"
  - "copyoffload"
```

**provider_types**, **plan_statuses** -- Project-specific enums (MTV example):

```yaml
provider_types:
  - "VSPHERE"
  - "RHV"
  - "OPENSTACK"
```

### environment.yaml

Platform and infrastructure requirements for test execution.

```yaml
platform:
  name: "OpenShift Container Platform"
  short_name: "OCP"
  cli_tools:
    - "oc"
    - "kubectl"
    - "virtctl"

cluster_requirements:
  topology: "Multi-node"
  cpu_virtualization: "Standard"
  min_worker_nodes: 2

version_constraints:
  ocp: "4.22+"
  cnv: "4.22+"

operators:
  - name: "OpenShift Virtualization"
    namespace: "openshift-cnv"
    csv_pattern: "kubevirt-hyperconverged"

network:
  cni: "OVN-Kubernetes"
  multus: true
  sriov_operator: "optional"

storage:
  default_class: "ocs-storagecluster-ceph-rbd"
  block_storage: true
  shared_storage: true
```

### pii_exceptions.yaml

Controls PII sanitization behavior. Names listed here are allowed in generated
documents without replacement.

```yaml
allowed_product_names:       # Product/vendor names that are not PII
  - "Red Hat"
  - "OpenShift"

allowed_project_names:       # Open source project names
  - "Kubernetes"
  - "KubeVirt"

allowed_technical_standards: # Technical standards and protocols
  - "SR-IOV"
  - "NVMe"

vendor_replacements:         # Generic replacements for actual vendor names
  gpu: "GPU Vendor"
  nic: "NIC Vendor"
  storage: "Storage Infrastructure Vendor"
```

## Feature Toggles Reference

Feature toggles are defined in `_defaults.yaml` and can be overridden per
project in `project.yaml`. Project values take precedence.

| Toggle | Default | Effect when `true` | Effect when `false` |
|--------|---------|-------------------|---------------------|
| `polarion` | `false` | Include Polarion markers in generated test stubs and tests | Omit Polarion markers |
| `unit_tests` | `false` | Informational only | Informational only |
| `tier1_tests` | `true` | Enable `/generate-go-tests`, include Go stubs in `/std-builder` | Block Go test generation |
| `tier2_tests` | `true` | Enable `/generate-python-tests`, include Python stubs in `/std-builder` | Block Python test generation |
| `stp_generation` | `true` | Enable `/stp-builder` | Block `/stp-builder` with early exit |
| `std_generation` | `true` | Enable `/std-builder` | Block `/std-builder` with early exit |
| `lsp_analysis` | `true` | Run regression-analyzer in STP pipeline, run ticket-context-analyzer in code generation | Skip LSP-based analysis |
| `pii_sanitization` | `true` | Run pii-sanitizer in document-formatter | Skip PII sanitization |

**Example:** MTV disables `tier1_tests` (Python-only project) and `polarion`
(uses pytest-jira instead). CNV enables `polarion` and both test tiers.

## Schema Validation

`_schema.yaml` defines validation rules that the project-resolver checks:

- **Required files** -- Every project must have `project.yaml`,
  `repositories.yaml`, `components.yaml`, `jira.yaml`, `environment.yaml`,
  and `pii_exceptions.yaml`
- **Optional files** -- `tier1.yaml` and `tier2.yaml` are only required when
  their corresponding feature toggle is `true`
- **Required fields** -- Each YAML file has required fields (e.g.,
  `project.yaml` must have `project_id` and `display_name`)
- **Toggle consistency** -- If `tier1_tests` is `true`, `tier1.yaml` must exist
  (and likewise for `tier2_tests` / `tier2.yaml`)

## Defaults Inheritance

`_defaults.yaml` provides shared defaults inherited by all projects:

- **feature_toggles** -- Default toggle values (projects override individual
  toggles, unset toggles inherit the default)
- **output_structure** -- File path patterns for all output types (STP, STD,
  stubs, tests) using `{JIRA_ID}` and `{feature}` placeholders
- **test_id_format** -- Pattern for test scenario IDs
  (`TS-{JIRA_ID}-{NUM:03d}`)
- **pii_rules** -- Default PII replacement values (customer names, IPs,
  hostnames, domains)
- **stp_defaults** -- Default STP document settings

Projects do not need to redefine these values unless they want different
behavior.

## Optional Directories

### patterns/

Contains YAML files with code patterns for the test generators:

- `tier1_patterns.yaml` -- Go/Ginkgo patterns (VM creation, assertions,
  network setup)
- `tier2_patterns.yaml` -- Python/pytest patterns (fixtures, assertions,
  resource creation)

Fresh LSP patterns extracted at runtime take priority over these historical
patterns.

### reference/

Contains example test files that generators use as style references:

- `tier1/` -- Example Go/Ginkgo test files
- `tier2/` -- Example Python/pytest test files

Each subdirectory should include a `README.md` explaining what the reference
files demonstrate.

### templates/

Contains templates for document and code generation:

- `stp/` -- STP markdown templates (`stp-template.md`,
  `section-requirements.md`)
- `std/` -- STD YAML templates (`std_template.yaml`,
  `std_template_comprehensive.yaml`, `supplemental.md`)
- `tier1/` -- Go test file templates (`.go.template` files)
- `tier2/` -- Python test file templates (`.py.template` files,
  `fixture_patterns.yaml`, `pytest_template.py`)
