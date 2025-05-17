# asciicast2script / script2asciicast

A pair of commands to convert between `asciinema`'s asciicasts and `script`'s typescript/timingfile.
asciicast2script supports both v2 and v3 asciicasts.
script2asciicast outputs asciicast v2 by default.
The `-v3` flag can be used to output asciicast v3.

## Installation

```
go install github.com/wk-y/asciicast2script/cmd/...@latest
```

## Usage examples

Conversion to script:
```
asciinema rec -c 'timeout 5 top -d 0.5' demo.cast
asciicast2script demo.cast
scriptreplay -t timingfile
```

Conversion to asciicast:
```
# record a script
script -c 'timeout 5 top -d 0.5' --timing=timingfile typescript
script2asciicast demo.cast
asciinema play demo.cast
```
