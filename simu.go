package main

import "math"

type Point struct {
	x, y float64
}

type Unit struct {
	Point
	id int
	r  float64
	vx float64
	vy float64
}

type Pod struct {
	Unit
	angle            float64
	nextCheckpointId int
	checked          int
	timeout          int
	partner          *Pod
	shield           bool
}

func (p *Point) distance(p2 Point) float64 {
	return (p.x-p2.x)*(p.x-p2.x) + (p.y-p2.y)*(p.y-p2.y)
}

func (pod *Pod) getAngle(p Point) float64 {
	d := pod.distance(p)
	dx := (p.x - pod.x) / d
	dy := (p.y - pod.y) / d

	a := math.Acos(dx) * 180 / math.Pi

	if dy < 0 {
		a = 360.0 - a
	}
	return a
}

func (pod *Pod) diffAngle(p Point) float64 {
	a := pod.getAngle(p)

	var right, left float64

	if pod.angle <= a {
		right = a - pod.angle
	} else {
		right = 360.0 - pod.angle + a
	}
	if pod.angle >= a {
		left = pod.angle - a
	} else {
		left = pod.angle + 360.0 - a
	}

	if right < left {
		return right
	}
	return left
}

func (pod *Pod) rotate(p Point) {
	a := pod.diffAngle(p)

	// On ne peut pas tourner de plus de 18° en un seul tour
	if a > 18.0 {
		a = 18.0
	} else if a < -18.0 {
		a = -18.0
	}

	pod.angle += a

	// L'opérateur % est lent. Si on peut l'éviter, c'est mieux.
	if pod.angle >= 360.0 {
		pod.angle = pod.angle - 360.0
	} else if pod.angle < 0.0 {
		pod.angle += 360.0
	}
}

func (pod *Pod) boost(thrust int) {
	// N'oubliez pas qu'un pod qui a activé un shield ne peut pas accélérer pendant 3 tours
	if pod.shield {
		return
	}

	// Conversion de l'angle en radian
	ra := pod.angle * math.Pi / 180.0

	// Trigonométrie
	pod.vx += math.Cos(ra) * float64(thrust)
	pod.vy += math.Sin(ra) * float64(thrust)
}

func (pod *Pod) move(t float64) {
	pod.x += pod.vx * t
	pod.y += pod.vy * t
}

func (pod *Pod) end() {
	pod.x = math.Round(pod.x)
	pod.y = math.Round(pod.y)
	pod.vx = math.Trunc(pod.vx * 0.85)
	pod.vy = math.Trunc(pod.vy * 0.85)

	// N'oubliez pas que le timeout descend de 1 chaque tour. Il revient à 100 quand on passe par un checkpoint
	pod.timeout -= 1
}

func (pod *Pod) play(p Point, thrust int) {
	pod.rotate(p)
	pod.boost(thrust)
	pod.move(1.0)
	pod.end()
}
