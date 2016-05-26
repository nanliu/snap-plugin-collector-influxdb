# snap collector plugin - influxdb

This plugin collects statistical and diagnostic information about each InfluxDB node. This information can be very useful to assist with troubleshooting and performance analysis of the database itself.

This plugin has ability to gather InfluxDB internal system monitoring information in response to the following commands: `SHOW STATS` and `SHOW DIAGNOSTICS`.
															
The plugin is used in the [snap framework] (http://github.com/intelsdi-x/snap).				

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
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

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

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
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-influxdb/blob/master/README.md#examples).

## Documentation

To learn more about influxDB System Monitoring, visit:

* [InfluxDB Server Monitoring Documentation] (https://docs.influxdata.com/influxdb/v0.9/administration/statistics/)
* [InfluxDB Server Monitoring README] (https://github.com/influxdata/influxdb/blob/master/monitor/README.md)
* blog post ["How to use the show stats command to monitor InfluxDB"](https://influxdata.com/blog/how-to-use-the-show-stats-command-and-the-_internal-database-to-monitor-influxdb/)

### Global Config

Global configuration file is described in snap's documentation. For this plugin section "influxdb" in "collector" specifing the following options needs to be added (see [exemplary configs file](examples/configs/snap-config-sample.json)):

Name 	  	| Data Type | Description
------------|-----------|-----------------------
"host" 		| string 	| hostname of InfluxDB http API
"port" 		| int	 	| port of InfluxDB http API (by default 8086)
"user" 		| string 	| user name
"password" 	| string 	| user password


### Collected Metrics

List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-influxdb/blob/master/METRICS.md).

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
Version: 4
Type: collector
Signed: false
Loaded Time: Fri, 26 Feb 2016 09:09:03 UTC
```
See all available metrics:
```
$ snapctl metric list

NAMESPACE                                                VERSIONS
/intel/influxdb/diagn/build/Branch                       2
/intel/influxdb/diagn/build/Commit                       2
/intel/influxdb/diagn/build/Version                      2
/intel/influxdb/diagn/network/hostname                   2
/intel/influxdb/diagn/runtime/GOARCH                     2
/intel/influxdb/diagn/runtime/GOMAXPROCS                 2
/intel/influxdb/diagn/runtime/GOOS                       2
/intel/influxdb/diagn/runtime/version                    2
/intel/influxdb/diagn/system/PID                         2
/intel/influxdb/diagn/system/currentTime                 2
/intel/influxdb/diagn/system/started                     2
/intel/influxdb/diagn/system/uptime                      2
/intel/influxdb/stat/engine/1/blks_write                 2
/intel/influxdb/stat/engine/1/blks_write_bytes           2
/intel/influxdb/stat/engine/1/blks_write_bytes_c         2
/intel/influxdb/stat/engine/1/points_write               2
/intel/influxdb/stat/engine/1/points_write_dedupe        2
/intel/influxdb/stat/engine/2/blks_write                 2
/intel/influxdb/stat/engine/2/blks_write_bytes           2
/intel/influxdb/stat/engine/2/blks_write_bytes_c         2
/intel/influxdb/stat/engine/2/points_write               2
/intel/influxdb/stat/engine/2/points_write_dedupe        2
/intel/influxdb/stat/httpd/auth_fail                     2
/intel/influxdb/stat/httpd/ping_req                      2
/intel/influxdb/stat/httpd/points_written_ok             2
/intel/influxdb/stat/httpd/query_req                     2
/intel/influxdb/stat/httpd/query_resp_bytes              2
/intel/influxdb/stat/httpd/req                           2
/intel/influxdb/stat/httpd/write_req                     2
/intel/influxdb/stat/httpd/write_req_bytes               2
/intel/influxdb/stat/runtime/Alloc                       2
/intel/influxdb/stat/runtime/Frees                       2
/intel/influxdb/stat/runtime/HeapAlloc                   2
/intel/influxdb/stat/runtime/HeapIdle                    2
/intel/influxdb/stat/runtime/HeapInUse                   2
/intel/influxdb/stat/runtime/HeapObjects                 2
/intel/influxdb/stat/runtime/HeapReleased                2
/intel/influxdb/stat/runtime/HeapSys                     2
/intel/influxdb/stat/runtime/Lookups                     2
/intel/influxdb/stat/runtime/Mallocs                     2
/intel/influxdb/stat/runtime/NumGC                       2
/intel/influxdb/stat/runtime/NumGoroutine                2
/intel/influxdb/stat/runtime/PauseTotalNs                2
/intel/influxdb/stat/runtime/Sys                         2
/intel/influxdb/stat/runtime/TotalAlloc                  2
/intel/influxdb/stat/shard/1/fields_create               2
/intel/influxdb/stat/shard/1/series_create               2
/intel/influxdb/stat/shard/1/write_points_ok             2
/intel/influxdb/stat/shard/1/write_req                   2
/intel/influxdb/stat/shard/2/fields_create               2
/intel/influxdb/stat/shard/2/series_create               2
/intel/influxdb/stat/shard/2/write_points_ok             2
/intel/influxdb/stat/shard/2/write_req                   2
/intel/influxdb/stat/write/point_req                     2
/intel/influxdb/stat/write/point_req_local               2
/intel/influxdb/stat/write/req                           2
/intel/influxdb/stat/write/write_ok                      2

```

Load file plugin for publishing:
```
$ snapctl plugin load $SNAP_DIR/build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 4
Type: publisher
Signed: false
Loaded Time: Fri, 26 Feb 2016 09:16:44 UTC
```

Create a task JSON file (exemplary files in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-influxdb/blob/master/examples/tasks/)):
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

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-influxdb/issues) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-influxdb/pulls).

## Community Support
This repository is one of **many** plugins in the **Snap Framework**: a powerful telemetry agent framework. 

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [snap Gitter channel](https://gitter.im/intelsdi-x/snap).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: 	[Izabella Raulin](https://github.com/IzabellaRaulin)
* Author: 	[Marcin Krolik](https://github.com/marcin-krolik)
