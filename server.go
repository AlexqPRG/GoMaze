package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Position struct { //создание структуры для хранения текущей позиции
	column int
	row    int
}

var isGame = false //переменная для отслеживания игры
freionfjiewnfjiernfjkwenjfnejkfejkfnjekn

var massive = [10][10]int{{1, 0, 0, 0, 0, 1, 0, 0, 0, 1}, //сам лабиринт (0 - проход, 1 - стена, 2 - выход)
	{1, 0, 0, 0, 1, 0, 0, 1, 1, 1},
	{1, 0, 1, 1, 1, 0, 0, 0, 1, 1},
	{1, 0, 1, 1, 1, 0, 0, 0, 1, 1},
	{1, 0, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 0, 1, 1},
	{1, 0, 0, 0, 1, 0, 1, 0, 1, 1},
	{1, 1, 0, 0, 0, 0, 1, 0, 1, 1},
	{1, 1, 0, 0, 0, 0, 1, 0, 0, 2}}

func main() {

	successMessage := "Вы сместились"
	failedMessage := "больше нельзя смещаться (вы достигли границы)"

	//Устанавливаем начальную позицию с помощью рандома
	// currentPos := Position{rand.Intn(10), rand.Intn(10)}
	currentPos := Position{0, 1} //устанавливаем изначальную позицию

	r := gin.Default()

	r.GET("/help", func(c *gin.Context) { //путь, по которому можно узнать основные команды
		c.JSON(200, gin.H{
			"Команды": "/go?where=up - пойти наверх /go?where=right - пойти направо   /go?where=down - пойти вниз   /go?where=left - пойти налево   /check - узнать обстановку вокруг",
		})
	})

	r.GET("/go", func(c *gin.Context) { //путь, для ходьбы по лабиринту

		if isGame {
			PrintMessage(c, "Вы нашли выход")
			return
		}

		where := c.Query("where")

		switch where {
		case "up":
			if currentPos.column-1 >= 0 && massive[currentPos.column-1][currentPos.row] != 1 {
				currentPos.column--
				PrintMessage(c, successMessage+" ВВЕРХ")
			} else {
				PrintMessage(c, "Вверх "+failedMessage)
			}

		case "down":
			if currentPos.column+1 <= 9 && massive[currentPos.column+1][currentPos.row] != 1 {
				currentPos.column++
				PrintMessage(c, successMessage+" ВНИЗ")
			} else {
				PrintMessage(c, "Вниз "+failedMessage)
			}

		case "right":
			if currentPos.row+1 <= 9 && massive[currentPos.column][currentPos.row+1] != 1 {
				currentPos.row++
				PrintMessage(c, successMessage+" ВПРАВО")
			} else {
				PrintMessage(c, "Вправо "+failedMessage)
			}

		case "left":
			if currentPos.row-1 >= 0 && massive[currentPos.column][currentPos.row-1] != 1 {
				currentPos.row--
				PrintMessage(c, successMessage+" ВЛЕВО")
			} else {
				PrintMessage(c, "Влево "+failedMessage)
			}

		default:
			c.JSON(200, gin.H{
				"message": "Такого действия нет",
			})
		}
		CheckExit(currentPos, c)

	})

	r.GET("/pos", func(c *gin.Context) { //узнаем текущую позицию

		message1 := fmt.Sprintf("Позиция по вертикали: %d, Позиция по горизонтали: %d", currentPos.column, currentPos.row)

		c.JSON(200, gin.H{
			"message": message1,
		})
	})

	r.GET("/check", func(c *gin.Context) { //путь, чтобы узнать объекты вокруг

		var up, right, down, left int = 3, 3, 3, 3

		if currentPos.column-1 >= 0 {
			up = massive[currentPos.column-1][currentPos.row]
		}

		if currentPos.row+1 <= 9 {
			right = massive[currentPos.column][currentPos.row+1]
		}

		if currentPos.column+1 <= 9 {
			down = massive[currentPos.column+1][currentPos.row]
		}

		if currentPos.row-1 >= 0 {
			left = massive[currentPos.column][currentPos.row-1]
		}

		names := Naming(up, right, down, left)

		c.JSON(200, gin.H{
			"СВЕРХУ": names[0],
			"СПРАВА": names[1],
			"СНИЗУ":  names[2],
			"СЛЕВА":  names[3],
		})
	})

	r.Run()

}

func Naming(numbers ...int) [4]string { //функция, которая возвращает названия объектов вокруг
	return_values := [4]string{}
	for i, number := range numbers {
		switch number {
		case 0:
			return_values[i] = "Проход"
		case 1:
			return_values[i] = "Стена"
		case 2:
			return_values[i] = "Выход (ПОБЕДА!)"
		case 3:
			return_values[i] = "Пустота (туда нельзя)"
		}
	}
	return return_values
}

func CheckExit(position Position, c *gin.Context) { //функция, которая проверяет не дошел ли пользователя до выхода
	if massive[position.column][position.row] == 2 {
		isGame = true

		c.JSON(200, gin.H{
			"MESSAGE": "YRA",
		})
	}
}

func PrintMessage(c *gin.Context, message string) { //функция, которая присылает ответы в JSON с message в роли содержимого
	c.JSON(200, gin.H{
		"message": message,
	})
}

