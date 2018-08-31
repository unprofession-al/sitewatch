![WebStalker](./logo.svg)

Watch websites for changes. Like Baywatch, but for websites. 

## If...

... You really consider running this in real live, notice that this documentation
is crap. Raise an issue and write some stuff about you. Do not forget to mention
your use case of WebStalker. We'll get things going for you...

## What it does

Ever refreshed a website ten times while waiting for some expected update in
order to figure out that nothing happend yet?

Let WebStalker waist its time to do so: Give it a couple of websites and WebStalker
checks if their content change (via md5 sum of its content). You'll only get 
notified if changes occure...

## Install

Since currently no binaries are provided you need to compile WebStalker by hand.
Go makes this easy: 

1) Install Go (https://golang.org/doc/)
2) Compile: `go get -u github.com/unprofession-al/webstalker`

That's it. 

## Configure

Create a config.yaml file and provide a list of websites you want to have stalked:

```
---
# check interval in seconds
interval: 300
debug: false
# overwrite this config file to store the hash of each site
store_hash: true
sites:
  Random String:
    url: https://www.random.org/strings/?num=10&len=10&digits=on&unique=on&format=html&rnd=new
    recipient: onlyreply@example.com
    template: There is another random string available
  Google Banner:
    url: https://www.google.com
    recipient: onlyreply@example.com
    # css selector of the element to watch
    # if not defined the whole document is consedered
    # use github.com/suntong/cascadia to test your selector
    selector: .hplogo
    template: Seems that google has a new banner
```

## Run

Get a help output via option `-h`:

```
webstalker -h
Usage of webstalker:
  -config string
    	path to the configuration file (default "config.yaml")
  -single
    	run only once (to be used when controlled via cron or simiar)
```


Run this is the directory where your config lives:

```
WEBSTALKER_NOTIFIER_SENDGRID="noreply@stalkingbastard.com SG.yG2dlva4R4KO8-ThisIsMySendGridKey" WEBSTALKER_NOTIFIER_STDOUT="YES" webstalker -config /path/to/config/file.yaml
```

Alternatively docker-compose to run `webstalker`:

```
version: '2'

services:
  webstalker:
    image: unprfssnl/webstalker
    restart: always
    volumes:
      - ./config.yaml:/config.yaml
      - /etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt
    container_name: webstalker
    command: webstalker -c /config.yaml
    environment:
    - "WEBSTALKER_NOTIFIER_SENDGRID=noreply@stalkingbastard.com SG.yG2dlva4R4KO8-ThisIsMySendGridKey"
    - "WEBSTALKER_NOTIFIER_STDOUT=true"

```

... and then...

```
docker-compose up
```
