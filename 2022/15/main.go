// https://adventofcode.com/2022/day/15

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

var inputs = [...]string{
	"Sensor at x=2, y=18: closest beacon is at x=-2, y=15",
	"Sensor at x=9, y=16: closest beacon is at x=10, y=16",
	"Sensor at x=13, y=2: closest beacon is at x=15, y=3",
	"Sensor at x=12, y=14: closest beacon is at x=10, y=16",
	"Sensor at x=10, y=20: closest beacon is at x=10, y=16",
	"Sensor at x=14, y=17: closest beacon is at x=10, y=16",
	"Sensor at x=8, y=7: closest beacon is at x=2, y=10",
	"Sensor at x=2, y=0: closest beacon is at x=2, y=10",
	"Sensor at x=0, y=11: closest beacon is at x=2, y=10",
	"Sensor at x=20, y=14: closest beacon is at x=25, y=17",
	"Sensor at x=17, y=20: closest beacon is at x=21, y=22",
	"Sensor at x=16, y=7: closest beacon is at x=15, y=3",
	"Sensor at x=14, y=3: closest beacon is at x=15, y=3",
	"Sensor at x=20, y=1: closest beacon is at x=15, y=3",
}

type Sensor struct {
	utils.Location2D[int]
	distToClosestBeacon int
}

type Beacon struct {
	utils.Location2D[int]
}

type System struct {
	sensors []Sensor
	beacons []Beacon
}

// tuningFrequency returns the tuning frequency of a beacon at the given location
func tuningFrequency(loc utils.Location2D[int], factor int) int {
	return loc.X*factor + loc.Y
}

// sensorRangesByRow returns the ranges of X coordinates in range of sensor among Y axis
func (system System) sensorRangesByRow(y int) []utils.Interval[int] {
	intervals := utils.Intervals[int]{}
	for _, sensor := range system.sensors {
		// compare the sensor's distance to its closest beacon to the distance with the sensor Y coordinate
		delta := sensor.distToClosestBeacon - utils.Abs(sensor.Y-y)
		// if we are in range (i.e delta >= 0), consider the X coordinates in this same range also in range
		if delta >= 0 {
			intervals.Insert(utils.NewInterval(sensor.X-delta, sensor.X+delta), 1)
		}
	}

	return intervals.Values()
}

// impossibleBeaconLocationsOnRow returns the number of X locations among Y axis where a beacon can not be
func (system System) impossibleBeaconLocationsOnRow(y int) (impossibleLocations int) {
	// get all X locations in Y sensors range
	for _, impossibleRange := range system.sensorRangesByRow(y) {
		impossibleLocations += impossibleRange.Len()
	}

	// subtract the number of beacons in Y range
	for _, beacon := range system.beacons {
		if beacon.Y == y {
			impossibleLocations--
		}
	}

	return
}

// getPossibleBeaconLocation returns the first possible location for a beacon to be, in the search area
func (system System) getPossibleBeaconLocation(searchArea utils.Interval[int]) utils.Location2D[int] {
	// iterate over every rows
	for y := searchArea.Min; y <= searchArea.Max; y++ {
		// get impossible X location range on this row
		impossibleLocations := system.sensorRangesByRow(y)

		// if there is more than 1 range, it means that the only possible location is between them
		if len(impossibleLocations) > 1 {
			return utils.NewLocation2D(utils.Min(impossibleLocations[0].Max, impossibleLocations[1].Max)+1, y)
		}
	}

	return utils.NewLocation2D(-1, -1)
}

// parseSystem parses input and returns the list of sensors and beacons in the system
func parseSystem(inputs []string) System {
	system := System{
		sensors: make([]Sensor, 0, len(inputs)),
		beacons: make([]Beacon, 0, len(inputs)),
	}

	beaconsMap := make(map[Beacon]bool, len(inputs))
	for _, input := range inputs {
		sensor := Sensor{}
		beacon := Beacon{}
		fmt.Sscanf(input, "Sensor at x=%v, y=%v: closest beacon is at x=%v, y=%v", &sensor.X, &sensor.Y, &beacon.X, &beacon.Y)

		// ensure beacon existence and unity
		if !beaconsMap[beacon] {
			// add beacon to the system
			beaconsMap[beacon] = true
			system.beacons = append(system.beacons, beacon)
		}

		// add sensor to the system
		sensor.distToClosestBeacon = sensor.ManhattanDist(beacon.Location2D)
		system.sensors = append(system.sensors, sensor)
	}

	return system
}

func main() {
	defer func(start time.Time) { fmt.Println("took:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init

	system := parseSystem(inputs[:])

	////////////////////////////////////////

	// 26
	fmt.Println("part1:", system.impossibleBeaconLocationsOnRow(10))

	////////////////////////////////////////

	// 56000011
	fmt.Println("part2:", tuningFrequency(system.getPossibleBeaconLocation(utils.NewInterval(0, 4000000)), 4000000))
}
