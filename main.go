package main

import (
	"fmt"
	"os"
)

type Vec2 struct {
	x, y int
}


type Checkpoint = Vec2
type Pod struct {
	pos Vec2
	vel Vec2
	acc Vec2
	lastPos Vec2
	lastVel Vec2
}
func updateMyPod(x, y int) {
	myPod.lastPos = myPod.pos
	myPod.lastVel = myPod.vel

	myPod.pos = Vec2{x, y}
	myPod.vel = Vec2{
		x: x - myPod.lastPos.x,
		y: y - myPod.lastPos.y,
	}
	myPod.acc = Vec2{
		x: myPod.vel.x - myPod.lastVel.x,
		y: myPod.vel.y - myPod.lastVel.y,
	}
}

var myPod = Pod{

}

var (
	boostConsumed bool

	Checkpoints []Checkpoint 
	target Checkpoint
	allCheckpointsSeen bool
	
)
///////////////////////////////////////////////
///////////////////////////////////////////////
func main() {
    for {

		var thurst int
		// var boost bool

		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle int
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)
		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)

		updateMyPod(x, y)

		updateAllCheckpointseen(Checkpoint{nextCheckpointX, nextCheckpointY})

		if !allCheckpointsSeen {
			addCheckpoint(nextCheckpointX, nextCheckpointY)
		}

		// if !boostConsumed && nextCheckpointDist > 4000 && nextCheckpointAngle < 25 {
		// 	boost = true
		// 	boostConsumed = true
		// }


		// thurst = 100 - norm(nextCheckpointAngle) * 3

		// if thurst < 0 {
		// 	thurst = 50
		// }

		fmt.Fprintf(os.Stderr, "Checkpoints: %v\n", Checkpoints)
		fmt.Fprintf(os.Stderr, "allseen: %t\n", allCheckpointsSeen)
		fmt.Fprintf(os.Stderr, "myPod: %v", myPod)
		// fmt.Fprintf(os.Stderr, "next: %d ; t: %d\n", nextCheckpointAngle, thurst)
		// fmt.Fprintf(os.Stderr, "bc: %t | b: %t | t: %d | next: %d", boostConsumed, boost, thurst, nextCheckpointDist)


		if nextCheckpointDist < 1500 && allCheckpointsSeen {
			target = getNextTarget(nextCheckpointX, nextCheckpointY)
		} else {
			target = Checkpoint{nextCheckpointX, nextCheckpointY}
		}
		thurst = 100

		
		if nextCheckpointDist < 1000 {
			thurst = 0
		}


		if nextCheckpointAngle > 5 {

		}

		if !boostConsumed {
			fmt.Printf("%d %d BOOST\n", target.x, target.y)
			boostConsumed = true
		} else {
			fmt.Printf("%d %d %d\n", target.x, target.y, thurst)
		}
    }
}
///////////////////////////////////////////////
///////////////////////////////////////////////
func norm(a int) int {
	n := abs(a)
	f := float64(n) / 180
	s := f * 99
	return int(s)
}

func abs (x int) int {
	if x < 0 {
		return -x
	}
	return x
}
///////////////////////////////////////////////
func has(s []Checkpoint, x Checkpoint) (int, bool) {
	for i, e := range s {
		if e == x {
			return i, true
		}
	}
	return -1, false
}
func addCheckpoint(x, y int) {
	c := Checkpoint{x, y}

	if _, b := has(Checkpoints, Checkpoint{x, y}); !b {
		Checkpoints = append(Checkpoints, c)
	}
}

func updateAllCheckpointseen (cCheckpoint Checkpoint) {
	if len(Checkpoints) > 1 && cCheckpoint == Checkpoints[0] {
		allCheckpointsSeen = true
	}
}

func getNextTarget(x, y int) Checkpoint {
	ci, _ := has(Checkpoints, Checkpoint{x, y})
	if ci == len(Checkpoints) - 1 {
		return Checkpoints[0]
	}
	return Checkpoints[ci + 1]
}
