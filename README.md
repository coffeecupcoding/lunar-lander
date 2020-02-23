# lunar-lander

This is a rewrite of the classic Lunar Lander game in several languages, based on the
BASIC source code from 101 BASIC Computer Games by David Ahl.  The original
can be found here:

https://www.atariarchives.org/basicgames/showpage.php?page=106

I wrote the initial Python version in honor of the 50th Anniversary of the Apollo 11 lunar
landing, and all of the versions as fun short projects to explore the basics of the languages.

I have taken some liberties with the text output of the program but have
tried to keep the same feel, and the calculations should be very similar.

All versions appear to work correctly, I wouldn't recommend trying
to put something on the Moon with any of them though.


## Notes on the BASIC code:

- The logic of the original took some thinking through... this would be a
classic example for 'GOTOs considered harmful'

- The variable W appears to be used for two different purposes

- I don't think line 400 does anything: V was just set to J so the IF in
line 390 takes precedence.  I'm not sure what it was intended to do, either

