# iota-iri_exporter

This is an implementation of the functionality of iota-prom-exporter in Go that is typically used by Iota IRI nodes for Prometheus/Grafana monitoring.
While creating the exporter in Go I decided to do this as an reimplementation, not as a straight port. The output is compatible with iota-prom-exporter, however, the exporter framwork is different. Interpretaton of the metrics is rebuild to match the iota-prom-exporter.

This exporter will create and use a database (Badger DB) to track transaction confirmation times. Note that this database will only be used if ZMQ messaging is enabled.

I started this project to port the key IRI metrics to an exporter program written in Go due to the following concerns with the existing iota-prom-exporter written in node.js:

1. One single (simple) executable that provides the key IRI status metrics for Prometheus. 
2. Only perform monitoring of the IRI node (no external metrics like BTC, ETH and Iota price).
3. Stability by using the same model as the stock node_exporter as provided by Prometheus.

As I am experimenting with building IRI nodes that use a minimum of resources, it became clear to me that I did not want the additional overhead of running a node.us program. In addition to that, the iota-node-exporter is also providing insight into external info like BTC, ETH and Iota market price which only adds unnecessary resources to the node.

It is my goal of this project to provide the same functionality as iota-prom-exporter with the option to switch off non-IRI exporters to match the iota-prom-exporter exports as needed.

1. iota-iri_exporter (work in progress): Export the main IRI metrics for consumption by Prometheus
  - [x] Pull getNodeInfo metrics.
  - [x] Pull getNeighbor metrics.
  - [x] Collect and Expose IRI ZMQ Metrics.
  - [x] Collect and Expose Transaction Confirmation Rate (duration) Metrics.
  - [x] Bitfinex: Export market prices for popular crypto.

With the option to remove specific metrics groups, federation of node monitoring should become more logical and efficient.  

# Use

Start the iota-iri_exporter program from the commandline.

example: `iota-iri_exporter --web.listen-address=":9311" --web.iri-path="http://myiotanode:14265"`

usage: iota-iri_exporter [\<flags\>]
```
Flags:
  --help                        Show context-sensitive help (also try --help-long and --help-man).
  --web.listen-address=":9311"  Address to listen on for web interface and telemetry.
  --web.telemetry-path="/metrics"  
                                Path under which to expose metrics.
  --web.iri-path="http://localhost:14265"  
                                URI of the IOTA IRI Node to scrape.
  --opt.zmq                     Enable ZMQ based metrics (database needed).
  --web.zmq-path="tcp://localhost:5556"  
                                URI of the IOTA IRI ZMQ Node to scrape.
  --db.database-path="./iotabadgerdb"  
                                Path for the database.
  --version                     Show application version.
  --log.level="info"            Only log messages with the given severity or above. Valid levels: [debug, info, warn,
                                error, fatal]
  --log.format="logger:stderr"  Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or
                                "logger:stdout?json=true"
```


Point your browser at http://localhost:9311/metrics

Node metrics should show.
