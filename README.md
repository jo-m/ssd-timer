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
