package refactor

import (
	"fmt"
	"strings"
)

// Format returns a human-readable string of the plan.
func Format(p Plan) string {
	if !p.HasChange() {
		return "NO_CHANGE"
	}
	var s strings.Builder
	_, _ = fmt.Fprintf(&s, "manifest:\n")
	_, _ = fmt.Fprintf(&s, "  resources: %s\n", p.Resources)
	_, _ = fmt.Fprintf(&s, "  patchesStrategicMerge: %s\n", p.PatchesStrategicMerge)
	_, _ = fmt.Fprintf(&s, "files:\n")
	for name, resources := range p.Create {
		_, _ = fmt.Fprintf(&s, "  + %s (%d)\n", name, len(resources))
	}
	for _, name := range p.Remove {
		_, _ = fmt.Fprintf(&s, "  - %s\n", name)
	}
	return s.String()
}
