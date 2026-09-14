package preprocessor

import (
	"strings"
	"testing"
)

func TestInjectScopeAttributeTagsTopLevelElements(t *testing.T) {
	content := `<p>one</p>
<section><p>two</p></section>
<img src="x.png">
`

	scoped := injectScopeAttribute(content, `data-__2_scope="1"`)

	expected := `<p data-__2_scope="1">one</p>
<section data-__2_scope="1"><p>two</p></section>
<img src="x.png" data-__2_scope="1">
`

	if scoped != expected {
		t.Errorf("unexpected scoped content:\n%s", scoped)
	}
}

func TestInjectScopeAttributeSkipsScriptAndStyleBlocks(t *testing.T) {
	content := `<style>p > b { color: red; }</style>
<script compiled>
	$value = "<p>not markup</p>";
</script>
<div>content</div>
`

	scoped := injectScopeAttribute(content, `data-__2_scope="2"`)

	expected := `<style>p > b { color: red; }</style>
<script compiled>
	$value = "<p>not markup</p>";
</script>
<div data-__2_scope="2">content</div>
`

	if scoped != expected {
		t.Errorf("unexpected scoped content:\n%s", scoped)
	}
}

func TestInjectScopeAttributeHandlesQuotedAngleBrackets(t *testing.T) {
	content := `<a title="a > b" href="/">link</a><br>`

	scoped := injectScopeAttribute(content, `data-__2_scope="3"`)

	if !strings.Contains(scoped, `<a title="a > b" href="/" data-__2_scope="3">`) {
		t.Errorf("the quoted \">\" broke the tag scan:\n%s", scoped)
	}

	if !strings.Contains(scoped, `<br data-__2_scope="3">`) {
		t.Errorf("expected the void element to be tagged:\n%s", scoped)
	}
}

func TestScopeComponentWrapsStyleBlocks(t *testing.T) {
	content := `<style>p { color: red; }</style>
<p>badge</p>
`

	scoped := scopeComponent(content, 4)

	expected := `<style>
@scope ([data-__2_scope="4"]) {
p { color: red; }
}
</style>
<p data-__2_scope="4">badge</p>
`

	if scoped != expected {
		t.Errorf("unexpected scoped component:\n%s", scoped)
	}
}

func TestHasStyleBlock(t *testing.T) {
	if hasStyleBlock("<p>no styles</p>") {
		t.Error("expected no style block")
	}

	if !hasStyleBlock("<STYLE>p { }</STYLE>") {
		t.Error("expected the case insensitive style block to be found")
	}
}
