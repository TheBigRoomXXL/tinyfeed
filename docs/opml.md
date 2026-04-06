# OPML Export

!!! info
    [OPML](https://opml.org/) is an open standard that has been used since the 2000s to aggregate feeds into collections. As such, OPML is very useful as an interoperable format between feed apps.


tinyfeed can be used to export a feed collection into OPML format instead of an HTML page. This is done using the custom template option:

```bash
tinyfeed --input feeds.txt --output index.opml --template opml
```

This will output a file with the following structure:

```xml
<opml version="2.0">
	<head>
		<title>Feed</title>
		<dateCreated>06 Apr 26 11:34 CEST</dateCreated>
		<docs>"https://opml.org/spec2.opml"</docs>
	</head>
	<body>
	
		<outline
			text="Redowan&#39;s Reflections" 
			title="Redowan&#39;s Reflections"
			xmlUrl="https://rednafi.com/index.xml"
			description="Recent content on Redowan&#39;s Reflections"
			htmlUrl="https://rednafi.com/"
			language="en"
			type="rss"
			version="RSS2"
		/>
	
		<outline
			text="Let&#39;s Discuss the Matter Further" 
			title="Let&#39;s Discuss the Matter Further"
			xmlUrl="https://rhodesmill.org/brandon/feed/"
			description="Thoughts and ideas from Brandon Rhodes"
			htmlUrl="https://rhodesmill.org/brandon/feed/"
			language="en"
			type="rss"
			version="RSS2"
		/>
	
		<outline
			text="Josh Collinsworth" 
			title="Josh Collinsworth"
			xmlUrl="https://joshcollinsworth.com/api/rss.xml"
			description="Josh Collinsworth&#39;s blog"
			htmlUrl="https://joshcollinsworth.com"
			language="en"
			type="rss"
			version="RSS2"
		/>
	
		<outline
			text="Sebastien Lovergne" 
			title="Sebastien Lovergne"
			xmlUrl="https://lovergne.dev/rss.xml"
			description="A humble tinkerer website about his interests"
			htmlUrl="https://lovergne.dev/"
			language="en-US"
			type="rss"
			version="RSS2"
		/>
	</body>
</opml>

```

Using OPML as input format is not yet supported but it is planned: [issue #29](https://github.com/TheBigRoomXXL/tinyfeed/issues/29)
