// https://adventofcode.com/2023/day/24

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Objects []Object

type Object struct {
	x, y, z    int
	vx, vy, vz int
}

func (objects Objects) getCollidingRock() Object {
	rock := Object{}

	// considering a rock R(x,y,z,vx,vy,vz) and a hail stone H(x,y,z,vx,vy,vz)
	// (1) R(t) = H(t)
	// (1) Rpi + t*Rvi = Hpi + t*Hvi
	// (1) t = (Hx-Rx)/(Rvx-Hvx) = (Hy-Ry)/(Rvy-Hvy) = (Hz-Rz)/(Rvz-Hvz)
	//
	// (1) for X and Y only:
	// (2) Ry*Rvx - Rx*Rvy = Hx*Hvy - Hy*Hvx + Ry*Hvx + Hy*Rvx - Hx*Rvy - Rx*Hvy
	//
	// (2) left part is constant for all H:
	// (3) Hix*Hivy - Hiy*Hivx + Ry*Hivx + Hiy*Rvx - Hix*Rvy - Rx*Hivy = Hjx*Hjvy - Hjy*Hjvx + Ry*Hjvx + Hjy*Rvx - Hjx*Rvy - Rx*Hjvy
	// (3) (Hjvy-Hivy)*Rx + (Hivx-Hjvx)*Ry + (Hiy-Hjy)*Rvx + (Hjx-Hix)*Rvy = Hjx*Hjvy - Hjy*Hjvx - Hix*Hivy + Hiy*Hivx
	n := 4
	Axy := make([][]int, n)
	Bxy := make([]int, n)
	for i := 0; i < n && i < len(objects)-1; i++ {
		Hi := objects[i]
		Hj := objects[i+1]
		Axy[i] = []int{Hj.vy - Hi.vy, Hi.vx - Hj.vx, Hi.y - Hj.y, Hj.x - Hi.x}
		Bxy[i] = Hj.x*Hj.vy - Hj.y*Hj.vx - Hi.x*Hi.vy + Hi.y*Hi.vx
	}
	Rxy := utils.GaussianEliminationSolver(Axy, Bxy)
	rock.x = Rxy[0]
	rock.y = Rxy[1]
	rock.vx = Rxy[2]
	rock.vy = Rxy[3]

	// using (3) with known X or Y gives:
	// (4) (Hjvz-Hivz)*Rx + (Hivx-Hjvx)*Rz + (Hiz-Hjz)*Rvx + (Hjx-Hix)*Rvz = Hjx*Hjvz - Hjz*Hjvx - Hix*Hivz + Hiz*Hivx
	// (4) (Hivx-Hjvx)*Rz + (Hjx-Hix)*Rvz = Hjx*Hjvz - Hjz*Hjvx - Hix*Hivz + Hiz*Hivx -(Hjvz-Hivz)*Rx - (Hiz-Hjz)*Rvx
	n = 2
	Az := make([][]int, n)
	Bz := make([]int, n)
	for i := 0; i < n && i < len(objects)-1; i++ {
		Hi := objects[i]
		Hj := objects[i+1]
		Az[i] = []int{Hi.vx - Hj.vx, Hj.x - Hi.x}
		Bz[i] = Hj.x*Hj.vz - Hj.z*Hj.vx - Hi.x*Hi.vz + Hi.z*Hi.vx - (Hj.vz-Hi.vz)*rock.x - (Hi.z-Hj.z)*rock.vx
	}
	Rz := utils.GaussianEliminationSolver(Az, Bz)
	rock.z = Rz[0]
	rock.vz = Rz[1]

	return rock
}

func (objects Objects) getIntersectionsCount(min, max int) int {
	count := 0
	for i := 0; i < len(objects); i++ {
		lhs := objects[i]
		for j := i + 1; j < len(objects); j++ {
			rhs := objects[j]
			det := rhs.vx*lhs.vy - rhs.vy*lhs.vx
			// no unique intersection
			if det == 0 {
				continue
			}

			dx := rhs.x - lhs.x
			dy := rhs.y - lhs.y
			u := dy*rhs.vx - dx*rhs.vy
			v := dy*lhs.vx - dx*lhs.vy
			// intersects in the past
			if (u > 0) != (det > 0) || (v > 0) != (det > 0) {
				continue
			}

			// intersects outside area
			f := float64(u) / float64(det)
			if (float64(lhs.vx)*f < float64(min-lhs.x) || float64(lhs.vx)*f > float64(max-lhs.x)) ||
				float64(lhs.vy)*f < float64(min-lhs.y) || float64(lhs.vy)*f > float64(max-lhs.y) {
				continue
			}

			count++
		}
	}
	return count
}

func parseHailStones(inputs []string) Objects {
	hailStones := make(Objects, len(inputs))
	for i, input := range inputs {
		fmt.Sscanf(input, "%v, %v, %v @ %v, %v, %v",
			&hailStones[i].x,
			&hailStones[i].y,
			&hailStones[i].z,
			&hailStones[i].vx,
			&hailStones[i].vy,
			&hailStones[i].vz,
		)
	}
	return hailStones
}

func main() {
	fmt.Println("--- 2023 Day 24: Never Tell Me The Odds ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	hailStones := parseHailStones(inputs)

	////////////////////////////////////////

	// 2
	fmt.Println("Part 1:", hailStones.getIntersectionsCount(7, 27))

	////////////////////////////////////////

	rock := hailStones.getCollidingRock()

	// 47
	fmt.Println("Part 2:", rock.x+rock.y+rock.z)
}
