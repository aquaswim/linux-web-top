package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	htmlEngine "github.com/gofiber/template/html/v2"
	"linux-web-top/data"
	"linux-web-top/util"
	"linux-web-top/websocket"
	"log"
	"os"
	"os/signal"
	"time"
)

var tickDuration time.Duration
var listenAddr string

func init() {
	var durationStr string
	flag.StringVar(&durationStr, "d", "3s", "tick duration")
	flag.StringVar(&listenAddr, "l", ":3000", "listen address")
	flag.Parse()
	d, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Panicln("Error parse tick duration", durationStr, "err:", err)
	}
	tickDuration = d
}

func main() {
	app := fiber.New(fiber.Config{
		Views: htmlEngine.New("./web", ".tmpl.html"),
	})

	app.Static("/public", "./web/public")
	app.Get("/", func(ctx *fiber.Ctx) error {
		//return ctx.SendFile("./web/index.tmpl.html")
		return ctx.Render("index", fiber.Map{
			"Hostname": util.GetHostname(),
		})
	})

	hub := websocket.NewHub()
	ticker := time.NewTicker(tickDuration)
	stop := make(chan bool)
	go func() {
		stat := data.Stat{}
		if err := util.PopulateNCpu(&stat); err != nil {
			log.Panicln("Failed to read /proc/stat:", err)
		}
		for {
			select {
			case <-ticker.C:
				if hub.IsNotEmpty() {
					if err := util.PopulateMemInfo(&stat); err != nil {
						log.Println("Failed to read memory info", err)
					}

					if err := util.PopulateCpuInfo(&stat); err != nil {
						log.Println("Failed to read cpu stats", err)
					}

					if err := util.PopulateNetworkBandwith(&stat); err != nil {
						log.Println("Failed to read network bandwith:", err)
					}
					hub.Broadcast(&stat)
				}
			case <-stop:
				log.Println("Stopping ticker")
				ticker.Stop()
			}
		}
	}()
	wsRoute := app.Group("/ws", websocket.Upgrade)
	wsRoute.Get("/stats", websocket.NewGetStatHandler(&hub))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Gracefully shutting down...")
		stop <- true
		_ = app.Shutdown()
	}()
	if err := app.Listen(listenAddr); err != nil {
		log.Panic(err)
	}
	log.Println("Running cleanup tasks...")
}
