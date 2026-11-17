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
//	@Summary		Health of the enabled engines and Japanese dictionaries, from real traffic
//	@Description	Without engine: every enabled engine and dictionary, always 200. With engine: that one only, re-checked first if it is failing, has had no request since start or none for 12 hours (at most every 2 minutes), and 503 while it is failing.
//	@Param			engine	query	string	false	"Engine or dictionary name (jisho, jpdb, moji, weblio)"
//	@Param			kind	query	string	false	"translate, tts or dict (default: the first one it has)"
//	@Router			/api/health [get]
func HandleHealth(c *fiber.Ctx) error {
	engine := utils.Sanitize(c.Query("engine"), "alpha")
	if engine == "" {
		engines := fiber.Map{}
		for _, id := range append(utils.EnabledEngines(), utils.JaDictSources()...) {
			engines[id] = libmozhi.EngineHealth(id)
		}
		return c.JSON(fiber.Map{"since": started, "engines": engines})
	}
	kinds := libmozhi.HealthKinds(engine)
	if len(kinds) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "unknown engine "+engine)
	}
	// A disabled engine is not a failing one.
	if !utils.EngineEnabled(engine) && !utils.JaDictEnabled(engine) {
		return c.JSON(fiber.Map{"status": "disabled"})
	}
	kind := c.Query("kind", kinds[0])
	h, ok := libmozhi.HealthCheck(engine, kind)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "engine "+engine+" has no "+kind)
	}
	if h.Status == "failing" {
		c.Status(fiber.StatusServiceUnavailable)
	}
	return c.JSON(h)
}
