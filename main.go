package tfinator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tf "github.com/hashicorp/terraform/terraform"
)

const planFileName = "plan.tfplan"

// DiffStat holds statistics on a Terraform diff
type DiffStat struct {
	add     int
	change  int
	destroy int
}

// DiffStats reports statistics on a Terraform plan: the number of
// resources which would be added, changed, or destroyed
func DiffStats(p *tf.Plan) DiffStat {
	s := DiffStat{}
	d := p.Diff
	if d.Empty() {
		return s
	}
	for _, m := range d.Modules {
		for _, rdiff := range m.Resources {
			switch rdiff.ChangeType() {
			case tf.DiffDestroyCreate:
				s.add++
				s.destroy++
			case tf.DiffDestroy:
				s.destroy++
			case tf.DiffCreate:
				s.add++
			default:
				s.change++
			}
		}
	}

	return s
}

// PlanStats runs "terraform plan" in a given directory, reads the resulting
// plan file, and returns a DiffStat representing the generated plan
func PlanStats(path string) (DiffStat, error) {
	return planStats(path, func(args ...string) error {
		output, err := exec.Command(args[0], args[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %v", output, err)
		}
		return nil
	})
}
func planStats(path string, runCommand func(args ...string) error) (DiffStat, error) {
	cmd := []string{"terraform", "plan", "-out", planFileName, path}
	if err := runCommand(cmd...); err != nil {
		return DiffStat{}, fmt.Errorf("couldn't run %q on %q: %v", strings.Join(cmd, " "), path, err)
	}
	file, err := os.Open(filepath.Join(path, planFileName))
	plan, err := tf.ReadPlan(file)
	if err != nil {
		return DiffStat{}, err
	}
	return DiffStats(plan), nil
}
