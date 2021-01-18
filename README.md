# Cave

Cave allows you to write fully functional interactive web applications with Go. You describe your templates and components in your Go program and Cave takes care of rendering your UI to the browser and pushing updates from user changes.

Cave is similar to [Phoenix LiveView](https://hexdocs.pm/phoenix_live_view/Phoenix.LiveView.html) although it was originally inspired by [Dash](https://plotly.com/dash/).

*Cave is Alpha software and incomplete. Don't take it seriously!*

## A Minimal Example

Let's walk through a minimal example to get a feel for the structure of Cave.

Here's a simple component:

```go
type SimpleComponent struct {
	Count    int
}

var _ cave.Renderer = new(ToDoApp) // this just ensures that we are implementing this interface

func (tda *SimpleComponent) Render() string {
	return `<div>{{ .Count }}</div>`
}
```

This component implements the `cave.Renderer` interface, which is the minimum requirement for a component. When rendering it simply outputs a div with a count value.

You can't do much with this component in its current form, but this is one of our core building blocks.

## Something More Complicated

Here's the cave version of a ToDo list app. It has a form input and adds items to the list when the form is submitted.

*If you just want to play around with the final result you [can try it out here](https://cave-demo.fly.dev/). The full source is also [here](./examples/to-do).*

```go
type ToDoApp struct {
	Items    []string
	ToDoList *ToDoList
}

var (
	_ cave.OnSubmiter = new(ToDoApp)
	_ cave.Renderer   = new(ToDoApp)
)

func (tda *ToDoApp) OnSubmit(name string, form map[string]string) {
	if name == "todo" {
		tda.Items = append(tda.Items, form["new"])
	}
}
func (tda *ToDoApp) Render() string {
	return `
	<div>
	<h3>TODO</h3>
	{{ render .ToDoList }}
	<form cave-submit=todo>
	  <label for="new-todo">
		What needs to be done?
	  </label>
	  <input type="text" name="new" />
	  <button>
		Add {{ len .Items | add 1 }}
	  </button>
	</form>
  </div>
`
}

type ToDoList struct {
	ToDoApp *ToDoApp
}

var _ cave.Renderer = new(ToDoList)

func (tdl *ToDoList) Render() string {
	return `
	<ul>
	 {{range .ToDoApp.Items }}
	 	<li>{{.}}</li>
	 {{end}}
	</ul>
`
}
```

We've introduced two things. The `cave.OnSubmiter` interface, and the curious tag `cave-submit`. Both of these things work together! When this component is rendered in the browser, Cave listens for submit events on our form and when they are made it shoots the form details over a websocket. The server then call `OnSubmit` on this component, computes the resulting changes in the HTML and shoots them back over a websocket. Outrageous!

## The Nitty Gritty

Now you're probably thinking "Rendering basic UI changes by pushing bytes over thousands of miles, count me in! How do I plug this thing into a server and start adding latency to my user experiences?". Well let's show you how!

Now that we have components we'll need to hook them up to a web server.

Let's create a new Cave, we'll call it cavern because that's cute. Let's put a layout and a component in our cavern.
```go
cavern := cave.New()
if err := cavern.AddTemplateFile("main", "layout.html"); err != nil {
	log.Fatal(err)
}
cavern.AddComponent("main", NewToDoApp)
```

`AddComponent` takes a `func() cave.Renderer` so that it can create a new component for every request. So we'll need to set up that function as well.

```go
func NewToDoApp() cave.Renderer {
	tda := &ToDoApp{Items: []string{"breathe"}}
	tld := &ToDoList{}
	tda.ToDoList = tld
	tld.ToDoApp = tda
	return tda
}
```

Layouts are html pages that render the html boilerplate we'll need outside of our components. A minimal example would be this:
```html
<!doctype html>
<html>
    <head><title>Hello</title></head>
    <body>
        {{ component "main" }}
        <script src="/bundle.js" type="application/javascript"></script>
    </body>
</html>
```
We need to load the javascript bundle that contains all the Cave goodies, and we need to mount the component "main" right where we want it.

Next we'll need to actually serve the page. I'm going to use [gin](https://github.com/gin-gonic/gin), but you could technically use anything that uses `http.ResponseWriter` and `*http.Response`.

```go
r.Use(func(c *gin.Context) {
	if _, ok := c.Request.URL.Query()["cavews"]; ok {
		cavern.ServeWS(c.Writer, c.Request)
		c.Abort()
	}
})
r.GET("/", func(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "text/html")
	_ = cavern.Render("main", c.Writer)
})
r.GET("/bundle.js", func(c *gin.Context) {
	cavern.ServeJS(c.Writer, c.Request)
})
```

A few things going on here:

1. We intercept all request with the query param `cavews` and assume they are websocket requests intended for cave.
2. We render out "main" layout at the root path.
3. We serve our bundle where out layout expects it.

That's it! Everything else is websocket magic, crazy hacks, and the strange feeling that we're making progress technically while regressing at the same time.
