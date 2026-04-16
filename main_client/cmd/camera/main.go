package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	_ "embed"
)

type Game struct{}

var (
	worldImg  *ebiten.Image
	cameraImg *ebiten.Image
	opt       *ebiten.DrawImageOptions
	camera    Camera
)

// Draw implements ebiten.Game.
func (g Game) Draw(screen *ebiten.Image) {
	cameraImg = worldImg.SubImage(camera.img).(*ebiten.Image)
	screen.DrawImage(cameraImg, opt)
}

type Camera struct {
	img    image.Rectangle
	width  int
	height int
}

// Layout implements ebiten.Game.
func (g Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}

// Update implements ebiten.Game.
func (g Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if camera.img.Min.X >= 0 && camera.img.Max.X <= 4000 {
			camera.img = camera.img.Add(image.Pt(camera.width, 0))
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if camera.img.Min.X > 0 {
			camera.img = camera.img.Sub(image.Pt(camera.width, 0))
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		camera.img = camera.img.Add(image.Pt(0, camera.height))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		camera.img = camera.img.Add(image.Pt(0, -camera.height))
	}
	return nil
}

func main() {
	var err error
	_ = err
	worldImg = ebiten.NewImage(4000, 4000)
	worldImg.Fill(color.RGBA{0, 0xff, 0, 0xff})
	w, h := ebiten.WindowSize()
	for j := range 4000 {
		worldImg.Set(4000, j, color.White)
		worldImg.Set(3999, j, color.White)
		worldImg.Set(3998, j, color.White)
		worldImg.Set(3997, j, color.White)
		worldImg.Set(0, j, color.RGBA{0xff, 0, 0, 0xff})
		worldImg.Set(1, j, color.RGBA{0xff, 0, 0, 0xff})
		worldImg.Set(2, j, color.RGBA{0xff, 0, 0, 0xff})
		worldImg.Set(3, j, color.RGBA{0xff, 0, 0, 0xff})
		worldImg.Set(4, j, color.RGBA{0xff, 0, 0, 0xff})
	}
	camera = Camera{
		img: image.Rectangle{
			Min: image.Point{
				X: 0,
				Y: 0,
			},
			Max: image.Point{
				X: w,
				Y: h,
			},
		},
		width: 200, height: 200,
	}
	opt = &ebiten.DrawImageOptions{}
	// for range 1000 {
	// 	worldImg.Set(rand.IntN(4000)+1, rand.IntN(4000), color.RGBA{uint8(rand.UintN(255)), uint8(rand.UintN(255)), uint8(rand.UintN(255)), 0xff})
	// 	worldImg.Set(rand.IntN(4000)-1, rand.IntN(4000), color.RGBA{uint8(rand.UintN(255)), uint8(rand.UintN(255)), uint8(rand.UintN(255)), 0xff})
	// 	worldImg.Set(rand.IntN(4000), rand.IntN(4000), color.RGBA{uint8(rand.UintN(255)), uint8(rand.UintN(255)), uint8(rand.UintN(255)), 0xff})
	// }

	ebiten.RunGame(Game{})
}
