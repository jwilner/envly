# envly

Drop in binary that loads environment variables from different locations.

S3:
```bash
ENVLY_URI=s3://my-bucket/path/to/envfile.env envly my-command
```

File:
```bash
ENVLY_URI=file:///tmp/my-testing-env.env envly my-command
```
