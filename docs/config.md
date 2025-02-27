# Config

Metrictank comes with an [example main config file](https://github.com/grafana/metrictank/blob/master/metrictank-sample.ini),
a [storage-schemas.conf file](https://github.com/grafana/metrictank/blob/master/scripts/config/storage-schemas.conf) and
a [storage-aggregation.conf file](https://github.com/grafana/metrictank/blob/master/scripts/config/storage-aggregation.conf)
an [index-rules.conf file](https://github.com/grafana/metrictank/blob/master/scripts/config/index-rules.conf)

The files themselves are well documented, but for your convenience, they are replicated below.  

Config values for the main ini config file can also be set, or overridden via environment variables.
They require the 'MT_' prefix.  Any delimiter is represented as an underscore.
Settings within section names in the config just require you to prefix the section header.

Examples:

```
MT_LOG_LEVEL: info                        # MT_<setting_name>
MT_CASSANDRA_WRITE_CONCURRENCY: 10        # MT_<setting_name>
MT_KAFKA_MDM_IN_DATA_DIR: /your/data/dir  # MT_<section_title>_<setting_name>
```

---


# Metrictank.ini


sample config for metrictank  
the defaults here match the default behavior.  
## misc ##

```
# instance identifier. must be unique. used in clustering messages, for naming queue consumers and emitted metrics.
instance = default
```

## data ##

```
# see https://github.com/grafana/metrictank/blob/master/docs/memory-server.md for more details
# forego persisting of first received (and typically incomplete) chunk
drop-first-chunk = false
# only ingest data for chunks that have a t0 equal or higher to the given timestamp. Specified per org. syntax: orgID:timestamp[,...]
ingest-from =
# max age for a chunk before to be considered stale and to be persisted to Cassandra
chunk-max-stale = 1h
# max age for a metric before to be considered stale and to be purged from in-memory ring buffer.
metric-max-stale = 3h
# Interval to run garbage collection job
gc-interval = 1h
# duration until when secondary nodes are considered to have enough data to be ready and serve requests.
# To prevent gaps in charts when running a cluster of nodes you need to either
# 1) have the nodes backfill data from Kafka (set via "offset" in the "kafka-mdm-in" config) or
# 2) set the warm-up-period to a value long enough to ensure all data received before the node started has been persisted in the store.
# See https://github.com/grafana/metrictank/blob/master/docs/clustering.md#priority-and-ready-state
warm-up-period = 1h
# org Id for publically (any org) accessible data
# leave at 0 to disable.
public-org = 0
```

## Profiling and logging ##

```
# see https://golang.org/pkg/runtime/#SetBlockProfileRate
block-profile-rate = 0
# 0 to disable. 1 for max precision (expensive!) see https://golang.org/pkg/runtime/#pkg-variables")
mem-profile-rate = 524288 # 512*1024
# heap profiletrigger: triggers a heap (memory) profile for diagnosis when usage threshold is breached
# recommended usage: set proftrigger-heap-thresh-rss such that it is much larger than "normal" usage, but lower
# then how much RAM capacity you have, so that a profile can be captured before the process gets killed by the OOM-killer
# inspect status frequency. set to 0 to disable
proftrigger-freq = 10s
# path to store triggered profiles
proftrigger-path = /tmp
# minimum time between triggered profiles
proftrigger-min-diff = 1h
# threshold for process RSS, the amount of RAM memory used. (0 to disable) (see "rss" on dashboard)
proftrigger-heap-thresh = 25000000000
# threshold for bytes allocated on heap (0 to disable) (see "allocated in heap" on dashboard)
# typically, this is not all that useful, "rss" above is what most people care about (and the heap uses less than rss),
# but this setting can help detect a large heap even if some of the memory is swapped out (and thus not accounted for in rss)
proftrigger-heap-thresh-heap = 0
# only log log-level and lower (read right to left: to the left is lower). panic|fatal|error|warning|info|debug
log-level = info
```

## request tracing via jaeger ##

```
[jaeger]
# Whether the tracer is enabled or not
enabled = false
# A comma separated list of name = value tracer level tags, which get added to all reported spans.
# The value can also refer to an environment variable using the format ${envVarName:default},
# where the :default is optional, and identifies a value to be used if the environment variable cannot be found
add-tags =
# the type of the sampler: const, probabilistic, rateLimiting, or remote
sampler-type = const
# The sampler parameter (number)
# - for "const" sampler, 0 or 1 for always false/true respectively
# - for "probabilistic" sampler, a probability between 0 and 1
# - for "rateLimiting" sampler, the number of spans per second
# - for "remote" sampler, param is the same as for "probabilistic"
#   and indicates the initial sampling rate before the actual one
#   is received from the mothership.
sampler-param = 1
# The HTTP endpoint when using the remote sampler
sampler-manager-addr = "http://jaeger:5778/sampling"
# The maximum number of operations that the sampler will keep track of.
# If an operation is not tracked, a default probabilistic
# sampler will be used rather than the per operation specific sampler.
sampler-max-operations = 0
# How often the remotely controlled sampler will poll jaeger-agent for the appropriate sampling strategy
sampler-refresh-interval = 10s
# The reporter's maximum queue size
reporter-max-queue-size = 0
# The reporter's flush interval
reporter-flush-interval = 10s
# Whether the reporter should also log the spans
reporter-log-spans = false
# HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces
collector-addr =
# Username to send as part of "Basic" authentication to the collector endpoint
collector-user =
# Password to send as part of "Basic" authentication to the collector endpoint
collector-password =
# UDP address of the agent to send spans to. (only used if collector-addr is empty)
agent-addr = localhost:6831
```

## metric data storage in cassandra ##

```
[cassandra]
# see https://github.com/grafana/metrictank/blob/master/docs/cassandra.md for more details
# enable the cassandra backend store plugin -- This setting is ignored and overridden (set to false) in query mode
enabled = true
# comma-separated list of hostnames to connect to
addrs = localhost
# keyspace to use for storing the metric data table
keyspace = metrictank
# desired write consistency (any|one|two|three|quorum|all|local_quorum|each_quorum|local_one
consistency = one
# how to select which hosts to query
# roundrobin                : iterate all hosts, spreading queries evenly.
# hostpool-simple           : basic pool that tracks which hosts are up and which are not.
# hostpool-epsilon-greedy   : prefer best hosts, but regularly try other hosts to stay on top of all hosts.
# tokenaware,roundrobin              : prefer host that has the needed data, fallback to roundrobin.
# tokenaware,hostpool-simple         : prefer host that has the needed data, fallback to hostpool-simple.
# tokenaware,hostpool-epsilon-greedy : prefer host that has the needed data, fallback to hostpool-epsilon-greedy.
host-selection-policy = tokenaware,hostpool-epsilon-greedy
# cassandra timeout
timeout = 1s
# max number of concurrent reads to cassandra
read-concurrency = 20
# max number of concurrent writes to cassandra
write-concurrency = 10
# max number of outstanding reads before reads will be dropped. This is important if you run queries that result in many reads in parallel
read-queue-size = 200000
# write queue size per cassandra worker. should be large engough to hold all at least the total number of series expected, divided by how many workers you have
write-queue-size = 100000
# how many times to retry a query before failing it
retries = 0
# size of compaction window relative to TTL
window-factor = 20
# if a read is older than this, it will be omitted, not executed
omit-read-timeout = 60s
# CQL protocol version. cassandra 3.x needs v3 or 4.
cql-protocol-version = 4
# enable the creation of the mdata keyspace and tables, only one node needs this
create-keyspace = true
# File containing the needed schemas in case database needs initializing
schema-file = /etc/metrictank/schema-store-cassandra.toml
# enable SSL connection to cassandra
ssl = false
# cassandra CA certficate path when using SSL
ca-path = /etc/metrictank/ca.pem
# host (hostname and server cert) verification when using SSL
host-verification = true
# enable cassandra user authentication
auth = false
# username for authentication
username = cassandra
# password for authentication
password = cassandra
# instruct the driver to not attempt to get host info from the system.peers table
disable-initial-host-lookup = false
# interval at which to perform a connection check to cassandra, set to 0 to disable.
connection-check-interval = 5s
# maximum total time to wait before considering a connection to cassandra invalid. This value should be higher than connection-check-interval.
connection-check-timeout = 30s
# Maximum chunkspan size used.
max-chunkspan = 24h
```

## Bigtable backend Store Settings ##

```
[bigtable-store]
# enable the bigtable backend store plugin -- This setting is ignored and overridden (set to false) in query mode
enabled = false
# Name of GCP project the bigtable cluster resides in
gcp-project = default
# Name of bigtable instance
bigtable-instance = default
# Name of bigtable table used for chunks
table-name = metrics
# Max number of chunks, per write thread, allowed to be unwritten to bigtable. Must be larger then write-max-flush-size
write-queue-size = 100000
# Max number of chunks in each batch write to bigtable
write-max-flush-size = 10000
# Number of writer threads to use
write-concurrency = 10
# Number concurrent reads that can be processed
read-concurrency = 20
# Maximum chunkspan size used.
max-chunkspan = 6h
# read timeout
read-timeout = 5s
# write timeout
write-timeout = 5s
# enable the creation of the table and column families
create-cf = true
```

## Retention settings ##

```
[retention]
# path to storage-schemas.conf file
schemas-file = /etc/metrictank/storage-schemas.conf
# path to storage-aggregation.conf file
aggregations-file = /etc/metrictank/storage-aggregation.conf
# enables/disables the enforcement of the future tolerance limitation
enforce-future-tolerance = true
# defines until how far in the future we accept datapoints. defined as a percentage fraction of the raw ttl of the matching retention storage schema
future-tolerance-ratio = 10
```

## instrumentation stats ##

```
[stats]
# enable sending graphite messages for instrumentation
enabled = true
# stats prefix (will add trailing dot automatically if needed)
# The default matches what the Grafana dashboard expects
# $instance will be replaced with the `instance` setting.
# note, the 3rd word describes the environment you deployed in.
prefix = metrictank.stats.default.$instance
# graphite address
addr = localhost:2003
# interval at which to send statistics
interval = 1
# timeout after which a write is considered not successful
timeout = 10s
# how many messages (holding all measurements from one interval. rule of thumb: a message is ~25kB) to buffer up in case graphite endpoint is unavailable.
# With the default of 20k you will use max about 500MB and bridge 5 hours of downtime when needed
buffer-size = 20000
```

## chunk cache ##

```
[chunk-cache]
# maximum size of chunk cache in bytes. 512 MB = (1024 ^ 2) * 512 = 536870912
# 0 disables cache
max-size = 536870912
```

## http api ##

```
[http]
# tcp address for metrictank to bind to for its HTTP interface
listen = :6060
# use gzip compression
gzip = true
# use HTTPS
ssl = false
# SSL certificate file
cert-file = /etc/ssl/certs/ssl-cert-snakeoil.pem
# SSL key file
key-file = /etc/ssl/private/ssl-cert-snakeoil.key
# lower resolution rollups will be used to try and keep requests below this number of datapoints. (0 disables limit)
max-points-per-req-soft = 1000000
# limit of number of datapoints a request can return. Requests that exceed this limit will be rejected. (0 disables limit)
max-points-per-req-hard = 20000000
# limit of number of series a request can operate on. Requests that exceed this limit will be rejected. (0 disables limit)
# note here we look at all lowlevel series (even if they will be merged or are equivalent), and can't accurately account for
# requests with duplicate or overlapping targets. See PR #1926 and #1929 for details
max-series-per-req = 250000
# require x-org-id authentication to auth as a specific org. otherwise orgId 1 is assumed
multi-tenant = true
# in case our /render endpoint does not support the requested processing, proxy the request to this graphite
fallback-graphite-addr = http://localhost:8080
# proxy to graphite when metrictank considers the request bad
proxy-bad-requests = true
# timezone for interpreting from/until values when needed, specified using [zoneinfo name](https://en.wikipedia.org/wiki/Tz_database#Names_of_time_zones) e.g. 'America/New_York', 'UTC' or 'local' to use local server timezone.
time-zone = local
# maximum number of concurrent threads for fetching data on the local node. Each thread handles a single series.
get-targets-concurrency = 20
# default limit for tagdb query results, can be overridden with query parameter "limit"
tagdb-default-limit = 100
# ratio of peer responses after which speculative querying (aka spec-exec) is used. Set to 1 to disable.
speculation-threshold = 1
# enable pre-normalization optimization
pre-normalization = true
# enable MaxDataPoints optimization (experimental)
mdp-optimization = false
# output query headers in logs
log-headers = false
```

## metric data inputs ##

```
[input]
# reject received metrics that have invalid input data (invalid utf8 or invalid tags)
reject-invalid-input = true
```

### carbon input (optional)

```
[carbon-in]
# This setting is ignored and overridden (set to false) in query mode
enabled = false
# tcp address
addr = :2003
# represents the "partition" of your data if you decide to partition your data.
partition = 0
```

### kafka-mdm input (optional, recommended)

```
[kafka-mdm-in]
# This setting is ignored and overridden (set to false) in query mode
enabled = false
# For incoming MetricPoint messages without org-id, assume this org id
org-id = 0
# tcp address (may be given multiple times as a comma-separated list)
brokers = kafka:9092
# Kafka version in semver format. All brokers must be this version or newer.
kafka-version = 2.0.0
# kafka topic (may be given multiple times as a comma-separated list)
topics = mdm
# offset to start consuming from. Can be oldest, newest or a time duration
# When using a duration but the offset request fails (e.g. Kafka doesn't have data so far back), metrictank falls back to `oldest`.
# the further back in time you go, the more old data you can load into metrictank, but the longer it takes to catch up to realtime data
offset = newest
# kafka partitions to consume. use '*' or a comma separated list of id's
partitions = *
# The number of metrics to buffer in internal and external channels
channel-buffer-size = 1000
# The minimum number of message bytes to fetch in a request
consumer-fetch-min = 1
# The default number of message bytes to fetch in a request
consumer-fetch-default = 32768
# The maximum amount of time the broker will wait for Consumer.Fetch.Min bytes to become available before it
consumer-max-wait-time = 1s
#The maximum amount of time the consumer expects a message takes to process
consumer-max-processing-time = 1s
# How many outstanding requests a connection is allowed to have before sending on it blocks
net-max-open-requests = 100
# Whether to enable TLS
tls-enabled = false
# Whether to skip TLS server cert verification
tls-skip-verify = false
# Client cert for client authentication (use with -tls-enabled and -tls-client-key)
tls-client-cert =
# Client key for client authentication (use with -tls-enabled and -tls-client-cert)
tls-client-key =
# Whether to enable SASL
sasl-enabled = false
# The SASL mechanism configuration (possible values: SCRAM-SHA-256, SCRAM-SHA-512)
sasl-mechanism =
# Username for client authentication (use with -sasl-enabled and -sasl-password)
sasl-username =
# Password for client authentication (use with -sasl-enabled and -sasl-user)
sasl-password =
```

## basic clustering settings ##

```
[cluster]
# Unique name of the cluster.  This node will only be able to join clusters with the same name.
name = metrictank
# The primary node writes data to cassandra. There should only be 1 primary node per shardGroup.
primary-node = true
# maximum priority before a node should be considered not-ready.
max-priority = 10
# TCP addresses of other nodes, comma separated. use this if you shard your data and want to query other nodes.
# If no port is specified, it is assumed the other nodes are using the same port this node is listening on.
peers =
# Operating mode of this node within the cluster. (dev|shard|query)
# * dev: gossip disabled. node is not aware of other nodes but can serve up all data it is aware of (from memory or from the store)
# * shard: gossip enabled. node receives data and participates in fan-in/fan-out if it receives queries but owns only a part of the data set and spec-exec if enabled.
# * query: gossip enabled. node receives no data and fans out queries to shard nodes (e.g. if you rather not query shard nodes directly)
mode = dev
# minimum number of shards that must be available for a query to be handled.
min-available-shards = 0
# How long to wait before aborting http requests to cluster peers and returning a http 503 service unavailable
http-timeout = 60s
# GOGC value to use when node is not ready.  Defaults to GOGC
# you can use this to set a more aggressive, latency-inducing GC behavior when the node is initializing and hungry for extra memory
# gc-percent-not-ready = 100
# duration until when the cluster topology can be considered up-to-date and this node to be ready to serve requests (when gossip enabled)
gossip-settle-period = 10s
```

## SWIM/gossip clustering settings ##

```
# for more details, see https://godoc.org/github.com/hashicorp/memberlist#Config
# all values correspond literally to the memberlist.Config options except where noted
[swim]
# config setting to use. If set to anything but manual, will override all other swim settings.
# Use manual|default-lan|default-local|default-wan. Note all our swim settings correspond to default-lan
# see:
# * https://godoc.org/github.com/hashicorp/memberlist#DefaultLANConfig
# * https://godoc.org/github.com/hashicorp/memberlist#DefaultLocalConfig
# * https://godoc.org/github.com/hashicorp/memberlist#DefaultWANConfig
use-config = manual
# binding TCP Address for UDP and TCP gossip (full ip/dns:port combo unlike memberlist.Config)
bind-addr = 0.0.0.0:7946
# advertised TCP address for UDP and TCP gossip (full ip/dns:port combo, or empty to use bind-addr)
# Useful for traversing NAT such as from inside docker
advertise-addr =
# timeout for establishing a stream connection with peers for a full state sync, and for stream reads and writes
tcp-timeout = 10s
# number of nodes that will be asked to perform an indirect probe of a node in the case a direct probe fails
indirect-checks = 3
# multiplier for number of retransmissions for gossip messages. Retransmits = RetransmitMult * log(N+1)
retransmit-mult = 4
# multiplier for determining when inaccessible/suspect node is delared dead. SuspicionTimeout = SuspicionMult * log(N+1) * ProbeInterval
suspicion-multi = 4
# multiplier for upper bound on detection time.  SuspicionMaxTimeout = SuspicionMaxTimeoutMult * SuspicionTimeout
suspicion-max-timeout-mult = 6
# interval between complete state syncs. 0 will disable state push/pull syncs
push-pull-interval = 30s
# interval between random node probes
probe-interval = 1s
# timeout to wait for an ack from a probed node before assuming it is unhealthy. This should be set to 99-percentile of network RTT
probe-timeout = 500ms
# turn off the fallback TCP pings that are attempted if the direct UDP ping fails
disable-tcp-pings = false
# will increase the probe interval if the node becomes aware that it might be degraded and not meeting the soft real time requirements to reliably probe other nodes.
awareness-max-multiplier = 8
# number of random nodes to send gossip messages to per GossipInterval
gossip-nodes = 3
# interval between sending messages that need to be gossiped that haven't been able to piggyback on probing messages. 0 disables non-piggyback gossip
gossip-interval = 200ms
# interval after which a node has died that we will still try to gossip to it. This gives it a chance to refute
gossip-to-the-dead-time = 30s
# message compression
enable-compression = true
# system's DNS config file. Override allows for easier testing
dns-config-path = /etc/resolv.conf
```

## clustering transports for tracking chunk saves between replicated node ##
### kafka as transport for clustering messages (recommended)

```
[kafka-cluster]
enabled = false
# tcp address (may be given multiple times as a comma-separated list)
brokers = kafka:9092
# Kafka version in semver format. All brokers must be this version or newer.
kafka-version = 2.0.0
# kafka topic (only one)
topic = metricpersist
# kafka partitions to consume. use '*' or a comma separated list of id's. Should match kafka-mdm-in's partitions.
partitions = *
# offset to start consuming from. Can be oldest, newest or a time duration
# When using a duration but the offset request fails (e.g. Kafka doesn't have data so far back), metrictank falls back to `oldest`.
# Should match your kafka-mdm-in setting
offset = newest
# Maximum time backlog processing can block during metrictank startup. Setting to a low value may result in data loss
backlog-process-timeout = 60s
# Whether to enable TLS
tls-enabled = false
# Whether to skip TLS server cert verification
tls-skip-verify = false
# Client cert for client authentication (use with -tls-enabled and -tls-client-key)
tls-client-cert =
# Client key for client authentication (use with -tls-enabled and -tls-client-cert)
tls-client-key =
# Whether to enable SASL
sasl-enabled = false
# The SASL mechanism configuration (possible values: SCRAM-SHA-256, SCRAM-SHA-512)
sasl-mechanism =
# Username for client authentication (use with -sasl-enabled and -sasl-password)
sasl-username =
# Password for client authentication (use with -sasl-enabled and -sasl-user)
sasl-password =
```

## metric metadata index ##
### in memory, cassandra-backed

```
[cassandra-idx]
# This setting is ignored and overridden (set to false) in query mode
enabled = true
# Cassandra keyspace to store metricDefinitions in.
keyspace = metrictank
# Cassandra table to store metricDefinitions in.
table = metric_idx
# Cassandra table to archive metricDefinitions in.
archive-table = metric_idx_archive
# comma separated list of cassandra addresses in host:port form
hosts = localhost:9042
#cql protocol version to use
protocol-version = 4
# write consistency (any|one|two|three|quorum|all|local_quorum|each_quorum|local_one
consistency = one
# cassandra request timeout. valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'
timeout = 1s
# number of concurrent connections to cassandra
num-conns = 10
# Max number of metricDefs allowed to be unwritten to cassandra
write-queue-size = 100000
#Interval at which the index should be checked for stale series. valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'
prune-interval = 3h
# Number of partitions to load concurrently on startup.
init-load-concurrency = 1
# synchronize index changes to cassandra. not all your nodes need to do this.
update-cassandra-index = true
#frequency at which we should update flush changes to cassandra. only relevant if update-cassandra-index is true. valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'. Setting to '0s' will cause instant updates.
update-interval = 4h
# enable SSL connection to cassandra
ssl = false
# cassandra CA certficate path when using SSL
ca-path = /etc/metrictank/ca.pem
# host (hostname and server cert) verification when using SSL
host-verification = true
# enable cassandra user authentication
auth = false
# username for authentication
username = cassandra
# password for authentication
password = cassandra
# enable the creation of the index keyspace and tables, only one node needs this
create-keyspace = true
# File containing the needed schemas in case database needs initializing
schema-file = /etc/metrictank/schema-idx-cassandra.toml
# instruct the driver to not attempt to get host info from the system.peers table
disable-initial-host-lookup = false
# interval at which to perform a connection check to cassandra, set to 0 to disable.
connection-check-interval = 5s
# maximum total time to wait before considering a connection to cassandra invalid. This value should be higher than connection-check-interval.
connection-check-timeout = 30s
```

### in-memory only

```
[memory-idx]
# This setting is ignored and overridden (set to false) in query mode
enabled = false
# enables/disables querying based on tags
tag-support = true
# enables/disables querying based on meta tags
meta-tag-support = false
# number of workers to spin up to evaluate tag queries
tag-query-workers = 5
# max runtime for a tag query. 0s disables limit (experimental: when hit, result is 200 OK with partial data)
tag-query-timeout = 0s
# size of regular expression cache in tag query evaluation
match-cache-size = 1000
# size of event queue in the meta tag enricher
meta-tag-enricher-queue-size = 100
# size of add metric event buffer in enricher
meta-tag-enricher-buffer-size = 10000
# how long to buffer enricher events before they must be processed
meta-tag-enricher-buffer-time = 5s
# path to index-rules.conf file
rules-file = /etc/metrictank/index-rules.conf
# maximum duration each second a prune job can lock the index.
max-prune-lock-time = 100ms
# use separate indexes per partition. experimental feature. See #1251, #1252
partitioned = false
# number of find expressions to cache (per org). 0 disables cache
find-cache-size = 1000
# size of queue for invalidating findCache entries
find-cache-invalidate-queue-size = 200
# max amount of invalidations to queue up in one batch.
find-cache-invalidate-max-size = 100
# max duration to wait building up a batch to invalidate.
find-cache-invalidate-max-wait = 5s
# amount of time to disable the findCache when the invalidate queue fills up.
find-cache-backoff-time = 60s
# enable buffering new metricDefinitions and writing them to the index in batches
write-queue-enabled = true
# maximum delay between flushing buffered metric writes to the index
write-queue-delay = 1s
# maximum number of metricDefinitions that can be added to the index in a single batch (approx)
write-max-batch-size = 5000
```

### Bigtable index

```
[bigtable-idx]
# This setting is ignored and overridden (set to false) in query mode
enabled = false
# Name of GCP project the bigtable cluster resides in
gcp-project = default
# Name of bigtable instance
bigtable-instance = default
# Name of bigtable table used for MetricDefs
table-name = metric_idx
# Max number of metricDefs allowed to be unwritten to bigtable. Must be larger then write-max-flush-size
write-queue-size = 100000
# Max number of metricDefs in each batch write to bigtable
write-max-flush-size = 10000
# Number of writer threads to use
write-concurrency = 5
# synchronize index changes to bigtable. not all your nodes need to do this.
update-bigtable-index = true
# frequency at which we should update the metricdefinition in bigtable, use 0s for instant updates
update-interval = 3h
# Interval at which the index should be checked for stale series.
prune-interval = 3h
# enable the creation of the table and column families
create-cf = true
```

### in memory, cassandra-backed

```
[cassandra-meta-record-idx]
enabled = true
# Cassandra keyspace to store metricDefinitions in.
keyspace = metrictank
# Cassandra table to store meta records.
meta-record-table = meta_records
# Cassandra table to store meta data of meta record batches.
meta-record-batch-table = meta_record_batches
# Interval at which to poll store for meta record updates.
meta-record-poll-interval = 10s
# Interval at which meta records of old batches get pruned.
meta-record-prune-interval = 24h
# The minimum age a batch of meta records must have to be pruned.
meta-record-prune-age = 72h
# comma separated list of cassandra addresses in host:port form
hosts = localhost:9042
#cql protocol version to use
protocol-version = 4
# write consistency (any|one|two|three|quorum|all|local_quorum|each_quorum|local_one
consistency = one
# cassandra request timeout. valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'
timeout = 1s
# number of concurrent connections to cassandra
num-conns = 10
# synchronize index changes to cassandra. not all your nodes need to do this.
update-cassandra-index = true
# enable SSL connection to cassandra
ssl = false
# cassandra CA certficate path when using SSL
ca-path = /etc/metrictank/ca.pem
# host (hostname and server cert) verification when using SSL
host-verification = true
# enable cassandra user authentication
auth = false
# username for authentication
username = cassandra
# password for authentication
password = cassandra
# enable the creation of the index keyspace and tables, only one node needs this
create-keyspace = true
# File containing the needed schemas in case database needs initializing
schema-file = /etc/metrictank/schema-idx-cassandra.toml
# instruct the driver to not attempt to get host info from the system.peers table
disable-initial-host-lookup = false
# interval at which to perform a connection check to cassandra, set to 0 to disable.
connection-check-interval = 5s
# maximum total time to wait before considering a connection to cassandra invalid. This value should be higher than connection-check-interval.
connection-check-timeout = 30s
```

### Bigtable index

```
[bigtable-meta-record-idx]
enabled = false
# Name of GCP project the bigtable cluster resides in
gcp-project = default
# Name of bigtable instance
bigtable-instance = default
# Table to store meta records
table-name = meta_records
# Table to store meta data of meta record batches.
batch-table-name = meta_record_batches
# Interval at which to poll store for meta record updates
poll-interval = 10s
# Interval at which meta records of old batches get pruned
prune-interval = 24h
# The minimum age a batch of meta records must have to be pruned
prune-age = 72h
# Synchronize meta record changes to bigtable. not all your nodes need to do this
update-records = true
# Enable the creation of the table and column families
create-cf = true
# Max number of metricDefs allowed to be unwritten to bigtable. Must be larger then write-max-flush-size
write-queue-size = 100000
# Max number of metricDefs in each batch write to bigtable
write-max-flush-size = 10000
# Number of writer threads to use
write-concurrency = 5
# synchronize index changes to bigtable. not all your nodes need to do this.
update-bigtable-index = true
# frequency at which we should update the metricdefinition in bigtable, use 0s for instant updates
update-interval = 3h
# Interval at which the index should be checked for stale series.
prune-interval = 3h
# enable the creation of the table and column families
create-cf = true
```

# index-rules.conf

```
# This config file controls when to prune metrics from the index
# Note:
# * This file is optional. If it is not present, we won't prune data
# * Patterns are regexes matched on the metric name (including tags) and tried from top to bottom. First match wins.
# * Anything not matched or resolving to max-stale=0 will not be pruned
# * Patterns are unanchored regular expressions; add '^' or '$' to match the beginning or end of a pattern
# * max-stale is a duration like 7d. if no data has been seen for this time window, it will be pruned. (compared against LastUpdate)
# * Valid units are s/sec/secs/second/seconds, m/min/mins/minute/minutes, h/hour/hours, d/day/days, w/week/weeks, mon/month/months, y/year/years

[default]
pattern = 
max-stale = 0
```

# storage-aggregation.conf

```
# This config file controls which summaries are created (using which consolidation functions) for your lower-precision archives, as defined in storage-schemas.conf
# It is an extension of http://graphite.readthedocs.io/en/latest/config-carbon.html#storage-aggregation-conf
# Note:
# * This file is optional. If it is not present, we will use avg for everything
# * Anything not matched also uses avg for everything
# * xFilesFactor is not honored yet.  What it is in graphite is a floating point number between 0 and 1 specifying what fraction of the previous retention level's slots must have non-null values in order to aggregate to a non-null value. The default is 0.5.
# * aggregationMethod specifies the functions used to aggregate values for the next retention level. Legal methods are avg/average, sum, min, max, and last. The default is average.
# Unlike Graphite, you can specify multiple, as it is often handy to have different summaries available depending on what analysis you need to do.
# When using multiple, the first one listed is the "primary" one, used for reading data unless another one is requested via consolidateBy().
# * the settings configured when metrictank starts are what is applied. So you can enable or disable archives by restarting metrictank.
#
# see https://github.com/grafana/metrictank/blob/master/docs/consolidation.md for related info.

[default]
pattern = .*
xFilesFactor = 0.1
aggregationMethod = avg,min,max
```

# storage-schemas.conf

```
# This config file sets up your retention rules.
# It is an extension of http://graphite.readthedocs.io/en/latest/config-carbon.html#storage-schemas-conf
# Note:
# * You can have 0 to N sections
# * The first match wins, starting from the top. If no match found, we default to single archive of minutely points, retained for 7 days in 2h chunks
# * The patterns are unanchored regular expressions, add '^' or '$' to match the beginning or end of a pattern.
# * When running a cluster of metrictank instances, all instances should have the same agg-settings.
# * Unlike whisper (graphite), the config doesn't stick: if you restart metrictank with updated settings, then those
# will be applied. The configured rollups will be saved by primary nodes and served in responses if they are ready.
# (note in particular that if you remove archives here, we will no longer read from them)
# * Retentions must be specified in order of increasing interval and retention
# * The reorderBuffer an optional buffer that temporarily keeps data points in memory as raw data and allows insertion at random order. The specified value is how many datapoints, based on the raw interval specified in the first defined retention, should be kept before they are flushed out. This is useful if the metric producers cannot guarantee that the data will arrive in order, but it is relatively memory intensive. If you are unsure whether you need this, better leave it disabled to not waste memory. When enabled, you can optionally via 'reorderBufferAllowUpdate' allow updating the value of data points already received (if the timestamp falls within the reorder buffer window).
#
# A given rule is made up of at least 3 lines: the name, regex pattern, retentions and optionally the reorder buffer size.
# The retentions line can specify multiple retention definitions. You need one or more, space separated.
#
# There are 2 formats for a single retention definition:
# 1) 'series-interval:count-of-datapoints'                   legacy and not easy to read
# 2) 'series-interval:retention[:chunkspan:numchunks:ready]' more friendly format with optionally 3 extra fields
#
#Series intervals and retentions are specified using the following suffixes:
#
#s - second
#m - minute
#h - hour
#d - day
#y - year
#
# The final 3 fields are specific to metrictank and if unspecified, use sane defaults.
# See https://github.com/grafana/metrictank/blob/master/docs/memory-server.md for more details
#
# chunkspan: duration of chunks. e.g. 10min, 30min, 1h, 90min...
# must be valid value as described here https://github.com/grafana/metrictank/blob/master/docs/memory-server.md#valid-chunk-spans
# Defaults to a the smallest chunkspan that can hold at least 100 points.
#
# numchunks: number of raw chunks to keep in in-memory ring buffer
# See https://github.com/grafana/metrictank/blob/master/docs/memory-server.md for details and trade-offs, especially when compared to chunk-cache
# which may be a more effective method to cache data and alleviate workload for cassandra.
# Defaults to 2
#
# ready: whether, or as of what data timestamp, the archive is ready for querying.
# This is useful if you recently introduced a new archive, but it's still being populated, so you want to control whether (or to which extent) the archive can be used for queries.
# It supports two syntaxes:
# * unix timestamp: the archive contains data as of this timestamp
# * boolean: (legacy): whether or not the archive is completely ready or not ready at all.
# Defaults to true
#
# Here's an example with multiple retentions:
# [apache_busyWorkers]
# pattern = ^servers\.www.*\.workers\.busyWorkers$
# retentions = 1s:1d:10min:1,1m:21d,15m:5y:2h:1:false
#
# This example has 3 retention definitions, the first and last override some default options (to use 10minutely and 2hourly chunks and only keep one of them in memory
# and the last rollup is marked as not ready yet for querying.

[default]
pattern = .*
retentions = 1s:35d:10min:7
# reorderBuffer = 20
# reorderBufferAllowUpdate = true
```

This file is generated by [config-to-doc](https://github.com/grafana/metrictank/blob/master/scripts/dev/config-to-doc.sh)

