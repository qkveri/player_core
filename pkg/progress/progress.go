package progress

import "fmt"

type Progress float64

const (
	Passed Progress = 1
)

const percentRatio = 0.01

func (p Progress) String() string { return fmt.Sprintf("%d%%", p.Percents()) }
func (p Progress) IsDone() bool   { return p == Passed }
func (p Progress) Percents() int  { return int(p / percentRatio) }
