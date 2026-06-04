package skip

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// skipTest skips the test by using Errorf which continues execution but marks as failed.
// Gauge doesn't have proper skip support - tests will appear as failed with skip message.
func skipTest(msg string) {
	// Log the skip message
	log.Printf("SKIPPING TEST: %s", msg)

	// Use Errorf so the test continues but is marked as failed with skip message
	// This is the best we can do with Gauge's limitations
	testsuit.T.Errorf("[SKIPPED] %s", msg)

	// Fail immediately to stop further step execution
	testsuit.T.Fail(fmt.Errorf("[SKIPPED] %s", msg))
}

// GetOCPMinorVersion returns the OCP cluster minor version as an integer (e.g. 20 for 4.20).
// Returns -1 if the version cannot be determined.
func GetOCPMinorVersion(cs *clients.Clients) int {
	cv, err := cs.ClusterVersion.Get(context.Background(), "version", metav1.GetOptions{})
	if err != nil {
		log.Printf("GetOCPMinorVersion: failed to get ClusterVersion: %v", err)
		return -1
	}

	parseMinor := func(version string) int {
		parts := strings.SplitN(version, ".", 3)
		if len(parts) >= 2 {
			if minor, err := strconv.Atoi(parts[1]); err == nil {
				return minor
			}
		}
		return -1
	}

	// Primary: Desired.Version is always set and reflects what the cluster is running/targeting.
	if v := cv.Status.Desired.Version; v != "" {
		if minor := parseMinor(v); minor != -1 {
			return minor
		}
	}

	// Fallback: first Completed history entry.
	for _, h := range cv.Status.History {
		if h.State == "Completed" {
			if minor := parseMinor(h.Version); minor != -1 {
				return minor
			}
		}
	}

	log.Printf("GetOCPMinorVersion: could not parse version from ClusterVersion (desired=%q)", cv.Status.Desired.Version)
	return -1
}

// IfOCPVersionGTE skips the current Gauge test if the cluster OCP minor version
// is greater than or equal to minMinor, with a message referencing a known bug.
//
// Example:
//
//	skip.IfOCPVersionGTE(cs, 20, "SRVKP-11139", "buildah-ns fails on OCP 4.20+")
func IfOCPVersionGTE(cs *clients.Clients, minMinor int, bugID, reason string) {
	minor := GetOCPMinorVersion(cs)

	// Log OCP version detection for debugging
	if minor == -1 {
		log.Printf("Skip check: Could not determine OCP version, will skip conservatively")
	} else {
		log.Printf("Skip check: Detected OCP 4.%d, threshold is 4.%d+ for bug %s", minor, minMinor, bugID)
	}

	if minor == -1 {
		// Could not determine OCP version — skip conservatively to avoid hitting
		// the known bug on an unknown cluster version.
		msg := fmt.Sprintf("[KNOWN BUG: %s] Test skipped - OCP version could not be determined (assuming >= 4.%d to avoid known failure)", bugID, minMinor)
		log.Printf("%s", msg)
		skipTest(msg)
		return
	}

	if minor >= minMinor {
		msg := fmt.Sprintf("[KNOWN BUG: %s] Test skipped on OCP 4.%d+ (cluster is 4.%d) - %s", bugID, minMinor, minor, reason)
		log.Printf("%s", msg)
		skipTest(msg)
	} else {
		log.Printf("Skip check: OCP 4.%d is below threshold 4.%d, test will run", minor, minMinor)
	}
}

// Bug skips the current Gauge test with a clear message referencing a known bug.
// Use this when a test is intentionally skipped due to a tracked product bug.
//
// Example:
//
//	skip.Bug("SRVKP-11139", "buildah-ns fails on OCP 4.20+ due to /proc/0/uid_map")
func Bug(bugID, reason string) {
	msg := fmt.Sprintf("[KNOWN BUG: %s] %s", bugID, reason)
	log.Printf("%s", msg)
	skipTest(msg)
}
