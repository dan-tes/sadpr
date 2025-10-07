package main

import (
	"fmt"
	"math"
	"math/rand"
)

type antCalc struct {
	current int
	past    []int
	fy      int
	route   []int
}
type CityPerm struct {
	distances [][]int
	cities    []string
}

func (f *CityPerm) getY(arr []int) int {
	// безопасно считаем длину пути для любого размера >=2
	if len(arr) < 2 {
		return math.MaxInt32 // очень плохо — длинна пути некорректна
	}
	sum := 0
	for i := 1; i < len(arr); i++ {
		past, next := arr[i-1], arr[i]
		sum += f.distances[past][next]
	}
	return sum
}
func (f *CityPerm) getCities(arr []int) []string {
	bestPos := make([]string, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		bestPos = append(bestPos, f.cities[int(arr[i])])
	}
	return bestPos
}
func (f *CityPerm) getD() int {
	return len(f.distances)
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
func generateArr(n int, m int) [][]float64 {
	arr := make([][]float64, n)
	for i := range arr {
		arr[i] = make([]float64, m)
		for j := range arr[i] {
			// небольшие положительные начальные феромоны
			arr[i][j] = 0.01 + rand.Float64()*0.01
		}
	}
	return arr
}
func generateAnts(nAnts int, nCities int) []*antCalc {
	arr := make([]*antCalc, nAnts)
	for i := range arr {
		// создаём past автоматически (все города кроме 0)
		past := make([]int, 0, nCities-1)
		for k := 1; k < nCities; k++ {
			past = append(past, k)
		}
		arr[i] = &antCalc{current: 0, past: past}
	}
	return arr
}
func nextVer(til []float64) int {
	r := rand.Float64()
	cum := 0.0
	for i, v := range til {
		cum += v
		if r <= cum {
			return i
		}
	}
	// на случай погрешностей вернём последний индекс
	return len(til) - 1
}
func main() {
	fmt.Println("dkdk")
	city := NewCityPerm()
	feramon := generateArr(len(city.distances), len(city.cities))

	// параметры (p — доля испарения, обычно в (0,1). p=1 — полностью убирает феромон!)
	alfa, beta, p := 3.0, 5.0, 0.1
	t := 0
	ants := generateAnts(10, len(city.cities))

	//инициализация глобального лучшего (копия первого муравья безопасно)
	globalBest := &antCalc{fy: math.MaxInt32, route: nil}

	for ; t < 1000; t++ {
		// каждый муравей строит маршрут
		for _, ant := range ants {
			ant.current = 0
			// восстановим past (вдруг предыдущая итерация изменила)
			ant.past = make([]int, 0, len(city.cities)-1)
			for k := 1; k < len(city.cities); k++ {
				ant.past = append(ant.past, k)
			}
			ant.route = []int{ant.current}

			for len(ant.past) != 0 {
				til := make([]float64, len(ant.past))
				cum := 0.0
				for i, next := range ant.past {
					eta := 1.0 / float64(city.distances[ant.current][next])
					num := math.Pow(feramon[ant.current][next], alfa) * math.Pow(eta, beta)
					til[i] = num
					cum += num
				}
				if cum == 0 {
					for i := range til {
						til[i] = 1.0 / float64(len(til))
					}
				} else {
					for i := range til {
						til[i] /= cum
					}
				}
				bestI := nextVer(til)
				feramon[ant.current][bestI] = feramon[ant.current][bestI]*(0.9) + 0.01
				nextCity := ant.past[bestI]
				ant.current = nextCity
				ant.route = append(ant.route, nextCity)
				ant.past = append(ant.past[:bestI], ant.past[bestI+1:]...)
			}
			// возвращение в начальную точку
			ant.route = append(ant.route, 0)
			ant.fy = city.getY(ant.route)
		}
		// испарение феромона
		for i := 0; i < city.getD(); i++ {
			for j := 0; j < city.getD(); j++ {
				feramon[i][j] *= (1.0 - p)
			}
		}

		// обновление по муравьям
		for _, ant := range ants {
			// защита на случай некорректного fy
			if ant.fy <= 0 {
				continue
			}
			delta := 0.1 / float64(ant.fy)
			for i := 1; i < len(ant.route); i++ {
				a, b := ant.route[i-1], ant.route[i]
				feramon[a][b] += delta
				feramon[b][a] += delta
			}
		}

		// добавляем дополнительно вклад лучшего в итерации
		bestAnt := ants[0]
		for _, ant := range ants {
			if ant.fy < bestAnt.fy {
				bestAnt = ant
			}
		}
		if bestAnt.fy > 0 {
			delta := 10.0 / float64(bestAnt.fy)
			for i := 1; i < len(bestAnt.route); i++ {
				a, b := bestAnt.route[i-1], bestAnt.route[i]
				feramon[a][b] += delta
				feramon[b][a] += delta
			}
		}

		// обновляем глобальный лучший (копируем маршрут)
		for _, ant := range ants {
			if ant.fy < globalBest.fy {
				routeCopy := make([]int, len(ant.route))
				copy(routeCopy, ant.route)
				globalBest = &antCalc{fy: ant.fy, route: routeCopy}
			}
		}

		fmt.Println("Итерация", t+1, "Текущее лучшее решение", globalBest.route, globalBest.fy)
	}

	fmt.Println("Лучшие решения по муравьям (последняя итерация):")
	for i, ant := range ants {
		fmt.Println("Муравей", i+1, ant.route, ant.fy)
	}
	fmt.Println("Лучшее решение (города):", city.getCities(globalBest.route), globalBest.fy)
}
