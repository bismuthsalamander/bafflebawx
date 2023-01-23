package playtonic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bismuthsalamander/bafflebawx/inceptor"
	"github.com/gin-gonic/gin"
)

type PlaytonicConfig struct {
	InceptorURL string
}

type BaseResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error,omitempty"`
}

type DiceResponse struct {
	BaseResponse
	RawDice []DieFace `json:"raw_dice"`
}

type RollRequest struct {
	Descriptor string `json:"descriptor"`
}

type RollResponse struct {
	DiceResponse
	Total RollResult `json:"total"`
}

type SkillCheckRequest struct {
	Difficulty int `json:"difficulty"`
	Modifier   int `json:"modifier"`
}

type SkillCheckResponse struct {
	DiceResponse
	CheckOutcome string
}

func success() BaseResponse {
	return BaseResponse{true, nil}
}

func failure(e error) BaseResponse {
	return BaseResponse{false, e}
}

func diceSuccess(r []DieFace) DiceResponse {
	return DiceResponse{success(), r}
}

func diceFailure(err error) DiceResponse {
	return DiceResponse{failure(err), make([]DieFace, 0)}
}

func skillCheckError(err error) SkillCheckResponse {
	return SkillCheckResponse{diceFailure(err), "FAILURE"}
}

func rollError(err error) RollResponse {
	return RollResponse{diceFailure(err), 0}
}

func healthCheck(rootUrl string) error {
	target := rootUrl + "/health"
	resp, err := http.Get(target)
	if err != nil {
		return fmt.Errorf("error getting /health page: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("/health returned HTTP status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var m inceptor.HealthResponse
	err = json.Unmarshal(body, &m)
	if err != nil {
		return fmt.Errorf("error unmarshaling response body: %v", err)
	}
	if !m.Success || m.Error != nil {
		return fmt.Errorf("health check failed: success %t error %v", m.Success, m.Error)
	}
	return nil
}

func Server(conf PlaytonicConfig) (*gin.Engine, error) {
	if e := healthCheck(conf.InceptorURL); e != nil {
		return nil, fmt.Errorf("error getting health check: %v", e)
	}
	r := gin.Default()
	r.POST("/skillcheck", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, skillCheckError(err))
		}
		var r SkillCheckRequest
		err = json.Unmarshal(body, &r)
		if err != nil {
			c.JSON(http.StatusBadRequest, skillCheckError(err))
			return
		}
		res, outcome, err2 := SkillCheck(r.Difficulty, r.Modifier)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, skillCheckError(err))
			return
		}
		c.JSON(http.StatusOK, SkillCheckResponse{diceSuccess(outcome.rawDice), res.String()})
	})
	r.POST("/rolldice", func(c *gin.Context) {
		var r RollRequest
		var err error
		var body []byte
		var rt RollType
		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, rollError(err))
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			c.JSON(http.StatusBadRequest, rollError(err))
			return
		}
		rt, err = ParseDescriptor(r.Descriptor)
		if err != nil {
			c.JSON(http.StatusBadRequest, rollError(err))
		}
		outcome, err2 := ExecuteRoll(rt)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, rollError(err))
			return
		}
		c.JSON(http.StatusOK, RollResponse{diceSuccess(outcome.rawDice), outcome.result})
	})
	return r, nil
}
