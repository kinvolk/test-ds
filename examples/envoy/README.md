Server listens on port 1234 and is serving sha256. Client has a
"default" server set up at port 2345.

Envoy needs to listen on port 2345 and forward that to port 1234.
