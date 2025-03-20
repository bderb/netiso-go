# NetISO Go

A simple Xbox 360 ISO game streaming server for use with Aurora dashboard.

## Description

This project is a recreation of the original NetISO.exe that was released for streaming Xbox 360 ISO files via TCP with a modified Aurora Dashboard. The server replicates the same functionality as the original server application, and extends it with additional functionality.

The main feature of this version of the application it is nearly identical in terms of performance while also being a cross platform option allowing for dedicated hosting options.

## Getting Started

To get started simply download the precompiled executable for your platform of choice and the server will start automatically. By default the directory where you run the executable will be where it searches for Xbox ISO game files, the server will only get Xbox Game disc files to stream.

### Dependencies

* Modified Xbox 360 running the custom Aurora Dashboard with support for this
* Windows/Linux machine on the same network as the Xbox 360

### Executing program

* Simply run the server via the Terminal with the selected parameters required, server with start streaming on localhost:4323.

## Help

The script has the following parameters that can be accessed using the `-h` flag while executing.

```
netiso -h

Usage of netiso:
  -d    Enable debug mode.
  -f string
        Path to the Xbox ISO files.
  -s    Logs network speed (mbps).
```

## License

This project is licensed under the [NAME HERE] License - see the LICENSE.md file for details