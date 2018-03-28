package tfinator

import (
	"testing"

	tf "github.com/hashicorp/terraform/terraform"
)

var (
	emptyPlan   tf.Plan = tf.Plan{}
	addAttrPlan tf.Plan = tf.Plan{
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
	requiresNewPlan tf.Plan = tf.Plan{
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
		plan     tf.Plan
		want     diffStat
	}{
		{
			"empty plan",
			emptyPlan,
			diffStat{change: 0},
		},
		{
			"add bar attribute",
			addAttrPlan,
			diffStat{change: 1},
		},
		{
			"requires new",
			requiresNewPlan,
			diffStat{add: 1, change: 0, destroy: 1},
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
