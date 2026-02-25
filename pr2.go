package main

import (
	"math"
	"net/http"
	"strconv"
)

func handlerPR2(w http.ResponseWriter, r *http.Request) {
	data := PageData{Res2: make(map[string]float64)}
	
	if r.Method == http.MethodPost {
		b, _ := strconv.ParseFloat(r.FormValue("B"), 64)
		n, _ := strconv.ParseFloat(r.FormValue("n"), 64)
		av, _ := strconv.ParseFloat(r.FormValue("av"), 64)
		fuel := r.FormValue("fuelType")

		var qri, ar, hv float64
		// Константи ПР2
		if fuel == "coal" {
			qri, ar, hv = 20.47, 25.2, 1.5
		} else if fuel == "oil" {
			qri, ar, hv = 39.48, 0.15, 0.0
		} else {
			data.Res2["Kem"], data.Res2["E"] = 0, 0
			data.HasResult = true
			parseTemplate("pr2.html").ExecuteTemplate(w, "layout", data)
			return
		}

		k_em := (math.Pow(10, 6) / qri) * av * (ar / (100 - hv)) * (1 - n)
		e_val := 0.000001 * k_em * qri * b

		data.Res2["Kem"] = k_em
		data.Res2["E"] = e_val
		data.HasResult = true
	}
	
	parseTemplate("pr2.html").ExecuteTemplate(w, "layout", data)
}