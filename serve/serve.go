package serve

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"codeberg.org/aryak/libmozhi"
	_ "codeberg.org/aryak/mozhi/docs"
	"codeberg.org/aryak/mozhi/pages"
	"codeberg.org/aryak/mozhi/public"
	"codeberg.org/aryak/mozhi/utils"
	"codeberg.org/aryak/mozhi/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	//"github.com/gofiber/fiber/v2/middleware/limiter"
	//	For debugging purposes
	//	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html"

	_ "github.com/joho/godotenv/autoload"
)

// @title			Mozhi API
// @version		1.0
// @description	API for Mozhi, the alternative-frontend for many translation engines.
// @license.name	AGPL 3.0
// @license.url	https://www.gnu.org/licenses/agpl-3.0.txt
// @BasePath		/api
func Serve(host, port string) {
	views := http.FS(views.GetFiles())
	engine := html.NewFileSystem(views, ".html")

	engine.AddFunc("newlinetobr", func(s string) template.HTML {
		return template.HTML(strings.ReplaceAll(strings.ReplaceAll(s, "\n", "<br>"), "\r", ""))
	})

	// Returns specific string if a string is in another slice
	engine.AddFunc("contains", func(data []string, s string, retval string) template.HTML {
		for _, val := range data {
			if val == s {
				return template.HTML(retval)
			}
		}
		return template.HTML("")
	})

	app := fiber.New(fiber.Config{
		Views:   engine,
		Prefork: false,
		AppName: "Mozhi",
		// kind of screwed up way to fix rate limits
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
		ProxyHeader:             fiber.HeaderXForwardedFor,
	})

	// For debugging purposes
	// app.Use(logger.New(logger.Config{
	// 	Format: "[${ip}]:${port} ${status} - ${method} ${path} ${queryParams}\n",
	// }))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	app.Use(recover.New())

	//ratelimiter := limiter.New(limiter.Config{
	//	Max:        5,
	//	Expiration: 5 * time.Minute,
	//	LimitReached: func(c *fiber.Ctx) error {
	//		return c.Status(429).Render("ratelimit_gt", fiber.Map{
	//			"Title":  "Rate limit exceeded",
	//		})
	//	},
	//})

	// add global headers
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "no-referrer")
		c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; script-src 'self' 'unsafe-inline'; frame-ancestors 'self'; form-action 'self'; base-uri 'self'; connect-src 'self';")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		return c.Next()
	})

	api := app.Group("/api")
	// Set CORS Header for all API endpoints
	api.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		return c.Next()
	})
	api.All("/translate", pages.HandleTranslate)
	api.All("/image", pages.HandleImg)
	api.Get("/source_languages", pages.HandleSourceLanguages)
	api.Get("/target_languages", pages.HandleTargetLanguages)
	api.Get("/engines", pages.HandleEngines)
	api.Get("/tts", pages.HandleTTS)
	api.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"fiberversion": fiber.Version,
			"goversion":    runtime.Version(),
		})
	})
	api.Get("/swagger/*", swagger.HandlerDefault) // default

	app.All("/", pages.HandleIndex)
	app.All("/switchlanguages", func(c *fiber.Ctx) error {
		engine := utils.Sanitize(utils.GetQueryOrFormValue(c, "engine"), "alpha")
		from := utils.GetQueryOrFormValue(c, "from")
		to := utils.GetQueryOrFormValue(c, "to")
		text := utils.GetQueryOrFormValue(c, "text")
		var swapText string
		if engine != "all" && engine != "some" && text != "" {
			translation, tlerr := libmozhi.Translate(engine, to, from, text)
			if tlerr == nil {
				swapText = translation.OutputText
			} else {
				swapText = text
			}
		}
		return c.Redirect("/?engine="+engine+"&from="+to+"&to="+from+"&text="+swapText+"&redirected=true", 301)
	})
	app.Get("/about", pages.HandleAbout)
	app.Use("/", filesystem.New(filesystem.Config{
		MaxAge: 2592000,
		Root:   http.FS(public.GetFiles()),
	}))

	hostVal, ok := os.LookupEnv("MOZHI_HOST")
	if !ok {
		hostVal = ""
	}
	if host != "" {
		hostVal = host
	}

	portVal, ok := os.LookupEnv("MOZHI_PORT")
	if !ok {
		portVal = "3000"
	}
	if port != "" {
		portVal = port
	}

	log.Fatal(app.Listen(hostVal + ":" + portVal))
}
