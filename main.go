package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/astaxie/beego"
)

var (
	configFile string
	config     Config
)

func init() {
	flag.StringVar(&configFile, "c", "config.json", "config file")
}

// Controller -
type Controller struct {
	beego.Controller
}

// Get -
func (c *Controller) Get() {
	uuid := c.Ctx.Input.Param(":uuid")
	var respond string
	source, exist := config.Valid[uuid]
	if !exist {
		log.Print("Invalid uuid: ", uuid)
		c.Ctx.WriteString("")
		return
	}

	urls := make([]string, 0)
	for _, s := range source {
		u, err := s.Parse()
		if err != nil {
			log.Print(err)
			continue
		}
		urls = append(urls, u...)
	}

	respond = strings.Join(urls, "\n")
	if strings.ToLower(c.Ctx.Input.Query("raw")) != "true" {
		respond = base64.StdEncoding.EncodeToString([]byte(respond))
	}

	c.Ctx.WriteString(respond)
	return
}

func main() {
	flag.Parse()

	f, err := os.Open(configFile)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	if err = json.NewDecoder(f).Decode(&config); err != nil {
		log.Fatal(err)
	}

	beego.Router("/:uuid", &Controller{})
	beego.Run(config.Address)
}
