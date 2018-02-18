# Installation Instructions

Note that these have not been verified yet. 
Some of the steps may be incorrect or incomplete


Prerequisites for building:
- Golang (Go 1.6 and 1.8 tested)
- libzmq3-dev (apt-get install libzmq3-dev)

Go prerequisites
- go get github.com/iotaledger/giota
- go get github.com/prometheus
- go get github.com/pebbe/zmq4

Get the iota-iri_exporter sources:
- go get github.com/maeck70/iota-iri_exporter


Steps:
1. `go install github.com/maeck70/iota-iri_exporter`
2. >> binary should be built into go/bin
3. Make sure iota-prom-exporter is not running `systemctl stop iota-prom-exporter`
4. Check the starter file for correct addresses and ports and start manually with `./run-exporter`
5. Since this exposes the iota-iri prometheus info on the same port as the iota-prom-exporter, the newmetrics should be visible in Grafana
6. Daemonize this program. There are more steps to do so which I wont go into. The script `copy-to-prometheus.sh` can be used to copt the binary to the /opt/prometheus folder and start the daemon.
 