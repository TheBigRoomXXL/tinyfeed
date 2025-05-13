# Styling

tinyfeed allows you to customize the style by linking to an external CSS stylesheet using
the `-s` or `--stylesheet` flag:

Using a direct link:

```bash
tinyfeed https://lovergne.dev/rss -s "https://feed.lovergne.dev/css/solarized.css" -o index.html
# add <link rel="stylesheet" type="text/css" href="http://localhost:8000/styles/solarized.css" , nonce="..."> to the webpage

```

Or a relative link:
```bash
tinyfeed https://lovergne.dev/rss -s "/solarized.css" -o index.html
# add <link rel="stylesheet" type="text/css" href="/solarized.css" , nonce="..."> to the webpage
```

!!! Info
    CSS is a declarative language used to describe the style of an HTML document. If you are not familiar with it, you can take a look at ["Write your first lines of CSS!" by Scrimba](https://scrimba.com/the-frontend-developer-career-path-c0j/~015?via=mdn) for an introduction and [CSSBed](https://www.cssbed.com/) for stylesheet examples.

## Color

The simplest form of styling is overriding the color to create an alternative theme. This
is the default colors used by tinyfeed:
```css
/* This define dark mode colors */
:root {
    --primary: #0AC8F5;     /* Unvisited links */
    --secondary: #D2F00F;   /* Visited links */
    --text: , #cfcfcf;      /* Title, description and metadata*/
    --border:  #c0c0c0;     /* Border for header, footer and feeds table */
    --background:  #1D1F21; /* Background color for the entire page */
}


/* This define light mode colors */
@media (prefers-color-scheme: light) {
    :root {
        --primary: #207D92;
        --secondary: #6A7A06;
        --text: #444444;
        --border: #333333;
        --background: #ffffff;
    }
}
```

You can then copy this definition in a stylesheet and tweak the color to your liking.
For example this is the stylesheet use to convert tinyfeed to a solarized theme:
```css
:root {
    --primary: #1696D8;
    --secondary: #d7b151;
    --text: #cfcfcf;
    --border: #c0c0c0;
    --background: #002C37;
}

@media (prefers-color-scheme: light) {
    :root {
        --text: #5c7075;
        --border: #EFE8D3;
        --background: #FFF6E1;
    }
}
```

**See solarized demo : [feed.lovergne.dev/demo/solarized](/demo/solarized.html)**
## Formatting

More generaly adding a stylesheet gives you access to the full power of CSS. For example adding 
```css
/* skipping color section... */

li {
    font-size: bold;
    margin-block-start: 0.5em;
}
``` 
would give you the compact layout of [lobste.rs](https://lobste.rs/). **See the lobste.rs skin demo:  [feed.lovergne.dev/demo/lobster](/demo/lobster.html)**

Or something more involved like:
```css
/* skipping color section... */

body {
    max-width: 796px;
}

h1 {
    color: #000000;
    font-size: 10pt;
    margin: 0.25em;
}

h2 {
    font-size: 10pt;
    font-weight: normal;
    font-family: Verdana, Geneva, sans-serif;
    display: inline;
}

header {
    background-color: #ff6600;
    border-bottom: none;
    margin-bottom: 0px;
}

main {
    background-color: #f6f6ef;
    padding-bottom: 1rem;
}

footer {
    border-top: 2px solid #ff6600;
    margin-top: 0px;
    background-color: #f6f6ef;
}

ol {
    margin-block: 0em;
    padding-top: 0.5em;
}

li {
    margin-block-start: 0.25em;
}

small {
    font-size: 8pt;
}
```
would reproduce the style and layout of [Hacker News](https://news.ycombinator.com/). **See the hacker news skin demo: [feed.lovergne.dev/demo/hackernews](/demo/hackernews.html)**

