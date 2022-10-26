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
	w.Resize(fyne.NewSize(300, 600))
	w.CenterOnScreen()

	label := widget.NewLabel("Введите размер матрицы ")
	entry := widget.NewEntry()
	label1 := widget.NewLabel("Введите матрицу ")
	entry1 := widget.NewEntry()
	label2 := widget.NewLabel("Векторы")
	entry2 := widget.NewEntry()
	label3 := widget.NewLabel("Приближение(только для метода Зейделя)")
	entry3 := widget.NewEntry()

	answer := widget.NewLabel("")

	scr := container.NewVScroll(answer)
	scr.SetMinSize(fyne.NewSize(300, 300))

	btn1 := widget.NewButton("Посчитать методом Гаусса", func() {
		gauss(answer, entry, entry1, entry2)
	})

	//-------------------------------------------------------------------------------------------------

	btn2 := widget.NewButton("Посчитать методом Зейделя", func() {
		seidel(answer, entry, entry1, entry2, entry3)
	})

	w.SetContent(container.NewVBox(
		label, entry,
		label1, entry1,
		label2, entry2,
		btn1,
		label3, entry3,
		btn2,
		scr,
	))

	w.ShowAndRun()

}

func (a matrix) printGauss(index []int, answer *widget.Label, b vector) {
	answer.Text = answer.Text + "Матрица\n"
	answer.SetText(answer.Text)
	for i := range a {
		for j := range a[i] {
			if a[i][index[j]] == -0 || math.Abs(a[i][index[j]]) < 0.00000001 {
				// необходимо чтобы избавиться от -0
				answer.Text = answer.Text + fmt.Sprintf("%9f ", 0.0)
			} else {
				answer.Text = answer.Text + fmt.Sprintf("%9f ", a[i][index[j]])
				answer.SetText(answer.Text)
			}
		}
		answer.Text = answer.Text + fmt.Sprintf("%9v ", "")
		if b[index[i]] == -0 || math.Abs(b[index[i]]) < 0.00000001 {
			answer.Text = answer.Text + fmt.Sprintf("%9f ", 0.0)
		} else {
			answer.Text = answer.Text + fmt.Sprintf("%9f ", b[index[i]])
		}
		answer.SetText(answer.Text)
		answer.Text = answer.Text + fmt.Sprint("\n")
		answer.SetText(answer.Text)
	}
	answer.Text = answer.Text + fmt.Sprint("\n")
	answer.Text = answer.Text + fmt.Sprint("-----------------------------")
	answer.Text = answer.Text + fmt.Sprint("\n")
	answer.SetText(answer.Text)
	answer.SetText(answer.Text)
}

func (a matrix) printSeidel(answer *widget.Label) {
	answer.Text = answer.Text + "Матрица\n"
	answer.SetText(answer.Text)
	for i := range a {
		for j := 0; j < len(a[i])-1; j++ {
			if a[i][j] == -0 || math.Abs(a[i][j]) < 0.00000001 {
				// необходимо чтобы избавиться от -0
				answer.Text = answer.Text + fmt.Sprintf("%9f ", 0.0)
			} else {
				answer.Text = answer.Text + fmt.Sprintf("%9f ", a[i][j])
				answer.SetText(answer.Text)
			}
		}
		answer.Text = answer.Text + fmt.Sprintf("%9v ", "")
		if a[i][len(a)] == -0 || math.Abs(a[i][len(a)]) < 0.00000001 {
			answer.Text = answer.Text + fmt.Sprintf("%9f ", 0.0)
		} else {
			answer.Text = answer.Text + fmt.Sprintf("%9f ", a[i][len(a)])
			answer.SetText(answer.Text)
			answer.Text = answer.Text + fmt.Sprint("\n")
			answer.SetText(answer.Text)
		}
	}
	answer.Text = answer.Text + fmt.Sprint("\n")
	answer.Text = answer.Text + fmt.Sprint("-----------------------------")
	answer.Text = answer.Text + fmt.Sprint("\n")
	answer.SetText(answer.Text)
	answer.SetText(answer.Text)
}

func gauss(answer *widget.Label, entry, entry1, entry2 *widget.Entry) {
	flagNoSolutions := false
	flagManySolutions := false
	answer.SetText("")
	n, err := strconv.Atoi(entry.Text)
	if err != nil {
		panic(err)
	}
	splitFunc := func(r rune) bool {
		return strings.ContainsRune(", \r", r)
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
	a.printGauss(index, answer, b)

	// прямой ход (Зануляю элементы под главной диагональю)
	for i := 0; i < len(a); i++ {

		// главный элемент, значение по умолчанию
		r := a[i][index[i]]

		// если главный элемент равен нулю, нужно найти другой
		// методом перестановки колонок в матрице
		if r == 0 {
			maxEl := 0.0
			maxId := 0

			for j := i; j < n; j++ {
				if math.Abs(a[j][i]) > maxEl {
					maxEl = math.Abs(a[j][i])
					maxId = j
				}
			}
			a[i], a[maxId] = a[maxId], a[i]
			b[i], b[maxId] = b[maxId], b[i]

			maxEl = 0
			maxId = 0

		}
		zero := 0.0
		// если главный элемент строки равен 0, метод гаусса не работает
		for k := 0; k < n; k++ {
			zero = a[k][k]
			if zero == 0.0 {
				if b[k] == 0.0 {
					flagManySolutions = true
					answer.SetText("система имеет множество решений")
				} else {
					answer.SetText("система не имеет решений")
					flagNoSolutions = true
					break
				}

			}
		}

		if flagNoSolutions {
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
		a.printGauss(index, answer, b)
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

	if !flagNoSolutions {
		answer.Text = answer.Text + "Вектор X\n"
		for i := 0; i < len(x); i++ {
			if x[index[i]] == 0.0 || math.Abs(x[index[i]]) < 0.00000001 {
				answer.Text = answer.Text + fmt.Sprintf("%9f ", 0.0)
			} else {
				answer.Text = answer.Text + fmt.Sprintf("%9f ", x[index[i]])
			}
		}
	}
	if flagManySolutions {
		answer.Text = answer.Text + fmt.Sprintf("Система имеет множество решений\n")
	}
	answer.SetText(answer.Text)
	flagNoSolutions = false
}

func seidel(answer *widget.Label, entry, entry1, entry2, entry3 *widget.Entry) {
	flag := false
	var eps float64
	answer.SetText("")
	n, err := strconv.Atoi(entry.Text)
	if err != nil {
		panic(err)
	}
	splitFunc := func(r rune) bool {
		return strings.ContainsRune(", \r", r)
	}

	numsStr := strings.FieldsFunc(entry1.Text, splitFunc)
	vecStr := strings.FieldsFunc(entry2.Text, splitFunc)
	epsStr := entry3.Text

	nums := make([]float64, n*n)
	vecs := make([]float64, n)

	eps, err = strconv.ParseFloat(epsStr, 64)

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

	for i := range a {
		a[i] = make([]float64, n+1)
	}

	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			a[i][j] = nums[cnt]
			cnt++
		}
		a[i][n] = vecs[i]
	}
	//----------------------------------------------
	a.printSeidel(answer)
	maxEl := 0.0
	maxId := 0

	r := a[0][0]

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if math.Abs(a[j][i]) > maxEl {
				maxEl = math.Abs(a[j][i])
				maxId = j
			}

		}
		a[i], a[maxId] = a[maxId], a[i]

		maxEl = 0
		maxId = 0
		r = a[i][i]
		if r == 0 {
			if a[i][len(a)] == 0 {
				answer.SetText("система имеет множество решений")
			} else {
				answer.SetText("система не имеет решений")
			}
			flag = true
			break
		}

	}
	if !flag {
		a.printSeidel(answer)

		// Введем вектор значений неизвестных на предыдущей итерации,
		// размер которого равен числу строк в матрице, т.е. size,
		// причем согласно методу изначально заполняем его нулями
		previousVariableValues := make(vector, n)
		for i := 0; i < n; i++ {
			previousVariableValues[i] = 0.0
		}

		for {
			// Введем вектор значений неизвестных на текущем шаге
			currentVariableValues := make(vector, n)
			// Посчитаем значения неизвестных на текущей итерации
			// в соответствии с теоретическими формулами
			for i := 0; i < n; i++ {
				// Инициализируем i-ую неизвестную значением
				// свободного члена i-ой строки матрицы
				currentVariableValues[i] = a[i][n]
				// Вычитаем сумму по всем отличным от i-ой неизвестным
				for j := 0; j < n; j++ {
					// При j < i можем использовать уже посчитанные
					// на этой итерации значения неизвестных
					if j < i {
						currentVariableValues[i] -= a[i][j] * currentVariableValues[j]
					}

					// При j > i используем значения с прошлой итерации
					if j > i {
						currentVariableValues[i] -= a[i][j] * previousVariableValues[j]
					}
				}
				currentVariableValues[i] /= a[i][i]
				if previousVariableValues[i] == 0.0 || math.Abs(previousVariableValues[i]) < 0.00000001 {
					previousVariableValues[i] = 0.0
				}
				answer.Text = answer.Text + "Вектора невязки\n"

				for i := 0; i < n; i++ {
					if previousVariableValues[i] == 0.0 || math.Abs(previousVariableValues[i]) < 0.00000001 {
						answer.Text = answer.Text + fmt.Sprintf("%9f ", 0.0)
					} else {
						answer.Text = answer.Text + fmt.Sprintf("%9f ", previousVariableValues[i])
					}

				}
				answer.Text = answer.Text + "\n"
				answer.SetText(answer.Text)
			}

			errEps := 0.0

			for i := 0; i < n; i++ {
				errEps += math.Abs(currentVariableValues[i] - previousVariableValues[i])
			}

			// Если необходимая точность достигнута, то завершаем процесс
			if errEps < eps {
				break
			}

			// Переходим к следующей итерации, так
			// что текущие значения неизвестных
			// становятся значениями на предыдущей итерации
			previousVariableValues = currentVariableValues
		}

		answer.Text = answer.Text + "Вектор X\n"

		for i := 0; i < n; i++ {
			if previousVariableValues[i] == 0.0 || math.Abs(previousVariableValues[i]) < 0.00000001 {
				answer.Text = answer.Text + fmt.Sprintf("%9f ", 0.0)
			} else {
				answer.Text = answer.Text + fmt.Sprintf("%9f ", previousVariableValues[i])
			}

		}
	}
	answer.SetText(answer.Text)
	flag = false
}
