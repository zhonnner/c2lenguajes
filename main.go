package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	game := Game{
		Screen: screen,
	}
	go game.Run()

	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)
	union := make(chan int)
	union2 := make(chan int)
	go func() {
		for {
			select {
			case <-done:
				game.Screen.Fini()
				os.Exit(0)
				return
			case t := <-ticker.C:
				fmt.Println(t)
				switch event := game.Screen.PollEvent().(type) {
				case *tcell.EventResize:
					game.Screen.Sync()
				case *tcell.EventKey:
					if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
						game.Screen.Fini()
						os.Exit(0)
					} else if event.Key() == tcell.KeyUp && game.snakeBody.Yspeed == 0 {
						go func() {
							union <- 0
							time.Sleep(200 * time.Millisecond)
							game.snakeBody.ChangeDir(-1, 0)
						}()
						go func() {
							x := <-union
							union2 <- x
							time.Sleep(400 * time.Millisecond)
							game.snakeBody2.ChangeDir(-1, 0)
						}()
						go func() {
							fmt.Println(<-union2)
							time.Sleep(600 * time.Millisecond)
							game.snakeBody3.ChangeDir(-1, 0)
						}()

					} else if event.Key() == tcell.KeyDown && game.snakeBody.Yspeed == 0 {
						go func() {
							union <- 0
							time.Sleep(200 * time.Millisecond)
							game.snakeBody.ChangeDir(1, 0)
						}()
						go func() {
							x := <-union
							union2 <- x
							time.Sleep(400 * time.Millisecond)
							game.snakeBody2.ChangeDir(1, 0)
						}()
						go func() {
							fmt.Println(<-union2)
							time.Sleep(600 * time.Millisecond)
							game.snakeBody3.ChangeDir(1, 0)
						}()
					} else if event.Key() == tcell.KeyLeft && game.snakeBody.Xspeed == 0 {
						go func() {
							union <- 0
							time.Sleep(200 * time.Millisecond)
							game.snakeBody.ChangeDir(0, -1)
						}()
						go func() {
							x := <-union
							union2 <- x
							time.Sleep(400 * time.Millisecond)
							game.snakeBody2.ChangeDir(0, -1)
						}()
						go func() {
							fmt.Println(<-union2)
							time.Sleep(600 * time.Millisecond)
							game.snakeBody3.ChangeDir(0, -1)
						}()
					} else if event.Key() == tcell.KeyRight && game.snakeBody.Xspeed == 0 {
						go func() {
							union <- 0
							time.Sleep(200 * time.Millisecond)
							game.snakeBody.ChangeDir(0, 1)
						}()
						go func() {
							x := <-union
							union2 <- x
							time.Sleep(400 * time.Millisecond)
							game.snakeBody2.ChangeDir(0, 1)
						}()
						go func() {
							fmt.Println(<-union2)
							time.Sleep(600 * time.Millisecond)
							game.snakeBody3.ChangeDir(0, 1)
						}()
					} else if event.Rune() == 'y' && game.GameOver {
						go game.Run()
					} else if event.Rune() == 'n' && game.GameOver {
						game.Screen.Fini()
						os.Exit(0)
					}
				}
			}
		}
	}()

	time.Sleep(20000 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")

}
