# pingback
Keep track of all your servers / nodes

# project status
[![Go Report Card](https://goreportcard.com/badge/github.com/rogierlommers/pingback)](https://goreportcard.com/report/github.com/rogierlommers/pingback) [![CircleCI](https://circleci.com/gh/rogierlommers/pingback/tree/master.svg?style=svg)](https://circleci.com/gh/rogierlommers/pingback/tree/master)

## run a server once

```
docker run --name=docker-pingback \
            -p 9005:8080 \
            -e "emailuser=youremailusername" \
            -e "emailpassword=youremailpass" \
            -e "emailsmtp=yoursmtp.server.com" \
            rogierlommers/pingback-server
```

This will start the server. Notification emails are sent to the provided (smtp server) address.

## run on each client

Create a crontab record which fires the client every x minutes:

```
pingback-linux-amd64 -mode client -server http://point-to-the-server >/dev/null 2>&1
```

## API
Open the server API endpoint `/api/ping` to view all your nodes.
