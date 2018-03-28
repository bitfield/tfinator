package tfinator

import tf "github.com/hashicorp/terraform/terraform"

// DiffStat holds statistics on a Terraform diff
type DiffStat struct {
	add     int
	change  int
	destroy int
}

// DiffStats reports statistics on a Terraform plan: the number of
// resources which would be added, changed, or destroyed
func DiffStats(p *tf.Plan) (DiffStat, error) {
	s := DiffStat{}
	d := p.Diff
	if d.Empty() {
		return s, nil
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

	return s, nil
}

/*
func PlanDir(path string) (diffStat, error) {
	return diffStat{}, nil
}
*/
