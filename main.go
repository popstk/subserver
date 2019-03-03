package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"github.com/astaxie/beego"
	"log"
	"os"
	"strings"
)

var (
	configFile string
	config     Config
)

func init() {
	log.SetFlags(log.Lshortfile|log.Ltime)
	flag.StringVar(&configFile, "c", "config.json", "config file")
}


func parseConfig(root string) ([]string, error) {
	urls := make([]string, 0)
	m := make(map[string]bool)

	q := make([]string, 0, 1)
	q = append(q, root)
	for len(q) > 0 {

		uuid := q[0]
		q = q[1:]
		m[uuid] = true

		source, exist := config.Valid[uuid]
		if !exist {
			return nil, errors.New("Invalid uuid: "+uuid)
		}

		for _, s := range source {
			if s.Type == "sub" {
				_, exist := m[s.Addr]
				if !exist {
					q = append(q, s.Addr)
				}

				continue
			}

			u, err := s.Parse()
			if err != nil {
				log.Print(err)
				continue
			}
			urls = append(urls, u...)
		}
	}

	return urls, nil
}

// Controller -
type Controller struct {
	beego.Controller
}


// Get -
func (c *Controller) Get() {
	uuid := c.Ctx.Input.Param(":uuid")

	log.Print("Get uuid = ", uuid)

	urls, err := parseConfig(uuid)
	if err != nil {
		log.Println(err)
		c.Ctx.WriteString("")
		return
	}

	valid := make([]string, 0, len(urls))
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u != "" {
			valid =append(valid, u)
		}
	}

	respond := strings.Join(valid, "\n")
	if strings.ToLower(c.Ctx.Input.Query("raw")) != "1" {
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
