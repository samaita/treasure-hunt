package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const ( // direction
	up = iota
	down
	right
	left
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
	// render terminal
	unix = iota + 1
	playground
)

const (
	delay_time       = 400  // time used to display a step taken for each movement before continuing
	pause_time       = 4000 // time used to display gimmick before full run exploration
	render_interface = unix // planned to be able to run in golang playground but got TIMEOUT instead!
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
	Range          map[int]int
}

type Treasure struct {
	Position [2]int
}

type TreasureMap struct {
	Size                         [2]int
	OriginalMapping              map[[2]int]int
	Mapping                      map[[2]int]int
	ListPossibleTreasureLocation map[[2]int]bool
	TreasureLocation             [2]int
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

	fmt.Println("Initial Condition, treasure hid in:", treasure.Position, "Wait for it..")
	time.Sleep(pause_time * time.Millisecond)

	treasureMap.setPossibleTreasure()
	treasureMap.render() // display map with possible treasure location

	fmt.Println("Now it's hidden! Let's go find it!")
	time.Sleep(pause_time * time.Millisecond)

	for !player.FoundTreasure {
		// player see unobstructed path, and determine which is treasure and which is path
		treasurePositionXY, listPathPosition := player.see(treasureMap)
		for _, pathPosition := range listPathPosition {
			treasureMap.setEntity(entity_path, pathPosition)
			treasureMap.updatePossibleTreasureLocation(listPathPosition)
		}

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
			treasureMap.clearPossibleTreasureLocation()
			treasureMap.setTreasureLocation(treasurePositionXY)
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
		Range:    make(map[int]int),
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

	// see all entity in x axis with same y axis / right direction ->
	treasurePosition, pathFound = checkMap(treasureMap, startX+1, startY, 1, axis_x)
	if treasureMap.OriginalMapping[treasurePosition] == entity_treasure {
		treasureFound = treasurePosition
	}
	listPathPosition = append(listPathPosition, pathFound...)
	p.Range[right] = len(pathFound)

	// see all entity in -x axis with same y axis / left direction <-
	treasurePosition, pathFound = checkMap(treasureMap, startX-1, startY, -1, axis_x)
	if treasureMap.OriginalMapping[treasurePosition] == entity_treasure {
		treasureFound = treasurePosition
	}
	listPathPosition = append(listPathPosition, pathFound...)
	p.Range[left] = len(pathFound)

	// see all entity in y axis with same x axis / up direction ^
	treasurePosition, pathFound = checkMap(treasureMap, startY+1, startX, 1, axis_y)
	if treasureMap.OriginalMapping[treasurePosition] == entity_treasure {
		treasureFound = treasurePosition
	}
	listPathPosition = append(listPathPosition, pathFound...)
	p.Range[up] = len(pathFound)

	// see all entity in -y axis with same x axis / down direction v
	treasurePosition, pathFound = checkMap(treasureMap, startY-1, startX, -1, axis_y)
	if treasureMap.OriginalMapping[treasurePosition] == entity_treasure {
		treasureFound = treasurePosition
	}
	listPathPosition = append(listPathPosition, pathFound...)
	p.Range[down] = len(pathFound)

	if treasureMap.OriginalMapping[treasureFound] == entity_treasure {
		p.FoundTreasure = true
	}

	// check possibility of path intersection with best probability to get the most explored map
	if p.DirectionTaken == up && p.Range[right] > p.Range[up] {
		p.DirectionTaken = right
	} else if p.DirectionTaken == right && p.Range[down] > p.Range[right] {
		p.DirectionTaken = down
	}

	return treasureFound, listPathPosition
}

// checkMap a shorthand to validate an unobstructed line of sight in original mapping, return treasure location, list of clear path in sight
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
		Size:                         size,
		Mapping:                      make(map[[2]int]int),
		OriginalMapping:              make(map[[2]int]int),
		ListPossibleTreasureLocation: make(map[[2]int]bool),
	}
}

// render display of the mapping, not the original mapping. It also print the info of list possible location of the treasure.
func (tm *TreasureMap) render() {
	var (
		treasureMapDrawPerLine, treasureMapDrawComplete, treasureMapAdditional string
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

	if len(tm.ListPossibleTreasureLocation) > 0 {
		for coordinate, possibleLocation := range tm.ListPossibleTreasureLocation {
			coordinateString := strconv.Itoa(coordinate[0]) + "," + strconv.Itoa(coordinate[1])
			if possibleLocation {
				treasureMapAdditional = treasureMapAdditional + fmt.Sprintf("{%s},", coordinateString)
			}
		}
		treasureMapDrawComplete = treasureMapDrawComplete + fmt.Sprintf("\nPossible treasure location: %s", treasureMapAdditional)
	}

	if tm.TreasureLocation != [2]int{} {
		coordinateString := strconv.Itoa(tm.TreasureLocation[0]) + "," + strconv.Itoa(tm.TreasureLocation[1])
		treasureMapDrawComplete = treasureMapDrawComplete + fmt.Sprintf("\nTreasure found at location: {%s}! Congratulation!", coordinateString)
	}

	renderToTerminal(treasureMapDrawComplete)
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
			tm.ListPossibleTreasureLocation[coordinate] = true
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

// setTreasureLocation mark the found treasure location
func (tm *TreasureMap) setTreasureLocation(treasurePositionXY [2]int) {
	tm.TreasureLocation = treasurePositionXY
}

// updatePossibleTreasureLocation keeping record of all possible treasure location
func (tm *TreasureMap) updatePossibleTreasureLocation(listPathPosition [][2]int) {
	// remove the possible treasure location if its a path
	for _, pathPosition := range listPathPosition {
		tm.ListPossibleTreasureLocation[pathPosition] = false
	}
}

// clearPossibleTreasureLocation empty the list of possible treasure location, usually used once the treasure found
func (tm *TreasureMap) clearPossibleTreasureLocation() {
	tm.ListPossibleTreasureLocation = make(map[[2]int]bool)
}

// renderToTerminal performing animated rendering to terminal.
// If you want to run in golang playground, change constant render_interface from unix to playground.
// But sadly still unable to run in playground
func renderToTerminal(output string) {
	switch render_interface {
	case unix:
		fmt.Println("\033[2J")
		fmt.Println(output)
	case playground:
		fmt.Printf("\x0c %s", output)
	}
	time.Sleep(delay_time * time.Millisecond)
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
