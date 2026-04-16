package asteroid

import (
	"fmt"
	"image/color"
	"time"

	"github.com/TaRosh/online_mover/game/entities"
	"github.com/TaRosh/online_mover/main_client/assets"
	"github.com/TaRosh/online_mover/main_client/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Asteroid struct {
	entities.Asteroid
	img    *ebiten.Image
	shader *ebiten.Shader
	opt    *ebiten.DrawImageOptions
	shOpt  *ebiten.DrawRectShaderOptions
	t      time.Time
}

func (a *Asteroid) Draw(screen *ebiten.Image, cam *camera.Camera) {
	defer a.opt.GeoM.Reset()
	defer a.shOpt.GeoM.Reset()
	bounds := a.img.Bounds()
	a.opt.GeoM.Translate(-float64(bounds.Dx()/2), -float64(bounds.Dy()/2))
	// swh := screen.Bounds()
	a.opt.GeoM.Translate(a.Position.X-cam.Pos.X, a.Position.Y-cam.Pos.Y)
	// a.opt.GeoM.Rotate(float64(a.Rotation) + math.Pi/2)
	a.img.DrawRectShader(bounds.Dx(), bounds.Dy(), a.shader, a.shOpt)
	// a.opt.GeoM.Translate(float64(swh.Dx())/2, float64(swh.Dy())/2)

	screen.DrawImage(a.img, a.opt)
	bond := a.Bounds()
	vector.StrokeRect(screen, float32(bond.Min.X-cam.Pos.X), float32(bond.Min.Y-cam.Pos.Y), float32(bond.Dx()), float32(bond.Dy()), 2, color.Black, false)
}

func (a *Asteroid) Update() {
	a.Asteroid.Update()
	a.shOpt.Uniforms["Time"] = float32(time.Since(a.t).Seconds())
}

func New(id uint32, r float64, c color.Color) (*Asteroid, error) {
	fmt.Println("MAKE ONE")
	img := ebiten.NewImage(int(r), int(r))
	// img.Fill(color.RGBA{0, 0xff, 0, 0})
	var err error
	t := time.Now()
	shaderOpt, err := setShaderOpt(float32(r / 2))
	if err != nil {
		return nil, err
	}
	shaderOpt.Uniforms["Time"] = t
	shader, err := assets.GetShader("shaders/asteroid.kage")
	if err != nil {
		return nil, err
	}
	a := Asteroid{
		Asteroid: *entities.NewAsteroid(id, r, r),
		opt:      &ebiten.DrawImageOptions{},
		img:      img,
		shader:   shader,
		shOpt:    shaderOpt,
		t:        t,
	}
	return &a, nil
}
