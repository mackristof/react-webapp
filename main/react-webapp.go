package main

import "github.com/gin-gonic/gin"
import "github.com/trustmaster/goflow"
import "fmt"

// struct Greeter used communication
type Greeter struct {
	flow.Component
	Name <-chan string
	Res  chan<- string
}

//
func (g *Greeter) OnName(name string) {
	greeting := fmt.Sprintf("Hello, %s!", name)
	g.Res <- greeting
}

type Logger struct {
	flow.Component
	Line <-chan string
	Res  chan<- string
}

func (p *Logger) OnLine(line string) {
	fmt.Println(line)
	p.Res <- line
}

type GreetingFlow struct {
	flow.Graph
}

func NewGreetingFlow() *GreetingFlow {
	n := new(GreetingFlow)
	n.InitGraphState()
	n.Add(new(Greeter), "greeter")
	n.Add(new(Logger), "logger")
	n.Connect("greeter", "Res", "logger", "Line")
	n.MapInPort("In", "greeter", "Name")
	n.MapOutPort("Out", "logger", "Res")
	return n
}

func main() {
	net := NewGreetingFlow()
	in := make(chan string)
	out := make(chan string)
	net.SetInPort("In", in)
	net.SetOutPort("Out", out)
	flow.RunNet(net)

	router := gin.Default()
	router.GET("/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		in <- name
		var resp struct {
			Title string `json:"greeter"`
			Value string
		}
		resp.Title = "golang"
		resp.Value = <-out
		c.JSON(200, resp)
	})
	router.Run(":8080")
	close(in)
	close(out)
	<-net.Wait()
}
