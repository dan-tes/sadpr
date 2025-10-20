package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Function === Твоя функция ===
type Function struct{}

func (f *Function) getY(arr []float64) float64 {
	if len(arr) != 2 {
		panic("arr must have length 2")
	}
	x, y := arr[0], arr[1]
	// Цель — минимизация, поэтому возвращаем отрицание (т.к. getY у тебя с минусом)
	return 0.26*(x*x+y*y) - 0.48*x*y
}

// === Координата (частица) ===
type cord struct {
	x []float64
	y float64
}

func (c cord) print(i int) {
	fmt.Printf("X%d = (", i+1)
	for j := 0; j < len(c.x); j++ {
		fmt.Printf(" %.10f", c.x[j])
	}
	fmt.Print(" )")
	fmt.Printf(" f()=%.25f\n", -c.y)
}

// === Вспомогательные функции ===
func randomInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Матрица вращения для 2D
func rotationMatrix(theta float64) [2][2]float64 {
	return [2][2]float64{
		{math.Cos(theta), -math.Sin(theta)},
		{math.Sin(theta), math.Cos(theta)},
	}
}

// === Алгоритм водоворота ===
func whirlpoolOptimize(
	f *Function,
	bounds [2][2]float64,
	nParticles int,
	maxIter int,
	alpha float64,
) cord {
	rand.Seed(time.Now().UnixNano())

	// Инициализация популяции
	population := make([]cord, nParticles)
	for i := range population {
		x := []float64{
			randomInRange(bounds[0][0], bounds[0][1]),
			randomInRange(bounds[1][0], bounds[1][1]),
		}
		population[i] = cord{x: x, y: f.getY(x)}
	}

	// Поиск лучшего
	best := population[0]
	for _, p := range population {
		if p.y < best.y {
			best = p
		}
	}

	// Основной цикл
	for iter := 0; iter < maxIter; iter++ {
		for i := range population {
			speed := 2.0
			theta := randomInRange(math.Pi*speed*0.75, speed*math.Pi)
			R := rotationMatrix(theta)

			dx := best.x[0] - population[i].x[0]
			dy := best.x[1] - population[i].x[1]

			// Вращение и сжатие
			newX := population[i].x[0] + alpha*(R[0][0]*dx+R[0][1]*dy)
			newY := population[i].x[1] + alpha*(R[1][0]*dx+R[1][1]*dy)

			// Ограничение по границам
			newX = math.Min(math.Max(newX, bounds[0][0]), bounds[0][1])
			newY = math.Min(math.Max(newY, bounds[1][0]), bounds[1][1])

			population[i].x = []float64{newX, newY}
			population[i].y = f.getY(population[i].x)

			// Обновляем лучший
			if population[i].y < best.y {
				best = cord{x: population[i].x, y: population[i].y}
			}
		}
		fmt.Printf("Итерация %d: лучший f() = %f в (%f, %f)\n",
			iter, -best.y, best.x[0], best.x[1])

	}

	return best
}

// === Тест ===
func main() {
	f := &Function{}
	bounds := [2][2]float64{{-5, 5}, {-5, 5}}

	best := whirlpoolOptimize(f, bounds, 30, 100, 0.5)

	fmt.Println("\nЛучшее найденное решение:")
	best.print(0)
}
