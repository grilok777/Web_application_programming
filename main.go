package main

import (
	"html/template"
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


	HasResult bool
}

func main() {
	// Функція для парсингу шаблонів з підтримкою Layout
	parse := func(file string) *template.Template {
		return template.Must(template.ParseFiles("D:/Mysor2/4kurs/Prow_web_2/Pr1/template/layout.html", "template/"+file))
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
				data.ResMazut["Ap"] = a * k_dry
				data.ResMazut["Qri"] = q_dafi*(100-w_val-data.ResMazut["Ap"])/100 - 0.025*w_val
			}
		}
		parse("pr1.html").ExecuteTemplate(w, "layout", data)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/pr1", http.StatusSeeOther)
	})

	println("Сервер запущено: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}