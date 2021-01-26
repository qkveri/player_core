package playlister

import (
	"errors"
	"fmt"
	"testing"

	"github.com/qkveri/player_core/pkg/domain"
)

func Test_intervalIndexBySeconds(t *testing.T) {
	svc := &service{
		musicData: &domain.MusicData{
			Intervals: []*domain.MusicDataInterval{
				{Start: 18000, End: 39600},
				{Start: 39600, End: 72000},
				{Start: 72000, End: 18000},
			},
		},
	}

	testCases := []struct {
		seconds int
		index   int
	}{
		{18000, 0},
		{39600, 1},
		{78000, 2},
		{0, 2},
		{17000, 2},

		{999999, 0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("seconds_%d", tc.seconds), func(t *testing.T) {
			if i := svc.intervalIndexBySeconds(tc.seconds); i != tc.index {
				t.Errorf("got %d, want %d", i, tc.index)
			}
		})
	}
}

func Test_getTrackIdByIntervalIndex(t *testing.T) {
	t.Run("interval empty", func(t *testing.T) {
		svc := &service{
			musicData: &domain.MusicData{
				Intervals: []*domain.MusicDataInterval{},
			},
		}

		if _, err := svc.getTrackIdByIntervalIndex(0); !errors.Is(err, intervalsEmpty) {
			t.Errorf("got %v, want %v", intervalsEmpty, err)
		}
	})
}
