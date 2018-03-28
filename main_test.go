package tfinator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tf "github.com/hashicorp/terraform/terraform"
)

var (
	emptyPlan   = &tf.Plan{}
	addAttrPlan = &tf.Plan{
		Diff: &tf.Diff{
			Modules: []*tf.ModuleDiff{
				&tf.ModuleDiff{
					Resources: map[string]*tf.InstanceDiff{
						"foo": &tf.InstanceDiff{
							Attributes: map[string]*tf.ResourceAttrDiff{
								"foo": &tf.ResourceAttrDiff{
									Old: "",
									New: "bar",
								},
							},
						},
					},
				},
			},
		},
	}
	destroyOnePlan = &tf.Plan{
		Diff: &tf.Diff{
			Modules: []*tf.ModuleDiff{
				&tf.ModuleDiff{
					Resources: map[string]*tf.InstanceDiff{
						"foo": &tf.InstanceDiff{
							Destroy: true,
						},
					},
				},
			},
		},
	}
	requiresNewPlan = &tf.Plan{
		Diff: &tf.Diff{
			Modules: []*tf.ModuleDiff{
				&tf.ModuleDiff{
					Resources: map[string]*tf.InstanceDiff{
						"foo": &tf.InstanceDiff{
							Destroy: true,
							Attributes: map[string]*tf.ResourceAttrDiff{
								"foo": &tf.ResourceAttrDiff{
									Old:         "",
									New:         "bar",
									RequiresNew: true,
								},
							},
						},
					},
				},
			},
		},
	}
	addOnePlan = &tf.Plan{
		Diff: &tf.Diff{
			Modules: []*tf.ModuleDiff{
				&tf.ModuleDiff{
					Resources: map[string]*tf.InstanceDiff{
						"foo": &tf.InstanceDiff{
							Attributes: map[string]*tf.ResourceAttrDiff{
								"foo": &tf.ResourceAttrDiff{
									Old:         "",
									New:         "bar",
									RequiresNew: true,
								},
							},
						},
					},
				},
			},
		},
	}
)

func mockRunCommand(args ...string) error {
	path := args[len(args)-1]
	var plan *tf.Plan
	switch {
	case strings.Contains(path, "add-one"):
		plan = addOnePlan
	default:
		plan = emptyPlan
	}
	outFile, err := os.Create(filepath.Join(path, planFileName))
	if err != nil {
		return err
	}
	if err := tf.WritePlan(plan, outFile); err != nil {
		return err
	}
	return nil
}

func TestDiffStat(t *testing.T) {
	var tests = []struct {
		planInfo string
		plan     *tf.Plan
		want     DiffStat
	}{
		{
			"no diff",
			emptyPlan,
			DiffStat{change: 0},
		},
		{
			"add one",
			addOnePlan,
			DiffStat{add: 1},
		},
		{
			"change one",
			addAttrPlan,
			DiffStat{change: 1},
		},
		{
			"requires new",
			requiresNewPlan,
			DiffStat{add: 1, change: 0, destroy: 1},
		},
		{
			"destroy one",
			destroyOnePlan,
			DiffStat{add: 0, change: 0, destroy: 1},
		},
	}
	for _, test := range tests {
		got := DiffStats(test.plan)
		if got != test.want {
			t.Errorf("DiffStats(%v) = %+v, want %+v ", test.planInfo, got, test.want)
		}
	}

}

func TestPlanStats(t *testing.T) {
	var tests = []struct {
		dir  string
		want DiffStat
	}{
		{"./testdata/nodiff", DiffStat{change: 0}},
		{"./testdata/add-one", DiffStat{add: 1}},
	}

	for _, test := range tests {
		got, err := planStats(test.dir, mockRunCommand)
		if err := os.Remove(filepath.Join(test.dir, planFileName)); err != nil {
			t.Fatalf("failed to remove test plan file %q from %q", planFileName, test.dir)
		}
		if err != nil {
			t.Errorf("planStats(%q): %s", test.dir, err)
		}
		if got != test.want {
			t.Errorf("planStats(%q) = %+v, want %+v ", test.dir, got, test.want)
		}
	}
}
