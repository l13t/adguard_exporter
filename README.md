# adguard_exporter

[![License Apache 2][badge-license]](LICENSE)

Fork from `povilasv/pihole_exporter`.

Modified to use the statistics AdGuardHome delivers.

## Overview

This Prometheus exporter checks your [AdGuard](https://github.com/AdguardTeam/AdGuardHome) statistics. Available metrics are:

- Average Responsetime
- DNS Queries
- Domains blocked

### Run docker container

```bash
docker run l13t/adguard_exporter:<tagname> -adguard http://<adguard_username>:<adguard_pwd>@<adguard_url>
```

## License

See [LICENSE](LICENSE) for the complete license.

## Changelog

A [changelog](ChangeLog.md) is available

[badge-license]: https://img.shields.io/badge/license-Apache2-green.svg?style=flat
