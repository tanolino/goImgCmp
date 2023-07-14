package proc

type CompareResult struct {
	jobs []*Job
	cmp  [][]float64
}

func NewCompareResult() *CompareResult {
	r := new(CompareResult)
	r.jobs = make([]*Job, 0)
	r.cmp = make([][]float64, 0)
	return r
}

func (t *CompareResult) Add(job *Job) {
	t.jobs = append(t.jobs, job)

	jobCount := len(t.jobs)
	cmp := make([]float64, jobCount-1)
	for i := 0; i+1 < jobCount; i++ {
		cmp[i] = job.CompareTo(t.jobs[i])
	}
	t.cmp = append(t.cmp, cmp)
}

func (t *CompareResult) AddArray(jobs []*Job) {
	for _, j := range jobs {
		t.Add(j)
	}
}

func (t *CompareResult) Diff(i1 int, i2 int) float64 {
	if i1 == i2 {
		return 1
	} else if i2 < i1 {
		return t.Diff(i2, i1)
	} else {
		return t.cmp[i2][i1]
	}
}

func (t *CompareResult) Len() int {
	return len(t.jobs)
}

func (t *CompareResult) At(i int) *Job {
	return t.jobs[i]
}

func (t *CompareResult) Del(i int) {
	// Remove the actual Job
	t.jobs = append(t.jobs[0:i], t.jobs[i+1:]...)
	t.cmp = append(t.cmp[0:i], t.cmp[i+1:]...)
	// Now remove the old compare results
	for j := i; j < len(t.cmp); j++ {
		t.cmp[j] = append(t.cmp[j][:i], t.cmp[j][i+1:]...)
	}
}
