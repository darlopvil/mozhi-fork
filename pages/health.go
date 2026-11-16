package pages

import (
	"time"

	"codeberg.org/aryak/libmozhi"
	"codeberg.org/aryak/mozhi/utils"
	"github.com/gofiber/fiber/v2"
)

// started is when this instance started: engines with no request since then
// are "unknown".
var started = time.Now()

// HandleHealth godoc
//
//	@Summary		Health of the enabled engines, from real traffic
//	@Description	Without engine: every enabled engine, always 200. With engine: that engine only, re-checked first if its last request failed (at most every 2 minutes), and 503 while it is failing.
//	@Param			engine	query	string	false	"Engine name"
//	@Param			kind	query	string	false	"translate (default) or tts"
//	@Router			/api/health [get]
func HandleHealth(c *fiber.Ctx) error {
	engine := utils.Sanitize(c.Query("engine"), "alpha")
	if engine == "" {
		engines := fiber.Map{}
		for _, id := range utils.EnabledEngines() {
			engines[id] = libmozhi.EngineHealth(id)
		}
		return c.JSON(fiber.Map{"since": started, "engines": engines})
	}
	if len(libmozhi.HealthKinds(engine)) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "unknown engine "+engine)
	}
	// A disabled engine is not a failing one.
	if !utils.EngineEnabled(engine) {
		return c.JSON(fiber.Map{"status": "disabled"})
	}
	kind := c.Query("kind", "translate")
	h, ok := libmozhi.HealthCheck(engine, kind)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "engine "+engine+" has no "+kind)
	}
	if h.Status == "failing" {
		c.Status(fiber.StatusServiceUnavailable)
	}
	return c.JSON(h)
}
