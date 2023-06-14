# rocketchat-term 

rocketchat-term is a minimal, lightweight terminal interface for reading rocketchat messages.
It streams incoming messages in realtime from every room/channel a user is a member of into one continous feed.

## Usage

The binaries found in bin can be executed directly without any commandline flags, or the project can be manually compiled with a go compiler.
It is best to put the binary in a place specified in your `$PATH` environment variable

## Authentication 

Authentication can be done via a standard username and password, an LDAP username and password or via a token.
For token authentication, see the example config in the config directory.
All connections are sent with TLS encryption.

rocketchat-term caches tokens it recieves in `~/.rocketchat-term/`

## Configuration

A configuration file isn't needed to run, but a configuration file can configure connection options, custom colouring, spacing and logging.
The config file should be a yaml file with path `~/.rocketchat-term/rocketchat-term.yaml` if no config is found then it uses defaults and asks for connection options.
See the config directory for an example 

## Colouring

By default, rocketchat-term uses Ansi-256 colouring as this is widely supported across terminals.
Full RGB colouring is supported for more modern terminals and can be specified in the config with hexcodes. 
Custom colouring is also supported for Ansi-256 colours.

More terminal colouring information can be found here [](https://en.wikipedia.org/wiki/ANSI_escape_code)
