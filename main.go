package tfinator

import (
	"fmt"
	"os"
	"path/filepath"

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
func PlanStats(path string, runCommand func(args ...string) error) (DiffStat, error) {
	if err := runCommand("terraform", "plan", path, "-out", planFileName); err != nil {
		return DiffStat{}, fmt.Errorf("couldn't run 'terraform plan' on %q: %v", path, err)
	}
	file, err := os.Open(filepath.Join(path, planFileName))
	plan, err := tf.ReadPlan(file)
	if err != nil {
		return DiffStat{}, err
	}
	return DiffStats(plan), nil
}
