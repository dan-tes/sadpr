package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type cord struct {
	x []float64
	y float64
}

type Function struct{}

func (f *Function) getY(arr []float64) float64 {
	if len(arr) != 2 {
		panic("arr must have length 2")
	}
	x, y := arr[0], arr[1]
	return -(0.26*(x*x+y*y) - 0.48*x*y)
}

func distance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

func (f *Function) localSearch(center []float64, delta float64, r *rand.Rand) []float64 {
	newX := make([]float64, len(center))
	for i := range center {
		newX[i] = center[i] + (r.Float64()*2-1)*delta
	}
	return newX
}

func sorted(arr []*cord) []*cord {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].y > arr[j].y
	})
	return arr
}

func (f *Function) run(tau int, r *rand.Rand) *cord {
	S := 10
	n := 2
	m := 2
	N := 5
	M := 3
	delta := 0.5

	// Шаг 1: разведчики
	var scouts []*cord
	for i := 0; i < S; i++ {
		x := []float64{20 * (r.Float64() - 0.5), 20 * (r.Float64() - 0.5)}
		scouts = append(scouts, &cord{x: x, y: f.getY(x)})
	}
	scouts = sorted(scouts)

	// Шаг 2: выделяем лучшие и перспективные области
	bestRegions := scouts[:n]
	promRegions := scouts[n : n+m]

	// Шаг 3: объединение близких областей
	allRegions := append([]*cord{}, bestRegions...)
	allRegions = append(allRegions, promRegions...)
	merged := make([]bool, len(allRegions))
	var regions []*cord
	for i := 0; i < len(allRegions); i++ {
		if merged[i] {
			continue
		}
		center := allRegions[i]
		for j := i + 1; j < len(allRegions); j++ {
			if distance(allRegions[i].x, allRegions[j].x) < delta {
				if allRegions[j].y > center.y {
					center = allRegions[j]
				}
				merged[j] = true
			}
		}
		regions = append(regions, center)
	}

	best := regions[0]
	bestCount := 0

	for iter := 0; ; iter++ {
		var newRegions []*cord

		for i, region := range regions {
			count := M
			if i < n {
				count = N
			}
			for k := 0; k < count; k++ {
				newX := f.localSearch(region.x, delta, r)
				c := &cord{x: newX, y: f.getY(newX)}
				if c.y > region.y {
					region = c
				}
			}
			newRegions = append(newRegions, region)
		}

		newRegions = sorted(newRegions)

		if newRegions[0].y > best.y {
			best = newRegions[0]
			bestCount = 0
		} else {
			bestCount++
			if bestCount >= tau {
				break
			}
		}

		regions = newRegions
		fmt.Printf("Итерация %d, текущий лучший: f=%.4f в (%.4f, %.4f)\n", iter, best.y, best.x[0], best.x[1])

	}

	return best
}

func main() {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	f := Function{}
	best := f.run(5, r)
	fmt.Printf("Лучшее решение: x=%.4f y=%.4f → f=%.4f\n", best.x[0], best.x[1], best.y)
}
