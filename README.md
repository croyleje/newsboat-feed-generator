# newsboat-feed-generator
A simple Go program to generate an _RSS_ feed from a _URL_.  It was originally
written to produce feeds for Rumble and YouTube channels but has been tested on
other sites that offer similar lists of content.  See the comments in
***main.go*** it's fairly self explanatory.

Due to changes in the way Rumble generates it's pages the _URLs_ have been
hard coded into the scraper so only the _channel name_ needs to be passed in
as an argument see examples below.

## Example Rumble Feeds
Rumble channels (feeds) have two basic formats either /user/channelName and
/c/channelName just pass the name of the channel via your _Newsboat URLs_
file and it will appear in your feeds list.

## Example Newsboat URLs
```
"exec:/path/to/file channelName"
```
