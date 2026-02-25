package main

import (
	"math"
	"net/http"
	"strconv"
)

// Нормальний закон розподілу
func gaussian(x, p, q float64) float64 {
	sqrtTwoPi := math.Sqrt(2 * math.Pi)
	exponent := -math.Pow(x-p, 2) / (2 * math.Pow(q, 2))
	return (1 / (q * sqrtTwoPi)) * math.Exp(exponent)
}

// Інтегрування (метод трапецій)
func integrate(numPoints int, p, q float64) float64 {
	abc_q := 5.0
	start := p - (p * abc_q / 100)
	end := p + (p * abc_q / 100)
	step := (end - start) / float64(numPoints)
	area := 0.0

	for x := start; x < end; x += step {
		y1 := gaussian(x, p, q)
		y2 := gaussian(x+step, p, q)
		area += (step / 2.0) * (y1 + y2)
	}
	return area
}

// Обробник для ПР №3
func handlerPR3(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Input3: make(map[string]float64),
		Res3:   make(map[string]float64),
	}

	if r.Method == http.MethodPost {
		pC, _ := strconv.ParseFloat(r.FormValue("pc"), 64)
		q1, _ := strconv.ParseFloat(r.FormValue("q1"), 64)
		q2, _ := strconv.ParseFloat(r.FormValue("q2"), 64)
		price, _ := strconv.ParseFloat(r.FormValue("pr"), 64)

		day := 24.0

		// --- Розрахунки для q1 (Наявне) ---
		area1 := integrate(200, pC, q1)
		prElec1 := area1 * day * pC
		fElec1 := (1 - area1) * day * pC
		received1 := prElec1 * price * 1000
		paid1 := fElec1 * price * 1000

		data.Res3["w1"] = area1 * 100
		data.Res3["pEl1"] = prElec1
		data.Res3["fEl1"] = fElec1
		data.Res3["profit1"] = received1 - paid1
		data.Res3["rec1"] = received1
		data.Res3["paid1"] = paid1

		// --- Розрахунки для q2 (Передбачуване) ---
		area2 := integrate(200, pC, q2)
		prElec2 := area2 * day * pC
		fElec2 := (1 - area2) * day * pC
		received2 := prElec2 * price * 1000
		paid2 := fElec2 * price * 1000

		data.Res3["w2"] = area2 * 100
		data.Res3["pEl2"] = prElec2
		data.Res3["fEl2"] = fElec2
		data.Res3["profit2"] = received2 - paid2
		data.Res3["rec2"] = received2
		data.Res3["paid2"] = paid2

		data.HasResult = true
	}

	parseTemplate("pr3.html").ExecuteTemplate(w, "layout", data)
}