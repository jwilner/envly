# envly

[![Build Status](https://travis-ci.org/jwilner/envly.svg?branch=master)](https://travis-ci.org/jwilner/envly)

Drop in binary that loads environment variables from different locations.

Intended as a simple secrets management solution within Docker; drop it in as an entrypoint, and your container will load the secrets from `ENVLY_URI` if defined. Supports multiple different backends; your container remains agnostic while you move it between deployments or platforms.

Plain bash:
```bash
echo "ENV_VAR=abcdef" >> my.env  # create an env file
aws s3 cp --sse AES256 my.env s3://envfiles/my.env  # store file encrypted in S3
export ENVLY_URI=s3://envfiles/my.env  # export URI
envly env | grep ENV_VAR  # loads ENV_VAR=abcdef from S3
```

Or within a Docker:
```bash
cat Dockerfile
# FROM debian
#
# ADD https://github.com/jwilner/envly/releases/download/v0.0.1/envly-s3-linux-amd64 /usr/local/bin/envly
#
# ENTRYPOINT envly
docker build -t envly-example .
docker run -e ENVLY_URI=s3://envfiles/my.env envly-example env
```

## Backends:

- S3: `ENVLY_URI=s3://my-bucket/path/to/envfile.env`
- Google Cloud Storage: `ENVLY_URI=gs://my-bucket/path/to/envfile.env`
- HTTP/HTTPS: `ENVLY_URI=https://example.com/path/to/envfile.env`
- File: `ENVLY_URI=file:///tmp/my-testing-env.env`
