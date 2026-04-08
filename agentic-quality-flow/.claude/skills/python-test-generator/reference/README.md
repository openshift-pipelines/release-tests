# Tier2 Python/pytest Test Reference Examples

This directory contains **real working tier2 test examples** from RedHatQE/openshift-virtualization-tests repository.

## Purpose

These examples serve as:
1. **Pattern library** - Real code patterns the generator should follow
2. **Quality benchmark** - Generated code should match these patterns
3. **Validation reference** - Compare generated tests against these examples
4. **Documentation** - Show what "good" tier2 tests look like

## Usage

The `python-test-generator` skill references these examples to:
- Validate its pattern detection logic
- Ensure generated code follows RedHatQE standards
- Learn new patterns as more examples are added
- Verify LSP analysis patterns are correctly applied

## Files

Each file is a complete, working test from RedHatQE/openshift-virtualization-tests:

### bgp_connectivity_test.py
- **Source:** RedHatQE/openshift-virtualization-tests PR #2881
- **Pattern:** BGP connectivity tests with CUDN VMs
- **Features:**
  - Module-level pytestmark with multiple markers
  - TCP server/client fixtures
  - Migration test with connectivity validation

**Key Patterns:**
```python
# Pattern 1: Module-level markers
pytestmark = [
    pytest.mark.bgp,
    pytest.mark.ipv4,
    pytest.mark.usefixtures("bgp_setup_ready"),
]

# Pattern 2: Test with migration
def test_connectivity_is_preserved_during_migration(tcp_server, tcp_client):
    migrate_vm_and_verify(vm=tcp_server.vm)
    assert is_tcp_connection(server=tcp_server, client=tcp_client)
```

### flat_overlay_connectivity_test.py
- **Source:** RedHatQE/openshift-virtualization-tests PR #2881
- **Pattern:** Flat overlay network connectivity tests
- **Features:**
  - Class-based test organization
  - MANDATORY Markers section in class docstring
  - Gating tests
  - Negative test patterns
  - Cross-namespace connectivity

**Key Patterns:**
```python
# Pattern 1: Class docstring with MANDATORY Markers section
class TestFlatOverlayConnectivity:
    """
    Tests for flat overlay network connectivity.

    Markers:           # ⬅️ MANDATORY SECTION
        - s390x
        - ipv4

    Preconditions:
        - Multi-network policy enabled
    """

# Pattern 2: Negative test
def test_no_connectivity(self, vm_a, vm_d_ip):
    """
    [NEGATIVE] Test that VMs on separate networks cannot communicate.
    """
    assert_no_ping(src_vm=vm_a, dst_ip=vm_d_ip)

# Pattern 3: Cross-namespace test
def test_connectivity_between_namespaces(self, nad1, nad2, vm_a, vm_e):
    assert nad1.name == nad2.name
    assert_ping_successful(src_vm=vm_a, dst_ip=lookup_iface_status_ip(...))
```

### service_manifest_test.py & service_virtctl_test.py
- **Source:** RedHatQE/openshift-virtualization-tests PR #2881
- **Pattern:** Service configuration tests
- **Features:**
  - Parametrized tests with indirect fixtures
  - IP family policy validation
  - Dual stack service tests

**Key Patterns:**
```python
# Pattern: Parametrization with marks
@pytest.mark.parametrize(
    "single_stack_service_ip_family, single_stack_service",
    [
        pytest.param("IPv4", "IPv4", marks=[pytest.mark.ipv4, pytest.mark.polarion("CNV-5789")]),
        pytest.param("IPv6", "IPv6", marks=[pytest.mark.ipv6, pytest.mark.polarion("CNV-12557")]),
    ],
    indirect=["single_stack_service"],
)
```

### online_resize_test.py
- **Source:** RedHatQE/openshift-virtualization-tests PR #2881
- **Pattern:** Storage online resize tests
- **Features:**
  - Context managers for resource cleanup
  - Sequential and simultaneous disk expansion
  - Snapshot integration
  - Migration after resize

**Key Patterns:**
```python
# Pattern 1: Context manager for resize
with wait_for_resize(vm=vm):
    expand_pvc(dv=dv, size_change=SMALLEST_POSSIBLE_EXPAND)

# Pattern 2: Snapshot workflow
with vm_snapshot(vm=vm, name="snapshot-before") as snap_before:
    expand_pvc(dv=dv, size_change=SMALLEST_POSSIBLE_EXPAND)
    with vm_snapshot(vm=vm, name="snapshot-after") as snap_after:
        with vm_restore(vm=vm, name=snap_before.name) as restored:
            check_file_unchanged(vm=restored)
```

### vm_lifecycle_restart_test.py
- **Source:** RedHatQE/openshift-virtualization-tests PR #2881
- **Pattern:** VM lifecycle tests
- **Features:**
  - Simple VM fixture pattern
  - Context manager usage
  - running_vm() helper

**Key Patterns:**
```python
# Pattern: Simple VM fixture
@pytest.fixture()
def vm_to_restart(unprivileged_client, namespace):
    name = "vm-to-restart"
    with VirtualMachineForTests(
        client=unprivileged_client,
        name=name,
        namespace=namespace.name,
        body=fedora_vm_body(name=name),
    ) as vm:
        running_vm(vm=vm)  # ⬅️ ALWAYS call running_vm
        yield vm
```

### vm_run_strategy_test.py
- **Source:** RedHatQE/openshift-virtualization-tests PR #2881
- **Pattern:** VM run strategy tests
- **Features:**
  - Complex parametrization with indirect fixtures
  - Dictionary-driven test logic
  - ResourceEditor for VM updates
  - Post-upgrade marker

**Key Patterns:**
```python
# Pattern 1: Module-level marker
pytestmark = pytest.mark.post_upgrade

# Pattern 2: Dictionary-driven logic
RUN_STRATEGY_DICT = {
    MANUAL: {
        "start": {"status": True, "run_strategy": MANUAL},
        "stop": {"status": None, "run_strategy": MANUAL},
    },
}

# Pattern 3: ResourceEditor
ResourceEditor(patches={vm: {"spec": {"runStrategy": run_strategy}}}).update()
```

### node_maintenance_test.py
- **Source:** RedHatQE/openshift-virtualization-tests PR #2881
- **Pattern:** Node maintenance and migration tests
- **Features:**
  - Context manager for node operations
  - Multiple VM fixtures with random names
  - Process preservation validation
  - Windows and Linux VM patterns

**Key Patterns:**
```python
# Pattern 1: Context manager for node operations
with node_mgmt_console(node=source_node, node_mgmt="drain"):
    check_migration_process_after_node_drain(client=client, vm=vm)

# Pattern 2: Random VM naming
name = f"vm-nodemaintenance-{random.randrange(99999)}"

# Pattern 3: Process preservation check (Windows)
pre_pid = start_and_fetch_processid_on_windows_vm(vm=vm, process_name=process_name)
# ... drain node ...
post_pid = fetch_pid_from_windows_vm(vm=vm, process_name=process_name)
assert post_pid == pre_pid
```

## Critical Patterns from LSP Analysis

All examples follow these **MANDATORY** patterns from LSP analysis:

1. ✅ **Class docstrings MUST include Markers section**
2. ✅ **Context managers for all resources** (`with` statements)
3. ✅ **running_vm() called after VM creation**
4. ✅ **Module-level pytestmark for common fixtures**
5. ✅ **No global fixture redefinition** in local conftest
6. ✅ **Proper import order** (logging, pytest, utilities, tests)

## Adding More Examples

When you receive new tier2 test examples:

1. Add the file to this directory
2. Update this README with:
   - File name
   - Source location (PR number)
   - Pattern description
   - Key patterns to learn
3. Update `../patterns/pattern_rules.yaml` if new patterns detected
4. Update `../templates/` if new template patterns needed
5. Update `docs/tier2_pattern_library.yaml` if major patterns found

## Structure Reference

```
python-test-generator/
├── SKILL.md                                # Workflow logic only
├── patterns/
│   ├── README.md
│   └── pattern_rules.yaml                  # Template selection
├── templates/
│   ├── localnet_cudn_test.py.template
│   ├── network_connectivity_test.py.template
│   └── ... (more templates)
└── reference/                              # THIS DIRECTORY
    ├── README.md                           # This file
    ├── bgp_connectivity_test.py
    ├── flat_overlay_connectivity_test.py
    ├── service_manifest_test.py
    ├── service_virtctl_test.py
    ├── online_resize_test.py
    ├── vm_lifecycle_restart_test.py
    ├── vm_run_strategy_test.py
    └── node_maintenance_test.py
```

## Quality Criteria

Examples in this directory must be:
- ✅ Complete working tests (pass `pytest --collect-only`)
- ✅ From RedHatQE/openshift-virtualization-tests repository
- ✅ Following RedHatQE coding standards
- ✅ Demonstrating tier2 test patterns validated by LSP analysis
- ✅ Include proper docstrings with Markers sections
