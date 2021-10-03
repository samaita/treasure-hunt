package main

import (
	"fmt"
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

	// direction
	up = iota
	right
	down
)

var (
	treasureMap              = make(map[[2]int]int)
	treasureMapOriginalState = make(map[[2]int]int)
	treasureMapSizeXY        = [2]int{8, 6}
	playerPositionXY         = [2]int{2, 2}
	playerPosibleMove        = [][2]int{}
	listCustomObstacle       = [][2]int{
		{3, 2},
		{3, 4},
		{4, 4},
		{5, 4},
		{5, 3},
		{7, 3},
	}
	directionOrder = []int{up, right, down}
	directionTaken = map[int]bool{
		up:    false,
		right: false,
		down:  false,
	}
)

func init() {
	generateMap()
	spawnPlayer()
	spawnTreasure()
}

func main() {
	var (
		treasureFound      bool
		treasurePositionXY [2]int
	)

	// print initial condition
	renderToTerminal()

	// hide the treasure
	setPossibleTreasure()

	// print map with hidden treasure
	renderToTerminal()

	for !treasureFound {
		treasurePositionXY, treasureFound = checkTreasure()
		if treasureFound {
			renderToTerminal()
			revealMap(treasurePositionXY)
			renderToTerminal()
			break
		} else {
			renderToTerminal()
			movePlayer()
			renderToTerminal()
		}

	}
}

// movePlayer update the coordinate of the player limited by predefined direction
func movePlayer() {

}

func revealMap(treasurePositionXY [2]int) {
	for coordinate := range treasureMap {
		if treasureMap[coordinate] == treasure && coordinate != treasurePositionXY {
			treasureMap[coordinate] = path
		}
	}
}

// checkTreasure check all unobstructed line of X & Y from player position
func checkTreasure() ([2]int, bool) {
	var (
		startX, startY               = playerPositionXY[0], playerPositionXY[1]
		maxSightX, maxSightY         = treasureMapSizeXY[0], treasureMapSizeXY[1]
		foundTreasure, foundObstacle = false, false
	)

	checkPositive, checkNegative := true, true
	for x := startX + 1; x <= maxSightX; x++ {
		currentSightXPositive := [2]int{x, startY}
		currentSightXNegative := [2]int{-x, startY}

		log.Println("x", currentSightXPositive, currentSightXNegative)

		if checkPositive {
			foundObstacle, foundTreasure = checkMapEntity(currentSightXPositive)
			if foundTreasure {
				return currentSightXPositive, foundTreasure
			}
			checkPositive = !foundObstacle
		}

		if checkNegative {
			foundObstacle, foundTreasure = checkMapEntity(currentSightXNegative)
			if foundTreasure {
				return currentSightXNegative, foundTreasure
			}
			checkNegative = !foundObstacle
		}

		if !checkPositive && !checkNegative {
			x = maxSightX
		}
	}

	checkPositive, checkNegative = true, true
	for y := startY + 1; y <= maxSightY; y++ {
		currentSightYPositive := [2]int{startX, y}
		currentSightYNegative := [2]int{startX, -y}

		log.Println("y", currentSightYPositive, currentSightYNegative)

		if checkPositive {
			foundObstacle, foundTreasure = checkMapEntity(currentSightYPositive)
			if foundTreasure {
				return currentSightYPositive, foundTreasure
			}
			checkPositive = !foundObstacle
		}

		if checkNegative {
			foundObstacle, foundTreasure = checkMapEntity(currentSightYNegative)
			if foundTreasure {
				return currentSightYNegative, foundTreasure
			}
			checkNegative = !foundObstacle
		}

		log.Println("y", checkPositive, checkNegative)
		log.Println("y", treasureMapOriginalState[currentSightYPositive], treasureMapOriginalState[currentSightYNegative])

		if !checkPositive && !checkNegative {
			y = maxSightY
		}
	}

	return [2]int{}, false
}

func checkMapEntity(position [2]int) (bool, bool) {
	switch treasureMapOriginalState[position] {
	case path:
		treasureMap[position] = path
	case treasure:
		return false, true
	case obstacle:
		return true, false
	default:
		return true, false
	}

	return false, false
}

func renderToTerminal() {
	fmt.Println("\033[2J")
	fmt.Println(drawTreasureMap())

	time.Sleep(1 * time.Second)
}

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
		treasurePositionXY                   [2]int
	)

	for true {
		rand.Seed(time.Now().UnixNano())

		treasurePositionX = rand.Intn(xMax-xMin) + xMin
		treasurePositionY = rand.Intn(yMax-yMin) + yMin
		treasurePositionXY = [2]int{treasurePositionX, treasurePositionY}

		if treasureMap[treasurePositionXY] == path {
			treasureMap[treasurePositionXY] = treasure
			break
		}
	}
}

// setPossibleTreasure hide the actual treasure coordinate to be exposed later
func setPossibleTreasure() {
	for coordinate := range treasureMap {
		treasureMapOriginalState[coordinate] = treasureMap[coordinate]
		if treasureMap[coordinate] == path {
			treasureMap[coordinate] = treasure
		}
	}
}

// spawnPlayer put the player a.k.a X on the treasure map
func spawnPlayer() {
	treasureMap[playerPositionXY] = player
}

// drawTreasureMap draw the map in a ready print a,.,,
func drawTreasureMap() string {
	var (
		treasureMapDrawX, treasureMapDrawY string
	)

	for y := 1; y <= treasureMapSizeXY[1]; y++ {
		treasureMapDrawX = ""
		if y < treasureMapSizeXY[1] {
			treasureMapDrawX = "\n"
		}
		for x := 1; x <= treasureMapSizeXY[0]; x++ {
			treasureMapDrawX = treasureMapDrawX + convertIntToEntity(treasureMap[[2]int{x, y}])
		}
		treasureMapDrawY = treasureMapDrawX + treasureMapDrawY

	}

	return treasureMapDrawY
}

// convertIntToEntity convert code constant of entity to a map drawn entity
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
