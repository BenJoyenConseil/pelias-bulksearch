package main

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/plugin/iriscontrol"
	"strconv"
	"fmt"
	"github.com/valyala/fasthttp"
)


//const HOST = "http://compute1:3100"
const HOST = "https://search.mapzen.com"

type Address struct {
	Text string 	`json:"text"`
	Size int	`json:"size"`
}


func (a *Address) EncodeForUrl() string {
	return "text=" + iris.HTMLEscape(a.Text) + "&size=" + strconv.Itoa(a.Size)
}

var c *fasthttp.Client

func main() {
	iris.Plugins.Add(iriscontrol.New(9090, map[string]string{
		"admin": "admin",
	}))
	iris.Set(iris.OptionGzip(true))
	iris.Set(iris.OptionReadBufferSize(65536))
	iris.Set(iris.OptionWriteBufferSize(65536))
	c = &fasthttp.Client{}

	iris.Get("/", Home)

	iris.Post("/v1/search", SearchBulk)

	iris.Listen(":8080")
}

func SearchBulk(ctx *iris.Context) {

	addresses := make([]*Address, 0)
	ctx.ReadJSON(&addresses)

	var geocoded string
	geocoded += "["
	for i, add := range addresses {
		url := HOST + "/v1/search?" + add.EncodeForUrl()
		statusCode, body, err := c.Get(nil, url)
		if statusCode != fasthttp.StatusOK {
			fmt.Println("Status : " + strconv.Itoa(statusCode))
		}
		if err != nil {
			fmt.Println(err)
		}
		geocoded += string(body)
		if i != len(addresses) -1 {
			geocoded += ","
		}
	}
	geocoded += "]"
	fmt.Println("Processed " + strconv.Itoa(len(addresses)) + " addresses")

	ctx.WriteString(geocoded)
}

func Home(ctx *iris.Context) {
	ctx.Write("%s", "Pelias Bulk-Search API\nAPI version : 1\nendpoint : POST -> /v1/search\nformat : \n[\n\t{\n\t\"text\":\"\",\n\t\"size\", 1\n\t}\n]")
}
