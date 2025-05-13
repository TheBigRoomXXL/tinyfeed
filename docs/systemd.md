# Systemd

!!! info 
    Systemd is a service manager for Linux that starts and manages services automatically. Each service is defined by a small text file, called a unit file, which containing instructions on how to run the service.

Here is a simple [systemd unit file](https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html) to run tinyfeed when your system start and update the page every 24 hours. With this setup you can edit the feeds list at `~/feeds.txt` and the output will be updated after the next run. 

```ini
# /etc/systemd/system/tinyfeed.service

[Unit]
Description=tinyfeed service
After=network.target

[Service]
Type=simple
Restart=always
User=<USER>
WorkingDirectory=/home/<USER>/
ExecStart=/usr/local/bin/tinyfeed --daemon -i feeds.txt -o index.html

[Install]
WantedBy=mutli-user.target
```

Whith this config tinyfeed will run with your permission and output `index.html` in your home directory. You might want instead to setup a dedicated user and point a webserver to its home directory for better isolation.

If you have SELinux enabled (in fedora for example) you will need to allow `systemd` to execute binaries in the `usr/local/bin` directory with the following commands:
```bash
sudo semanage fcontext -a -t bin_t /usr/local/bin 
sudo chcon -Rv -u system_u -t bin_t /usr/local/bin 
sudo restorecon -R -v /usr/local/bin
```

