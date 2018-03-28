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
	return planStats(path, func(verb string, args ...string) error {
		cmdArgs := append([]string{}, verb)
		cmdArgs = append(cmdArgs, args...)
		cmd := exec.Command("terraform", cmdArgs...)
		cmd.Dir = path
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %v", output, err)
		}
		return nil
	})
}
func planStats(path string, runTFCommand func(verb string, args ...string) error) (DiffStat, error) {
	planPath := filepath.Join(path, planFileName)
	if err := runTFCommand("init", path); err != nil {
		return DiffStat{}, fmt.Errorf("couldn't run 'terraform init' on %q: %v", path, err)
	}
	cmdLine := []string{"-out", planPath, path}
	if err := runTFCommand("plan", cmdLine...); err != nil {
		return DiffStat{}, fmt.Errorf("couldn't run 'terraform plan %s' on %q: %v", strings.Join(cmdLine, " "), path, err)
	}
	file, err := os.Open(planPath)
	if err != nil {
		return DiffStat{}, fmt.Errorf("couldn't open plan file for reading: %v", err)
	}
	plan, err := tf.ReadPlan(file)
	if err != nil {
		return DiffStat{}, fmt.Errorf("couldn't parse plan file %q: %v", planPath, err)
	}
	return DiffStats(plan), nil
}
