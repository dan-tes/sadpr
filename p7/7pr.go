package main

import (
	"fmt"
	_ "math"
	"math/rand"
	"time"
)

// Функция Матьяса (формула 1.2.1)
func matyas(x, y float64) float64 {
	return 0.26*(x*x+y*y) - 0.48*x*y
}

// Структура особи (индивида)
type Individual struct {
	x, y    float64
	fitness float64
}

const (
	popSize       = 100
	nGenerations  = 200
	xMin, xMax    = -10.0, 10.0
	baseCrossover = 0.8
	baseMutation  = 0.2
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Генетический алгоритм" +
		"" +
		"")
	// Инициализация популяции
	pop := make([]Individual, popSize)
	for i := range pop {
		pop[i].x = rand.Float64()*(xMax-xMin) + xMin
		pop[i].y = rand.Float64()*(xMax-xMin) + xMin
		pop[i].fitness = matyas(pop[i].x, pop[i].y)
	}

	for gen := 0; gen < nGenerations; gen++ {
		// Сортировка по приспособленности (чем меньше, тем лучше)
		sortByFitness(pop)

		best := pop[0]
		avgFit := averageFitness(pop)

		// Новое поколение
		newPop := make([]Individual, 0, popSize/2)

		// Селекция (оставляем лучших)
		survivors := pop[:popSize/2]

		for len(newPop) < popSize/2 {
			// Выбор родителей
			p1 := survivors[rand.Intn(len(survivors))]
			p2 := survivors[rand.Intn(len(survivors))]

			// Адаптивные коэффициенты
			pc := adaptiveCrossover(p1.fitness, avgFit)
			pm := adaptiveMutation(p1.fitness, avgFit)

			// Кроссовер
			child := crossover(p1, p2, pc)

			// Мутация
			child = mutate(child, pm)

			// Проверка границ
			child.x = clamp(child.x, xMin, xMax)
			child.y = clamp(child.y, xMin, xMax)

			child.fitness = matyas(child.x, child.y)
			newPop = append(newPop, child)
		}

		// Новая популяция = элита + потомки
		pop = append(survivors, newPop...)

		if gen%20 == 0 {
			fmt.Printf("Поколение %d: f(%.4f, %.4f) = %.6f\n",
				gen, best.x, best.y, best.fitness)
		}
	}

	best := pop[0]
	fmt.Printf("\nЛучшее решение: x=%e, y=%e, f(x,y)=%e\n",
		best.x, best.y, best.fitness)
}

// --- Вспомогательные функции ---

func sortByFitness(pop []Individual) {
	for i := 0; i < len(pop); i++ {
		for j := i + 1; j < len(pop); j++ {
			if pop[j].fitness < pop[i].fitness {
				pop[i], pop[j] = pop[j], pop[i]
			}
		}
	}
}

func averageFitness(pop []Individual) float64 {
	sum := 0.0
	for _, ind := range pop {
		sum += ind.fitness
	}
	return sum / float64(len(pop))
}

func adaptiveCrossover(fitness, avg float64) float64 {
	//if fitness < avg {
	//	return baseCrossover * 0.5
	//}
	return baseCrossover + rand.Float64()*0.2
}

func adaptiveMutation(fitness, avg float64) float64 {
	//if fitness < avg {
	//	return baseMutation * 0.5
	//}
	return baseMutation + rand.Float64()*0.3
}

func crossover(p1, p2 Individual, rate float64) Individual {
	if rand.Float64() > rate {
		return p1
	}
	alpha := rand.Float64()
	return Individual{
		x: alpha*p1.x + (1-alpha)*p2.x,
		y: alpha*p1.y + (1-alpha)*p2.y,
	}
}

func mutate(ind Individual, rate float64) Individual {
	if rand.Float64() < rate {
		ind.x += rand.NormFloat64() * 0.2
	}
	if rand.Float64() < rate {
		ind.y += rand.NormFloat64() * 0.2
	}
	return ind
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
