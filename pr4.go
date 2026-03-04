package main

import (
	"math"
	"net/http"
	"strconv"
)

func handlerPR4(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Input4:   make(map[string]float64),
		Res4:     make(map[string]float64),
		Res4Text: make(map[string]string),
	}

	if r.Method == http.MethodPost {
		calcType := r.FormValue("calc_type")
		data.HasResult = true

		if calcType == "task1" {
			// ЗАВДАННЯ 1: Вибір кабелів
			ik, _ := strconv.ParseFloat(r.FormValue("Ik"), 64)
			tf, _ := strconv.ParseFloat(r.FormValue("Tf"), 64)
			pm, _ := strconv.ParseFloat(r.FormValue("Pm"), 64)
			tm, _ := strconv.ParseFloat(r.FormValue("Tm"), 64)
			mat, _ := strconv.ParseInt(r.FormValue("materl"), 10, 64)
			is, _ := strconv.ParseInt(r.FormValue("isolat"), 10, 64)

			unorm := 10.0
			ct := 92.0

			// 1. Струми
			im := (pm / 2) / (math.Sqrt(3) * unorm)
			impa := 2 * im
			data.Res4["Im"] = im
			data.Res4["Impa"] = impa

			// 2. Економічний переріз
			table := [][]float64{
				{2.5, 2.1, 1.8}, {1.3, 1.1, 1.0}, {3.0, 2.5, 2.0},
				{1.6, 1.4, 1.2}, {3.5, 3.1, 2.7}, {1.9, 1.7, 1.6},
			}
			inY := 0
			if tm > 1000 && tm < 3000 {
				inY = 0
			} else if tm >= 3000 && tm < 5000 {
				inY = 1
			} else if tm >= 5000 {
				inY = 2
			}

			j := table[mat+is][inY]
			data.Res4["Sek"] = im / j

			// 3. Переріз
			data.Res4["S"] = math.Ceil((ik * math.Sqrt(tf)) / ct)

			data.Res4Text["t1"] = "Для внутрізаводської мережі вибираємо броньовані кабелі з паперовою ізоляцією в алюмінієвій оболонці типу ААБ."
			data.Res4Text["t2"] = "Вибираємо кабель ААБ 10 3х25 з допустимим струмом Ідоп = 90 А."

		} else if calcType == "task2" {
			// ЗАВДАННЯ 2: Струми КЗ ГПП
			p, _ := strconv.ParseFloat(r.FormValue("P"), 64)
			uCh, _ := strconv.ParseFloat(r.FormValue("Uch"), 64)

			sNom := 6.3
			uK := 10.5

			xc := math.Pow(uCh, 2) / p
			xt := (uK / 100.0) * math.Pow(uCh, 2) / sNom
			x := xc + xt

			data.Res4["R"] = x
			data.Res4["I"] = uCh / (math.Sqrt(3) * x)

		} else if calcType == "task3" {
			// ЗАВДАННЯ 3: Струми КЗ підстанції
			rn, _ := strconv.ParseFloat(r.FormValue("Rn"), 64)
			xn, _ := strconv.ParseFloat(r.FormValue("Xn"), 64)
			rmin, _ := strconv.ParseFloat(r.FormValue("Rmin"), 64)
			xmin, _ := strconv.ParseFloat(r.FormValue("Xmin"), 64)

			uV, uN := 115.0, 11.0
			xT, kPr := 233.0, 0.009
			rl, xl := 7.91, 4.49

			// Опори
			xShN := xn + xT
			xShMin := xmin + xT

			zSh := math.Sqrt(math.Pow(xShN, 2) + math.Pow(rn, 2))
			zShMin := math.Sqrt(math.Pow(xShMin, 2) + math.Pow(rmin, 2))

			zShN := math.Sqrt(math.Pow(xShN*kPr, 2) + math.Pow(rn*kPr, 2))
			zShMinN := math.Sqrt(math.Pow(xShMin*kPr, 2) + math.Pow(rmin*kPr, 2))

			zSumN := math.Sqrt(math.Pow(xl+xShN*kPr, 2) + math.Pow(rl+rn*kPr, 2))
			zSumMinN := math.Sqrt(math.Pow(xl+xShMin*kPr, 2) + math.Pow(rl+rmin*kPr, 2))

			// Струми
			data.Res4["i3Sh"] = (uV * 1000) / (math.Sqrt(3) * zSh)
			data.Res4["i3ShMin"] = (uV * 1000) / (math.Sqrt(3) * zShMin)
			data.Res4["i2Sh"] = data.Res4["i3Sh"] * math.Sqrt(3) / 2
			data.Res4["i2ShMin"] = data.Res4["i3ShMin"] * math.Sqrt(3) / 2

			data.Res4["i3ShN"] = (uN * 1000) / (math.Sqrt(3) * zShN)
			data.Res4["i3ShMinN"] = (uN * 1000) / (math.Sqrt(3) * zShMinN)
			data.Res4["i2ShN"] = data.Res4["i3ShN"] * math.Sqrt(3) / 2
			data.Res4["i2ShMinN"] = data.Res4["i3ShMinN"] * math.Sqrt(3) / 2

			data.Res4["i3LN"] = (uN * 1000) / (math.Sqrt(3) * zSumN)
			data.Res4["i3LMinN"] = (uN * 1000) / (math.Sqrt(3) * zSumMinN)
			data.Res4["i2LN"] = data.Res4["i3LN"] * math.Sqrt(3) / 2
			data.Res4["i2LMinN"] = data.Res4["i3LMinN"] * math.Sqrt(3) / 2
		}
	}

	parseTemplate("pr4.html").ExecuteTemplate(w, "layout", data)
}