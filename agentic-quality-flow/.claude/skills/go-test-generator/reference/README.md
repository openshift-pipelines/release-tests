# Tier1 Go/Ginkgo Test Reference Examples

This directory contains **real working tier1 test examples** from kubevirt/kubevirt repository.

## Purpose

These examples serve as:
1. **Pattern library** - Real code patterns the generator should follow
2. **Quality benchmark** - Generated code should match these patterns
3. **Validation reference** - Compare generated tests against these examples
4. **Documentation** - Show what "good" tier1 tests look like

## Usage

The `go-test-generator` skill references these examples to:
- Validate its pattern detection logic
- Ensure generated code follows kubevirt standards
- Learn new patterns as more examples are added

## Files

Each file is a complete, working test from kubevirt/kubevirt:

### passt_binding_test.go
- **Source:** kubevirt/kubevirt `tests/network/`
- **Pattern:** Passt network binding plugin tests
- **Features:**
  - DescribeTable with IPv4/IPv6 parametrization
  - Ordered contexts with BeforeAll/AfterAll
  - TCP/UDP connectivity tests
  - Migration tests
  - Multiple VMI setup with client/server pattern

**Key Patterns to Learn:**
```go
// Pattern 1: DescribeTable for IPv4/IPv6
DescribeTable("connectivity", func(ipFamily k8sv1.IPFamily) {
    Entry("[IPv4]", k8sv1.IPv4Protocol),
    Entry("[IPv6]", k8sv1.IPv6Protocol),
})

// Pattern 2: Ordered context with resource sharing
Context("TCP with port specification", Ordered, decorators.OncePerOrderedCleanup, func() {
    BeforeAll(func() {
        // Setup expensive resources once
    })
    // Multiple It blocks share resources
})

// Pattern 3: Helper function extraction
func waitUntilVMIsReady(loginTo console.LoginToFunction, vmis ...*v1.VirtualMachineInstance) {
    for idx, vmi := range vmis {
        *vmis[idx] = *libwait.WaitUntilVMIReady(vmi, loginTo, ...)
    }
}

// Pattern 4: Custom assertions
func assertSourcePodContainersTerminate(labelSelector, fieldSelector string, vmi *v1.VirtualMachineInstance) bool {
    return Eventually(func() k8sv1.PodPhase {
        // Check condition
    }, 30*time.Second).Should(Equal(k8sv1.PodSucceeded))
}
```

## Adding More Examples

When you receive new tier1 test examples from colleagues:

1. Add the file to this directory
2. Update this README with:
   - File name
   - Source location
   - Pattern description
   - Key patterns to learn
3. Update `../patterns/pattern_rules.yaml` if new patterns detected
4. Update `../templates/` if new template patterns needed

## Structure Reference

```
go-test-generator/
├── SKILL.md                          # Workflow logic only
├── patterns/
│   └── pattern_rules.yaml            # Pattern detection rules
├── templates/
│   ├── network_connectivity_test.go.template
│   └── basic_vmi_test.go.template
└── reference/                        # THIS DIRECTORY
    ├── README.md                     # This file
    └── passt_binding_test.go         # Working example
```

## Quality Criteria

Examples in this directory must be:
- ✅ Complete working tests (compile with Bazel)
- ✅ From kubevirt/kubevirt repository
- ✅ Following kubevirt coding standards
- ✅ Demonstrating tier1 test patterns
- ✅ Well-documented with comments
