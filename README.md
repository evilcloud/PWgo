![](assets/background.jpg)

# PWgo

#### Go implementation of PW pass gen

## What is it?

This is a password generator, based somewhat on principles of xkcd password strength comic

![](assets/password_strength.png)

[link](https://xkcd.com/936/) to the original comic

It generates passwords, usernames, and for fun, two word combos, that may or may not mean something.

This was a personal project, later filled with features requested by a close person. In my view this is by far the best password generator I have personally seen. Why? Because once in a while we are forced to enter the passwords manually, despite our password managers placing them into the clipboard, and when this starts, I get extremely annoyed to constantly shift attention between my password manager and the input field, grabbing only couple of random alphanumeric characters at a time. Here once or two views is sufficient to retype a 40 charachter long string. All this without sacrifising the safety.

### So, I don't need a password manager then?

It's 2020 and anyone who is not making a full use of password manager right now may as well post their SSN and credit card number as the banner on their Facebook. Of course you need a password manager, but you also can't reuse the passwords. This is why you need a password generator. Making the passwords yourself will only make you repeat the same guessable strings of data.

Of course password managers offer password generatiors, but I personally find them pretty dumb and often uncomfortable to use, especially when the password needs to be copied manually. It's amazing how often this is necessary.

## Real words wut?

The app uses three dictionaries:

**Adjectives**: is a list of English language adjectives. Taken from [hugsy](https://gist.github.com/hugsy/8910dc78d208e40de42deb29e62df913)

**Nouns** is a 6'775 entries strong list of English language nouns. Taken from [The Great Noun List](http://www.desiquintans.com/nounlist)

**Profanities** is a list of profanities, heavily edited -- I have tried as much as I could to remove prejudices and leave only insults. Many have slipped though and stayed. The list is in the state of cleaning. In any case, I advise against using profanities option.



## How to use it

The app sits in the menubar and once the icon is clicked, it will generate the following data:

1. Username -- consists of a merged adjective and a noun. All spaces and punctuation are omitted. Words are capitalised and then concatenated
2. Password -- consists of 4, 6 or 8 words, one special character and one number (excluding zero, to avoid confusion). The words in the password are random, capitalized. The numeral and special character are set at the end of the password as default (can be changed in the settings)
3. Words of Wisdom -- just a fun experiment where a random adjective and noun are juxtaposed

All items mentioned above are copied into the clipboard when clicked on.

Password is automatically added to the clipboard once the menu icon is pressed.

The app is written in pure Go. Compiled and tested on macOS. I assume it can be compiled and run on Windows and Linux as well, but I haven't tried.

## How to use it?

![](assets/Screenshot.png)

### Last copy-clicked

This menu shows you your last copy-clicked username and password (as well as indicating the time passed since the click). Can be very useful, depending on your workflow.

Username or password will be copied to clipboard when clicked. The data will stay unchanged until a new username / password respectively is clicked, or app restarts.

### Words of Wisdom

A fun juxtaposition of an adjective and noun. Can be copied into clipboard, but isn't saved anywhere.

### Settings

1. **Length** defines the length of the password in words
2. **Additional security** will place the number and special charachter randomly between words in the password, instead of placing them at the end. Arguably this is somewhat safer, but I personally see no huge upside.
3. **NSFW** adds a dictionary of insults and offensive words to the existing dictionary. Theoretically this strengthens the password strength by adding into the pool of predetermined words, but some of the terms truly shouldn't be there. Use at your own risk, or better don't use at all.
4. **Sailor-redneck mode** weakest, most immature mode of username and password generation -- this option removes adjectives and nouns dictionaries and leaves only profanities. Done purely for novelty purposes. Not advisable to use at all. Ever.


## What's next?

*Nothing*. Hisorically the app in its current form has appeared as a port, and later improvemnt of its original Python version. I guess there will be some minor updates, but essentially the app will be rewritten in Swift or Objective-C some time soon. Go is beautiful, but it has its limitations with these type of applications.

Currently the app is very stable, so mostly the updates will fall in either "optimisation" or "additional features" camps.
