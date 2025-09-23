package main

import (
	"fmt"
	"math"
	"math/rand"
)

type particle struct {
	X []float64
	y float64
	v []float64
	D int
}

// Значение функции
func getY(arr []float64) float64 {
	if len(arr) != 2 {
		panic("arr must have length p2")
	}
	x, y := arr[0], arr[1]
	return 0.26*(x*x+y*y) - 0.48*x*y
}

func newParticle() *particle {
	D := 2
	generate := func(getNum func() float64) []float64 {
		arr := make([]float64, 0)
		for i := 0; i < D; i++ {
			arr = append(arr, getNum())
		}
		return arr
	}
	x, v := generate(func() float64 { return rand.Float64() * 10 }), generate(func() float64 { return 0.0 })
	return &particle{
		X: x,
		v: v,
		D: D,
		y: math.MaxFloat64,
	}
}
func (p *particle) updateX() {
	for i := 0; i < p.D; i++ {
		p.X[i] = math.Min(math.Max(p.X[i]+p.v[i], -10), 10)
	}
}
func (p *particle) updateV(bestY float64) {
	for i := 0; i < p.D; i++ {
		r1, r2 := rand.Float64(), rand.Float64()
		p.v[i] += 2*r1*(p.y-p.X[i]) + 2*r2*(bestY-p.X[i])
	}
}
func generateSwarm() []*particle {
	arr := make([]*particle, 0)
	for i := 0; i < 100; i++ {
		arr = append(arr, newParticle())
	}
	return arr
}

func main() {
	rand.Seed(64)
	global := func(bestY float64, j int, arr []*particle) float64 {
		//fmt.Println(bestY)
		return bestY
	}
	local := func(bestY float64, j int, arr []*particle) float64 {
		localMin := math.MaxFloat64
		num := 1
		for i := j - num; i < j+num; i++ {
			var ind int
			if i < 0 {
				ind = len(arr) + i - 1
			} else {
				ind = i % len(arr)
			}
			localMin = math.Min(localMin, arr[ind].y)
		}
		return localMin
	}
	getDecision := func(updateV func(bestY float64, j int, arr []*particle) float64) ([]float64, float64) {
		maxIter := 100
		var particles = generateSwarm()
		bestX, bestY := make([]float64, 2), math.MaxFloat64
		for i := 0; i < maxIter; i++ {
			for j := range particles {
				if i%5 == 0 {
					fmt.Print("Итерация № ", i)
					st := ""
					for _, x := range particles[j].X {
						st += fmt.Sprintf(", %.2f", x) // округление до 2 знаков после запятой
					}
					st1 := ""
					for _, x := range particles[j].v {
						st1 += fmt.Sprintf(", %.2f", x) // округление до 2 знаков после запятой
					}
					fmt.Printf(" Частица %d координаты: %s скорость %s, Y: %.2f \n", j, st[2:], st1[2:], getY(particles[j].X))
				}
				if getY(particles[j].X) < particles[j].y {
					particles[j].y = getY(particles[j].X)
					if bestY > particles[j].y {
						bestY = particles[j].y
						copy(bestX, particles[j].X)
					}
				}
			}
			for j := range particles {
				particles[j].updateV(updateV(bestY, j, particles))
				particles[j].updateX()
			}
		}

		fmt.Println("Лучшее решение ", bestX, ":", bestY)

		return bestX, bestY
	}
	fmt.Println("Глобальный роевой алгоритм")
	getDecision(global)
	fmt.Println("Локальный роевой алгоритм")
	getDecision(local)

}
