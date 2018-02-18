# Simple script for copying the built iota-iri_exporter program from the go/bin directory to /opt/prometheus
# This script is designed for Iota IRI installations using the 'Playbook" script'

systemctl stop iota-iri_exporter
cp /home/marcel/go/bin/iota-iri_exporter /opt/prometheus/iota-iri_exporter
chown prometheus:prometheus /opt/prometheus/iota-iri_exporter
systemctl start iota-iri_exporter
/opt/prometheus/iota-iri_exporter --version
