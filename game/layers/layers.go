package layers

import quadtree "github.com/TaRosh/online_mover/quad_tree"

const (
	LayerPlayer   quadtree.Layer = 1 << 0
	LayerBullet   quadtree.Layer = 1 << 1
	LayerAsteroid quadtree.Layer = 1 << 2
)
