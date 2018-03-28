package tfinator

import (
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
	addInstancePlan = &tf.Plan{
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
			addInstancePlan,
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
		got, err := DiffStats(test.plan)
		if err != nil {
			t.Errorf("DiffStats(%v): %s", test.planInfo, err)
		}
		if got != test.want {
			t.Errorf("DiffStats(%v) = %+v, want %+v ", test.planInfo, got, test.want)
		}
	}

}

/*
func TestPlanDir(t *testing.T) {
	var tests = []struct {
		dir  string
		want diffStat
	}{
		{"./testdata/nodiff", diffStat{change: 0}},
		{"./testdata/change1", diffStat{change: 1}},
	}

	for _, test := range tests {
		got, err := PlanDir(test.dir)
		if err != nil {
			t.Errorf("PlanDir(%q): %s", test.dir, err)
		}
		if got != test.want {
			t.Errorf("PlanDir(%q) = %v, want %v ", test.dir, got, test.want)
		}
	}
}

*/
