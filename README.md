# Wolfram Rule Simulator

This program displays a simulated Wolfram rule.

### Prerequisites

Download and install [termbox-go](http://github.com/nsf/termbox-go) using: ```go get github.com/nsf/termbox-go```

### Usage

Simply call ```go run wolframsimulator.go``` to compile and run.

Arguments determine program behavior as follows:

* ```-rule``` (```-r```) specifies the rule to simulate.
* ```-width``` (```-w```) specifies the width of the simulation.
* ```-height``` (```-h```) specifies the height of the simulation.

Default argument values:

* ```rule =``` 30
* ```width =``` width of terminal
* ```height =``` height of terminal
