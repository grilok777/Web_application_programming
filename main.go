package main

import (
	"html/template"
	"net/http"
)

// Спільна структура даних для всіх шаблонів
type PageData struct {
	Input1, Res1         map[string]float64
	InputMazut, ResMazut map[string]float64
	Input2, Res2         map[string]float64
	Input3, Res3         map[string]float64 // Для сонячних станцій
	HasResult            bool
}

// Глобальна функція парсингу шаблонів (її бачитимуть інші файли)
func parseTemplate(file string) *template.Template {
	return template.Must(template.ParseFiles("templates/layout.html", "templates/"+file))
}

func main() {
	// Маршрути (роутинг)
	http.HandleFunc("/pr1", handlerPR1) // Функцію handlerPR1 можете винести в pr1.go
	http.HandleFunc("/pr2", handlerPR2) // Функцію handlerPR2 можете винести в pr2.go
	http.HandleFunc("/pr3", handlerPR3) // Ця функція знаходиться у файлі pr3.go

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/pr1", http.StatusSeeOther)
	})

	println("Сервер запущено: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}