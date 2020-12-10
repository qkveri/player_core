package domain

type (
	Ad struct {
		ID    int
		Title string
		Times AdTimes
		Track Track
	}

	AdTimes struct {
		Days  []int
		Times []int
	}
)
