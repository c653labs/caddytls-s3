# caddytls S3
**Work in progress**

[Caddy](https://caddyserver.com) TLS plugin that will use [S3](https://aws.amazon.com/s3/) for storing certificate files.

This plugin is useful when you either want to persist certificate files outside of the operating systems disk or share certificates between multiple Caddy servers.

## Environment variables
- CADDY_S3_BUCKET
- CADDY_S3_REGION
