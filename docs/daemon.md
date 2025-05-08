# Daemon Mode

!!! info
    A daemon is a background process that runs independently of user interaction, typically used to perform tasks at scheduled intervals or in response to specific events. 

In the context of Tinyfeed, the daemon mode allows the tool to periodically update the generated HTML page with the latest content from the specified feeds. This is particullary usefull to integrate tinyfeed with service manager such as [systemd](/systemd), [OpenRC](/openrc) or even [Docker](docker). For details about those use case, please see the corresponding workflow section of the documentation.

To use Tinyfeed in daemon mode, you can use the `-D / --daemon` flag along with the `-o / --output` flag to specifcy which file to keep updated.

For example:
```bash
tinyfeed --daemon -i feeds.txt -o index.html
```

By default the page will be updated once every day (1440 minutes). You can change with the `-I / --interval` flag to specify the duration (in minutes) between each update. For example:
```bash
tinyfeed --daemon -i feeds.txt -o index.html -I 720
```
When ajusting the update interval please keep in mind that tinyfeed does not yet implement any [cache mecanism](https://github.com/TheBigRoomXXL/tinyfeed/issues/11) and as such any short interval will create a probably [unwanted](https://rachelbythebay.com/w/2024/12/17/packets/) load on the feed provider.

!!! opinion
    While most media feeds bombard us with last minutes new, I concider that a once a day update is actually a fine rythme that encourage a healthier consumption of new. Maybe don't adjust the interval time down if you don't have a specific need. 


