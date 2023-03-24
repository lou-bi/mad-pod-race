package main

import (
	"fmt"
	"math"
)

type Vec2 struct{ x, y int }

var cpm = CheckpointManager{}
var a = 1
var boostConsumed bool

var firstShot bool

func main() {
	for {

		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle int
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)
		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)

		ccp := Vec2{nextCheckpointX, nextCheckpointY}

		cpm.add(ccp)
		cpm.updateCPId(ccp)
		cpm.lap(ccp)
		cpm.updateAllCPSeen(ccp)

		absNCPA := math.Abs(float64(nextCheckpointAngle))

		thrust := 100

		// if !cpm.allCPSeen {

		if nextCheckpointAngle < a && nextCheckpointAngle > -a {
			thrust = 100
		} else {
			a := (1.0 - (absNCPA / 90.0))
			b := (float64(nextCheckpointDist) / (2.0 * 600.0))

			if a < 0 {
				a = 0.0
			}

			thrust = int(100.0 * a * b)
		}

		if thrust > 100 {
			thrust = 100
		}

		var target Vec2

		if cpm.allCPSeen {
			target = cpm.bcps[cpm.currentCPId]
		} else {
			target = cpm.computeInitialBestTraj(ccp)
		}

		if int(absNCPA) < a && nextCheckpointDist > 6000 {
			fmt.Println(target.x, target.y, "BOOST")
			boostConsumed = true
		} else {
			fmt.Println(target.x, target.y, thrust)
		}
	}
}

type CheckpointManager struct {
	cps         []Vec2
	bcps        []Vec2
	currentLap  int
	currentCPId int
	bestBoost   Vec2
	allCPSeen   bool
}

func (cpm *CheckpointManager) has(cp Vec2) bool {
	for _, c := range cpm.cps {
		if c == cp {
			return true
		}
	}
	return false
}

func (cpm *CheckpointManager) add(cp Vec2) {
	if !cpm.has(cp) {
		cpm.cps = append(cpm.cps, cp)
	}
}

func (cpm *CheckpointManager) updateCPId(cp Vec2) {
	for i, c := range cpm.cps {
		if c == cp {
			cpm.currentCPId = i
		}
	}
}

func (cpm *CheckpointManager) lap(cp Vec2) {
	if cpm.currentCPId != 0 && cp == cpm.cps[0] {
		cpm.currentLap++
		cpm.currentCPId = 0
	}
}

func (cpm *CheckpointManager) updateAllCPSeen(cp Vec2) {
	if cpm.allCPSeen == false && len(cpm.cps) > 1 && cpm.cps[0] == cp {
		cpm.allCPSeen = true
		cpm.computeBestBoost()
		cpm.computeBestTraj()
	}
}

func (cpm *CheckpointManager) computeBestBoost() {
	longest := 0
	var best Vec2
	for i := 0; i < len(cpm.cps)-1; i++ {
		a := cpm.cps[i]
		b := cpm.cps[i+1]
		dx := (a.x - b.x) * (a.x - b.x)
		dy := (a.y - b.y) * (a.y - b.y)
		d := dx + dy

		if d > longest {
			longest = d
			best = cpm.cps[i+1]
		}
	}
	cpm.bestBoost = best
}
func (cpm *CheckpointManager) computeBestTraj() {
	checkpointRadius := 400.0
	for i, cp := range cpm.cps {
		var previous int
		var next int

		if i == 0 {
			previous = len(cpm.cps) - 1
			next = i + 1
		} else if i == len(cpm.cps)-1 {
			previous = i - 1
			next = 0
		} else {
			previous = i - 1
			next = i + 1
		}
		a := cpm.cps[previous]
		b := cpm.cps[next]

		ab := Vec2{x: b.x - a.x, y: b.y - a.y}
		ap := Vec2{x: cp.x - a.x, y: cp.y - a.y}
		orientation := ab.x*ap.y - ab.y*ap.x

		var midx float64
		var midy float64

		if orientation < 0 {
			midx = float64(a.y-b.y) / 2
			midy = float64(b.x-a.x) / 2
		} else {
			midx = float64(b.y-a.y) / 2
			midy = float64(a.x-b.x) / 2
		}

		norm := math.Sqrt(math.Pow(midx, 2) + math.Pow(midy, 2))
		normalized := []float64{midx / norm, midy / norm}

		scale := []float64{normalized[0] * checkpointRadius, normalized[1] * checkpointRadius}

		bestPoint := Vec2{
			x: int(float64(cp.x) + scale[0]),
			y: int(float64(cp.y) + scale[1]),
		}

		cpm.bcps = append(cpm.bcps, bestPoint)
	}
}
func (cpm *CheckpointManager) computeInitialBestTraj(cp Vec2) Vec2 {
	checkpointRadius := 600.0

	midx := 8000.0
	midy := 4500.0

	dx := float64(cp.x) - midx
	dy := float64(cp.y) - midy

	magnitude := math.Sqrt(dx*dx + dy*dy)
	if magnitude > 0 {
		dx /= magnitude
		dy /= magnitude
	}

	return Vec2{
		x: int(midx + dx*(magnitude-checkpointRadius)),
		y: int(midy + dy*(magnitude-checkpointRadius)),
	}
}
