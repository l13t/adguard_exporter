# adguard_exporter

[![License Apache 2][badge-license]](LICENSE)

Fork from `povilasv/pihole_exporter`.

Modified to use the statistics AdGuardHome delivers.

# Overview 

This Prometheus exporter checks your [AdGuard](https://github.com/AdguardTeam/AdGuardHome) statistics. Available metrics are:

-   Average Responsetime
-   DNS Queries
-   Domains blocked

## Docker Deployment

-   Build Image:

    docker build -t adguard-exporter .

-   Start Container

    docker run -d -p 9311:9311 adguard-exporter -adguard http://192.168.1.5

## License

See [LICENSE](LICENSE) for the complete license.

## Changelog

A [changelog](ChangeLog.md) is available

[badge-license]: https://img.shields.io/badge/license-Apache2-green.svg?style=flat
