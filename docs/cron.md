# Cron

!!! info
    
    cron is the time-based job scheduler in Unix-like computer operating systems. cron enables users to schedule jobs (commands or shell scripts) to run periodically at certain times, dates or intervals. It is commonly used to automate system maintenance or administration.  
    *- from [Wikipedia](https://en.wikipedia.org/wiki/Cron)*

## Crontab

To integrate tinyfeed with cron start with opening your crontab with `crontab -e`, then simply add the following line

```cron
0 6 * * * tinyfeed -i /path/to/feeds.txt -o /path/to/index.html
```

tinyfeed  will run daily at 6 AM to update your `index.html`.

Since `crond` executes tinyfeed with your user permissions, ensure you have read and write access to both `feeds.txt` and `index.html`. If you need to run it as a different user, use the `-u` flag with `crontab`, which requires root access.

If you want to change the frequency you need to update `0 6 * * *` with a new cron definition. You might use [crontab.guru](https://crontab.guru/) if you are unfamilliar with this format.

## Alterting

If your server is set up with a mail service (such as [Postfix](https://www.postfix.org/)), you can add `MAILTO="your.email@domain.com"` before the tinyfeed line. This setup ensures you receive an email notification whenever tinyfeed encounters an error.

By default, tinyfeed outputs warnings to `stderr`, which can also trigger an email. To prevent this, you can use the `-q` or `--quiet` flag.
