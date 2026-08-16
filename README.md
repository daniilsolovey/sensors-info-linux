# sensors-info-linux

Small Linux desktop utility that shows current system information in a desktop notification.

## Information displayed

The notification is grouped into four sections.

**TIME**

- local time from the current Linux system timezone;
- current date;
- Moscow time.

**NETWORK**

- Wi-Fi network name;
- network ping;
- VPN status.

**SYSTEM**

- CPU temperature;
- average and maximum instantaneous CPU frequency across logical CPUs;
- used RAM, total RAM, and usage percentage.

**POWER**

- battery percentage;
- battery charging/discharging state.

## Example

![sensors-info-linux notification](docs/sensors-info-linux-notification.png)

## Requirements

- Linux
- `notify-send` (libnotify)
- a notification daemon, for example dunst
- `nmcli` (NetworkManager) for VPN status
- lm-sensors data for CPU temperature (`coretemp`)

## Build and run

```bash
go build -o sensors-info-linux
./sensors-info-linux
```

The notification is titled `System info` and closes after 5 seconds.

## CPU frequency

The utility reads current CPU frequencies from `/proc/cpuinfo` (`cpu MHz` lines).

`avg` is the average instantaneous frequency across all logical CPUs.

`max` is the highest instantaneous frequency among all logical CPUs during the current `/proc/cpuinfo` read. It is not the maximum frequency observed over a historical time interval.

Example:

```text
Frequency  1150(avg MHz) / 2100(max MHz)
```

## RAM usage

RAM usage is calculated using the amount of memory currently available to applications:

```text
used = total - available
```

The notification shows usage percentage, used memory, and total memory.

Example:

```text
RAM  (19%) 5.9(used GB) / 31.0(total GB)
```

## CPU temperature

CPU temperature is read from lm-sensors chip `coretemp-isa-0000`, sensor `Core 0`.

## Local time

Local time is obtained using Go's `time.Now()`, so it follows the timezone currently configured in Linux. Time is shown as `HH:MM`.

The current timezone can be checked with:

```bash
timedatectl
```

Moscow time is `Europe/Moscow`.

## Network

Wi-Fi name comes from the current wireless interface. If no network is found, the notification shows `disconnected`.

Ping is a single ICMP request to `www.google.com` and is shown in milliseconds. ICMP ping usually needs `CAP_NET_RAW` or root.

VPN status is taken from NetworkManager:

```bash
nmcli con show --active
```

The status is `Connected` if an active connection contains `tun0`, otherwise `Disconnected`.

## Desktop notifications

Notifications are displayed using `notify-send` and are intended to be used with a notification daemon such as dunst.

For a larger system-info notification, a dunst configuration such as the following can be used:

```ini
font = Ubuntu Mono 13
width = 420
height = (0, 700)
```

Reload dunst after changing its configuration:

```bash
dunstctl reload
```
