#!/usr/bin/env bash

set -eu

CONNECTION_STRING=${1}
HOSTNAME=$(hostname | tr -d '\n')
AUDIT_LOG_PATH=$(sudo find /var -name audit.log | grep 'merged' | tr -d '\n')

# Install JRE and logstash (this sets up some important config).
apt install -y openjdk-11-jre-headless
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add -
echo "deb https://artifacts.elastic.co/packages/7.x/apt stable main" | tee -a /etc/apt/sources.list.d/elastic-7.x.list
apt-get update && apt-get install -y logstash

# Need to install specific version of logstash to avoid a JRuby bug.
mkdir -p /tmp/downloads
pushd /tmp/downloads
wget https://artifacts.elastic.co/downloads/logstash/logstash-7.3.2.tar.gz
tar xvf logstash-7.3.2.tar.gz
rm -rf /usr/share/logstash
mv logstash-7.3.2 /usr/share/logstash
popd

# Install mongodb output plugin.
/usr/share/logstash/bin/logstash-plugin install logstash-output-mongodb

# Add logstash conf.
cat << EOF > /etc/logstash/conf.d/logstash.conf
input {
  file {
    path => "${AUDIT_LOG_PATH}"
    mode => "tail"
    start_position => "beginning"
  }
}

filter {
  json {
    source => "message"
    remove_field => [ "annotations" ]
  }
}

output {
  mongodb {
    id => "${HOSTNAME}"
    uri => "${CONNECTION_STRING}"
    database => "kubelet"
    collection => "audit-logs"
  }
}
EOF

# Add systemd service.
cat << EOF > /etc/systemd/system/logstash.service
[Unit]
Description=logstash

[Service]
Type=simple
User=root
Group=root
# Load env vars from /etc/default/ and /etc/sysconfig/ if they exist.
# Prefixing the path with '-' makes it try to load, but if the file doesn't
# exist, it continues onward.
EnvironmentFile=-/etc/default/logstash
EnvironmentFile=-/etc/sysconfig/logstash
ExecStart=/usr/share/logstash/bin/logstash "--path.settings" "/etc/logstash"
Restart=always
WorkingDirectory=/
Nice=19
LimitNOFILE=16384

[Install]
WantedBy=multi-user.target
EOF

# Start logstash.
systemctl daemon-reload
systemctl restart logstash
