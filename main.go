package tfinator

import tf "github.com/hashicorp/terraform/terraform"

type diffStat struct {
	add     int
	change  int
	destroy int
}

func DiffStats(p tf.Plan) (diffStat, error) {
	s := diffStat{}
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

func PlanDir(path string) (diffStat, error) {
	return diffStat{}, nil
}
