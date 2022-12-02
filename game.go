package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
)

type Game struct {
	Screen     tcell.Screen
	snakeBody  SnakeBody
	snakeBody2 SnakeBody
	snakeBody3 SnakeBody
	FoodPos    Part
	Score      int
	GameOver   bool
}

func worker(id int) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func drawParts(s tcell.Screen, snakeParts []Part, foodPos Part, snakeStyle tcell.Style, foodStyle tcell.Style) {
	s.SetContent(foodPos.X, foodPos.Y, '\u25CF', nil, foodStyle)
	for _, part := range snakeParts {
		s.SetContent(part.X, part.Y, ' ', nil, snakeStyle)
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, text string) {
	row := y1
	col := x1
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func checkCollision(parts []Part, otherPart Part) bool {
	for _, part := range parts {
		if part.X == otherPart.X && part.Y == otherPart.Y {
			return true
		}
	}
	return false
}

func (g *Game) UpdateFoodPos(width int, height int) {
	g.FoodPos.X = rand.Intn(width)
	g.FoodPos.Y = rand.Intn(height)
	if g.FoodPos.Y == 1 && g.FoodPos.X < 10 {
		g.UpdateFoodPos(width, height)
	}
}

func (g *Game) Run() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.Screen.SetStyle(defStyle)
	width, height := g.Screen.Size()
	g.snakeBody.ResetPos(width, height)
	g.snakeBody2.ResetPos(width+30, height+10)
	g.snakeBody3.ResetPos(width-30, height-10)
	g.UpdateFoodPos(width, height)
	g.GameOver = false
	g.Score = 0
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)
	for {
		longerSnake := false
		g.Screen.Clear()
		if checkCollision(g.snakeBody.Parts[len(g.snakeBody.Parts)-1:], g.FoodPos) {
			g.UpdateFoodPos(width, height)
			longerSnake = true
			g.Score++
		}
		g.snakeBody.Update(width, height, longerSnake)
		g.snakeBody2.Update(width, height, longerSnake)
		g.snakeBody3.Update(width, height, longerSnake)
		drawParts(g.Screen, g.snakeBody.Parts, g.FoodPos, snakeStyle, defStyle)
		drawParts(g.Screen, g.snakeBody2.Parts, g.FoodPos, snakeStyle, defStyle)
		drawParts(g.Screen, g.snakeBody3.Parts, g.FoodPos, snakeStyle, defStyle)
		str1 := strconv.Itoa(g.snakeBody.Parts[0].X) + "," + strconv.Itoa(g.snakeBody.Parts[0].Y)
		str2 := strconv.Itoa(g.snakeBody.Parts[1].X) + " ," + strconv.Itoa(g.snakeBody.Parts[1].Y)
		str4 := strconv.Itoa(g.snakeBody.Parts[2].X) + " ," + strconv.Itoa(g.snakeBody.Parts[2].Y)
		str3 := str1 + " " + str2 + " " + str4
		drawText(g.Screen, 1, 1, 8+len(str3), 1, str3)
		time.Sleep(40 * time.Millisecond)
		g.Screen.Show()
	}
	g.GameOver = true
	drawText(g.Screen, width/2-20, height/2, width/2+20, height/2, "Game Over, Score: "+strconv.Itoa(g.Score)+", Play Again? y/n")
	g.Screen.Show()
}
