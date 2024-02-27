# Dead Switch

This is a simple dead switch written in Go. You basically just set the parameters in the command-line and it will load a site where you can enter a password every X days to verify yourself.

The deadswitch will send out a warning at halfway and 3/4 to remind the user to reset the switch. It will also send an email every time the timer is reset.

## Usage

You will need a Gmail account with 2FA and an app password. [This tutorial](https://support.google.com/mail/answer/185833?hl=en) will show you how to get this set up.

You will need Golang to compile: go build

And then you just have to run the program with the parameters. For example:

./deadswitch -days 180 -owner youremail@gmail.com -message "If you're receiving this message, I may be in trouble - here is some information" -recipient person1@examplemail.com -recipient mum@dontbesad.com -key "your gmail application password" -auth myverificationpassword

**The order of flags is not important**

By default, the program will start a server on 127.0.0.1:3451 - this will point to an auth page where you can enter your verification password to reset the dead switch/prevent it from triggering. If the user fails to reset the dead switch by the time limit, it will send out the deadswitch message (as long as credentials are correct and server is still running).

It's recommended that you try out a short test dead switch first (eg. 1 day) before setting a real one to make sure everything is running well. 

It currently saves a file for persistence - **this file has your auth/key in plaintext** so be careful with it. The file basically enables running the deadswitch without flags - you can set a script to automatically run ./deadswitch on boot to make it resistant to resets.

## Flags
-days - Number of days before the dead switch triggers 

-owner - The 'from' email

-message - The message - this must with enclosed within quotes

-recipient - Add recipients - you can add as many as needed

-key - Your Google application password 

-auth - Set a password to verify yourself

-file - Add a file (must be in same directory)

-port - Change server port (default is 3451)