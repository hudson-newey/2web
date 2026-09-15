package preprocessor

import (
	"strings"
	"testing"
)

func TestExtractSlotContent(t *testing.T) {
	children := `
  <h2 slot="title">My Card</h2>
  <div slot="footer"><span>nested</span></div>
  <p>body content</p>
`

	named, defaultContent := extractSlotContent(children)

	if named["title"] != "<h2>My Card</h2>" {
		t.Errorf("expected the title slot content without the slot attribute, got %q", named["title"])
	}

	if named["footer"] != `<div><span>nested</span></div>` {
		t.Errorf("expected the nested footer content, got %q", named["footer"])
	}

	if !strings.Contains(defaultContent, "<p>body content</p>") {
		t.Errorf("expected the default content, got %q", defaultContent)
	}

	if strings.Contains(defaultContent, "slot=") {
		t.Errorf("the default content must not contain slot assignments, got %q", defaultContent)
	}
}

func TestExtractSlotContentWithoutSlots(t *testing.T) {
	named, defaultContent := extractSlotContent(`<p>only default</p>`)

	if len(named) != 0 {
		t.Errorf("expected no named slots, got %v", named)
	}

	if !strings.Contains(defaultContent, "<p>only default</p>") {
		t.Errorf("expected the default content, got %q", defaultContent)
	}
}

func TestSpliceSlots(t *testing.T) {
	component := `<header><slot name="title">fallback title</slot></header>
<p><slot>fallback body</slot></p>`

	// Filled slots.
	spliced := spliceSlots(component, `<h2 slot="title">My Card</h2><p>body</p>`)

	if strings.Contains(spliced, "<slot") {
		t.Errorf("expected every slot to be replaced:\n%s", spliced)
	}

	if !strings.Contains(spliced, "<header><h2>My Card</h2></header>") {
		t.Errorf("expected the named slot content:\n%s", spliced)
	}

	if !strings.Contains(spliced, "<p>body</p>") {
		t.Errorf("expected the default slot content:\n%s", spliced)
	}

	// Unfilled slots render their fallback content.
	spliced = spliceSlots(component, "")

	if !strings.Contains(spliced, "fallback title") || !strings.Contains(spliced, "fallback body") {
		t.Errorf("expected the fallback content:\n%s", spliced)
	}
}

func TestMatchedElementLengthHandlesNesting(t *testing.T) {
	content := `<div slot="x"><div>inner</div></div>after`

	length := matchedElementLength(content, "div")

	if content[:length] != `<div slot="x"><div>inner</div></div>` {
		t.Errorf("unexpected element match: %q", content[:length])
	}
}
