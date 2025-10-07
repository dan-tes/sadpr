package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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
func (c cord) print(i int) {
	fmt.Printf("X%d = (", i+1)
	for j := 0; j < len(c.x); j++ {
		fmt.Printf(" %.4f", c.x[j])
	}
	fmt.Print(")")
	fmt.Printf("f()=%.4f\n", -c.y)
}
func distance(a, b []float64) float64 {
	//fmt.Printf("a = %.4f, b = %.4f |", a, b)
	sum := 0.0
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	//fmt.Printf("%.4f\n", math.Sqrt(sum))
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
func show(allPoints []*cord, i int) {
	// Визуализация точек
	p := plot.New()
	p.Title.Text = "Визуализация точек"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	pts := make(plotter.XYs, len(allPoints))
	for i, pt := range allPoints {
		pts[i].X = pt.x[0]
		pts[i].Y = pt.x[1]
	}

	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	p.Add(s)

	if err := p.Save(10*vg.Inch, 10*vg.Inch, "points"+strconv.Itoa(i)+".png"); err != nil {
		panic(err)
	}
}
func (f *Function) run(tau int, r *rand.Rand) *cord {
	S := 10
	n := 2
	m := 2
	N := 5
	M := 3
	deltaFind, deltaSpawn := 0.5, 0.5

	// Шаг 1: разведчики
	var scouts []*cord
	for i := 0; i < S; i++ {
		x := []float64{20 * (r.Float64() - 0.5), 20 * (r.Float64() - 0.5)}
		scout := &cord{x: x, y: f.getY(x)}
		scouts = append(scouts, scout)
		//scout.print(i)
	}
	scouts = sorted(scouts)

	allRegions := scouts
	merged := make([]bool, len(allRegions))
	var regions []*cord
	for i := 0; i < len(allRegions); i++ {
		if merged[i] {
			continue
		}
		center := allRegions[i]
		for j := i + 1; j < len(allRegions); j++ {
			//fmt.Print(i, j)
			if distance(allRegions[i].x, allRegions[j].x) < deltaFind {
				if allRegions[j].y > center.y {
					center = allRegions[j]
				}
				merged[j] = true
			}
		}
		regions = append(regions, center)
		if len(regions) >= n+m {
			break
		}
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
				newX := f.localSearch(region.x, deltaSpawn, r)
				c := &cord{x: newX, y: f.getY(newX)}
				if c.y > region.y {
					region = c
				}
			}
			newRegions = append(newRegions, region)
		}

		newRegions = sorted(newRegions)
		fmt.Println("Центры регионов")
		for i, region := range newRegions {
			region.print(i)
		}
		if newRegions[0].y > best.y {
			best = newRegions[0]
			bestCount = 0
		} else {
			bestCount++
			if bestCount >= tau {
				break
			}
		}
		deltaSpawn = deltaFind * ((30000 - float64(iter)) / 30000)
		//fmt.Println(deltaSpawn)
		regions = newRegions
		fmt.Printf("Итерация %d, текущий лучший: f=%.8f в (%.4f, %.4f)\n", iter, -best.y, best.x[0], best.x[1])
		//show(regions, iter)
	}
	return best
}

func main() {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	f := Function{}
	best := f.run(10000, r)
	fmt.Printf("Лучшее решение: x=%.4f y=%.4f → f=%.13f\n", best.x[0], best.x[1], -best.y)
}
