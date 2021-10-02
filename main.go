package main

import (
	"log"
	"math/rand"
	"time"
)

/*

Treasure Hunt!

y
|
########
#......#
#.###..#
#...#.##
#X#....#
######## -- x

# represents an obstacle.
. represents a clear path.
X represents the player’s starting position.

A treasure is hidden within one of the clear path points, and the user must find it.
User have long range of vision, it able to see left/right/up/down whenever the path is unobstructed.
User restricted to only move in 3 direction, up then right then down.
The location of the treasure is located randomly every runtime.
All possible location of the treasure marked as $.

*/

const (
	// entity in map
	path = iota
	treasure
	player
	obstacle
)

var (
	treasureMap        = make(map[[2]int]int)
	treasureMapSizeXY  = [2]int{8, 6}
	playerPositionXY   = [2]int{2, 2}
	listCustomObstacle = [][2]int{
		{3, 2},
		{3, 4},
		{4, 4},
		{5, 4},
		{5, 3},
		{7, 3},
	}
)

func init() {
	generateMap()
	spawnPlayer()
	spawnTreasure()
	renderTreasureMap()
}

func main() {}

func generateMap() {
	generateMapObstacleDefault()
	generateMapObstacleCustom()
}

// generateMapObstacleDefault putting the obstacle on the outer sandbox
func generateMapObstacleDefault() {
	for y := 1; y <= treasureMapSizeXY[1]; y++ {
		for x := 1; x <= treasureMapSizeXY[0]; x++ {
			switch true {
			case x == 1, y == 1, x == treasureMapSizeXY[0], y == treasureMapSizeXY[1]:
				treasureMap[[2]int{x, y}] = obstacle
			default:
				treasureMap[[2]int{x, y}] = path
			}
		}
	}
}

// generateMapObstacleCustom putting the obstacle by predefined location
func generateMapObstacleCustom() {
	for _, customObstacle := range listCustomObstacle {
		treasureMap[customObstacle] = obstacle
	}
}

// spawnTreasure put the treasure in the map randomly. It require several loop to ensure the treasure located on a clear path
func spawnTreasure() {
	var (
		xMin, xMax                           = 1, treasureMapSizeXY[0]
		yMin, yMax                           = 1, treasureMapSizeXY[1]
		treasurePositionX, treasurePositionY int
	)

	for true {
		rand.Seed(time.Now().UnixNano())

		treasurePositionX = rand.Intn(xMax-xMin) + xMin
		treasurePositionY = rand.Intn(yMax-yMin) + yMin

		if treasureMap[[2]int{treasurePositionX, treasurePositionY}] == path {
			treasureMap[[2]int{treasurePositionX, treasurePositionY}] = treasure
			break
		}
	}
}

// spawnPlayer put the player a.k.a X on the treasure map
func spawnPlayer() {
	treasureMap[playerPositionXY] = player
}

// renderTreasureMap display the map in the terminal
func renderTreasureMap() {
	var (
		treasureMapDrawX, treasureMapDrawY string
	)

	for y := 1; y <= treasureMapSizeXY[1]; y++ {
		treasureMapDrawX = "\n"
		for x := 1; x <= treasureMapSizeXY[0]; x++ {
			treasureMapDrawX = treasureMapDrawX + convertIntToEntity(treasureMap[[2]int{x, y}])
		}
		treasureMapDrawY = treasureMapDrawX + treasureMapDrawY

	}

	log.Println(treasureMapDrawY)
}

func convertIntToEntity(code int) string {
	switch code {
	case path:
		return "."
	case obstacle:
		return "#"
	case player:
		return "X"
	case treasure:
		return "$"
	default:
		return "."
	}
}