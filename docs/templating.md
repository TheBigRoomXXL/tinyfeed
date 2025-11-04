# Templating

Internally, tinyfeed uses a Golang HTML template to generate its webpage. You can provide your own template for page generation using the `-t` / `--template` flag:
```bash
tinyfeed -i feeds.txt -o index.html -t my-template.html
```
For an example template, check out the [built-in one](https://github.com/TheBigRoomXXL/tinyfeed/blob/main/internal/built-in). To learn about HTML+Go templates, check the [official documentation](https://pkg.go.dev/html/template).

## Securing against Cross-Site-Scripting  🛡️

### TL;DR

When doing your templating it is heavily recommended to use the following [Content Security Policy (CSP)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CSP) to avoid [Cross-Site-Scripting (or XSS)](https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/XSS) attack:

```html
<meta http-equiv="Content-Security-Policy" content="
    default-src 'none';
    base-uri 'self';
    img-src 'self' data: ;
    style-src 'nonce-{{.Metadata.nonce}}' ;
    script-src 'nonce-{{.Metadata.nonce}}' 'strict-dynamic' ;   
">
```

### Understanding the Cross-Site-Scripting risk and how to avoid it

The content of RSS, Atom, and JSON feeds **MUST** be treated as uncontrolled user input. These feeds may include various types of unusual and potentially harmful payloads. For example, if you have the following section in your template:
```html
<ol>
    {{range .Items}}
    <li>
        <h2>{{.Title}}</h2>
        <p>{{.Description}}</p>
    </li>
    {{end}}
</ol>
```

Then if a feed has the following description:
```html
<meta http-equiv="refresh" content="1;URL=https://example.com">
```

It will inject a `<meta>` tag in your HTML which will be interpreted as a redirect header, and the result will be that any user who opens your webpage will be redirected to `example.com`. This is what we call a [Cross-Site-Scripting (or XSS)](https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/XSS) attack. This can not only break your page but also **put your users at risk**.

!!! info
    This example is actually a real case encountered during the development of tinyfeed. The twist is that it was not even designed to be harmful. My guess is that it was designed to redirect readers opening an article in a feed reader to the canonical link of the article. However, this show that even if you trust the authors of the feeds you choose to not be voluntarily harmful, they can still break your webpage if you don't manage their inputs carefully.

While escaping the provided inputs helps secure the content (and tinyfeed does so with the [Go HTML package](https://pkg.go.dev/html)),  it is insufficient. Cross-Site-Scripting attacks can take a [tremendous number of forms](https://github.com/swisskyrepo/PayloadsAllTheThings/tree/master/XSS%20Injection) and often uses complex evasion techniques to bypass filters and escaping mechanisms. To really be secure from injection, we need to use an awesome feature provided by the browser: [Content Security Policy (CSP)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CSP). CSP enables us to specify which resources (such as scripts, images, and stylesheets) are permitted to load on our webpage. While it is a powerfull tool, it can be somewhat complex to configure. Therefore, I suggest reviewing the [MDN Guide](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CSP) to better understand its usage. If you are not comfortable with the subject, you can start by using the CSP tinyfeed already uses and adapt it if needed:
```html
<meta http-equiv="Content-Security-Policy" content="
    default-src 'none';
    base-uri 'self';
    img-src 'self' data: ;
    style-src 'nonce-{{.Metadata.nonce}}' ;
    script-src 'nonce-{{.Metadata.nonce}}' 'strict-dynamic' ;
    ">
```

Here is a breakdown of what it does:
- `default-src 'none';`: Disables the loading of all resources, including any JavaScript.
- `base-uri 'self';`: stop a [<base> element](https://developer.mozilla.org/fr/docs/Web/HTML/Reference/Elements/base) from changing the base URI of relative link.
- `img-src 'self' data: ;` Enable loading images from the same domain or from [data urls](https://developer.mozilla.org/en-US/docs/Web/URI/Reference/Schemes/data)
- `style-src 'nonce-{{.Metadata.nonce}}';`: Enable loading stylesheets if they contain a given nonce. This nonce is generated randomly and securely at each execution of tinyfeed.
- `script-src 'nonce-{{.Metadata.nonce}}' 'strict-dynamic' ;`: Enable loading script if they contain a given nonce (same as stylesheets). `strict-dynamic` enable trusted script to load there own script, this make scripts with dependencies easy to use.

After tweaking the CSP for your needs, I recommend auditing it with Google's [CSP Evaluator](https://csp-evaluator.withgoogle.com/).

## Available data

Inside your template, you will have access to tinyfeed's data through the following struct:

```go
type data struct {
    Metadata    map[string]string
    Items       []Item // it's *gofeed.Item with a FeedName field added.
    Feeds       []*gofeed.Feed
    Stylesheets []string
    Scripts     []string
}
```

Metadata is an arbitrary map, there is currently no backward compatibility promise on it's content (but we keep it stable as much as possible).

Metadata currently holds the following keys: `name`, `description`, `nonce`, `day`, `datetime`. For the most recent list of keys, refer to `printHTML` in [main.go](https://github.com/TheBigRoomXXL/tinyfeed/blob/main/internal/main.go). 

You will also have access to two functions for formatting:
```go
// Provide the publication date of an Item with the format "2006-01-02"
func publication(item *gofeed.Item) string

// Return the hostname of an Item with "www" trimmed
func domain(item *gofeed.Item) string
```
