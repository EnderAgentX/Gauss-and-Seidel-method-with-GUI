package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"math"
	"strconv"
	"strings"
)

/**
  пр 1
	3
	3,-9,3,2,-4,4,1,8,-18
	-18,-10,35
	Ответ: 1, 2, -1

 пр 2
  3
  1,-3,1,0,0,-2,0,11,-20
  -6,2,42
  Ответ: 1, 2, -1
*/

type matrix [][]float64
type vector []float64

func main() {
	newApp := app.New()
	w := newApp.NewWindow("Метод Гаусса")
	w.Resize(fyne.NewSize(1820, 980))
	w.CenterOnScreen()

	label := widget.NewLabel("Введите размер матрицы ")
	entry := widget.NewEntry()
	label1 := widget.NewLabel("Введите матрицу ")
	entry1 := widget.NewEntry()
	label2 := widget.NewLabel("Векторы")
	entry2 := widget.NewEntry()

	answer := widget.NewLabel("")

	scr := container.NewVScroll(answer)
	scr.SetMinSize(fyne.NewSize(400, 150))

	btn := widget.NewButton("Посчитать", func() {
		flag := false
		answer.SetText("")
		n, err := strconv.Atoi(entry.Text)
		if err != nil {
			panic(err)
		}

		splitFunc := func(r rune) bool {
			return strings.ContainsRune(", ", r)
		}

		numsStr := strings.FieldsFunc(entry1.Text, splitFunc)
		vecStr := strings.FieldsFunc(entry2.Text, splitFunc)

		nums := make([]float64, n*n)
		vecs := make([]float64, n)

		for i := 0; i < len(numsStr); i++ {
			nums[i], err = strconv.ParseFloat(numsStr[i], 64)
		}
		for i := 0; i < len(vecStr); i++ {
			vecs[i], err = strconv.ParseFloat(vecStr[i], 64)
		}

		if err != nil {
			panic(err)
		}

		a := make(matrix, n)
		b := make(vector, n)

		for i := range a {
			a[i] = make([]float64, n)
		}

		cnt := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				a[i][j] = nums[cnt]
				cnt++
			}
		}
		for i := 0; i < n; i++ {
			b[i] = vecs[i]
		}

		// индекс, определяет порядок колонок в матрице
		index := make([]int, len(a))
		for i := range index {
			index[i] = i
		}

		// отображаем исходные данные
		a.dump(index, answer)
		b.dump(answer)

		// прямой ход (Зануляю элементы под главной диагональю)
		for i := 0; i < len(a); i++ {

			// главный элемент, значение по умолчанию
			r := a[i][index[i]]

			// если главный элемент равен нулю, нужно найти другой
			// методом перестановки колонок в матрице
			if r == 0 {
				var kk int

				// двигаемся вправо от диагонального элемента, для поиска максимального по модулю элемента
				for k := i; k < len(a); k++ {
					if math.Abs(a[i][index[k]]) > r {
						kk = k
					}
				}

				// если удалось найти главный элемент
				if kk > 0 {
					// меняем местами колонки, так чтобы главный элемент встал в диагональ матрицы
					index[i], index[kk] = index[kk], index[i]
				}

				// получаем главный элемента, текущей строки из диагонали
				r = a[i][index[i]]
			}

			// если главный элемент строки равен 0, метод гаусса не работает
			if r == 0 {
				if b[i] == 0 {
					answer.SetText("система имеет множество решений")
				} else {
					answer.SetText("система не имеет решений")
				}
				flag = true
				break
			}

			// деление элементов текущей строки, на главный элемент
			for j := 0; j < len(a[i]); j++ {
				a[i][index[j]] /= r
			}
			b[i] /= r

			// вычитание текущей строки из всех ниже расположенных строк с занулением I - ого элемента в каждой из них
			for k := i + 1; k < len(a); k++ {
				r = a[k][index[i]]
				for j := 0; j < len(a[i]); j++ {
					a[k][index[j]] = a[k][index[j]] - a[i][index[j]]*r
				}
				b[k] = b[k] - b[i]*r
			}

			// отображаем дамп матрицы A и вектора B
			a.dump(index, answer)
			b.dump(answer)
		}

		var x vector = make(vector, len(b))

		// обратный ход (выражаем x1,x2,x3,xn)
		for i := len(a) - 1; i >= 0; i-- {
			// Задается начальное значение элемента x[I].
			x[i] = b[i]

			// Корректируется искомое значение x[I].
			// В цикле по J от I+1 до N (в случае, когда I=N, этот шаг не выполняется) производятся вычисления x[I]:=  x[I] - x[J]* A[I, J].
			for j := i + 1; j < len(a); j++ {
				x[i] = x[i] - (x[j] * a[i][index[j]])
			}
		}

		if !flag {
			answer.Text = answer.Text + "Вектор X\n"
			for i := 0; i < len(x); i++ {
				answer.Text = answer.Text + fmt.Sprintf("%9.2v ", x[index[i]])
			}
		}

		answer.SetText(answer.Text)
		flag = false
	})

	//box1 := container.NewVBox(label, entry)
	//box2 := container.NewVBox(label1, entry1)
	//box3 := container.NewVBox(label2, entry2, btn)

	w.SetContent(container.NewVBox(
		/*
			box1,
			box2,
			box3,
			answer,
			scr,
		*/
		label, entry,
		label1, entry1,
		label2, entry2, btn,
		answer,
		scr,
	))

	w.ShowAndRun()

}

// отображение дампа матрицы

func (a matrix) dump(index []int, answer *widget.Label) {
	answer.Text = answer.Text + "Матрица\n"
	answer.SetText(answer.Text)
	for i := range a {
		for j := range a[i] {
			if a[i][index[j]] == 0 {
				// необходимо чтобы избавиться от -0
				answer.Text = answer.Text + fmt.Sprintf("%9.2f ", 0.0)
				//fmt.Printf("[0] ")
			} else {
				answer.Text = answer.Text + fmt.Sprintf("%9.2f ", a[i][index[j]])

				//fmt.Printf("[%v] ", a[i][index[j]])
				answer.SetText(answer.Text)
			}
		}
		answer.Text = answer.Text + fmt.Sprint("\n")
		answer.SetText(answer.Text)
	}
	answer.Text = answer.Text + fmt.Sprint("\n")
	answer.SetText(answer.Text)
}

// отображение дампа вектора
func (b vector) dump(answer *widget.Label) {
	answer.Text = answer.Text + "Вектор\n"
	answer.SetText(answer.Text)

	for i := 0; i < len(b); i++ {
		answer.Text = answer.Text + fmt.Sprintf("%9.2v ", b[i])
		answer.SetText(answer.Text)
	}

	answer.Text = answer.Text + fmt.Sprint("\n")
	answer.Text = answer.Text + fmt.Sprint("-----------------------------")
	answer.Text = answer.Text + fmt.Sprint("\n")
	answer.SetText(answer.Text)
}
