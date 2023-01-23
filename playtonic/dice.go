package playtonic

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bismuthsalamander/bafflebawx/inceptor"
)

//TODO: collapse SideCount and RollResult into one type?

type DiceCount uint8
type RollResult int32
type SideCount uint16
type RollDescriptor string
type Modifier int16
type ResultRule uint8

const RR_SUM = 0
const RR_HIGHEST = 1

const SUFFIXES = "^*"

type RollType struct {
	numDice       DiceCount
	numSides      SideCount
	modifier      Modifier
	resultRule    ResultRule
	explodingDice bool
}

type RollOutcome struct {
	result  RollResult
	rawDice []RollResult
}

func (rt RollType) isMaxDie(r RollResult) bool {
	return r == RollResult(rt.numSides)
}

func dieRoll(sides SideCount) (RollResult, error) {
	n, err := inceptor.Uint64()
	if err != nil {
		return 0, err
	}
	return RollResult((n % uint64(sides)) + 1), nil
}

func ParseDescriptor(desc string) (RollType, error) {
	rt := RollType{1, 0, 0, RR_SUM, false}
	d := desc
	if len(d) < 1 {
		return rt, fmt.Errorf("descriptor cannot be empty")
	}
	for strings.Contains(SUFFIXES, d[len(d)-1:]) {
		if d[len(d)-1] == '^' {
			rt.resultRule = RR_HIGHEST
		} else if d[len(d)-1] == '*' {
			rt.explodingDice = true
		}
		d = d[:len(d)-1]
	}
	parts := strings.Split(d, "d")
	if len(parts) != 2 {
		return rt, fmt.Errorf("descriptor %s contains %d instances of char 'd'; expected 1", desc, len(parts))
	}
	if len(parts[0]) != 0 {
		dicecount, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return rt, fmt.Errorf("could not parse dice count '%s' from descriptor %s", parts[0], desc)
		}
		if dicecount > math.MaxUint8 || dicecount <= 0 {
			return rt, fmt.Errorf("dice count %d (max: %d) in descriptor %s", dicecount, math.MaxUint8, desc)
		}
		rt.numDice = DiceCount(dicecount)
	}
	mod_index := strings.IndexAny(parts[1], "+-")
	if mod_index != -1 {
		mod_str := parts[1][mod_index:]
		parts[1] = parts[1][:mod_index]
		if mod_str != "" {
			mod_magnitude, err2 := strconv.ParseUint(mod_str[1:], 10, 64)
			if err2 != nil {
				return RollType{}, fmt.Errorf("could not parse modifier magnitude '%s' (from )descriptor %s)", mod_str[1:], desc)
			}
			rt.modifier = Modifier(mod_magnitude)
			if mod_str[0] == '-' {
				rt.modifier *= -1
			}
			if rt.modifier > math.MaxInt16 || rt.modifier < math.MinInt16 {
				return RollType{}, fmt.Errorf("modifier %d outside of range [%d, %d] (descriptor %s)", rt.modifier, math.MaxInt16, math.MinInt16, desc)
			}
		}
	}
	sides, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return RollType{}, fmt.Errorf("error parsing side count '%s' (descriptor %s)", parts[1], desc)
	}
	if sides > math.MaxUint16 || sides <= 0 {
		return RollType{}, fmt.Errorf("side count %d outside of range [1, %d] (descriptor %s)", sides, math.MaxUint16, desc)
	}
	rt.numSides = SideCount(sides)
	return rt, nil
}

func ExecuteRoll(rt RollType) (RollOutcome, error) {
	var res RollResult = 0
	//theoretically, the loop may be faster for non-exploding dice if we made
	//rawDice with the correct number of slots, but we'd have to either write
	//two separate versions of the loop or add some ugly logic to check the
	//slice's size on each iteration.
	rawDice := make([]RollResult, 0)
	for i := 0; DiceCount(i) < rt.numDice; i++ {
		d, err := dieRoll(rt.numSides)
		if err != nil {
			return RollOutcome{}, fmt.Errorf("error rolling die: %v", err)
		}
		rawDice = append(rawDice, d)
		if rt.explodingDice && rt.isMaxDie(d) {
			var extraDie RollResult = RollResult(rt.numSides)
			for rt.isMaxDie(extraDie) {
				extraDie, err = dieRoll(rt.numSides)
				if err != nil {
					return RollOutcome{}, fmt.Errorf("error rolling die: %v", err)
				}
				rawDice = append(rawDice, extraDie)
				d += extraDie
			}
		}
		if rt.resultRule == RR_SUM {
			res += d
		} else if rt.resultRule == RR_HIGHEST {
			if d > res {
				res = d
			}
		}
	}
	res += RollResult(rt.modifier)
	return RollOutcome{res, rawDice}, nil
}
