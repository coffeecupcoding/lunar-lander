# lunar-lander-rust
Rewrite of the classic lunar lander game in Rust

This is based on on the Python version I created
(see https://github.com/coffeecupcoding/lunar-lander-python).  I wanted to write something small
in Rust, and that project was fresh in my mind.

## Known issues

- It's not very good Rust, as this is a first pass

- It could use some refactoring, particularly to split up run_game() ~~and
make the input code a separate function~~ Done

- As for the Python version, some of the formulas are suspect, ~~and it has
the same end-game bug (though it's expressed a bit differently)~~ Fixed
after finding and fixing the bug in the Python version

## Future work

Well...

~~If I add anything, it'll be bounds/error checking in the input, the
program ought not fail no matter what you type.~~

