# snap collector plugin - influxdb

This plugin collects statistical and diagnostic information about each InfluxDB node. This information can be very useful to assist with troubleshooting and performance analysis of the database itself.

This plugin gather InfluxDB internal system monitoring information in response to commands `SHOW STATS` and `SHOW DIAGNOSTICS`.
															
The plugin is used in the [snap framework] (http://github.com/intelsdi-x/snap).				

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Global Config](#global-config)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements

- Linux system
- InfluxDB (version 0.9 or higher)

### Installation

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-influxdb																				

Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-influxdb.git
```

Build the snap influxdb collector plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage

* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)


## Documentation

To learn more about influxDB System Monitoring, visit:

* [InfluxDB Server Monitoring doc] (https://docs.influxdata.com/influxdb/v0.9/administration/statistics/)
* [InfluxDB Server Monitoring README] (https://github.com/influxdata/influxdb/blob/master/monitor/README.md)
* blog post ["How to use the show stats command to monitor InfluxDB"](https://influxdata.com/blog/how-to-use-the-show-stats-command-and-the-_internal-database-to-monitor-influxdb/)

### Global Config

Global configuration files are described in snap's documentation. For this plugin section "influxdb" in "collector" specifing the following options needs to be added (see [exemplary configs file](examples/configs/snap-config-sample.json)):

Name 	  	| Data Type | Description
------------|-----------|-----------------------
"host" 		| string 	| hostname of InfluxDB http API
"port" 		| int	 	| port of InfluxDB http API (by default 8086)
"user" 		| string 	| user name
"password" 	| string 	| user password


### Collected Metrics

This plugin has the ability to gather:

a) all **diagnostic information** of InfluxDB system itself, represented by the metrics with prefix `/intel/influxdb/diagn/`

b) all **statistical information** of InfluxDB system itself, represented by the metrics with prefix `/intel/influxdb/stat/`
                                                                                                
Metric Name | Description
------------ | -------------
/intel/influxdb/diagn/build/Branch |
/intel/influxdb/diagn/build/Commit |
/intel/influxdb/diagn/build/Version |
/intel/influxdb/diagn/network/hostname |
/intel/influxdb/diagn/runtime/GOARCH |
/intel/influxdb/diagn/runtime/GOMAXPROCS |
/intel/influxdb/diagn/runtime/GOOS |
/intel/influxdb/diagn/runtime/version |
/intel/influxdb/diagn/system/PID |
/intel/influxdb/diagn/system/currentTime |
/intel/influxdb/diagn/system/started |
/intel/influxdb/diagn/system/uptime |
| |
/intel/influxdb/stat/engine/<engine_id>/blks_write |
/intel/influxdb/stat/engine/<engine_id>/blks_write_bytes |
/intel/influxdb/stat/engine/<engine_id>/blks_write_bytes_c |
/intel/influxdb/stat/engine/<engine_id>/points_write |
/intel/influxdb/stat/engine/<engine_id>/points_write_dedupe |
/intel/influxdb/stat/httpd/auth_fail |
/intel/influxdb/stat/httpd/ping_req |
/intel/influxdb/stat/httpd/points_written_ok |
/intel/influxdb/stat/httpd/query_req |
/intel/influxdb/stat/httpd/query_resp_bytes |
/intel/influxdb/stat/httpd/req |
/intel/influxdb/stat/httpd/write_req |
/intel/influxdb/stat/httpd/write_req_bytes |
/intel/influxdb/stat/runtime/Alloc |
/intel/influxdb/stat/runtime/Frees |
/intel/influxdb/stat/runtime/HeapAlloc |
/intel/influxdb/stat/runtime/HeapIdle |
/intel/influxdb/stat/runtime/HeapInUse |
/intel/influxdb/stat/runtime/HeapObjects |
/intel/influxdb/stat/runtime/HeapReleased |
/intel/influxdb/stat/runtime/HeapSys |
/intel/influxdb/stat/runtime/Lookups |
/intel/influxdb/stat/runtime/Mallocs |
/intel/influxdb/stat/runtime/NumGC |
/intel/influxdb/stat/runtime/NumGoroutine |
/intel/influxdb/stat/runtime/PauseTotalNs |
/intel/influxdb/stat/runtime/Sys |
/intel/influxdb/stat/runtime/TotalAlloc |
/intel/influxdb/stat/shard/<shard_id>/fields_create |
/intel/influxdb/stat/shard/<shard_id>/series_create |
/intel/influxdb/stat/shard/<shard_id>/write_points_ok |
/intel/influxdb/stat/shard/<shard_id>/write_req |
/intel/influxdb/stat/write/point_req |
/intel/influxdb/stat/write/point_req_local |
/intel/influxdb/stat/write/req |
/intel/influxdb/stat/write/write_ok |

The list of available metrics might be vary depending on the influxdb version or the system configuration.

Diagnostics information are gathered only once at the beggining of collecting process, because they are constant during running the influxdb process.

In task manifest there are declaration of metrics names which will be collected and an interval (see [exemplary task manifest] (examples/tasks/influxdb-file.json)). By default metrics are gathered once per second.

There is a possibility to use asterisk (*) in metrics names (in all levels), for example:
 - `/intel/influxdb/*` means that all available metrics representing statistical and diagnostic information will be collected
 - `/intel/influxdb/stat/*` means that all metrics representing statistical information will be collected
 - `/intel/influxdb/stat/httpd/*` means that only metrics for module `httpd` will be collected (i.e.: auth_fail, req, ping_req, query_req, query_resp_bytes, write_req, write_req_bytes, points_written_ok)


### Examples

Example of running snap influxdb collector and writing data to file.

Run the snap daemon:
```
$ snapd -l 1 -t 0 --config $SNAP_INFLUXDB_COLLECTOR_PLUGIN_DIR/examples/configs/snap-config-sample.json
```

Load snap influxdb collector plugin:
```
$ snapctl plugin load $SNAP_INFLUXDB_COLLECTOR_PLUGIN_DIR/build/rootfs/snap-plugin-collector-influxdb
Plugin loaded
Name: influxdb
Version: 1
Type: collector
Signed: false
Loaded Time: Fri, 26 Feb 2016 09:09:03 UTC
```
See all available metrics:
```
$ snapctl metric list

NAMESPACE                                                VERSIONS
/intel/influxdb/*                                        1
/intel/influxdb/diagn/*                                  1
/intel/influxdb/diagn/build/*                            1
/intel/influxdb/diagn/build/Branch                       1
/intel/influxdb/diagn/build/Commit                       1
/intel/influxdb/diagn/build/Version                      1
/intel/influxdb/diagn/network/*                          1
/intel/influxdb/diagn/network/hostname                   1
/intel/influxdb/diagn/runtime/*                          1
/intel/influxdb/diagn/runtime/GOARCH                     1
/intel/influxdb/diagn/runtime/GOMAXPROCS                 1
/intel/influxdb/diagn/runtime/GOOS                       1
/intel/influxdb/diagn/runtime/version                    1
/intel/influxdb/diagn/system/*                           1
/intel/influxdb/diagn/system/PID                         1
/intel/influxdb/diagn/system/currentTime                 1
/intel/influxdb/diagn/system/started                     1
/intel/influxdb/diagn/system/uptime                      1
/intel/influxdb/stat/*                                   1
/intel/influxdb/stat/engine/*                            1
/intel/influxdb/stat/engine/1/*                          1
/intel/influxdb/stat/engine/1/blks_write                 1
/intel/influxdb/stat/engine/1/blks_write_bytes           1
/intel/influxdb/stat/engine/1/blks_write_bytes_c         1
/intel/influxdb/stat/engine/1/points_write               1
/intel/influxdb/stat/engine/1/points_write_dedupe        1
/intel/influxdb/stat/engine/2/*                          1
/intel/influxdb/stat/engine/2/blks_write                 1
/intel/influxdb/stat/engine/2/blks_write_bytes           1
/intel/influxdb/stat/engine/2/blks_write_bytes_c         1
/intel/influxdb/stat/engine/2/points_write               1
/intel/influxdb/stat/engine/2/points_write_dedupe        1
/intel/influxdb/stat/httpd/*                             1
/intel/influxdb/stat/httpd/auth_fail                     1
/intel/influxdb/stat/httpd/ping_req                      1
/intel/influxdb/stat/httpd/points_written_ok             1
/intel/influxdb/stat/httpd/query_req                     1
/intel/influxdb/stat/httpd/query_resp_bytes              1
/intel/influxdb/stat/httpd/req                           1
/intel/influxdb/stat/httpd/write_req                     1
/intel/influxdb/stat/httpd/write_req_bytes               1
/intel/influxdb/stat/runtime/*                           1
/intel/influxdb/stat/runtime/Alloc                       1
/intel/influxdb/stat/runtime/Frees                       1
/intel/influxdb/stat/runtime/HeapAlloc                   1
/intel/influxdb/stat/runtime/HeapIdle                    1
/intel/influxdb/stat/runtime/HeapInUse                   1
/intel/influxdb/stat/runtime/HeapObjects                 1
/intel/influxdb/stat/runtime/HeapReleased                1
/intel/influxdb/stat/runtime/HeapSys                     1
/intel/influxdb/stat/runtime/Lookups                     1
/intel/influxdb/stat/runtime/Mallocs                     1
/intel/influxdb/stat/runtime/NumGC                       1
/intel/influxdb/stat/runtime/NumGoroutine                1
/intel/influxdb/stat/runtime/PauseTotalNs                1
/intel/influxdb/stat/runtime/Sys                         1
/intel/influxdb/stat/runtime/TotalAlloc                  1
/intel/influxdb/stat/shard/*                             1
/intel/influxdb/stat/shard/1/*                           1
/intel/influxdb/stat/shard/1/fields_create               1
/intel/influxdb/stat/shard/1/series_create               1
/intel/influxdb/stat/shard/1/write_points_ok             1
/intel/influxdb/stat/shard/1/write_req                   1
/intel/influxdb/stat/shard/2/*                           1
/intel/influxdb/stat/shard/2/fields_create               1
/intel/influxdb/stat/shard/2/series_create               1
/intel/influxdb/stat/shard/2/write_points_ok             1
/intel/influxdb/stat/shard/2/write_req                   1
/intel/influxdb/stat/write/*                             1
/intel/influxdb/stat/write/point_req                     1
/intel/influxdb/stat/write/point_req_local               1
/intel/influxdb/stat/write/req                           1
/intel/influxdb/stat/write/write_ok                      1

```

Load file plugin for publishing:
```
$ snapctl plugin load $SNAP_DIR/build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 26 Feb 2016 09:16:44 UTC
```

Create a task JSON file (exemplary file in examples/tasks/influxdb-file.json):  
```json

{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/influxdb/stat/httpd/*": {},
				"/intel/influxdb/stat/write/*": {},
				"/intel/influxdb/stat/runtime/Alloc": {},
				"/intel/influxdb/stat/runtime/Frees": {},
				"/intel/influxdb/diagn/system/PID" : {}
            },
            "config": {},
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_influxdb_internal_monitoring"
                    }
                }
            ]
        }
    }
}
```

Create a task:
```
$ snapctl task create -t $SNAP_INFLUXDB_COLLECTOR_PLUGIN_DIR/examples/tasks/influxdb-file.json
Using task manifest to create task
Task created
ID: d4392f17-11f0-4f64-8701-2708f432b50a
Name: Task-d4392f17-11f0-4f64-8701-2708f432b50a
State: Running
```

See sample output from `snapctl task watch <task_id>`

```
$ snapctl task watch d4392f17-11f0-4f64-8701-2708f432b50a

Watching Task (d4392f17-11f0-4f64-8701-2708f432b50a):
NAMESPACE                                        DATA                    TIMESTAMP                                       SOURCE
/intel/influxdb/diagn/system/PID                 16191                   2016-02-26 09:25:47.353886681 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/httpd/auth_fail             63                      2016-02-26 09:25:47.354068464 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/httpd/ping_req              14                      2016-02-26 09:25:47.354073486 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/httpd/points_written_ok     6.69802551e+08          2016-02-26 09:25:47.354137322 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/httpd/query_req             2310                    2016-02-26 09:25:47.354132981 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/httpd/query_resp_bytes      2.0678075e+07           2016-02-26 09:25:47.354113672 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/httpd/req                   6.824003e+06            2016-02-26 09:25:47.354028193 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/httpd/write_req             6.821614e+06            2016-02-26 09:25:47.354095848 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/httpd/write_req_bytes       7.1455502846e+10        2016-02-26 09:25:47.354105454 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/runtime/Alloc               1.07946848e+08          2016-02-26 09:25:47.354209458 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/runtime/Frees               1.5858776371e+10        2016-02-26 09:25:47.354412092 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/write/point_req             6.82720753e+08          2016-02-26 09:25:47.354579209 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/write/point_req_local       6.82720753e+08          2016-02-26 09:25:47.354619612 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/write/req                   7.000631e+06            2016-02-26 09:25:47.354624219 +0000 UTC         node-25.domain.tld
/intel/influxdb/stat/write/write_ok              7.000823e+06            2016-02-26 09:25:47.354609029 +0000 UTC         node-25.domain.tld

```
(Keys `ctrl+c` terminate task watcher)


These data are published to file and stored there (in this example in /tmp/published_influxdb_internal_monitoring).

Stop task:
```
$ snapctl task stop d4392f17-11f0-4f64-8701-2708f432b50a
Task stopped:
ID: d4392f17-11f0-4f64-8701-2708f432b50a
```

### Roadmap

There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-influxdb/issues).

## Community Support
This repository is one of **many** plugins in the **Snap Framework**: a powerful telemetry agent framework. To reach out on other use cases, visit:

* [Snap Gitter channel] (https://gitter.im/intelsdi-x/snap)

The full project is at http://github.com:intelsdi-x/snap.

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
List authors, co-authors and anyone you'd like to mention

* Author: 	[Izabella Raulin](https://github.com/IzabellaRaulin)

**Thank you!** Your contribution is incredibly important to us.
