package kemono

import (
	tshtml "gabe565.com/transsmute/internal/html"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func (p *Post) TemplateDescription() g.Group {
	return g.Group{
		g.Iff(p.Content != "", func() g.Node {
			return g.Group{
				html.H3(g.Text("Content")),
				g.Raw(tshtml.FormatHR(p.Content, true)),
			}
		}),

		g.Iff(len(p.Attachments) != 0, func() g.Node {
			return g.Group{
				html.H3(g.Text("Files")),
				g.Map(p.Attachments, func(a *Attachment) g.Node {
					return html.P(
						html.A(
							html.Href(a.URL().String()),
							func() g.Node {
								if a.IsImage() {
									return html.Img(
										html.Src(a.ThumbURL().String()),
										html.Alt(a.Name),
										html.Title(a.Name),
									)
								}
								return g.Text(a.Name)
							}(),
						),
						g.Iff(a.IsVideo(), func() g.Node {
							return html.Video(html.Controls(),
								html.Source(html.Src(a.URL().String())),
							)
						}),
					)
				}),
			}
		}),

		g.Iff(p.Embed.Subject != "" && p.Embed.URL != "", func() g.Node {
			return html.P(
				html.A(html.Href(p.URL().String()),
					html.Strong(
						g.Text(p.Embed.Subject),
					),
				),
			)
		}),

		g.Iff(len(p.Tags) != 0, func() g.Node {
			return g.Group{
				html.H3(g.Text("Tags")),
				html.P(
					g.Map(p.Tags, func(t string) g.Node {
						return g.Group{
							g.Raw("\n"),
							html.A(html.Href(p.Creator.TagURL(t).String()),
								g.Text(t),
							),
						}
					}),
				),
			}
		}),
	}
}
