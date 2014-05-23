# Startup Speeddating Ticker

This is the Ticker used for the EC Startup Speeddating.
It is written in [Go](http://golang.org/) and uses Wecksockets for real-Time communication.

After starting the server, everyone accessing the root URL will see the same synchronized ticker. Type `/admin` in the address bar to get to the admin interface after entering the password.

This way you can let the timer run from on a Laptop connected to a projector, and remote-control the timer with your smartphone or from another computer.

Build:

    go build
    # and run
    ./ssd-ticker -p your_secret

Run using `foreman` (env vars can be edited in `.env`):

    go build
    foreman start

Just run using `go run`:

    make run
    # or
    go run *.go -p your_secret

Build and install in gopath:

    git clone git@github.com:jo-m/ssd-timer.git $GOPATH/src/github.com/jo-m/ssd-timer
    cd $GOPATH/src/github.com/jo-m/ssd-timer
    go get

## Deploy to heroku
    heroku apps:create ssd-timer; \
    heroku labs:enable --app ssd-timer websockets; \
    heroku config:add --app ssd-timer BUILDPACK_URL=https://github.com/kr/heroku-buildpack-go.git; \
    heroku config:add --app ssd-timer ADMIN_PASSWORD=<password>; \
    git push heroku master

## Check logs
    heroku logs --app ssd-timer --tail
