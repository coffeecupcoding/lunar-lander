# lunar-lander

This is a rewrite of the classic Lunar Lander game in Python, based on the
BASIC source code from 101 BASIC Computer Games by David Ahl.  The original
can be found here:

https://www.atariarchives.org/basicgames/showpage.php?page=106

I have taken some liberties with the text output of the program but have
tried to keep the same feel, and the calculations should be very similar.

It appears to work correctly for typical usage, I wouldn't recommend trying
to put something on the Moon with it though.



## Some known issues with this code:

- It's not very Pythonic

- There are no options to set initial parameters, and the intro text uses
static values

- The overall design is not particularly consistent internally (it's
essentially version 1.2, and work on it stopped when replay worked)

- The calculations in the program are copied from the original; I don't
fully understand some parts and may also have mis-read some of the BASIC

- Even where the formulas are correct, there may be differences in calculation
due to different implementations of floating point numbers

- There's no input checking currently

- The output could be improved, particularly to fit on 40 column screens

- There are no unit tests (NASA would not sign off on this code)

- There are likely bugs (in edge cases, of course)

- It likely uses more memory than any computer had at the time the original
was written


## Other notes:

- The logic of the original took some thinking through... this would be a
classic example for 'GOTOs considered harmful'

- The variable W appears to be used for two different purposes.

- I don't think line 400 does anything: V was just set to J so the IF in
line 390 takes precedence.  I'm not sure what it was intended to do, either.

