# Welcome 

![banner](img/banner.svg)

**tinyfeed** is a CLI tool that generate a static HTML page from a collection of feeds.

It's dead simple, no database, no config file, just a CLI and some HTML 

Give it a list of RSS, Atom or JSON feeds urls and it will generate a single HTML page for it. Then you can effortlessly set it up in `crond,` `systemd` or `openrc` and voilà, you’ve got yourself a webpage that aggregates your favorite feeds.

## Feature

- RSS, Atom and JSON feeds are all supported thanks to the awesome 
[gofeed library](https://github.com/mmcdole/gofeed)
- Highly customizable, especially with the ability to use external stylesheets and templates.
- Dark / Light theme based on system preference
- The generated page is lightweight and fully accessible.
- Supports a daemon mode to re-generate the output periodically.


## Live demos

**Default tinyfeed: [feed.lovergne.dev/demo](/demo)**  
**With solarize theme : [feed.lovergne.dev/demo/solarized](/demo/solarized.html)**  
**With lobste.rs skin : [feed.lovergne.dev/demo/lobster](/demo/lobster.html)**  
**With hacker new skin : [feed.lovergne.dev/demo/hackernew](/demo/hackernews.html)**  


## Screenshots

![screenshots of feed.lovergne.dev](img/screenshots.png)

Visited links are in yellow, unvisited in blue. 


## Feedback, help or bug report

If you need anything related to this project, whether it's just giving feedback, requesting a feature, or simply asking for help to understand something, open an [issue]((https://github.com/TheBigRoomXXL/tinyfeed/issues)) on the official [repository](https://github.com/TheBigRoomXXL/tinyfeed/).

You have created a page with tinyfeed and you want to share it? You can open a merge request or an issue to add it to the demo section.

*If you want to contribute something other than a demo, please open an issue first so that we can collaborate efficiently.*


## Acknowledgement

Thank you [@MariaLetta](https://github.com/MariaLetta) for the awesome [free-gophers-pack ](https://github.com/MariaLetta/free-gophers-pack) which I adapted for the banner.

The project was heavily inspired by the awesomely simple [tinystatus](https://github.com/bderenzo/tinystatus) and message boards like Lobste.rs and Hacker News.
