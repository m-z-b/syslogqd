# syslogqd: A simple Syslog server 

Command line syslog server written in Go for investigating / debugging IoT and similar devices.

## Installation
If you have Go installed, run:
```
go install github.com/m-z-b/syslogqd@latest
```
or download one of the releases: pre-built binaries for Windows and Linux/amd64 are included with each release. 
Simply rename the appropriate binary and copy it to somewhere in your search path. 

## Usage
```
$ syslogqd
syslogqd V1.0 listening on port 514 for severity >= debug
Use Ctrl-C to exit
2022-06-06T13:44:58Z 192.168.1.49: shellyplus1-7c87ce72ad58 274 33503.945 2 2|mgos_http_server.c:180 0x3ffd6e40 HTTP connection from 192.168.1.204:56508
2022-06-06T13:44:58Z 192.168.1.49: shellyplus1-7c87ce72ad58 275 33504.261 2 2|mg_rpc.c:314 shelly.getdeviceinfo via WS_in 192.168.1.204:56508
2022-06-06T13:45:30Z 192.168.1.49: shellyplus1-7c87ce72ad58 276 33536.661 2 2|mg_rpc.c:314 shelly.getconfig via WS_in 192.168.1.204:56508 user admin
...
```
By default syslogqd will listen for UDP and TCP syslog messages on the default port of 514, and the messages will be written to standard output. Note that you must have root privileges on Linux to listen to ports 1-1023; on Windows you may need to enable firewall access. 

`syslogqd -help` will show full usage information. There are options to:
 - listen on a different port
 - ignore lower-severity messages
 - save a copy of the output to a file
 - suppress output to stdout


## Timestamps, facilities, and severities

Note that simple devices may not obey any relevant RFCs. *syslogqd* is therefore designed to be tolerant 
of poorly formatted messages, messages without timestamps, etc. 

The output format is a timestamp followed by an IP address and optional severity/facility, followed by the contents of the syslog message. 

If a timestamp can be found in the syslog message, it will be extracted from the message and used as the message time: otherwise the time the message was received will be used. Message timestamps are shown in UTC.

If a severity/facility is found in the message it will be extracted and converted from `<number>` to a severity/facility string. 

## Contributing

Suggestions and pull requests are welcome. 

Please remember that *syslogqd* is intended as a simple tool for viewing and recording syslog messages from IoT and other devices: additional features should suit this intended usage.

# License

syslogqd is licensed under the Apache License, Version 2.0.
