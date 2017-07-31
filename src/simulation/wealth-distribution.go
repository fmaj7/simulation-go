package main

import (
	"math/rand"
	"flag"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"net/http"
	"log"
	"strconv"
)

// Attempt to simulate conclusion from https://mp.weixin.qq.com/s/gOHPDLkdvvskq06qj_QRgg
// where the social wealth are guaranteed to diverge in an extreme way

func Simulate(population int, numOfSim int) map[int]int {
	m := make(map[int]int)

	// initialize everyone
	for i := 0; i < population; i++ {
		m[i] = 100
	}

	for i := 0; i < numOfSim * population; i++ {
		m[rand.Intn(population)] -= 1
		m[rand.Intn(population)] += 1
	}

	return m
}

func DrawChart(res http.ResponseWriter, req *http.Request) {

	population, _ := strconv.Atoi(req.URL.Query().Get("population"))
	simulation, _ := strconv.Atoi(req.URL.Query().Get("simulation"))

	result := Simulate(population, simulation)

	bars := make([]chart.Value, len(result))

	for i := 0; i < len(result); i++ {
		bars[i] = chart.Value{Value: float64(result[i]), Label: fmt.Sprintf("%v", i)}
	}

	fmt.Printf("bars: %v\n", bars)

	sbc := chart.BarChart{
		Width: 1024,
		Height: 512,
		BarWidth: 5,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: bars,
	}

	res.Header().Set("Content-Type", "image/png")
	err := sbc.Render(chart.PNG, res)
	if err != nil {
		fmt.Printf("\nError rendering chart: %v\n", err)
	}
}

func ServeChart() {
	listenPort := fmt.Sprintf(":8080")
	fmt.Printf("Listening on %s\n", listenPort)
	http.HandleFunc("/", DrawChart)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}

func Cli() {
	population := flag.Int("p", 100, "population of the simulation")
	numOfSim := flag.Int("n", 17000, "number of the simulation")

	flag.Parse()

	result := Simulate(*population, *numOfSim)
	fmt.Println(result)
}

func main() {
	ServeChart()
}
