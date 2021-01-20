package cave

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubcomponentBasic(t *testing.T) {
	tcp := NewTestComponentParent()
	tcp.(*TestComponentParent).Count2.Count = 2

	lc, err := newLiveComponent(tcp)
	if err != nil {
		t.Fatal(err)
	}
	{
		var sb strings.Builder
		if err := lc.render(tcp, &sb); err != nil {
			log.Fatal(err)
		}
		// correct components in the correct order
		// subcomponent indexing starts at zero
		assert.Equal(t, sb.String(), "<div cave-subcomponent=\"0\"><div>5</div></div>\n\t<div cave-subcomponent=\"1\"><div>2</div></div>\n\t<div cave-subcomponent=\"2\"><div>0</div></div>")
	}
	tcp.(*TestComponentParent).order = true
	{
		var sb strings.Builder
		if err := lc.render(tcp, &sb); err != nil {
			log.Fatal(err)
		}
		// we reverse the components, but the rendering still
		// uses the correct indexes
		assert.Equal(t, sb.String(), "<div cave-subcomponent=\"2\"><div>0</div></div>\n\t\t<div cave-subcomponent=\"1\"><div>2</div></div>\n\t\t<div cave-subcomponent=\"0\"><div>5</div></div>")
	}
}
