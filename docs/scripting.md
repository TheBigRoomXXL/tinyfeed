# Scripting

tinyfeed allows you to inject javascript by linking to an javascript file with the `-S` or `--script` flag:

Using a direct link:

```bash
tinyfeed -S "https://feed.lovergne.dev/scripts/top.js" -o index.html https://lovergne.dev/rss
# add <script type="module" src="https://feed.lovergne.dev/scripts/top.js" nonce="..."></script> to the webpage

```

Or a relative link:
```bash
tinyfeed -S "/scripts/top.js" -o index.html https://lovergne.dev/rss
# add 	<script type="module" src="/scripts/top.js" nonce="..."></script> to the webpage
```

Adding a script is ideal if you want to extend tinyfeed with with more interactive features. 


## Example: Back to Top button

Adding a button with a callback is the kind of lightweight customization you can easily achieve with JavaScript.

Add the following `top.js` file to the directory where you output your `index.html`:

```js
function onScroll() {
    if (document.body.scrollTop > 20 || document.documentElement.scrollTop > 20) {
        button.style.display = "block";
    } else {
        button.style.display = "none";
    }
}

function GoToTop() {
    document.body.scrollTop = 0; // For Safari
    document.documentElement.scrollTop = 0; // For Chrome, Firefox, IE and Opera
}

// Create and stylize the button
const button = document.createElement('button');
button.title = "Go to the top of the page";
button.onclick = GoToTop;
button.style.display = 'none';
button.style.position = 'fixed';
button.style.bottom = '30px';
button.style.right = '30px';
button.style.zIndex = '99';
button.style.cursor = 'pointer';
button.style.width = '3em';
button.style.height = '3em';
button.style.backgroundColor = 'current';
button.style.mask = 'size: 100 % 100 %';
button.style.maskRepeat = 'no';
button.style.maskImage = `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23000' d='M22 12a10 10 0 0 1-10 10A10 10 0 0 1 2 12A10 10 0 0 1 12 2a10 10 0 0 1 10 10M7.4 15.4l4.6-4.6l4.6 4.6L18 14l-6-6l-6 6z'/%3E%3C/svg%3E")`;
button.style.webkitMask = 'size: 100 % 100 %';
button.style.webkitMaskRepeat = 'no';
button.style.webkitMaskImage = `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23000' d='M22 12a10 10 0 0 1-10 10A10 10 0 0 1 2 12A10 10 0 0 1 12 2a10 10 0 0 1 10 10M7.4 15.4l4.6-4.6l4.6 4.6L18 14l-6-6l-6 6z'/%3E%3C/svg%3E")`;

// Add it to the body
document.body.appendChild(button);
window.onscroll = onScroll;
```

Then generate your `index.html` with `tinyfeed -S "top.js" -o index.html https://lovergne.dev/rss`. The resulting page should have "Back To Top" button appear on the botton right when you start scrolling.

**Check out the demo with that exact script:  [feed.lovergne.dev/demo/back-to-top](/demo/back-to-top.html)**
