package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

// Функція для форматування чисел у вигляді "X.X*10^-Y" (аналог JS strPower)
func strPower(temp float64) string {
	if temp >= 0.01 {
		return fmt.Sprintf("%.4f", temp)
	}
	power := 0
	for temp*math.Pow(10, float64(power)) < 1 {
		power++
	}
	val := temp * math.Pow(10, float64(power))
	return fmt.Sprintf("%.1f*10^-%d", val, power)
}

func handlerPR5(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Input5:   make(map[string]float64),
		Res5:     make(map[string]float64),
		Res5Text: make(map[string]string),
	}

	if r.Method == http.MethodPost {
		data.HasResult = true
		calcType := r.FormValue("calc_type")

		if calcType == "task1" {
			// Таблиця даних 
			table := [][]float64{
				{0.015, 100, 1, 43},
				{0.02, 80, 1, 28},
				{0.005, 60, 0.5, 10},
				{0.05, 60, 0.5, 10},
				{0.01, 30, 0.1, 30},
				{0.02, 15, 0.33, 15},
				{0.01, 15, 0.33, 15},
				{0.03, 2, 0.167, 5},
				{0.05, 4, 0.33, 10},
				{0.1, 160, 0.5, 0},
				{0.1, 50, 0.5, 0},
				{0.07, 10, 0.167, 35},
				{0.02, 8, 0.167, 35},
				{0.02, 10, 0.167, 35},
				{0.03, 44, 1, 9},
				{0.005, 17.5, 1, 9},
			}

			// Зчитуємо 16 інпутів у масив
			allEq := make([]float64, 16)
			for i := 0; i < 16; i++ {
				val, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("n%d", i)), 64)
				allEq[i] = val
			}

			// 1. Частота відмов одноколової системи
			tempWos := 0.0
			for i := 0; i < len(table); i++ {
				tempWos += allEq[i] * table[i][0]
			}
			data.Res5["wOs"] = tempWos

			// 2. Середня тривалість відновлення
			temp := 0.0
			for i := 0; i < len(table); i++ {
				temp += allEq[i] * table[i][0] * table[i][1]
			}
			tempT := temp / tempWos
			data.Res5["t"] = tempT

			// 3. Коефіцієнт аварійного простою
			tempKe := tempWos * tempT / 8760.0
			data.Res5Text["kE"] = strPower(tempKe)

			// 4. Коефіцієнт планового простою
			tempKpMax := 0.0
			for i := 0; i < len(table); i++ {
				if allEq[i] > 0 {
					kpVal := table[i][2] * table[i][3]
					if kpVal > tempKpMax {
						tempKpMax = kpVal
					}
				}
			}
			tempKp := 1.2 * tempKpMax / 8760.0
			data.Res5Text["kP"] = strPower(tempKp)

			// 5 & 6. Частота відмов двоколової системи
			tempWdk := 2 * tempWos * (tempKe + tempKp)
			tempWds := tempWdk + 0.02
			data.Res5Text["wDk"] = strPower(tempWdk)
			data.Res5Text["wDs"] = strPower(tempWds)

		} else if calcType == "task2" {
			pE, _ := strconv.ParseFloat(r.FormValue("priceE"), 64)
			pP, _ := strconv.ParseFloat(r.FormValue("pricePl"), 64)
			tE, _ := strconv.ParseFloat(r.FormValue("timeE"), 64)
			tP, _ := strconv.ParseFloat(r.FormValue("timePl"), 64)

			// Константи
			T := 6451.0
			P_const := 5.12
			w_const := 0.01

			// Математичне сподівання збитків
			MlossEm := w_const * tE * P_const * T
			MlossPl := tP * P_const * T
			loss := pE*MlossEm + pP*MlossPl

			data.Res5["loss"] = loss
		}
	}

	parseTemplate("pr5.html").ExecuteTemplate(w, "layout", data)
}