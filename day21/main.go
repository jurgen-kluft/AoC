package main

/*
--- Day 21: RPG Simulator 20XX ---

Little Henry Case got a new video game for Christmas. It's an RPG, and he's stuck
on a boss. He needs to know what equipment to buy at the shop. He hands you the
controller.

In this game, the player (you) and the enemy (the boss) take turns attacking.
The player always goes first. Each attack reduces the opponent's hit points by
at least 1. The first character at or below 0 hit points loses.

Damage dealt by an attacker each turn is equal to the attacker's damage score
minus the defender's armor score. An attacker always does at least 1 damage.
So, if the attacker has a damage score of 8, and the defender has an armor score
of 3, the defender loses 5 hit points. If the defender had an armor score of 300,
the defender would still lose 1 hit point.

Your damage score and armor score both start at zero. They can be increased by
buying items in exchange for gold. You start with no items and have as much gold
as you need. Your total damage or armor is equal to the sum of those stats from
all of your items. You have 100 hit points.

Here is what the item shop is selling:

Weapons:    Cost  Damage  Armor
Dagger        8     4       0
Shortsword   10     5       0
Warhammer    25     6       0
Longsword    40     7       0
Greataxe     74     8       0

Armor:      Cost  Damage  Armor
Leather      13     0       1
Chainmail    31     0       2
Splintmail   53     0       3
Bandedmail   75     0       4
Platemail   102     0       5

Rings:      Cost  Damage  Armor
Damage +1    25     1       0
Damage +2    50     2       0
Damage +3   100     3       0
Defense +1   20     0       1
Defense +2   40     0       2
Defense +3   80     0       3

You must buy exactly one weapon; no dual-wielding. Armor is optional, but you
can't use more than one. You can buy 0-2 rings (at most one for each hand).
You must use any items you buy. The shop only has one of each item, so you can't
buy, for example, two rings of Damage +3.

For example, suppose you have 8 hit points, 5 damage, and 5 armor, and that the
boss has 12 hit points, 7 damage, and 2 armor:

The player deals 5-2 = 3 damage; the boss goes down to 9 hit points.
The boss deals 7-5 = 2 damage; the player goes down to 6 hit points.
The player deals 5-2 = 3 damage; the boss goes down to 6 hit points.
The boss deals 7-5 = 2 damage; the player goes down to 4 hit points.
The player deals 5-2 = 3 damage; the boss goes down to 3 hit points.
The boss deals 7-5 = 2 damage; the player goes down to 2 hit points.
The player deals 5-2 = 3 damage; the boss goes down to 0 hit points.
In this scenario, the player wins! (Barely.)

You have 100 hit points. The boss's actual stats are in your puzzle input.
What is the least amount of gold you can spend and still win the fight?
*/

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func iterateOverLinesInTextFile(filename string, action func(string, int)) {
	// Open the file.
	f, _ := os.Open(filename)
	defer f.Close()

	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)

	// Loop over all lines in the file and print them.
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		action(line, lineNumber)
		lineNumber++
	}
}

type stats struct {
	HitPoints int32
	Damage    int32
	Armor     int32
}

func (s *stats) deserialize(str string) {
	parts := strings.Split(str, ":")
	if len(parts) == 2 {
		value, _ := strconv.ParseInt(strings.Trim(parts[1], " "), 10, 32)
		if parts[0] == "Hit Points" {
			s.HitPoints = int32(value)
		} else if parts[0] == "Damage" {
			s.Damage = int32(value)
		} else if parts[0] == "Armor" {
			s.Armor = int32(value)
		}
	}
}

func (s *stats) print() {
	fmt.Printf("Stats ============================================================ \n")
	fmt.Printf("    HitPoints : %v\n", s.HitPoints)
	fmt.Printf("    Damage : %v\n", s.Damage)
	fmt.Printf("    Armor : %v\n", s.Armor)
}

type weapon struct {
	Name   string
	Cost   int32
	Damage int32
	Armor  int32
}

type armor struct {
	Name   string
	Cost   int32
	Damage int32
	Armor  int32
}

type ring struct {
	Name   string
	Cost   int32
	Damage int32
	Armor  int32
}

type shop struct {
	weapons []weapon
	armors  []armor
	rings   []ring
}

func (s *shop) deserialize(str string) {
	parts := strings.Split(str, ",")
	if len(parts) == 5 {
		index := 0
		ctgr := strings.Trim(parts[index], " ")
		index++
		name := strings.Trim(parts[index], " ")
		index++
		cost, _ := strconv.ParseInt(strings.Trim(parts[index], " "), 10, 32)
		index++
		damg, _ := strconv.ParseInt(strings.Trim(parts[index], " "), 10, 32)
		index++
		armr, _ := strconv.ParseInt(strings.Trim(parts[index], " "), 10, 32)
		index++

		ctgr = strings.ToLower(ctgr)
		name = strings.ToLower(name)

		switch ctgr {
		case "weapon":
			w := weapon{Name: name, Cost: int32(cost), Damage: int32(damg), Armor: int32(armr)}
			s.weapons = append(s.weapons, w)
			break
		case "armor":
			a := armor{Name: name, Cost: int32(cost), Damage: int32(damg), Armor: int32(armr)}
			s.armors = append(s.armors, a)
			break
		case "ring":
			r := ring{Name: name, Cost: int32(cost), Damage: int32(damg), Armor: int32(armr)}
			s.rings = append(s.rings, r)
			break
		}

	}
}

func (s *shop) print() {
	fmt.Printf("Shop ============================================================ \n")

	fmt.Printf("Weapons : %v\n", len(s.weapons))
	for _, w := range s.weapons {
		fmt.Printf("    Name : %v\n", w.Name)
		fmt.Printf("      Cost : %v\n", w.Cost)
		fmt.Printf("      Damage : %v\n", w.Damage)
		fmt.Printf("      Armor : %v\n", w.Armor)
	}

	fmt.Printf("Armors : %v\n", len(s.armors))
	for _, a := range s.armors {
		fmt.Printf("    Name : %v\n", a.Name)
		fmt.Printf("      Cost : %v\n", a.Cost)
		fmt.Printf("      Damage : %v\n", a.Damage)
		fmt.Printf("      Armor : %v\n", a.Armor)
	}

	fmt.Printf("Rings : %v\n", len(s.rings))
	for _, r := range s.rings {
		fmt.Printf("    Name : %v\n", r.Name)
		fmt.Printf("      Cost : %v\n", r.Cost)
		fmt.Printf("      Damage : %v\n", r.Damage)
		fmt.Printf("      Armor : %v\n", r.Armor)
	}
}

func readInputFromFile(filename string) (shopx *shop, statx *stats) {
	statx = &stats{}
	shopx = &shop{weapons: make([]weapon, 0), armors: make([]armor, 0), rings: make([]ring, 0)}

	object := ""

	computator := func(text string, line int) {
		if text == "#Boss" {
			object = "boss"
		} else if text == "#Shop" {
			object = "shop"
		} else {
			if object == "boss" {
				statx.deserialize(text)
			} else if object == "shop" {
				shopx.deserialize(text)
			}
		}

	}
	iterateOverLinesInTextFile(filename, computator)

	return shopx, statx
}

func findOptimumThingsToBuyToBeatBoss(player *stats, boss *stats, shop *shop) {

	// Rules:
	// - You must and can only buy one weapon
	// - You optionally can be one armor
	// - You optionally can buy maximum 2 rings
	// - The shop only has one of each item

	// Run all the combinations:
}

func main() {
	player := &stats{HitPoints: 100, Damage: 1, Armor: 0}
	shop, boss := readInputFromFile("input.text")
	shop.print()
	boss.print()
	findOptimumThingsToBuyToBeatBoss(player, boss, shop)
}
