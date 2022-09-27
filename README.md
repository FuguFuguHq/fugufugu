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

prints a report of privacy pages and all external scripts and images it found.

```
+--------------+----------------------+
| PRIVACY PAGE | TITLE                |
+--------------+----------------------+
| /privacy/    | Datenschutzerklärung |
+--------------+----------------------+
+--------------------------------+------------------+------------+---------+--------+-------+-----+
| SITE                           | COMPANY          | PRODUCT    | COUNTRY | SCRIPT | IMAGE | CSS |
+--------------------------------+------------------+------------+---------+--------+-------+-----+
| scripts.simpleanalyticscdn.com | Simple Analytics | Analytics  | EU      | Yes    |       |     |
| t5972a59c.emailsys1a.net       | rapidmail        | Newsletter | EU      | Yes    |       |     |
+--------------------------------+------------------+------------+---------+--------+-------+-----+
Summary https://www.amazingcto.com: 38 pages - External resources: 2 scripts | 0 images | 0 css
```

## Verbose

Verbose mode

`./fugufugu -url https://www.amazingcto.com -verbose`

will print what fugufugu is currently doing.

## Cookies

FuguFugu will not check for cookies in resources by default to speed up checking.

`-cookie` will enable cookie checking

`./fugufugu -url https://www.amazingcto.com -cookie`

## Max Pages

FuguFugu will by default check 10.000 pages. `-max` sets a new maximum for pages.

This checks only 10 pages:

`./fugufugu -url https://www.amazingcto.com -max 10`
