package assets

import "github.com/hajimehoshi/ebiten/v2"

var shaderCache = map[string]*ebiten.Shader{}

func loadShader(path string) (*ebiten.Shader, error) {
	data, err := FS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ebiten.NewShader(data)
}

func GetShader(path string) (*ebiten.Shader, error) {
	shader, exist := shaderCache[path]
	if exist {
		return shader, nil
	}
	shader, err := loadShader(path)
	if err != nil {
		return nil, err
	}
	shaderCache[path] = shader
	return shader, nil
}
