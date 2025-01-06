package util

import (
	"strconv"
)

func Map[TSource, TTarget any](source []TSource, trans func(TSource) TTarget) []TTarget {
	target := make([]TTarget, 0, len(source))
	for _, s := range source {
		target = append(target, trans(s))
	}
	return target
}

func MapStringToInt(source []string) []int {
	target := make([]int, 0, len(source))
	for _, s := range source {
		v, _ := strconv.Atoi(s)
		target = append(target, v)
	}
	return target
}
