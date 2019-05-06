# envly

[![Build Status](https://travis-ci.org/jwilner/envly.svg?branch=master)](https://travis-ci.org/jwilner/envly)

Drop in binary that loads environment variables from different locations. 

Intended as a simple secrets management solution within Docker.

S3:
```bash
ENVLY_URI=s3://my-bucket/path/to/envfile.env envly my-command
```

GS:
```bash
ENVLY_URI=gs://my-bucket/path/to/envfile.env envly my-command
```

HTTP/HTTPS:
```bash
ENVLY_URI=https://example.com/path/to/envfile.env envly my-command
```

File:
```bash
ENVLY_URI=file:///tmp/my-testing-env.env envly my-command
```
