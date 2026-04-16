package assets

import (
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed  img shaders
var FS embed.FS

var imageCache = map[string]*ebiten.Image{}

func loadImage(path string) (image.Image, error) {
	f, err := FS.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func loadEbitenImage(path string) (*ebiten.Image, error) {
	img, err := loadImage(path)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func GetImage(path string) (*ebiten.Image, error) {
	img, exist := imageCache[path]
	if exist {
		return img, nil
	}
	img, err := loadEbitenImage(path)
	if err != nil {
		return nil, err
	}
	imageCache[path] = img
	return img, nil
}
