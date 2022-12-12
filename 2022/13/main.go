// https://adventofcode.com/2022/day/13

package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	smaller = iota - 1
	equal
	greater
)

type Data interface {
	fmt.Stringer

	unmarshal(interface{}) Data
	compare(Data) int
}

type Integer int

// String returns string representation of an Integer data
func (integer Integer) String() string {
	return fmt.Sprint(int(integer))
}

// toPacket converts an Integer data to a Packet data
func (integer *Integer) toPacket() *Packet {
	return &Packet{integer}
}

// unmarshal unmarshals an interface (assuming float64) into an Integer data
func (integer *Integer) unmarshal(data interface{}) Data {
	*integer = Integer(data.(float64))
	return integer
}

// compare compares this Integer data with an other, and returns whether there are equals, greater or smaller
func (left Integer) compare(right Data) int {
	switch right := right.(type) {
	case *Integer:
		// left side is smaller, so inputs are in the right order
		if left < *right {
			return smaller
		}

		// right side is smaller, so inputs are not in the right order
		if left > *right {
			return greater
		}

	case *Packet:
		// mixed types: convert left to Packet and retry comparison
		return left.toPacket().compare(right)
	}

	// continue checking
	return equal
}

type Packet []Data

// String returns string representation of a Packet data
func (packet Packet) String() (s string) {
	values := make([]string, len(packet))
	for i, data := range packet {
		values[i] = data.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

// UnmarshalJSON parses the JSON-encoded raw data and stores the result in the Packet data
func (packet *Packet) UnmarshalJSON(raw []byte) error {
	var data []interface{}
	json.Unmarshal(raw, &data)
	packet.unmarshal(data)
	return nil
}

// unmarshal unmarshals an interface (assuming []any) into a Packet data
// a Packet is compose of multiple data, so every data contained is marshalled and added into it
func (packet *Packet) unmarshal(data interface{}) Data {
	for _, data := range data.([]interface{}) {
		switch data.(type) {
		// marshal into a new Integer data
		case float64:
			*packet = append(*packet, new(Integer).unmarshal(data))
		// marshal into a new Packet data
		case []interface{}:
			*packet = append(*packet, new(Packet).unmarshal(data))
		}
	}
	return packet
}

// compare compares this Packet data with an other, and returns whether there are equals, greater or smaller
func (left Packet) compare(right Data) int {
	switch right := right.(type) {
	case *Packet:
		n, m := Integer(len(left)), Integer(len(*right))

		// compare each value of each list
		for i := Integer(0); i < n && i < m; i++ {
			if comp := left[i].compare((*right)[i]); comp != equal {
				return comp
			}
		}

		// continue checking by comparing list lengths
		// if the lists are the same length and no comparison makes a decision about the order, continue checking the next part of the input
		return n.compare(&m)

	case *Integer:
		// mixed types: convert right to Packet and retry comparison
		return left.compare(right.toPacket())
	}

	// continue checking
	return equal
}

// getDecoderKey adds divider packets, sort the data and afterward returns the decoder key by locating the divider packets
func (packet Packet) getDecoderKey(dividers ...Integer) (key int) {
	// add divider packets
	for _, divider := range dividers {
		divider := divider
		packet = append(packet, &Packet{&Packet{&divider}})
	}

	// sort packets
	slices.SortFunc(packet, Data.compare)

	// retrieve dividers
	key = 1
	for i, data := range packet {
		for j, divider := range dividers {
			if data.compare(&divider) == equal {
				key *= (i + 1)
				// in case we locate the last divider, return the key
				if len(dividers) <= 1 {
					return key
				}
				// after finding a divider, remove it from the list
				dividers = append(dividers[:j], dividers[j+1:]...)
				break
			}
		}
	}

	return -1
}

type Pairs [][2]Packet

// flat returns a one level list of data composed by all Packets paired
func (pairs Pairs) flat() []Data {
	packets := make([]Data, 0, len(pairs)*2)
	for _, pair := range pairs {
		pair := pair // copy
		packets = append(packets, &pair[0], &pair[1])
	}
	return packets
}

// rightOrderedPairs returns the number of right ordered pairs (i.e. first element is smaller than second)
func (pairs Pairs) rightOrderedPairs() (rightOrdered int) {
	for i, pair := range pairs {
		if pair[0].compare(&pair[1]) == smaller {
			rightOrdered += i + 1
		}
	}
	return
}

// parsePairs parses input and returns a list of Packet pair
func parsePairs(inputs []string) (pairs Pairs) {
	for i := 0; i < len(inputs)-1; i += 2 {
		var pair [2]Packet
		json.Unmarshal([]byte(inputs[i]), &pair[0])
		json.Unmarshal([]byte(inputs[i+1]), &pair[1])
		pairs = append(pairs, pair)

		// skip empty lines
		if i < len(inputs)-2 && inputs[i+2] == "" {
			i++
		}
	}
	return
}

func main() {
	fmt.Println("--- 2022 Day 13: Distress Signal ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	pairs := parsePairs(inputs)

	////////////////////////////////////////

	// 13
	fmt.Println("Part 1:", pairs.rightOrderedPairs())

	////////////////////////////////////////

	// 140
	fmt.Println("Part 2:", Packet(pairs.flat()).getDecoderKey(2, 6))
}
