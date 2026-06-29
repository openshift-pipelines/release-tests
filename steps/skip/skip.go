package skip

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/skip"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

// Step: Skip if OCP >= <minVersion> (Bug: <bugID>)
// Examples:
//   - Skip if OCP >= "4.20" (Bug: "SRVKP-11139")
//   - Skip if OCP >= "20" (Bug: "SRVKP-11139")
var _ = gauge.Step("Skip if OCP >= <minVersion> (Bug: <bugID>)", func(minVersion, bugID string) {
	cs := store.Clients()

	minMinor := parseOCPVersion(minVersion)
	if minMinor == -1 {
		// Could not parse version - skip conservatively with bug reference
		log.Printf("Warning: Could not parse OCP version '%s', skipping test conservatively for bug %s", minVersion, bugID)
		skip.Bug(bugID, fmt.Sprintf("Test skipped on OCP %s+ (version parse failed)", minVersion))
		return
	}

	skip.IfOCPVersionGTE(cs, minMinor, bugID, fmt.Sprintf("Known failure on OCP 4.%d+", minMinor))
})

// parseOCPVersion parses version string in formats: "4.20", "20", etc.
// Returns the minor version number or -1 if parsing fails.
func parseOCPVersion(version string) int {
	version = strings.TrimSpace(version)

	// Try format "4.20" or "4.20.1"
	if strings.HasPrefix(version, "4.") {
		parts := strings.Split(version, ".")
		if len(parts) >= 2 {
			if minor, err := strconv.Atoi(parts[1]); err == nil {
				return minor
			}
		}
	}

	// Try format "20" (just the minor version)
	if minor, err := strconv.Atoi(version); err == nil {
		return minor
	}

	return -1
}
