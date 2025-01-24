package crncl

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type AsideBlockExtension struct{}

func (ext AsideBlockExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(AsideBlockParser{}, 500),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(AsideBlockRenderer{}, 500),
	))
}

var AsideBlockNodeKind = ast.NewNodeKind("AsideBlock")

var _ ast.Node = &AsideBlockNode{}

type AsideBlockNode struct {
	ast.BaseBlock
	Title string
}

func (n *AsideBlockNode) Dump(source []byte, level int) {
	m := map[string]string{
		"Title": n.Title,
	}
	ast.DumpHelper(n, source, level, m, nil)
}

func (n *AsideBlockNode) Kind() ast.NodeKind {
	return AsideBlockNodeKind
}

var _ parser.BlockParser = AsideBlockParser{}

type AsideBlockParser struct{}

// Trigger returns a list of characters that triggers Parse method of
// this parser.
// If Trigger returns a nil, Open will be called with any lines.
func (p AsideBlockParser) Trigger() []byte {
	return []byte{':'}
}

// Open parses the current line and returns a result of parsing.
//
// Open must not parse beyond the current line.
// If Open has been able to parse the current line, Open must advance a reader
// position by consumed byte length.
//
// If Open has not been able to parse the current line, Open should returns
// (nil, NoChildren). If Open has been able to parse the current line, Open
// should returns a new Block node and returns HasChildren or NoChildren.
func (p AsideBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if len(line) < 3 {
		return nil, parser.NoChildren // clearly not open to aside block
	}
	if !bytes.Equal(line[:3], []byte(":::")) {
		return nil, parser.NoChildren // if not expected open, move on
	}
	reader.Advance(len(line))

	return &AsideBlockNode{
		Title: strings.TrimSpace(string(line[3:])),
	}, parser.HasChildren
}

// Continue parses the current line and returns a result of parsing.
//
// Continue must not parse beyond the current line.
// If Continue has been able to parse the current line, Continue must advance
// a reader position by consumed byte length.
//
// If Continue has not been able to parse the current line, Continue should
// returns Close. If Continue has been able to parse the current line,
// Continue should returns (Continue | NoChildren) or
// (Continue | HasChildren)
func (p AsideBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	if len(line) >= 3 && bytes.Equal(line[:3], []byte(":::")) { // in the case the first 3 bytes are ":::", we are done
		reader.Advance(len(line))
		return parser.Continue | parser.NoChildren
	}

	return parser.Continue | parser.HasChildren
}

// Close will be called when the parser returns Close.
func (p AsideBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	//no-op
}

// CanInterruptParagraph returns true if the parser can interrupt paragraphs,
// otherwise false.
func (p AsideBlockParser) CanInterruptParagraph() bool {
	return false
}

// CanAcceptIndentedLine returns true if the parser can open new node when
// the given line is being indented more than 3 spaces.
func (p AsideBlockParser) CanAcceptIndentedLine() bool {
	return false
}

type AsideBlockRenderer struct{}

func (r AsideBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(AsideBlockNodeKind, r.render)
}

func (r AsideBlockRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	asideBlock := node.(*AsideBlockNode)
	if entering {
		tpl, err := template.New("aside").Parse(`<aside class="bg-blue-100 p-4"><h3>{{ . }}</h3>`)
		if err != nil {
			return ast.WalkStop, fmt.Errorf("parsing aside template: %w", err)
		}
		err = tpl.Execute(w, asideBlock.Title)
		if err != nil {
			return ast.WalkStop, fmt.Errorf("parsing aside template: %w", err)
		}
	} else {
		w.WriteString("</aside>")
	}

	return ast.WalkContinue, nil
}
