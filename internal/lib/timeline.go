package lib

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type ITimeLine interface {
	AddWork(work *Work)
	RemoveWork(workId uuid.UUID)
	IsBusy(start, end time.Time) bool
	GetWorksInRange(start, end time.Time) []*Work
	CheckOverlapping() bool
}

type TimeLine struct {
	Works []*Work
}

func (timeline *TimeLine) sortWorkByStart() {
	slices.SortFunc(timeline.Works, func(w1, w2 *Work) int {
		if w1.Start.Equal(w2.Start) {
			return 0
		} else if w1.Start.Before(w2.Start) {
			return -1
		} else {
			return 1
		}
	})

}

func (timeline *TimeLine) AddWork(work *Work) {
	timeline.Works = append(timeline.Works, work)
	timeline.sortWorkByStart()
}

func (timeline *TimeLine) RemoveWork(workId uuid.UUID) {
	for i, w := range timeline.Works {
		if *w.ReferenceId == workId || w.Id == workId {
			timeline.Works = append(timeline.Works[:i], timeline.Works[i+1:]...)
		}
	}

	timeline.sortWorkByStart()
}

func (timeline *TimeLine) IsBusy(start, end time.Time) bool {
	for _, work := range timeline.Works {
		if Between(work.Start, start, end) || Between(work.End, start, end) {
			return true
		}
	}
	return false
}

func (timeline *TimeLine) GetWorksInRange(start, end time.Time) []*Work {
	works := make([]*Work, 0)
	for _, work := range timeline.Works {
		if Between(work.Start, start, end) || Between(work.End, start, end) {
			works = append(works, work)
		}
	}
	return works
}

// TODO : Fix this function
func (timeline *TimeLine) CheckOverlapping() bool {
	for i := 0; i < len(timeline.Works); i++ {
		for j := i + 1; j < len(timeline.Works); j++ {
			if timeline.Works[i].Start.Before(timeline.Works[j].End) && timeline.Works[i].End.After(timeline.Works[j].Start) {
				return true
			}
		}
	}

	return false
}

func NewTimeLine() ITimeLine {
	return &TimeLine{
		Works: make([]*Work, 0),
	}
}

// utils

func BeforeOrEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.Before(t2)
}

func AfterOrEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.After(t2)
}

func Between(t, start, end time.Time) bool {
	return AfterOrEqual(t, start) && BeforeOrEqual(t, end)
}
