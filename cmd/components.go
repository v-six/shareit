package main

import (
	"net/http"
	"time"

	"github.com/3di-clockwork/devops-test/cmd/types"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

const (
	title = "Share It!"
)

func Home() g.Node {
	return page(Article(
		Header(Strong(g.Raw("🔗 Share a new file"))),
		Form(Action("/files"), Method("post"), EncType("multipart/form-data"),
			Input(Name("file"), Type("file"), Required()),
			FieldSet(
				Legend(g.Raw("How long shall we keep your file?")),
				Label(Input(Type("radio"), Name("ttl"), Value("1"), Checked()), g.Raw("1 Hour")),
				Label(Input(Type("radio"), Name("ttl"), Value("24")), g.Raw("1 Day")),
				Label(Input(Type("radio"), Name("ttl"), Value("720")), g.Raw("1 Month")),
			),
			Input(Type("submit"), Value(title)),
		),
	))
}

func FileNotFound() g.Node {
	return page(Article(
		Header(Strong(g.Raw("❌ Error"))),
		Span(g.Raw("This file does not exist")),
	))
}

func FileDetail(r *http.Request, content *types.Content) g.Node {
	baseUrl := Config.PublicURL
	if len(baseUrl) == 0 {
		scheme := r.URL.Scheme
		if len(scheme) == 0 {
			scheme = "http"
		}
		baseUrl = scheme + "://" + r.Host
	}

	return page(Article(
		Header(Strong(g.Raw("📦 "+content.Filename))),
		Table(TBody(
			Tr(Td(g.Raw("Filename")), Td(g.Raw(content.Filename))),
			Tr(Td(g.Raw("Size")), Td(g.Rawf("%d kB", content.Size>>10))),
			Tr(Td(g.Raw("Expiry")), Td(g.Raw(content.Expiry.Format(time.Stamp)))),
			Tr(Td(g.Raw("Link")), Td(Code(g.Rawf("%s/files/%s", baseUrl, content.ID)))),
		)),
		Form(Action("/files/"+string(content.ID)+"/raw"), Method("get"), Target("_blank"),
			Input(Type("submit"), Value("Download!")),
		),
	))
}

func page(mainContent ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.classless.min.css")),
		},
		Body: []g.Node{header(), Main(mainContent...), footer()},
	})
}

func header() g.Node {
	return Header(A(Href("/"), H1(g.Raw(title))))
}

func footer() g.Node {
	return Footer(g.Rawf("made with ❤️ by the <code>clockwork</code> team"))
}
