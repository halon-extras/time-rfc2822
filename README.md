# RFC2822 time plugin

This plugin allows for creating a time in RFC2822 format based on an optional unix time and location (time zone).

## Installation

Follow the [instructions](https://docs.halon.io/manual/comp_install.html#installation) in our manual to add our package repository and then run the below command.

### Ubuntu

```
apt-get install halon-extras-time-rfc2822
```

### RHEL

```
yum install halon-extras-time-rfc2822
```

## Exported functions

### time_rfc2822([unixtime, location])

Create a time in RFC2822 format. An optional unix time can be provided, otherwise the current time will be used. An optional location (time zone) can be provided, otherwise the current location will be used.

**Params**

- unixtime `number` - The unix time (optional)
- location `string` - The location (optional)

**Returns**

On success it will return a string that contains the time in RFC2822 format. On error an exception will be thrown.

**Example**

```
import { time_rfc2822 } from "extras://time-rfc2822";
echo time_rfc2822();
echo time_rfc2822(1645012689.519825);
echo time_rfc2822(time(), "Europe/Stockholm");
```