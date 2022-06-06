# syslogqd: A simple Syslog server 

A simple command line syslog server written in go for investigating / debugging IoT and similar devices.

Some IoT devices, VoIP devices, etc. can be configured to write debugging information to a syslog server. While testing
a device, it can be useful to have a convenient syslog server that can be run as needed from the command line to
either show log entries as they occur, or write them to a file.

Note that simple devices often don't obey any relevant RFCs. **syslogqd** is therefore designed to be tolerant 
of poorly formatted messages, messages without timestamps, etc. 

## Usage
```
syslogqd [options]
```
By default syslogqd will listen for UDP and TCP syslog messages on the default port of 514, and the messages will be written to standard output. Note that you must be have root privileges (Linux) or enable firewall access (Windows) to listen on ports 1-1023. 

Running
```
syslogqd -help
```
Will generate full usage information. There are options to:
 - listen on a different port
 - ignore lower-severity messages
 - save a copy of the output to a file
 - suppress output to stdout

## Timestamps, facilities, and severities

The output format is a timestamp followed by an IP address and optional severity/facility, followed by the contents of the syslog message. 

If a timestamp can be found in the syslog message, it will be extracted from the message and used as the message time: otherwise the current time will be used. Message timestamps are shown in UTC.

If a severity/facility is found in the message it will be extracted and converted from `<number>` to a severity/facility string. 

Here's some sample output:
```
2022-06-06T13:44:58Z 192.168.1.49: shellyplus1-7c87ce72ad58 274 33503.945 2 2|mgos_http_server.c:180 0x3ffd6e40 HTTP connection from 192.168.1.204:56508
2022-06-06T13:44:58Z 192.168.1.49: shellyplus1-7c87ce72ad58 275 33504.261 2 2|mg_rpc.c:314 shelly.getdeviceinfo via WS_in 192.168.1.204:56508
2022-06-06T13:45:30Z 192.168.1.49: shellyplus1-7c87ce72ad58 276 33536.661 2 2|mg_rpc.c:314 shelly.getconfig via WS_in 192.168.1.204:56508 user admin

```

# Build and Install

Assuming you have a Go compiler installed, simply execute `go install github.com/m-z-b/syslogqd@latest` to download, 
build, and install the latest version.

If you do not have a Go compiler installed, pre-built binaries for Windows and Linux/amd64 are included with each release. 
Simply rename the appropriate binary and place it in your path. 


# Contributing

Suggestions and pull requests are welcome. 

Please remember that syslogqd is intended as a simple tool for viewing and recording syslog entries from IoT and other devices: additional features should suit this intended usage.

# License

syslogqd is licensed under the Apache License, Version 2.0.
