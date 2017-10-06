# Wolfram Rule Simulator

This program is designed to print a simulation of one or many wolfram rules in the terminal.

#### Prerequisites

This program uses STTY to detect terminal window size.

#### Usage

Simply call ```go run wolframrules.go``` to compile and run. Arguments determine program behavior as follows:

..* Calling with no arguments will print all 256 Wolfram rules across the full width of the terminal for the full height of the terminal, displaying a new rule every second.

..* Calling with at least one argument will have the first be interpretted as a rule ID.

..* Calling with at least two arguments will have the second be interpretted as the number of lines for which the rule will be simulated.

..* Calling with at least three arguments will ahve the third be interprestted as the width of the rows to be simulated.
