# Wolfram Rule Simulator

This program is designed to print a simulation of one or many Wolfram rules in the terminal.

#### Usage

Simply call ```go run wolframrules.go``` to compile and run. Arguments determine program behavior as follows:

```-rule``` specifies rule to simulate. 
```-rows``` specifies the number of rows for which to simulate a rule.
```-width``` specifies row width of simulation.

#### Default Argument Values

Not specifying ```rule``` will cause the program to print all 256 rules sequentially, one every second.
Not specifying ```rows``` will set it to the terminal height (accounting for space needed to print rule label at the top).
Not specifying ```width``` will set it to the width of the terminal.