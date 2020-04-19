package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"

	"github.com/petderek/gogofer"
)

var (
	port = flag.Int("p", 7070, "the port to listen on")
	host = flag.String("h", "127.0.0.1", "the host to listen on")
)

const (
	gophermap = `iInfo!
iThis is an example of a gophermap file. gogoferd currently supports
ithe following types locally: 1,0,i. images, binaries, and sound will
iall be supported someday
0Lorem Ipsum	ipsum.txt`

	ipsum = `90's +1 bespoke tilde quinoa. Adaptogen lomo disrupt leggings succulents iPhone
brooklyn lo-fi cronut ramps pug four dollar toast aesthetic taxidermy woke.
Mixtape quinoa small batch selvage kombucha, subway tile sriracha bushwick fam
kinfolk hashtag microdosing. Echo park health goth authentic, etsy fixie la
croix 3 wolf moon jean shorts deep v pickled thundercats hella slow-carb. Twee
gentrify heirloom master cleanse scenester try-hard keytar, succulents cronut
banjo blog. Try-hard biodiesel shaman, meggings woke PBR&B iPhone craft beer
YOLO coloring book franzen 90's. Leggings hexagon butcher tote bag gochujang 
edison bulb sustainable fashion axe pickled pok pok narwhal typewriter subway                                           
tile squid.

You probably haven't heard of them single-origin coffee crucifix poutine four
dollar toast seitan live-edge. Stumptown af raclette air plant, copper mug
polaroid hell of keytar kitsch semiotics bitters messenger bag intelligentsia
actually. Bespoke pop-up yr butcher everyday carry cloud bread VHS health goth
coloring book retro pabst polaroid mumblecore XOXO. Austin bitters banjo, irony
selvage umami ethical jianbing quinoa portland. Readymade tumeric photo booth
chambray direct trade tbh polaroid kombucha. Tofu master cleanse before they
sold out affogato fanny pack irony portland taxidermy forage biodiesel hoodie
VHS roof party. Umami portland tacos prism forage brunch, typewriter four
dollar toast.`
)

// StaticGopherMapHandler shows an example of a custom handler using the gophermap type.
// It returns the GopherMap for an empty request, or some ipsum for a request to
// /ipsum.txt
type StaticGopherMapHandler struct{}

func (s *StaticGopherMapHandler) Serve(request gogofer.Selector) gogofer.Response {
	switch request.Path {
	case "":
		return gogofer.NewGopherMap(bytes.NewBufferString(gophermap), *host, *port)
	case "ipsum.txt":
		return &gogofer.StaticTextResponse{
			Message: []byte(ipsum),
		}
	default:
		return nil
	}
}

func main() {
	flag.Parse()
	server := &gogofer.Server{
		Addr:    fmt.Sprintf("%s:%d", *host, *port),
		Handler: &StaticGopherMapHandler{},
	}

	log.Fatal(server.ListenAndServe())
}
