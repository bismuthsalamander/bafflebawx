package playtonic

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bismuthsalamander/bafflebawx/inceptor"
)

type DiceCount uint8
type RollResult int32
type DieFace uint16
type RollDescriptor string
type Modifier int16
type ResultRule uint8 //this would be an enum in C

const RR_SUM = 0     //result of a roll is the sum of the dice
const RR_HIGHEST = 1 //result of a roll is the highest die rolled

const SUFFIXES = "^*" //^ means highest die, not sum; * means exploding dice

type RollType struct {
	numDice       DiceCount
	numSides      DieFace
	modifier      Modifier
	resultRule    ResultRule
	explodingDice bool
}

type RollOutcome struct {
	result  RollResult
	rawDice []DieFace
}

func (rt RollType) isMaxDie(r DieFace) bool {
	return r == rt.numSides
}

func dieRoll(sides DieFace) (DieFace, error) {
	n, err := inceptor.Uint64()
	if err != nil {
		return 0, err
	}
	return DieFace((n % uint64(sides)) + 1), nil
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
			return rt, fmt.Errorf("could not parse dice count '%s' (descriptor %s)", parts[0], desc)
		}
		if dicecount > math.MaxUint8 || dicecount <= 0 {
			return rt, fmt.Errorf("dice count %d out of range [1, %d] (descriptor %s)", dicecount, math.MaxUint8, desc)
		}
		rt.numDice = DiceCount(dicecount)
	}
	mod_index := strings.IndexAny(parts[1], "+-")
	if mod_index != -1 {
		mod_str := parts[1][mod_index:]
		parts[1] = parts[1][:mod_index]
		if mod_str != "" {
			mod, err2 := strconv.ParseInt(mod_str[1:], 10, 64)
			if err2 != nil {
				return RollType{}, fmt.Errorf("could not parse modifier magnitude '%s' (descriptor %s)", mod_str[1:], desc)
			}
			if mod_str[0] == '-' {
				mod *= -1
			}
			//We could also manually compare to MinInt16 and MaxInt16, but this
			//approach is a little more robust if Modifier later changes type
			rt.modifier = Modifier(mod)
			if mod != int64(rt.modifier) {
				return RollType{}, fmt.Errorf("modifier %s outside of range [%d, %d] (descriptor %s)", mod_str[1:], math.MinInt16, math.MaxInt16, desc)
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
	rt.numSides = DieFace(sides)
	return rt, nil
}

func ExecuteRoll(rt RollType) (RollOutcome, error) {
	var res RollResult = 0
	//theoretically, the loop may be faster for non-exploding dice if we made
	//rawDice with the correct number of slots, but we'd have to either write
	//two separate versions of the loop or add some ugly logic to check the
	//slice's size on each iteration.
	rawDice := make([]DieFace, 0)
	for i := 0; DiceCount(i) < rt.numDice; i++ {
		d, err := dieRoll(rt.numSides)
		if err != nil {
			return RollOutcome{}, fmt.Errorf("error rolling die: %v", err)
		}
		rawDice = append(rawDice, d)
		if rt.explodingDice && rt.isMaxDie(d) {
			extraDie := rt.numSides
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
			res += RollResult(d)
		} else if rt.resultRule == RR_HIGHEST {
			if RollResult(d) > res {
				res = RollResult(d)
			}
		}
	}
	res += RollResult(rt.modifier)
	return RollOutcome{res, rawDice}, nil
}
