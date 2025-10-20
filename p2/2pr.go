package main

import (
	"fmt"
	"math"
	"math/rand"
)

// Интерфейс задачи
type problem interface {
	generateNewDecision(arr []float64, t float64, iter int) []float64
	getD() int
	getY(arr []float64) float64
	run(maxIter int, T0 float64, current []float64)
}

type Function struct{}

// Генерация нового решения (распределение Коши)
func (f *Function) generateNewDecision(x []float64, t float64, k int) []float64 {
	newX := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		u := rand.Float64()*10 - 5
		newX[i] = (t / math.Pow(math.Pow(u-x[i], 2)+math.Pow(t, 2), 3/2)) / math.Pow(math.Pi, 2)
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
			{0, 240, 360, 181, 303, 242, 195}, // Москва
			{240, 0, 110, 40, 114, 12, 103},   // Суздаль
			{360, 110, 0, 131, 82, 94, 81},    // Иваново
			{181, 40, 131, 0, 143, 102, 111},  // Владимир
			{303, 114, 82, 143, 0, 97, 99},    // Кострома
			{242, 12, 94, 102, 97, 0, 180},    // Ярославль
			{195, 103, 81, 111, 99, 180, 0},   // Ростов
		},
		cities: []string{"Москва", "Суздаль", "Иваново", "Владимир", "Кострома", "Ярославль", "Ростов"},
	}
}

func (f *CityPerm) generateNewDecision(arr []float64, _ float64, iter int) []float64 {
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
	T := T0
	best := make([]float64, len(current))
	copy(best, current)
	bestVal := currentVal

	for k := 1; k <= maxIter; k++ {
		// охлаждение
		T = T * 0.99

		// новый кандидат
		candidate := problem.generateNewDecision(current, T, k)
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
		fmt.Println(T)
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
	function.run(100, 200, []float64{5, -5})
}
