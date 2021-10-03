# Treasure Hunt (in Terminal, Unix Only)
You got a map from a mysterious guy. Your pirate-sense now tingling, thrill for a hunt!

The map is look like this.

```
y
|
########
#......#
#.###..#
#...#.##
#X#....#
######## -- x
```
Wait, some information was left, how kind of him!
```
# represents an obstacle.
. represents a clear path
X represents your starting position.
```
In this hunt, you'll act as a player. Some rules applied:
- A treasure is hidden within one of the clear entity_path points, and we must find it.

- Player have long range of vision, it able to see left/right/up/down whenever the path is unobstructed.

- Player restricted to only move in 3 direction, up then right then down.

- The location of the entity_treasure is located randomly every runtime. Uh-oh, hopefully you'll still found it.

- All possible location of the treasure marked as $, let's eliminate the possibility during the exploring.

# For Dev / Tester
Minimum Requirement(s) to run:
- Unix
- Golang v1.14+ installed

To run the program clone/copy the code and type this command:
```
go run main.go
```
## Sample Output
### Initial Condition
```
########
#......#
#.###..#
#...#.##
#X#...$#
########
Initial Condition, treasure hid in: [7 2] Wait for it..
```
It will pause 4s at default, before continuing.
### Exploration
```
########
#.$$$$$#
#X###$$#
#...#$##
#.#$$$$#
########
Possible treasure location: {6,5},{4,5},{7,5},{5,2},{4,2},{3,5},{5,5},{6,2},{7,4},{6,3},{6,4}
```
### Final Output
```
########
#......#
#.###..#
#...#.##
#.#..X$#
########
Treasure found at location: {7,2}! Congratulation!
```

Have fun!

# FAQ
## I want to change the treasure location, how to do it?
Even if it generate randomly, you can still set the location manually. Go to this method
```
func (t *Treasure) randomizePosition(sizeX, sizeY int) {
...
}
```
Replace entire method to 
```
func (t *Treasure) randomizePosition(sizeX, sizeY int) {
	t.Position = [2]int{x,y}
    return
}
```
Make sure the x & y wont be higher than the mapSize, and make sure you didn't put it into clear path. I'd suggest to use `[2]int{7,2}` for full experience!

## The map updated too fast, how to slow it down?
Update `delay_time` from 400 to your liking, The time unit used is millisecond.
