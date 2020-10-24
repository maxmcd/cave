package cave

type Renderer interface {
	Render(ctx Context, ui UI) (err error)
}

type Context struct {
}
type UI struct {
}

func DIV(f ...interface{}) func(ctx Context, ui UI) (err error) {
	return func(ctx Context, ui UI) (err error) {

		return nil
	}
}

func Render(ctx Context, ui UI) (err error) {
	return nil
}

type Cave struct {
	renderer Renderer
}

func (c *Cave) NewHook() Hook {
	return Hook{}
}

type Hook struct {
	v interface{}
	c *Cave
}

func (h *Hook) SetValue()
