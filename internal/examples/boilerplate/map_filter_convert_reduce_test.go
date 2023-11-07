package boilerplate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/sort"
)

func Test_Map_Filter_Conver_Reduce(t *testing.T) {

	type Employer struct {
		name   string
		salary int
	}

	employers := map[string][]Employer{
		"internal": {{"Alice", 100}, {"Bob", 90}},
		"external": {{"Chris", 125}},
		"foreing":  {{"Mari", 99}},
	}

	noForeings := filter.Values(employers, func(employers []Employer) bool {
		return slice.Has(employers, func(e Employer) bool { return e.name != "Mari" })
	})

	assert.Equal(t, slice.Of("external", "internal"), sort.Asc(map_.Keys(noForeings)))

	var (
		getSalary                 = func(e Employer) int { return e.salary }
		getDepartmentAndSalarySum = func(department string, e []Employer) (string, int) {
			return department, slice.ConvertAndReduce(e, getSalary, op.Sum[int])
		}
	)

	departmentSalary := map_.Convert(noForeings, getDepartmentAndSalarySum)

	assert.Equal(t, 2, len(departmentSalary))
	assert.Equal(t, 190, departmentSalary["internal"])
	assert.Equal(t, 125, departmentSalary["external"])

}
