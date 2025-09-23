package main

import (
	"fmt"
	"math"
	"math/rand"
)

// Интерфейс задачи
type problem interface {
	generateNewDecision(arr []float64, t float64) []float64
	getD() int
	getY(arr []float64) float64
	run(maxIter int, T0 float64, current []float64)
}

type Function struct{}

// Генерация нового решения (распределение Коши)
func (f *Function) generateNewDecision(x []float64, t float64) []float64 {
	newX := make([]float64, len(x))
	copy(newX, x)

	for i := 0; i < len(x); i++ {
		u := rand.Float64() // равномерное [0,1)
		step := t * math.Tan(math.Pi*(u-0.5))
		newX[i] += step
	}
	//newX[0], newX[1] = 10*rand.Float64(), rand.Float64()*10
	return newX
}

// Размерность задачи
func (f *Function) getD() int {
	return 2
}

// Значение функции
func (f *Function) getY(arr []float64) float64 {
	if len(arr) != 2 {
		panic("arr must have length p2")
	}
	x, y := arr[0], arr[1]
	return 0.26*(x*x+y*y) - 0.48*x*y
}
func (f *Function) run(maxIter int, T0 float64, current []float64) {
	best, val, current, currentVal := simulatedAnnealing(maxIter, T0, f, current)
	fmt.Println("Текущее решение:", current, "значение:", currentVal)

	fmt.Println("Лучшее решение:", best, "значение:", val)
}

type CityPerm struct {
	distances [][]int
	cities    []string
}

func NewCityPerm() *CityPerm {
	return &CityPerm{
		distances: [][]int{
			// Москва Суздаль Иваново Владимир Кострома Ярославль Ростов
			{0, 240, 260, 181, 303, 242, 195}, // Москва
			{220, 0, 110, 40, 114, 102, 103},  // Суздаль
			{250, 100, 0, 131, 82, 94, 81},    // Иваново
			{176, 40, 135, 0, 143, 102, 111},  // Владимир
			{300, 110, 85, 141, 0, 97, 99},    // Кострома
			{234, 100, 93, 117, 99, 0, 80},    // Ярославль
			{190, 100, 87, 118, 97, 81, 0},    // Ростов
		},
		cities: []string{"Москва", "Суздаль", "Иваново", "Владимир", "Кострома", "Ярославль", "Ростов"},
	}
}

func (f *CityPerm) generateNewDecision(arr []float64, t float64) []float64 {
	if len(arr) != 8 || arr[0] != 0 || arr[len(arr)-1] != 0 {
		panic("arr must have length 8")
	}
	i, j := 0, 0
	for i == j {
		i, j = rand.Intn(len(arr)-2)+1, rand.Intn(len(arr)-2)+1
	}
	arr[i], arr[j] = arr[j], arr[i]
	return arr
}

func (f *CityPerm) getD() int {
	return 1
}
func (f *CityPerm) getY(arr []float64) float64 {
	sum := 0.0
	if len(arr) != 8 {
		panic("arr must have length 8")
	}
	for i := 1; i < len(arr); i++ {
		past, next := int(arr[i-1]), int(arr[i])
		sum += float64(f.distances[past][next])
	}
	return sum
}

func (f *CityPerm) run(maxIter int, T0 float64, current []float64) {
	best, val, current, currentVal := simulatedAnnealing(maxIter, T0, f, current)
	var bestPos, currentPos []string
	for i := 0; i < len(best); i++ {
		bestPos = append(bestPos, f.cities[int(best[i])])
	}
	for i := 0; i < len(best); i++ {
		currentPos = append(currentPos, f.cities[int(best[i])])
	}
	fmt.Println("Текущее решение:", currentPos, "значение:", currentVal)

	fmt.Println("Лучшее решение:", bestPos, "значение:", val)
}

// Имитация отжига
func simulatedAnnealing(maxIter int, T0 float64, problem problem, current []float64) ([]float64, float64, []float64, float64) {
	// начальное решение
	currentVal := problem.getY(current)

	best := make([]float64, len(current))
	copy(best, current)
	bestVal := currentVal

	for k := 1; k <= maxIter; k++ {
		// охлаждение
		T := T0 / math.Pow(float64(k), 1.0/float64(problem.getD()))

		// новый кандидат
		candidate := problem.generateNewDecision(current, T)
		// ограничиваем в пределах [-p2,p2]

		candidateVal := problem.getY(candidate)
		delta := candidateVal - currentVal

		if delta < 0 {
			current, currentVal = candidate, candidateVal
		} else {
			a, b := rand.Float64(), math.Exp(-delta/T)
			if a < b {
				current, currentVal = candidate, candidateVal
			}
		}
		fmt.Println("Итерация ", k, "из ", maxIter, " Текущее решение", current, "Текущий результат", currentVal)
		// обновляем лучшее найденное
		if currentVal < bestVal {
			copy(best, current)
			bestVal = currentVal
		}
	}

	return best, bestVal, current, currentVal
}

func main() {
	cityPerm := NewCityPerm()
	cityPerm.run(10000, 200, []float64{0, 1, 2, 3, 4, 5, 6, 0})
	function := &Function{}
	function.run(10000, 70, []float64{5, -5})
}
