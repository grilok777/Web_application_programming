package main

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
)

// Структура для передачі даних у шаблони
type PageData struct {
	// Для ПР1 Завдання 1
	Input1 map[string]float64
	Res1   map[string]float64
	// Для ПР1 Завдання 2 (Мазут)
	InputMazut map[string]float64
	ResMazut   map[string]float64
	// Для ПР2
	Input2 map[string]float64
	Res2   map[string]float64
	// Для ПР3
	Input3 map[string]float64
	Res3   map[string]float64

	HasResult bool
}

// Константи та матриці для ПР3
var (
	KActP_X = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 14, 16, 18, 20, 25, 30, 35, 40, 50, 60, 80, 100}
	KActP_Y = []float64{0.1, 0.15, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}
	KActP   = [][]float64{
		{8.00, 5.33, 4.00, 2.67, 2.00, 1.60, 1.33, 1.14, 1},
		{6.22, 4.33, 3.39, 2.45, 1.98, 1.60, 1.33, 1.14, 1},
		{4.06, 2.89, 2.31, 1.74, 1.45, 1.34, 1.22, 1.14, 1},
		{3.24, 2.35, 1.91, 1.47, 1.25, 1.21, 1.22, 1.14, 1},
		{2.84, 2.09, 1.72, 1.35, 1.16, 1.16, 1.08, 1.03, 1},
		{2.64, 1.96, 1.62, 1.28, 1.14, 1.13, 1.06, 1.01, 1},
		{2.49, 1.86, 1.54, 1.23, 1.12, 1.1, 1.04, 1},
		{2.37, 1.78, 1.48, 1.19, 1.1, 1.08, 1.02, 1},
		{2.27, 1.71, 1.43, 1.16, 1.09, 1.07, 1.01, 1},
		{2.18, 1.65, 1.39, 1.13, 1.07, 1.05, 1},
		{2.04, 1.56, 1.32, 1.08, 1.05, 1.03, 1},
		{1.94, 1.49, 1.27, 1.05, 1.02, 1},
		{1.85, 1.43, 1.23, 1.02, 1},
		{1.78, 1.39, 1.19, 1},
		{1.72, 1.35, 1.16, 1},
		{1.60, 1.27, 1.10, 1},
		{1.51, 1.21, 1.05, 1},
		{1.44, 1.16, 1},
		{1.40, 1.13, 1},
		{1.30, 1.07, 1},
		{1.25, 1.03, 1},
		{1.16, 1},
		{1},
	}
	KActPdepartment = [][]float64{
		{8.00, 5.33, 4.00, 2.67, 2.00, 1.60, 1.33, 1.14},
		{5.01, 3.44, 2.69, 1.90, 1.52, 1.24, 1.11, 1},
		{2.40, 2.17, 1.80, 1.42, 1.23, 1.14, 1.08, 1},
		{2.28, 1.73, 1.46, 1.19, 1.06, 1.04, 1.00, 0.97},
		{1.31, 1.12, 1.02, 1.00, 0.98, 0.96, 0.94, 0.93},
		{1.20, 1.00, 0.96, 0.95, 0.94, 0.93, 0.92, 0.91},
		{1.10, 0.97, 0.91, 0.90, 0.90, 0.90, 0.90, 0.90},
		{0.80, 0.80, 0.80, 0.85, 0.85, 0.85, 0.90, 0.90},
		{0.75, 0.75, 0.75, 0.75, 0.75, 0.80, 0.85, 0.85},
		{0.65, 0.65, 0.65, 0.70, 0.70, 0.75, 0.80, 0.80},
	}
)

func main() {
	parse := func(file string) *template.Template {
		return template.Must(template.ParseFiles("template/layout.html", "template/"+file))
	}

	// Маршрут для ПР №1
	http.HandleFunc("/pr1", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Input1:     make(map[string]float64),
			InputMazut: make(map[string]float64),
			Res1:       make(map[string]float64),
			ResMazut:   make(map[string]float64),
		}

		if r.Method == http.MethodPost {
			data.HasResult = true
			calcType := r.FormValue("calc_type")

			if calcType == "fuel" {
				h, _ := strconv.ParseFloat(r.FormValue("H"), 64)
				c, _ := strconv.ParseFloat(r.FormValue("C"), 64)
				s, _ := strconv.ParseFloat(r.FormValue("S"), 64)
				n, _ := strconv.ParseFloat(r.FormValue("N"), 64)
				o, _ := strconv.ParseFloat(r.FormValue("O"), 64)
				w_val, _ := strconv.ParseFloat(r.FormValue("W"), 64)
				a, _ := strconv.ParseFloat(r.FormValue("A"), 64)

				k_dry := 100 / (100 - w_val)
				k_flam := 100 / (100 - w_val - a)

				data.Res1["Kdry"] = k_dry
				data.Res1["Kflam"] = k_flam
				data.Res1["Hc"] = h * k_dry
				data.Res1["Cc"] = c * k_dry
				data.Res1["Sc"] = s * k_dry
				data.Res1["Nc"] = n * k_dry
				data.Res1["Oc"] = o * k_dry
				data.Res1["Ac"] = a * k_dry

				data.Res1["Hg"] = h * k_flam
				data.Res1["Cg"] = c * k_flam
				data.Res1["Sg"] = s * k_flam
				data.Res1["Ng"] = n * k_flam
				data.Res1["Og"] = o * k_flam

				q_p_h := (339*c + 1030*h - 108.8*(o-s) - 25*w_val) / 1000
				data.Res1["Qph"] = q_p_h
				data.Res1["Qch"] = (q_p_h + 0.025*w_val) * 100 / (100 - w_val)
				data.Res1["Qgh"] = (q_p_h + 0.025*w_val) * 100 / (100 - w_val - a)

			} else if calcType == "mazut" {
				c, _ := strconv.ParseFloat(r.FormValue("C"), 64)
				h, _ := strconv.ParseFloat(r.FormValue("H"), 64)
				o, _ := strconv.ParseFloat(r.FormValue("O"), 64)
				s, _ := strconv.ParseFloat(r.FormValue("S"), 64)
				q_dafi, _ := strconv.ParseFloat(r.FormValue("Qdafi"), 64)
				w_val, _ := strconv.ParseFloat(r.FormValue("W"), 64)
				a, _ := strconv.ParseFloat(r.FormValue("A"), 64)
				v, _ := strconv.ParseFloat(r.FormValue("V"), 64)
				
				k_dry := (100 - w_val) / 100
				data.ResMazut["Vp"] = v * k_dry 

				k_flam := (100 - w_val - a) / 100
				data.ResMazut["Cp"] = c * k_flam
				data.ResMazut["Hp"] = h * k_flam
				data.ResMazut["Op"] = o * k_flam
				data.ResMazut["Sp"] = s * k_flam
				data.ResMazut["Ap"] = a * (100 - w_val) / 100
				data.ResMazut["Qri"] = q_dafi*(100-w_val-data.ResMazut["Ap"])/100 - 0.025*w_val
			}
		}
		parse("pr1.html").ExecuteTemplate(w, "layout", data)
	})

	// Маршрут для ПР №2
	http.HandleFunc("/pr2", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{Res2: make(map[string]float64)}
		if r.Method == http.MethodPost {
			b, _ := strconv.ParseFloat(r.FormValue("B"), 64)
			n, _ := strconv.ParseFloat(r.FormValue("n"), 64)
			av, _ := strconv.ParseFloat(r.FormValue("av"), 64)
			fuel := r.FormValue("fuelType")

			var qri, ar, hv float64
			if fuel == "coal" {
				qri, ar, hv = 20.47, 25.2, 1.5
			} else if fuel == "oil" {
				qri, ar, hv = 39.48, 0.15, 0.0
			} else {
				data.Res2["Kem"] = 0
				data.Res2["E"] = 0
				data.HasResult = true
				parse("pr2.html").ExecuteTemplate(w, "layout", data)
				return
			}

			k_em := (math.Pow(10, 6) / qri) * av * (ar / (100 - hv)) * (1 - n)
			e_val := 0.000001 * k_em * qri * b

			data.Res2["Kem"] = k_em
			data.Res2["E"] = e_val
			data.HasResult = true
		}
		parse("pr2.html").ExecuteTemplate(w, "layout", data)
	})

	// Маршрут для ПР №3
	http.HandleFunc("/pr3", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Input3: make(map[string]float64),
			Res3:   make(map[string]float64),
		}

		if r.Method == http.MethodPost {
			pShlif, _ := strconv.ParseFloat(r.FormValue("p"), 64)
			kUse, _ := strconv.ParseFloat(r.FormValue("k"), 64)
			qPiv, _ := strconv.ParseFloat(r.FormValue("q"), 64)

			tempK := 75.16 + 40*kUse + 0.15*4*pShlif
			tempPn := 376.0 + 4*pShlif
			gK := tempK / tempPn
			data.Res3["grK"] = gK

			tempP2 := 13192.0 + 4*math.Pow(pShlif, 2)
			eN := math.Ceil(math.Pow(tempPn, 2) / tempP2)
			data.Res3["efecN"] = eN

			// Знаходимо індекси у масивах (без падіння програми)
			inX, inY := 0, 0
			for i, v := range KActP_X {
				if v >= eN {
					inX = i
					break
				}
			}
			for i, v := range KActP_Y {
				if v >= gK {
					inY = i
					break
				}
			}
			
			// Захист від виходу за межі масиву (Index Out of Bounds)
			if inX >= len(KActP) {
				inX = len(KActP) - 1
			}
			if inY >= len(KActP[inX]) {
				inY = len(KActP[inX]) - 1
			}
			
			pK := KActP[inX][inY]
			data.Res3["actK"] = pK

			pP := pK * tempK
			data.Res3["actP"] = pP
			reactEshr := 66.926 + 0.798*pShlif + 40*kUse + 10.8*qPiv
			pQ := pK * reactEshr
			data.Res3["reactP"] = pQ
			data.Res3["fullP"] = math.Sqrt(pP*pP + pQ*pQ)
			data.Res3["grI"] = pP / 0.38

			kPh := 3*tempK + 232
			hP := 3*tempPn + 440
			uKd := kPh / hP
			data.Res3["useK"] = uKd
			
			tempAll := 3*tempP2 + 48800
			epN := math.Ceil(math.Pow(hP, 2) / tempAll)
			data.Res3["efecNAll"] = epN

			idxX := 0
			if epN <= 10 {
				idxX = 1
			} else if epN <= 25 {
				idxX = 3
			} else {
				idxX = 5
			}
			
			pKAll := KActPdepartment[idxX][7] 
			data.Res3["actKAll"] = pKAll

			data.Res3["actP038"] = pKAll * kPh
			data.Res3["reactP038"] = pKAll * (reactEshr*3 + 120)
			data.Res3["fullP038"] = math.Sqrt(math.Pow(data.Res3["actP038"], 2) + math.Pow(data.Res3["reactP038"], 2))
			data.Res3["grI038"] = data.Res3["actP038"] / 0.38

			data.HasResult = true
		}
		parse("pr3.html").ExecuteTemplate(w, "layout", data)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/pr1", http.StatusSeeOther)
	})

	println("Сервер запущено: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}