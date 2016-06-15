package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type bucket struct {
	data []int
}

func newBucket() *bucket {
	return &bucket{
		data: make([]int, 64),
	}
}

func (l *bucket) Record(d time.Duration) {
	order := 64
	for ; order >= 0; order-- {
		if int64(d)&(1<<uint(order)) != 0 {
			break
		}
	}
	l.data[order] += 1
}

func (l *bucket) String() string {
	total := 0
	max := 0
	for i := 0; i < 63; i++ {
		total += l.data[i]
		if l.data[i] > max {
			max = l.data[i]
		}
	}

	maxWidth := 20
	var r []string

	last := 0
	for i := 0; i < 63; i++ {
		width := maxWidth * l.data[i] / max
		// skip first
		if width == 0 && len(r) == 0 {
			continue
		}
		ns := math.Pow(2, float64(i)+1.0)
		order := math.Pow(10, math.Floor(math.Log10(ns)))
		dur := time.Duration(math.Floor(0.5+ns/order)*order) * time.Nanosecond
		r = append(r, fmt.Sprintf("%5v  %s\n", dur, strings.Repeat("X", width)))
		if width > 0 {
			last = len(r)
		}
	}
	return strings.Join(r[0:last], "")
}
