# iota-iri_exporter

This is an implementation of the iota-prom-exporter in Go that is typically used by Iota IRI nodes.
I started this project to port the key IRI metrics to an exporter program written in Go due to the following concerns with the existing iota-prom-exporter written in node.js:

1. One single (simple) executable that provides the input for Prometheus. 
2. Only perform monitoring of the IRI node (no external info like BTC, ETH and Iota price).
3. Stability by using the same model as the stock node_exporter as provided by Prometheus.

As I am experimenting with building IRI nodes that use a minimum of resources, it became clear to me that I did not want the additional overhead of running a node.us program. In addition to that, the iota-node-exporter is also providing insight into external info like BTC, ETH and Iota market price which only adds unnecessary resources to the node.

It is my goal of this project to provide the same functionality as iota-prom-exporter in a different package with the option to run multiple exporters to match the iota-prom-exporter exports.

1. iota-iri_exporter (work in progress): Export the main IRI metrics for consumption by Prometheus
2. iota-tangle_exporter (planned): Export tangle metrics that pertain to the whole tangle, not this particular node 
3. bitfinex_exporter (planned): Export market prices for popular crypto

With the proposed metrics exporting breakdown, federation of node monitoring should become more efficient and logical.  

# Use

Start the iota-iri_exporter program from the commandline.
Point your browser at http:\\localhost:9311\metrics

Node metrics should show.
