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
	sum := 0
	//if len(arr) != 8 {
	//	panic("arr must have length 8")
	//}
	for i := 1; i < len(arr); i++ {
		past, next := int(arr[i-1]), int(arr[i])
		sum += f.distances[past][next]
	}
	return sum
}
func (f *CityPerm) getD() int {
	return len(f.distances)
}
func (a *antCalc) getPast() []int {
	return a.past
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
func generateArr(n int, m int) [][]float64 {
	arr := make([][]float64, n)
	for i := range arr {
		arr[i] = make([]float64, m)
		for j := range arr[i] {
			arr[i][j] = rand.Float64() / 10
		}
	}
	return arr
}
func generateAnts(n int) []*antCalc {
	arr := make([]*antCalc, n)
	for i := range arr {
		arr[i] = &antCalc{current: 0, past: []int{1, 2, 3, 4, 5, 6}}
	}
	return arr
}
func nextVer(til []float64, summ float64) int {
	r := rand.Float64() * summ
	cum := 0.0
	bestI := 0
	for i, v := range til {
		cum += v
		if r <= cum {
			bestI = i
			break
		}
	}
	return bestI
}
func main() {
	rand.Seed(61)
	city := NewCityPerm()
	feramon := generateArr(len(city.distances), len(city.cities))
	alfa, beta, p, t := 2.0, 1.0, 0.5, 0
	ants := generateAnts(10)
	globalBest := ants[0]
	for {
		for _, ant := range ants {
			ant.current = 0
			ant.past = []int{1, 2, 3, 4, 5, 6}
			ant.route = []int{ant.current} // сохраняем маршрут

			for len(ant.past) != 0 {
				til := make([]float64, len(ant.past))
				cum := 0.0
				for i, next := range ant.past {
					eta := 1.0 / float64(city.distances[ant.current][next])
					num := math.Pow(feramon[ant.current][next], alfa) * math.Pow(eta, beta)
					til[i] = num
					cum += num
				}
				bestI := nextVer(til, cum)
				nextCity := ant.past[bestI]
				ant.current = nextCity
				ant.route = append(ant.route, nextCity)
				ant.past = append(ant.past[:bestI], ant.past[bestI+1:]...)
			}
			ant.route = append(ant.route, 0)
			ant.fy = city.getY(ant.route)
		}
		// испарение
		for i := 0; i < city.getD(); i++ {
			for j := 0; j < city.getD(); j++ {
				feramon[i][j] *= (1 - p)
			}
		}

		// обновление по муравьям
		for _, ant := range ants {
			for i := 1; i < len(ant.route); i++ {
				a, b := ant.route[i-1], ant.route[i]
				feramon[a][b] += 1.0 / float64(ant.fy)
				feramon[b][a] += 1.0 / float64(ant.fy) // симметрия
			}
		}
		bestAnt := ants[0]
		for _, ant := range ants {
			if ant.fy < bestAnt.fy {
				bestAnt = ant
			}
		}
		for i := 1; i < len(bestAnt.route); i++ {
			a, b := bestAnt.route[i-1], bestAnt.route[i]
			feramon[a][b] += 5.0 / float64(bestAnt.fy) // усиление глобально лучшего
			feramon[b][a] += 5.0 / float64(bestAnt.fy)
		}
		t += 1
		if t == 10000 {
			break
		}
		for _, ant := range ants {
			if ant.fy < globalBest.fy {
				a := *ant       // разворачиваем указатель в значение
				globalBest = &a // сохраняем указатель на копию
			}

		}
	}
	for _, ant := range ants {
		fmt.Println(ant.route, ant.fy)
	}
	fmt.Println("best", globalBest.route, globalBest.fy)

}
