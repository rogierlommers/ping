# pingback
Keep track of all your servers / nodes

## run a server once

```
docker run --name=docker-pingback \
            -p 9005:8080 \
            -e "emailuser=your@gmail" \
            -e "emailpassword=yourpass" \
            rogierlommers/pingback-server
```

This will start the server. Notification emails are sent to the provided (gmail) address.

## run on each client

Create a crontab record which fires the client every x minutes:

```
pingback-linux-amd64 -mode client -server http://point-to-the-server >/dev/null 2>&1
```

## API
Open the server API endpoint `/api/ping` to view all your nodes.
