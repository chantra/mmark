package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mparser"
	"github.com/mmarkdown/mmark/xml"
)

func TestMmark(t *testing.T) {
	// open all *.md files and test them
	dir := "testdata"
	testFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not read %s: %q", dir, err)
	}
	for _, f := range testFiles {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) != ".md" {
			continue
		}
		base := f.Name()[:len(f.Name())-3]

		doTest(t, base)
	}
}

var ext = parser.CommonExtensions | parser.HeadingIDs | parser.AutoHeadingIDs | parser.Footnotes |
	parser.OrderedListStart | parser.Attributes | parser.Mmark

func doTest(t *testing.T, basename string) {
	p := parser.NewWithExtensions(ext)
	cwd := mparser.NewCwd()
	p.Opts = parser.ParserOptions{
		ParserHook:    mparser.TitleHook,
		ReadIncludeFn: cwd.ReadInclude,
	}
	opts := xml.RendererOptions{
		Flags: xml.CommonFlags | xml.XMLFragment,
	}
	renderer := xml.NewRenderer(opts)

	filename := filepath.Join("testdata", basename+".md")
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't open '%s', error: %v\n", filename, err)
		return
	}

	filename = filepath.Join("testdata", basename+".xml")
	expected, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't open '%s', error: %v\n", filename, err)
	}
	expected = bytes.TrimSpace(expected)

	actual := markdown.ToHTML(input, p, renderer)
	actual = bytes.TrimSpace(actual)
	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("\n    [%#v]\nExpected[%s]\nActual  [%s]",
			basename+".md", expected, actual)
	}
}
