# OpenRC

!!! Info
    OpenRC is an [init system](https://github.com/OpenRC/openrc/blob/master/service-script-guide.md) designed for Unix-like OS. Each service is defined by a small text file, called a unit file, which containing instructions on how to run the service. It's notably the default init system for Alpine and Gentoo.


To create an OpenRC service, you just need to add an [init file](https://wiki.alpinelinux.org/wiki/Writing_Init_Scripts) at `/etc/init.d/tinyfeed`.
Then, you can enable the service with `rc-update add tinyfeed default` and start or stop it with `rc-service tinyfeed <COMMAND>`. Below, you can find the most minimal service file that will enable you to start and stop tinyfeed with the feed list and rendered page located at `/etc/tinyfeed`.

```ini
#!/sbin/openrc-run

depend() {
	need net
	use dns 
}

command="/usr/local/bin/tinyfeed"
command_args="--daemon -i /etc/tinyfeed/feeds.txt -o /etc/tinyfeed/index.html"
command_background=true
pidfile="/run/${RC_SVCNAME}.pid"
```
For more advanced patterns (like running as your user instead of root), you can check out the official [OpenRC documentation](https://github.com/OpenRC/openrc/blob/master/service-script-guide.md) on the subject or [Gentoo wiki](https://wiki.gentoo.org/wiki/OpenRC/openrc-init) and [Alpine user handbook](https://docs.alpinelinux.org/user-handbook/0.1a/Working/openrc.html)
