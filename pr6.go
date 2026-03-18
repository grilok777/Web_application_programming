package main

import (
	"math"
	"net/http"
	"strconv"
)

// Константи та матриці
var (
	KActP_X = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 14, 16, 18, 20, 25, 30, 35, 40, 50, 60, 80, 100}
	KActP_Y = []float64{0.1, 0.15, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}

	KActP_X_department = []float64{1, 2, 3, 4, 5}
	KActP_Y_department = []float64{0.1, 0.15, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7}	

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


func handlerPR6(w http.ResponseWriter, r *http.Request) {
	data := PageData{Input3: make(map[string]float64), Res3: make(map[string]float64)}

		if r.Method == http.MethodPost {
			pShlif, _ := strconv.ParseFloat(r.FormValue("p"), 64)
			kUse, _ := strconv.ParseFloat(r.FormValue("k"), 64)
			qPiv, _ := strconv.ParseFloat(r.FormValue("q"), 64)

			// 1. Груповий КВ
			tempK := 75.16 + 40*kUse + 0.15*4*pShlif
			tempPn := 376.0 + 4*pShlif
			gK := tempK / tempPn
			data.Res3["grK"] = gK

			// 2. Ефективна кількість
			tempP2 := 13192.0 + 4*math.Pow(pShlif, 2)
			eN := math.Ceil(math.Pow(tempPn, 2) / tempP2)
			data.Res3["efecN"] = eN

			// 3. Розрахунковий коефіцієнт (ТОЧНА ІНТЕРПОЛЯЦІЯ)
			inY := 0
			minDiff := math.MaxFloat64
			for i, val := range KActP_Y {
				diff := math.Abs(val - gK)
				if diff < minDiff {
					minDiff = diff
					inY = i
				}
			}

			inX := -1
			for i, val := range KActP_X {
				if val == eN {
					inX = i
					break
				}
			}

			var pK float64
			if inX != -1 {
				safeY := inY
				if safeY >= len(KActP[inX]) { safeY = len(KActP[inX]) - 1 }
				pK = KActP[inX][safeY]
			} else {
				eNafter, eNbefore := math.MaxFloat64, -math.MaxFloat64
				for _, curr := range KActP_X {
					if curr > eN && curr < eNafter { eNafter = curr }
					if curr < eN && curr > eNbefore { eNbefore = curr }
				}

				inXafter, inXbefore := 0, 0
				for i, val := range KActP_X {
					if val == eNafter { inXafter = i }
					if val == eNbefore { inXbefore = i }
				}

				safeYafter := inY
				if safeYafter >= len(KActP[inXafter]) { safeYafter = len(KActP[inXafter]) - 1 }
				pKafter := KActP[inXafter][safeYafter]

				safeYbefore := inY
				if safeYbefore >= len(KActP[inXbefore]) { safeYbefore = len(KActP[inXbefore]) - 1 }
				pKbefore := KActP[inXbefore][safeYbefore]

				pK = pKbefore + (eN-eNbefore)/(eNafter-eNbefore)*(pKafter-pKbefore)
			}
			data.Res3["actK"] = pK

			// 4-7. Навантаження ШР
			pP := pK * tempK
			data.Res3["actP"] = pP
			reactEshr := 66.926 + 0.798*pShlif + 40*kUse + 10.8*qPiv
			pQ := pK * reactEshr
			data.Res3["reactP"] = pQ
			data.Res3["fullP"] = math.Sqrt(pP*pP + pQ*pQ)
			data.Res3["grI"] = pP / 0.38

			// 8-10. Цех в цілому
			kPh := 3*tempK + 232
			hP := 3*tempPn + 440
			uKd := kPh / hP
			data.Res3["useK"] = uKd

			tempAll := 3*tempP2 + 48800
			epN := math.Ceil(math.Pow(hP, 2) / tempAll)
			data.Res3["efecNAll"] = epN

			// Точний пошук індексів цеху 
			uKdRounded2 := math.Round(uKd*100) / 100
			uKdRounded1 := math.Round(uKd*10) / 10

			inYd := -1
			for i, val := range KActP_Y_department {
				if math.Abs(val-uKdRounded2) < 0.001 { inYd = i; break }
			}
			if inYd == -1 {
				for i, val := range KActP_Y_department {
					if math.Abs(val-uKdRounded1) < 0.001 { inYd = i; break }
				}
			}
			if inYd == -1 { inYd = 7 }

			inXd := -1
			for i, val := range KActP_X_department {
				if val == epN { inXd = i; break }
			}
			if inXd == -1 {
				if epN <= 8 { inXd = 5 } else if epN <= 10 { inXd = 6 } else if epN <= 25 { inXd = 7 } else if epN <= 50 { inXd = 8 } else { inXd = 9 }
			}

			safeYd := inYd
			if safeYd >= len(KActPdepartment[inXd]) { safeYd = len(KActPdepartment[inXd]) - 1 }
			pKd := KActPdepartment[inXd][safeYd]
			data.Res3["actKAll"] = pKd

			// 11-14. Шини 0.38
			data.Res3["actP038"] = pKd * kPh
			data.Res3["reactP038"] = pKd * (reactEshr*3 + 120)
			data.Res3["fullP038"] = math.Sqrt(math.Pow(data.Res3["actP038"], 2) + math.Pow(data.Res3["reactP038"], 2))
			data.Res3["grI038"] = data.Res3["actP038"] / 0.38

			data.HasResult = true
		}
		parseTemplate("pr6.html").ExecuteTemplate(w, "layout", data)
}