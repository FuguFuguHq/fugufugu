![FuguFugu](https://github.com/FuguFuguHq/fugufugu/blob/main/Logo.png?raw=true)

# FuguFugu

FuguFugu is a tool to check your website for external scripts and images. With a large website and many developers it's 
often unclear what external resources your website uses. If you're in the EU, each external usage needs to be documented
in the privacy statement. Each external resource can log the user IP and could - depending on the browser - set cookies.
Each external script can undermine the security of your user.

Fugu is a fish that has some very toxic parts. Just like a nice website where only one bad script can be toxic and kill you.

## Build

Build FuguFugu:

`go build`

should create a `fugufugu` executable.

## Run 

To check [https://www.amazingcto.com](https://www.amazingcto.com) (my site ;-) for external scripts:

`./fugufugu -url https://www.amazingcto.com`

prints a report of all external scripts and images it found.

```
+--------------------------------+------------------+---------+--------+-------+-----+--------+
| SITE                           | COMPANY          | COUNTRY | SCRIPT | IMAGE | CSS | COOKIE |
+--------------------------------+------------------+---------+--------+-------+-----+--------+
| t5972a59c.emailsys1a.net       | rapidmail        | EU      | Yes    |       |     |        |
| scripts.simpleanalyticscdn.com | Simple Analytics | EU      | Yes    |       |     |        |
| datenschutz-generator.de       |                  |         |        | Yes   |     |        |
+--------------------------------+------------------+---------+--------+-------+-----+--------+
```

**Only use on your own website!**