package playtonic

import "fmt"

type SkillCheckResult byte

const SUCCESS = SkillCheckResult('s')
const FAILURE = SkillCheckResult('f')
const CRITICAL_SUCCESS = SkillCheckResult('S')
const CRITICAL_FAILURE = SkillCheckResult('F')

func (r SkillCheckResult) String() string {
	switch r {
	case SUCCESS:
		return "success"
	case FAILURE:
		return "failure"
	case CRITICAL_SUCCESS:
		return "critical success"
	case CRITICAL_FAILURE:
		return "critical failure"
	default:
		return "unknown"
	}
}

func SkillCheck(difficulty int, modifiers int) (SkillCheckResult, RollOutcome, error) {
	descriptor := "1d20"
	if modifiers > 0 {
		descriptor += "+" + fmt.Sprintf("%d", modifiers)
	} else if modifiers < 0 {
		descriptor += "-" + fmt.Sprintf("%d", modifiers*-1)
	}
	rt, err1 := ParseDescriptor(descriptor)
	if err1 != nil {
		return FAILURE, RollOutcome{}, err1
	}
	outcome, err := ExecuteRoll(rt)
	if err != nil || err1 != nil {
		return FAILURE, outcome, err
	}
	if outcome.rawDice[0] == 1 {
		return CRITICAL_FAILURE, outcome, nil
	} else if outcome.rawDice[0] == 20 {
		return CRITICAL_SUCCESS, outcome, nil
	}
	if outcome.result >= RollResult(difficulty) {
		return SUCCESS, outcome, nil
	}
	return FAILURE, outcome, nil
}
