package main

import (
	"fmt"
	"strconv"
	"strings"
)

// distributor divides the work between workers and interacts with other goroutines.
// whole program is the distributor function for now
// does the game logic need to be in a different function?
func distributor(p golParams, d distributorChans, alive chan []cell) {

	// Create the 2D slice to store the world.
	// Two of these, one for source, one for target
	world := make([][]byte, p.imageHeight)
	for i := range world {
		world[i] = make([]byte, p.imageWidth)
	}
	// Create the 2D slice to store the new world.
	newWorld := make([][]byte, p.imageHeight)
	for i := range world {
		newWorld[i] = make([]byte, p.imageWidth)
	}

	// Request the io goroutine to read in the image with the given filename. //these lines chnage to output, send to outpuVal
	d.io.command <- ioInput
	d.io.filename <- strings.Join([]string{strconv.Itoa(p.imageWidth), strconv.Itoa(p.imageHeight)}, "x")

	// The io goroutine sends the requested image byte by byte, in rows.
	for y := 0; y < p.imageHeight; y++ {
		for x := 0; x < p.imageWidth; x++ {
			val := <-d.io.inputVal
			if val != 0 {
				fmt.Println("Alive cell at", x, y)
				world[y][x] = val
			}
		}
	}

	// Calculate the new state of Game of Life after the given number of turns.
	for turns := 0; turns < p.turns; turns++ {
		for y := 0; y < p.imageHeight; y++ {
			for x := 0; x < p.imageWidth; x++ {
				var sum = 0
				var maxWidth = p.imageWidth - 1
				var maxHeight = p.imageHeight - 1

				// if image width or height is 0 or max -> edge cases switch statement better?
				if (x == 0) || (y == 0) || (x == maxWidth) || (y == maxHeight) {

					// edge cases
					// is declaring these necessary
					var yplus = y + 1
					var yminus = y - 1
					var xplus = x + 1
					var xminus = x - 1

					if y == 0 {
						yplus = y + 1
						yminus = maxHeight
					}

					if y == maxHeight {
						yplus = 0
						yminus = y - 1
					}

					if x == 0 {
						xplus = x + 1
						xminus = maxWidth
					}

					if x == maxWidth {
						xplus = 0
						xminus = x - 1
					}

					if world[yminus][xminus] == 0xFF {
						sum = sum + 1
					} // sum++
					if world[yminus][x] == 0xFF {
						sum = sum + 1
					}
					if world[yminus][xplus] == 0xFF {
						sum = sum + 1
					}

					if world[y][xminus] == 0xFF {
						sum = sum + 1
					}
					if world[y][xplus] == 0xFF {
						sum = sum + 1
					}

					if world[yplus][xminus] == 0xFF {
						sum = sum + 1
					}
					if world[yplus][x] == 0xFF {
						sum = sum + 1
					}
					if world[yplus][xplus] == 0xFF {
						sum = sum + 1
					}

					// If your not on the edge?
				} else {

					// Intellij says that the conditions are always true
					// What case are you trying to cover here
					for vertical := -1; vertical < 2; y++ {
						for horizontal := -1; horizontal < 2; x++ {
							if world[y+vertical][x+horizontal] == 0xFF {
								sum = sum + 1
							} else {
								// is this necessary?
								sum = sum
							}
						}
					}
				}

				if sum < 2 && (world[y][x] == 0xFF) {
					world[y][x] = newWorld[y][x] ^ 0xFF
				} else if sum == 2 || sum == 3 && (world[y][x] == 0xFF) {

				} else if sum == 3 && world[y][x] != 0xFF {
					world[y][x] = newWorld[y][x] ^ 0xFF
				} else if sum > 3 && (world[y][x] == 0xFF) {
					world[y][x] = newWorld[y][x] ^ 0xFF
				} else {

				}
				d.io.inputVal <- newWorld[y][x]
			}
		}
	}

	// Create an empty slice to store coordinates of cells that are still alive after p.turns are done.
	var finalAlive []cell
	// Go through the world and append the cells that are still alive.

	for y := 0; y < p.imageHeight; y++ {
		for x := 0; x < p.imageWidth; x++ {
			if world[y][x] != 0 {
				finalAlive = append(finalAlive, cell{x: x, y: y})
			}
		}
	}

	// Make sure that the Io has finished any output before exiting.
	d.io.command <- ioCheckIdle
	<-d.io.idle

	// Return the coordinates of cells that are still alive.
	alive <- finalAlive
}
