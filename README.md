![](assets/background.jpg)

# PWgo

#### Go implementation of PW pass gen

## What is it?

PWgo is a password generator, based somewhat on principles of xkcd password strength comic

![](assets/password_strength.png)
[link](https://xkcd.com/936/) to the original

The app sits in the menubar and upon call (by click) will generate the following data:

1. Username -- consists of a merged adjective and a noun. All spaces and punctuation are omitted. Words are capitalised and then concatenated
2. Password -- consists of 4, 6 or 8 words, one special character and one number (excluding zero, to avoid confusion). The words in the password are random, capitalized. The numeral and special character are set at the end of the password as default (can be changed in the settings)
3. Words of Wisdom -- just a fun experiment where a random adjective and noun are juxtaposed

All items mentioned above are copied into the clipboard when clicked on.

Password is automatically added to the clipboard once the menu icon is pressed.

The app is written in pure Go. Compiled and tested on macOS. I assume it can be compiled and run on Windows and Linux as well, but I haven't tried.
