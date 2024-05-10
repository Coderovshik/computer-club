/*
Compclub manages the work of the computer club.
Output is written to standard output.
Path to computer club configuration .txt file has to be specified.

Usage:

	compclub <config>

Configuration file structure:

	<number of desks in the club>
	<worktime start> <worktime end>
	<per hour pay for desk usage>
	<event 1 occured at> <event 1 id> <event 1 body>
	<event 2 occured at> <event 2 id> <event 2 body>
	...
	<event N occured at> <event N id> <event N body>

The first line contains the number of tables as a positive integer.
The second line specifies the start and end times of the computer club,
separated by a space.
The third line specifies the cost of an hour in the computer club as an integer
positive number.
Then a list of incoming events is specified, separated by line breaks. Inside
lines, a single space is used as a separator between elements.

  - Client names are a combination of characters from the alphabet [a..z, 0..9, _, -].
  - Time is specified in 24-hour format with a colon as a separator XX:XX, leading zeros are required for input and output (for example 15:03 or 08:09).
  - Each table has its own number from 1 to N, where N is the total number of tables indicated in configurations.
  - All events occur sequentially in time. (event time N+1) ≥ (time events N).

In case any of these requirements are violated, program stops execution and outputs an error.
*/
package main

import (
	"fmt"
	"os"
)

func main() {
	inputPath, err := GetInputPath()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	f, err := GetInputFile(inputPath)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer f.Close()

	c, err := ParseClub(f)
	if err != nil {
		fmt.Printf("%s %s", f.Name(), err)
		os.Exit(1)
	}

	c.Run()
}
