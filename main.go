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

# represents anentity_obstacle.
. represents a clear entity_path.
X represents the entity_player’s starting position.

A entity_treasure is hidden within one of the clear entity_path points, and the user must find it.
User have long range of vision, it able to see left/right/up/down whenever the entity_path is unobstructed.
User restricted to only move in 3 direction, up then right then down.
The location of the entity_treasure is located randomly every runtime.
All possible location of the entity_treasure marked as $.

*/
const ( // direction
	up = iota
	right
	down
	stuck
)

const (
	// entity in map
	entity_path = iota + 1
	entity_treasure
	entity_player
	entity_obstacle
)

const (
	// axis
	axis_x = iota
	axis_y
)

const (
	delay_time = 500 // time used to display a step taken for each movement before continuing
)

var (
	mapSize             = [2]int{8, 6} // define maximum size of the map [X,Y]
	playerStartPosition = [2]int{2, 2} // define default location of the player
	listCustomObstacle  = [][2]int{
		{3, 2},
		{3, 4},
		{4, 4},
		{5, 4},
		{5, 3},
		{7, 3},
	}
)

type Player struct {
	Position       [2]int
	DirectionTaken int
	FoundTreasure  bool
}

type Treasure struct {
	Position [2]int
}

type TreasureMap struct {
	Size            [2]int
	OriginalMapping map[[2]int]int
	Mapping         map[[2]int]int
}

func main() {
	player := NewPlayer()
	treasure := NewTreasure()
	treasureMap := NewTreasureMap(mapSize)

	treasureMap.createMap(listCustomObstacle)
	treasureMap.setEntity(entity_player, player.Position)
	for true {
		treasure.randomizePosition(mapSize[0], mapSize[1])
		if treasureMap.setEntity(entity_treasure, treasure.Position) {
			break
		}
	}
	treasureMap.render() // display initial condition with treasure

	treasureMap.setPossibleTreasure()
	treasureMap.render() // display map with possible treasure location

	for !player.FoundTreasure {
		// player see unobstructed path, and determine which is treasure and which is path
		treasurePositionXY, listPathPosition := player.see(treasureMap)
		for _, pathPosition := range listPathPosition {
			treasureMap.setEntity(entity_path, pathPosition)
		}
		treasureMap.render()

		if !player.FoundTreasure {
			// keep moving until found the treasure
			newPosition, _ := player.move(treasureMap.Mapping)
			oldPosition := player.Position

			// stop, when player cannot move any longer
			if newPosition == oldPosition {
				break
			}

			// move the player into new position, put path on the older position
			treasureMap.setEntity(entity_path, oldPosition)
			treasureMap.setEntity(entity_player, newPosition)
			treasureMap.render()

			// update player position
			player.setPosition(newPosition)
		} else {
			treasureMap.revealMap(treasurePositionXY)
			treasureMap.render()
			break
		}

	}
}

// NewPlayer creating a new player with initial position
func NewPlayer() Player {
	return Player{
		Position: playerStartPosition,
	}
}

// setPosition update the value of player position
func (p *Player) setPosition(newPosition [2]int) {
	p.Position = newPosition
}

// move update the coordinate of the entity_player limited by predefined direction
func (p *Player) move(treasureMap map[[2]int]int) ([2]int, bool) {

	if p.DirectionTaken == up {
		newPlayerPositionXY := [2]int{p.Position[0], p.Position[1] + 1}
		if treasureMap[newPlayerPositionXY] == entity_obstacle {
			p.DirectionTaken = right
		} else {
			return newPlayerPositionXY, true
		}
	}

	if p.DirectionTaken == right {
		newPlayerPositionXY := [2]int{p.Position[0] + 1, p.Position[1]}
		if treasureMap[newPlayerPositionXY] == entity_obstacle {
			p.DirectionTaken = down
		} else {
			return newPlayerPositionXY, true
		}
	}

	if p.DirectionTaken == down {
		newPlayerPositionXY := [2]int{p.Position[0], p.Position[1] - 1}
		if treasureMap[newPlayerPositionXY] == entity_obstacle {
			p.DirectionTaken = stuck
		} else {
			return newPlayerPositionXY, true
		}
	}

	return p.Position, false
}

// see check all unobstructed line of X & Y from entity_player position
func (p *Player) see(treasureMap TreasureMap) ([2]int, [][2]int) {
	var (
		startX, startY                  = p.Position[0], p.Position[1]
		treasurePosition, treasureFound [2]int
		listPathPosition, pathFound     [][2]int
	)

	// see all entity in x axis with same y axis
	treasurePosition, pathFound = checkMap(treasureMap, startX+1, startY, 1, axis_x)
	if treasureMap.OriginalMapping[treasurePosition] == entity_treasure {
		treasureFound = treasurePosition
	}
	listPathPosition = append(listPathPosition, pathFound...)

	// see all entity in -x axis with same y axis
	treasurePosition, pathFound = checkMap(treasureMap, startX-1, startY, -1, axis_x)
	if treasureMap.OriginalMapping[treasurePosition] == entity_treasure {
		treasureFound = treasurePosition
	}
	listPathPosition = append(listPathPosition, pathFound...)

	// see all entity in y axis with same x axis
	treasurePosition, pathFound = checkMap(treasureMap, startY+1, startX, 1, axis_y)
	if treasureMap.OriginalMapping[treasurePosition] == entity_treasure {
		treasureFound = treasurePosition
	}
	listPathPosition = append(listPathPosition, pathFound...)

	// see all entity in -y axis with same x axis
	treasurePosition, pathFound = checkMap(treasureMap, startY-1, startX, -1, axis_y)
	if treasureMap.OriginalMapping[treasurePosition] == entity_treasure {
		treasureFound = treasurePosition
	}
	listPathPosition = append(listPathPosition, pathFound...)

	if treasureMap.OriginalMapping[treasureFound] == entity_treasure {
		p.FoundTreasure = true
	}

	return treasureFound, listPathPosition
}

// checkMap a shorthand to validate an unobstructed line of sight in original mapping
func checkMap(treasureMap TreasureMap, startAxis int, staticAxis int, addValue int, typeAxis int) ([2]int, [][2]int) {
	var (
		check            = true
		treasurePosition [2]int
		pathPosition     [][2]int
		currentPosition  [2]int
	)
	for check {
		if typeAxis == axis_x {
			currentPosition = [2]int{startAxis, staticAxis}
		} else {
			currentPosition = [2]int{staticAxis, startAxis}
		}

		if check {
			switch treasureMap.OriginalMapping[currentPosition] {
			case entity_path:
				pathPosition = append(pathPosition, currentPosition)
			case entity_treasure:
				treasurePosition = currentPosition
			case entity_obstacle:
				check = false
			default:
				check = false
			}
		}
		startAxis += addValue
	}

	return treasurePosition, pathPosition
}

// NewTreasure creating a new blank
func NewTreasure() Treasure {
	return Treasure{}
}

// randomizePosition put the entity_treasure in the map randomly. It require several loop to ensure the entity_treasure located on a clear entity_path
func (t *Treasure) randomizePosition(sizeX, sizeY int) {
	var (
		xMin, xMax                           = 1, sizeX
		yMin, yMax                           = 1, sizeY
		treasurePositionX, treasurePositionY int
		treasurePositionXY                   [2]int
	)

	rand.Seed(time.Now().UnixNano())

	treasurePositionX = rand.Intn(xMax-xMin) + xMin
	treasurePositionY = rand.Intn(yMax-yMin) + yMin
	treasurePositionXY = [2]int{treasurePositionX, treasurePositionY}

	t.Position = treasurePositionXY
}

// NewTreasureMap creating a new blank treasure map
func NewTreasureMap(size [2]int) TreasureMap {
	return TreasureMap{
		Size:            size,
		Mapping:         make(map[[2]int]int),
		OriginalMapping: make(map[[2]int]int),
	}
}

// render display of the mapping, not the original mapping. It also print the info of list possible location of the treasure.
func (tm *TreasureMap) render() {
	var (
		treasureMapDrawPerLine, treasureMapDrawComplete string
	)

	for y := 1; y <= tm.Size[1]; y++ {
		treasureMapDrawPerLine = ""
		if y < tm.Size[1] {
			treasureMapDrawPerLine = "\n"
		}
		for x := 1; x <= tm.Size[0]; x++ {
			treasureMapDrawPerLine = treasureMapDrawPerLine + convertIntToEntity(tm.Mapping[[2]int{x, y}])
		}
		treasureMapDrawComplete = treasureMapDrawPerLine + treasureMapDrawComplete

	}

	// fmt.Println("\033[2J")
	fmt.Println(treasureMapDrawComplete)
	time.Sleep(delay_time * time.Millisecond)
}

// generateMapObstacleDefault putting theentity_obstacle on the outer sandbox
func (tm *TreasureMap) generate() {
	for y := 1; y <= tm.Size[1]; y++ {
		for x := 1; x <= tm.Size[0]; x++ {
			switch true {
			case x == 1, y == 1, x == tm.Size[0], y == tm.Size[1]:
				tm.Mapping[[2]int{x, y}] = entity_obstacle
			default:
				tm.Mapping[[2]int{x, y}] = entity_path
			}
		}
	}
}

// generateMapObstacleCustom putting theentity_obstacle by predefined location
func (tm *TreasureMap) addObstacle(listCustomObstacle [][2]int) {
	for _, customObstacle := range listCustomObstacle {
		tm.Mapping[customObstacle] = entity_obstacle
	}
}

// createMap generate a sandbox, obstacle in it boundaries, and custom obstacle inside
func (tm *TreasureMap) createMap(obstacle [][2]int) {
	tm.generate()
	tm.addObstacle(obstacle)
}

// setEntity put the entity in a position within the map
func (tm *TreasureMap) setEntity(entity int, position [2]int) bool {
	switch tm.Mapping[position] {
	case entity_obstacle:
		return false
	case entity_path:
		if entity == entity_treasure || entity == entity_path || entity == entity_player {
			tm.Mapping[position] = entity
			return true
		}
	case entity_treasure:
		if tm.OriginalMapping[position] == entity_treasure {
			return false
		} else {
			tm.Mapping[position] = entity
			return true
		}
	case entity_player:
		if entity == entity_path {
			tm.Mapping[position] = entity
			return true
		}
	default:
		return false
	}

	return false
}

// setPossibleTreasure hide the actual entity_treasure coordinate to be exposed later
func (tm *TreasureMap) setPossibleTreasure() {
	for coordinate := range tm.Mapping {
		tm.OriginalMapping[coordinate] = tm.Mapping[coordinate]
		if tm.Mapping[coordinate] == entity_path {
			tm.Mapping[coordinate] = entity_treasure
		}
	}
}

// revealMap unhide all unexplored possible treasure
func (tm *TreasureMap) revealMap(treasurePositionXY [2]int) {
	for coordinate := range tm.Mapping {
		if tm.Mapping[coordinate] == entity_treasure && coordinate != treasurePositionXY {
			tm.Mapping[coordinate] = entity_path
		}
	}
}

// convertIntToEntity convert code constant of entity to a map drawn entity
func convertIntToEntity(code int) string {
	switch code {
	case entity_path:
		return "."
	case entity_obstacle:
		return "#"
	case entity_player:
		return "X"
	case entity_treasure:
		return "$"
	default:
		return "."
	}
}
