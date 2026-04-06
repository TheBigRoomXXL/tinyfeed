<opml version="2.0">
	<head>
		<title>{{.Metadata.name}}</title>
		<dateCreated>{{.Metadata.datetimeRFC822}}</dateCreated>
		<docs>"https://opml.org/spec2.opml"</docs>
	</head>
	<body>
	{{range .Feeds}}
		<outline
			text="{{.Title}}" 
			title="{{.Title}}"
			xmlUrl="{{.FeedLink}}"
			{{ if .Description }}description="{{.Description}}"{{end}}
			{{ if .Link }}htmlUrl="{{.Link}}"{{end}}
			{{ if .Language }}language="{{.Language}}"{{end}}
			{{ if .FeedType }}type="{{.FeedType}}"{{end}}
			{{ if .FeedVersion }}version="{{opmlRSSVersion .}}"{{end}}
		/>
	{{end}}
	</body>
</opml>
