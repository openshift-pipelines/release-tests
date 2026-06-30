---
name: pii-sanitizer
description: Sanitize PII and sensitive data from STP documents
model: claude-opus-4-6
---

# PII Sanitizer Skill

**Phase:** Post-Processing
**User-Invocable:** true

## Purpose

Sanitize Personally Identifiable Information (PII) and sensitive data from STP documents.

## When to Use

- Invoked by **document-formatter** subagent during post-processing
- Can be invoked standalone by users via `/pii-sanitizer`

## Input

```yaml
document: |
  # STP content with potentially sensitive data
  VM 'acme-corp-prod-db-01' on node 'ocp-worker-acme-3.acme.internal'
  failed migration. IP: 10.42.15.87, User: jsmith@acme.com
  PVC: pvc-acme-database-prod, Namespace: acme-production
```

## Output Format

```yaml
sanitized_document: |
  # STP content with sanitized data
  VM 'database-vm' on node 'worker-node-1'
  failed migration. IP: 192.0.2.10, User: testuser@example.com
  PVC: pvc-example, Namespace: test-namespace

sanitization_summary:
  ips_replaced: 3
  hostnames_replaced: 5
  emails_replaced: 2
  customer_names_replaced: 4
  vendor_names_replaced: 1
  credentials_found: 0
  total_replacements: 15
```

## Project-Specific PII Exceptions

Read `{project_context.config_dir}/pii_exceptions.yaml` for project-specific PII exceptions (e.g., allowed vendor/product names). The sanitization rules below are generic and apply to all projects.

## Data Categories to Sanitize

### Customer Information

| Original | Replacement |
|:---------|:------------|
| Customer names | `<customer>`, `Example Corp`, `ACME Inc` |
| Account IDs | `<account-id>` |
| Organization names | `<organization>` |

### User Identifiers

| Original | Replacement |
|:---------|:------------|
| Usernames | `testuser`, `admin-user`, `<username>` |
| Email addresses | `user@example.com`, `admin@example.org` |
| Employee IDs | `<user-id>` |

### Credentials

**NEVER include credentials in output:**
- Passwords → `<password>`
- API keys → `<api-key>`
- Tokens → `<token>`
- Certificates → `<certificate>`
- Secrets → `<secret>`

### Network Information

| Original | Replacement Range |
|:---------|:------------------|
| IP addresses | RFC 5737: `192.0.2.x`, `198.51.100.x`, `203.0.113.x` |
| MAC addresses | `00:00:5E:00:53:xx` |
| Hostnames | `<hostname>`, `worker-node-1`, `master-node-1` |
| FQDNs | `example.com`, `example.org`, `example.net` |

### Infrastructure Names

| Original | Replacement |
|:---------|:------------|
| VM names | `test-vm`, `fedora-vm`, `windows-vm` |
| Pod names | `pod-example` |
| Namespace names | `test-namespace`, `example-namespace` |
| PVC names | `test-pvc`, `pvc-example` |
| Storage classes | `storageclass-example` |
| NIC/Bridge names | `nic-example`, `br-example` |
| Cluster names | `cluster-example` |
| Node names | `node-example`, `worker-node-1` |

### File Paths

| Original | Replacement |
|:---------|:------------|
| `/home/jsmith/...` | `/home/<user>/...` |
| `/data/acme/...` | `/data/<customer>/...` |

### UUIDs

| Original | Replacement |
|:---------|:------------|
| Specific UUIDs | `<uuid>`, `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |

## Vendor Name Sanitization

**Never use specific vendor names (except Red Hat products and open source projects).**

| Vendor Category | Replace With |
|:----------------|:-------------|
| Virtualization (VMware, Hyper-V) | Virtualization Infrastructure Vendor |
| Network (Cisco, Juniper) | Network Infrastructure Vendor |
| Storage (NetApp, Dell EMC) | Storage Infrastructure Vendor |
| Cloud (AWS, Azure, GCP) | Cloud Infrastructure Provider |
| Hardware (Dell, HP) | Hardware Vendor |
| GPU (NVIDIA, AMD GPU) | GPU Vendor |
| NIC (Mellanox, Broadcom) | NIC Vendor |
| Backup (Veeam, Commvault) | Backup/DR Vendor |

**Exceptions (allowed):**
- Red Hat products: RHEL, Fedora, OpenShift, Ansible
- Open source projects: Kubernetes, KubeVirt, CDI, Prometheus
- Technical standards: SR-IOV, NVMe, iSCSI, NFS
- CPU tech references: Intel VT-x, AMD-V

## Sanitization Rules

### IP Address Replacement

1. Identify pattern: `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
2. Check if already documentation range (192.0.2.x, 198.51.100.x, 203.0.113.x)
3. If not, replace with sequential documentation IPs:
   - First IP → 192.0.2.1
   - Second IP → 192.0.2.2
   - etc.

### Hostname Replacement

1. Identify FQDN patterns: `*.company.com`, `*.internal`
2. Replace with generic: `worker-node-1.example.com`
3. Keep role indicators: `worker`, `master`, `compute`

### Email Replacement

1. Identify pattern: `[^@\s]+@[^@\s]+\.[^@\s]+`
2. Replace with: `user@example.com`, `admin@example.org`
3. Preserve role if evident: `admin@...` → `admin@example.com`

## Example Transformation

**Before:**
```
The migration of VM 'customer-prod-db-01' from node
'ocp4-worker-1.acme-corp.internal' (IP: 10.42.15.87) to
'ocp4-worker-2.acme-corp.internal' (IP: 10.42.15.88) failed.

User jsmith@acme-corp.com reported the issue. Using NetApp
Trident for storage provisioning.
```

**After:**
```
The migration of VM 'database-vm' from node
'worker-node-1.example.com' (IP: 192.0.2.1) to
'worker-node-2.example.com' (IP: 192.0.2.2) failed.

User testuser@example.com reported the issue. Using Storage
Provider for storage provisioning.
```

## Verification Checklist

Before returning sanitized document:
- [ ] No real customer names or identifiers
- [ ] No real IP addresses (except RFC 5737 ranges)
- [ ] No real hostnames or FQDNs (except example.com/org/net)
- [ ] No credentials, tokens, or secrets
- [ ] No real usernames or email addresses
- [ ] All infrastructure names are generic
- [ ] No third-party vendor names (except allowed exceptions)
- [ ] All vendor references use generic categories

## When in Doubt

**If uncertain whether data is sensitive: sanitize it.**

It is better to use generic names than risk exposing sensitive information.
