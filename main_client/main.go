package main

import (
	"github.com/TaRosh/online_mover/main_client/game"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	g := game.Game{}
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}
	g.Init()
	g.Run()
}
