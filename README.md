# standup

Notify a HipChat room with your daily status for standup.

## Setup

```bash
$ git clone git@github.com:danriti/standup.git
$ cd standup/cmd/standup-server
$ docker build -t standup ./
$ export HIPCHAT_TOKEN=<token>
$ export HIPCHAT_ROOM_ID=<room-id>
$ docker run -e HIPCHAT_TOKEN -e HIPCHAT_ROOM_ID -d -P standup:latest
$ docker ps -a
```
