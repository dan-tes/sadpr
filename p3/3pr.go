package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Структура частицы
type particle struct {
	X      []float64 // текущие координаты
	V      []float64 // текущая скорость
	Y      float64   // Значение функции в текущей позиции
	PBest  []float64 // лучшая позиция частицы
	PBestY float64   // значение функции в лучшей позиции
	D      int       // размерность
}

func (p *particle) print() {
	fmt.Println(p.X, p.V, p.Y, p.PBest, p.PBestY, p.D)
}

// Функция, которую минимизируем
func getY(arr []float64) float64 {
	if len(arr) != 2 {
		panic("arr must have length 2")
	}
	x, y := arr[0], arr[1]
	return 0.26*(x*x+y*y) - 0.48*x*y
}

// Создание новой частицы
func newParticle() *particle {
	D := 2
	x := []float64{rand.Float64()*10 - 5, rand.Float64()*10 - 5} // [-5, 5]
	v := []float64{rand.Float64()*2 - 1, rand.Float64()*2 - 1}   // [-1, 1]

	y := getY(x)

	return &particle{
		X:      x,
		V:      v,
		Y:      y,
		PBest:  append([]float64(nil), x...), // копия текущей позиции
		PBestY: y,
		D:      D,
	}
}

// Обновление позиции частицы
func (p *particle) updateX() {
	for i := 0; i < p.D; i++ {
		p.X[i] += p.V[i]
		if p.X[i] > 5 {
			p.X[i] = 5
		}
		if p.X[i] < -5 {
			p.X[i] = -5
		}
	}
}

// Обновление скорости частицы
func (p *particle) updateV(gBest []float64) {
	w := 0.5 // коэффициент инерции
	c1 := rand.Float64()
	c2 := rand.Float64()

	for i := 0; i < p.D; i++ {
		r1 := rand.Float64()
		r2 := rand.Float64()
		p.V[i] = w*p.V[i] +
			c1*r1*(p.PBest[i]-p.X[i]) +
			c2*r2*(gBest[i]-p.X[i])
	}
}

// Генерация роя
func generateSwarm(count int) []*particle {
	swarm := make([]*particle, 0, count)
	for i := 0; i < count; i++ {
		swarm = append(swarm, newParticle())
	}
	return swarm
}

// Основной алгоритм PSO
func getDecision(getYY func([]*particle, int) *particle, numParticles int, maxIter int) ([]float64, float64) {
	swarm := generateSwarm(numParticles)

	// Инициализация глобального минимума
	gBest := make([]float64, 2)
	gBestY := math.MaxFloat64

	for _, p := range swarm {
		if p.Y < gBestY {
			gBestY = p.Y
			copy(gBest, p.X)
		}
	}

	// Основной цикл
	for iter := 0; iter < maxIter; iter++ {
		for i, p := range swarm {
			p.updateV(getYY(swarm, i).X)
			p.updateX()
			p.Y = getY(p.X)
			if p.Y < p.PBestY {
				p.PBestY = p.Y
				copy(p.PBest, p.X)
			}
			if p.Y < gBestY {
				gBestY = p.Y
				copy(gBest, p.X)
			}
		}
		if iter%10 == 0 {
			fmt.Printf("Итерация %3d | лучший Y = %.6f | позиция = [%.4f, %.4f]\n",
				iter, gBestY, gBest[0], gBest[1])
		}

	}

	fmt.Printf("\nЛучшее решение: f(%e, %e) = %e\n", gBest[0], gBest[1], gBestY)
	return gBest, gBestY
}

// Точка входа
func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=== Роевой глобальный алгоритм  ===")
	global := func(a []*particle, _ int) *particle {
		best := a[0]
		for _, p := range a {
			if best.Y > p.Y {
				best = p
			}
		}
		return best
	}
	getDecision(global, 100, 200)
	fmt.Println("=== Роевой локальный алгоритм  ===")
	local := func(a []*particle, current int) *particle {
		best := a[0]
		lambda := 1
		for i := current - lambda; i < current+lambda; i++ {
			cur := min(max(i, 0), len(a)-1)
			if best.Y > a[cur].Y {
				best = a[cur]
			}
		}
		return best
	}
	getDecision(local, 100, 200)
}
