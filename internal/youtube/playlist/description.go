package playlist

import (
	tshtml "gabe565.com/transsmute/internal/html"
	ythtml "gabe565.com/transsmute/internal/youtube/html"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func (i *Item) TemplateDescription(embed bool) g.Group {
	desc := i.Description
	desc = tshtml.Escape(desc)
	desc = tshtml.FormatHR(desc, false)
	desc = tshtml.FormatURLs(desc)
	desc = ythtml.FormatHashtags(desc)
	desc = ythtml.FormatTimestamps(i.ResourceId.VideoId, desc)
	desc = tshtml.NL2BR(desc)

	return g.Group{
		g.Iff(embed, func() g.Node {
			return html.P(
				html.IFrame(
					html.Type("text/html"),
					html.Width("640"),
					html.Height("390"),
					g.Attr("frameborder", "0"),
					html.Src(i.EmbedURL().String()),
				),
			)
		}),

		g.Raw(desc),
	}
}
