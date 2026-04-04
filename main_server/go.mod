module github.com/TaRosh/online_mover/main_server

go 1.26.1

require (
	github.com/TaRosh/online_mover/game v0.0.0-00010101000000-000000000000
	github.com/TaRosh/online_mover/network v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/quasilyte/gmath v0.0.0-20250817142619-e0a8c6ee09b3
)

require github.com/TaRosh/online_mover/udp v0.0.0-00010101000000-000000000000 // indirect

replace github.com/TaRosh/online_mover/game => ../game

replace github.com/TaRosh/online_mover/udp => ../udp

replace github.com/TaRosh/online_mover/network => ../network
