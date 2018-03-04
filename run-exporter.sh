# Simple script starting the iota-iri_exporter program from the command line
# Modify as needed for your install

~/go/bin/iota-iri_exporter --log.level="debug" --web.listen-address=":9311" --web.iri-path="http://localhost:14265" --web.zmq-path="tcp://localhost:5556" --no-zmq
