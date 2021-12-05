## Tomato Timer CLI
Timer that can be used via command line interface.

### Instruction
To create executable file
`` make build ``

To run executable
`` ./cmd/tomato-timer/tomato-timer -s 10``

The above timer will be off in 10 second

To create package with related resources
`` make pack ``

Package will be located 
`` ./target/tomato-timer/ ``

### CLI Timer 

Possible command options

`` tomato-timer --help ``


`` tomato-timer --hours 1 --minutes 30 --seconds 45 `` or `` tomato-timer -h 1 -m 30 -s 45 ``
