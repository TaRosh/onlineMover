package bullet

import (
	"image/color"
	"math"
	"strconv"
	"time"

	"github.com/TaRosh/online_mover/game/entities"
	"github.com/TaRosh/online_mover/main_client/assets"
	"github.com/TaRosh/online_mover/main_client/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/ganim8/v2"
)

type Bullet struct {
	entities.Bullet
	animation *ganim8.Animation
	img       *ebiten.Image
	opt       *ebiten.DrawImageOptions
}

func (b *Bullet) Update() {
	b.Bullet.Update()
	b.animation.Update()
}

func (b *Bullet) Draw(screen *ebiten.Image, cam *camera.Camera) {
	b.opt.GeoM.Reset()
	b.animation.Draw(screen, ganim8.DrawOpts(b.Position.X-cam.Pos.X, b.Position.Y-cam.Pos.Y, float64(b.Rotation)+math.Pi/2))
	// b.opt.GeoM.Translate(b.Position.X, b.Position.Y)
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(int(b.ID)), int(b.Position.X+5), int(b.Position.Y+5))
	// screen.DrawImage(b.img, b.opt)
}

func NewBullet(id uint32, c color.Color) *Bullet {
	img, err := assets.GetImage("img/bullets/bullet.png")
	if err != nil {
		panic(err)
	}
	g32 := ganim8.NewGrid(32, 32, 96, 32)
	anim := ganim8.New(img, g32.Frames("1-3", 1), 100*time.Millisecond)
	b := Bullet{
		Bullet: *entities.NewBullet(id, 32, 32),
		// img:       bulletI,
		opt:       &ebiten.DrawImageOptions{},
		animation: anim,
	}
	return &b
}
