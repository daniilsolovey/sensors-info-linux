# sensors-info-linux

Small Linux desktop utility that displays current system information in a desktop notification.

## Information displayed

- local time from the current Linux system timezone;
- current date;
- Moscow time;
- Wi-Fi network name;
- network ping;
- VPN status;
- CPU temperature;
- average and maximum instantaneous CPU frequency across logical CPUs;
- used, total and available RAM;
- RAM usage percentage;
- battery percentage;
- battery charging/discharging state.

## Example

![sensors-info-linux notification](docs/sensors-info-linux-notification.png)

## CPU frequency

The utility reads the current CPU frequencies from `/proc/cpuinfo`.

`AVG` is the average instantaneous frequency across all logical CPUs.

`MAX` is the highest instantaneous frequency among all logical CPUs during the current `/proc/cpuinfo` read. It is not the maximum frequency observed over a historical time interval.

Example:

```text
Frequency   AVG 1150 MHz / MAX 2100 MHz
```

## RAM usage

RAM usage is calculated using the amount of memory currently available to applications:

```text
used = total - available
```

The notification shows used memory, total memory, usage percentage and available memory.

Example:

```text
RAM         5.9 / 31.0 GB (19%)
Available   25.1 GB
```

## Local time

Local time is obtained using Go's `time.Now()`, so it follows the timezone currently configured in Linux.

The current timezone can be checked with:

```bash
timedatectl
```

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
