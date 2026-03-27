# Tier2 Pattern Detection

This directory contains pattern detection configuration for tier2 Python/pytest test generation.

## Files

- **pattern_rules.yaml** - Simplified pattern detection rules (references main library)

## Main Pattern Library

The comprehensive tier2 pattern library is located at:
```
/docs/tier2_pattern_library.yaml
```

This file contains:
- 6 fixture categories (vm_instance, network, localnet, storage, migration, chaos)
- 290+ analyzed test files
- 85+ conftest.py patterns
- Complete LSP analysis from RedHatQE/openshift-virtualization-tests repository

## Pattern Detection Flow

1. Read `docs/tier2_pattern_library.yaml` (comprehensive patterns)
2. Read `patterns/pattern_rules.yaml` (template selection logic)
3. Match scenario against patterns
4. Select appropriate template
5. Generate code using selected template
