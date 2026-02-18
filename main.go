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
	Input2    map[string]float64
	Res2      map[string]float64

	HasResult bool
}

func main() {
	// Функція для парсингу шаблонів з підтримкою Layout
	parse := func(file string) *template.Template {
		return template.Must(template.ParseFiles("templates/layout.html", "templates/"+file)) 
	}

	// Маршрут для ПР №1 (Калькулятор палива та мазуту)
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
				// Логіка Завдання 1 (Паливо)
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
				// Логіка Завдання 2 (Мазут)
				c, _ := strconv.ParseFloat(r.FormValue("C"), 64)
				h, _ := strconv.ParseFloat(r.FormValue("H"), 64)
				o, _ := strconv.ParseFloat(r.FormValue("O"), 64)
				s, _ := strconv.ParseFloat(r.FormValue("S"), 64)
				q_dafi, _ := strconv.ParseFloat(r.FormValue("Qdafi"), 64)
				w_val, _ := strconv.ParseFloat(r.FormValue("W"), 64)
				a, _ := strconv.ParseFloat(r.FormValue("A"), 64)
				v, _ := strconv.ParseFloat(r.FormValue("V"), 64)
				k_dry := (100 - w_val) / 100
				data.ResMazut["Vp"] = v * k_dry // Перерахунок ванадію на робочу масу

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
			// Константи з вашого звіту ПР2
			if fuel == "coal" {
				qri, ar, hv = 20.47, 25.2, 1.5
			} else if fuel == "oil" {
				qri, ar, hv = 39.48, 0.15, 0.0
			} else {
				data.Res2["Kem"], data.Res2["E"] = 0, 0
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/pr1", http.StatusSeeOther)
	})

	println("Сервер запущено: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}