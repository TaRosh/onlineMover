module github.com/TaRosh/online_mover/main_client

go 1.26.1

require (
	github.com/TaRosh/online_mover/game v0.0.0-00010101000000-000000000000
	github.com/TaRosh/online_mover/network v0.0.0-00010101000000-000000000000
	github.com/hajimehoshi/ebiten/v2 v2.9.9
	github.com/joho/godotenv v1.5.1
	github.com/quasilyte/gmath v0.0.0-20250817142619-e0a8c6ee09b3
	github.com/stretchr/testify v1.11.1
	github.com/yohamta/ganim8/v2 v2.1.30
)

require (
	github.com/TaRosh/online_mover/udp v0.0.0-00010101000000-000000000000 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20221126150942-6ab00d035af9 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/TaRosh/online_mover/game => ../game

replace github.com/TaRosh/online_mover/network => ../network

replace github.com/TaRosh/online_mover/udp => ../udp
