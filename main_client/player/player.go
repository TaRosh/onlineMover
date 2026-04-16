package player

import (
	"image/color"
	"math"
	"strconv"

	_ "embed"

	"github.com/TaRosh/online_mover/game/entities"
	"github.com/TaRosh/online_mover/main_client/assets"
	"github.com/TaRosh/online_mover/main_client/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	entities.Player
	img    *ebiten.Image
	shader *ebiten.Shader
	opt    *ebiten.DrawImageOptions
	shOpt  *ebiten.DrawRectShaderOptions
}

func (p *Player) Draw(screen *ebiten.Image, cam *camera.Camera) {
	defer p.opt.GeoM.Reset()
	defer p.shOpt.GeoM.Reset()
	bounds := p.img.Bounds()
	p.opt.GeoM.Translate(-float64(bounds.Dx()/2), -float64(bounds.Dy()/2))
	// angle := p.Velocity.Angle() + math.Pi/2
	p.opt.GeoM.Rotate(float64(p.Rotation) + math.Pi/2)
	p.Clamp(cam)

	p.opt.GeoM.Translate(p.Position.X-cam.Pos.X, p.Position.Y-float64(cam.Pos.Y))
	bond := p.Bounds()
	vector.StrokeRect(screen, float32(bond.Min.X-cam.Pos.X), float32(bond.Min.Y-cam.Pos.Y), float32(bond.Dx()), float32(bond.Dy()), 2, color.Black, false)

	// p.opt.GeoM.Translate(p.Position.X-float64(cam.Rect.Min.X), p.Position.Y-float64(cam.Rect.Min.Y))
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(int(p.ID)), int(p.Position.X+5), int(p.Position.Y+5))

	// p.shOpt.GeoM.Translate(p.Position.X, p.Position.Y)
	p.img.DrawRectShader(bounds.Dx(), bounds.Dy(), p.shader, p.shOpt)
	screen.DrawImage(p.img, p.opt)
}

func (p *Player) Clamp(cam *camera.Camera) {
	minX := cam.Pos.X
	maxX := cam.Pos.X + cam.Width

	minY := cam.Pos.Y
	maxY := cam.Pos.Y + cam.Height

	if p.Position.X-p.Width/2 < minX {
		p.Position.X = minX + p.Width/2
	}
	if p.Position.X+p.Width/2 > maxX {
		p.Position.X = maxX - p.Width/2
	}

	if p.Position.Y-p.Height/2 < minY {
		p.Position.Y = minY + p.Height/2
	}
	if p.Position.Y+p.Height/2 > maxY {
		p.Position.Y = maxY - p.Height/2
	}
}

func New(id entities.PlayerID, c color.Color) (*Player, error) {
	img := ebiten.NewImage(48, 48)
	var err error

	// get 4 image for ship
	shaderOpt, err := setShaderOpt()
	if err != nil {
		return nil, err
	}
	shader, err := assets.GetShader("shaders/player.kage")
	if err != nil {
		return nil, err
	}

	p := Player{
		Player: *entities.NewPlayer(id, 32, 32),
		opt:    &ebiten.DrawImageOptions{},
		img:    img,
		shader: shader,
		shOpt:  shaderOpt,
	}
	return &p, nil
}
