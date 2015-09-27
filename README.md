# LiFX CLI

Cross-platform LiFX command line interface. Get info, set color and power state. This code has been developed and maintained by Ven at September 2015.

## Requirements
Golang 1.5 

https://golang.org/dl/

## Installation

```bash
go get github.com/2tvenom/golifx
go get github.com/2tvenom/lifx
```

## Example
Get help
```bash
lifx --help
```

Lookup bulbs
```bash
lifx --lookup
```
Output:
```
MAC: d0:73:d5:01:90:d7
IP: 192.168.0.2:56700
Power state: false
```

Get color state
```bash
lifx --lookup --color
```
Output:
```
MAC: d0:73:d5:01:90:d7
IP: 192.168.0.2:56700
Label: Ven LiFX
Power state: true
Color:
HUE: 52000
Saturation: 0
Brightness: 32336
Kelvin: 6196
```

Json format output:
```bash
lifx --lookup --color --json
```
Output:
```javascript
[{"color":{"brightness":32336,"hue":52000,"kelvin":6196,"saturation":0},"ip":{"IP":"192.168.0.2","Port":56700,"Zone":""},"label":"Ven LiFX","mac":"d0:73:d5:01:90:d7","power_state":true}]
```

Set color state:
```bash
lifx --lookup --color --json
```
Output:
```javascript
[{"color":{"brightness":32336,"hue":52000,"kelvin":6196,"saturation":0},"ip":{"IP":"192.168.0.2","Port":56700,"Zone":""},"label":"Ven LiFX","mac":"d0:73:d5:01:90:d7","power_state":true}]
```

Turn off:
```bash
lifx --bulb d0:73:d5:01:90:d7 --off
```

Turn on:
```bash
lifx --bulb d0:73:d5:01:90:d7 --on
```

Set blue color:
```bash
lifx --bulb d0:73:d5:01:90:d7 --hue=36240 --saturation=65535 --brightness=64580 --kelvin=3505
```

Set purple color:
```bash
lifx --bulb d0:73:d5:01:90:d7 --hue=49719 --saturation=65535 --brightness=64580 --kelvin=3505
```

Set red color:
```bash
lifx --bulb d0:73:d5:01:90:d7 --hue=64489 --saturation=63482 --brightness=65535 --kelvin=3500	
```

## Links
 - LiFX protocol specification http://lan.developer.lifx.com/
 - Community https://community.lifx.com/c/developing-with-lifx

## Licence
[WTFPL](http://www.wtfpl.net/)