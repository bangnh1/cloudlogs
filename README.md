# cloudlogs
Retrieve logs from AWS Cloudwatch

## Local build

```
$ go build  -o cloudlogs cmd/main/main.go
```

## How to use:

- Provide credentials

```
export AWS_ACCESS_KEY_ID=<<YOUR_AWS_ACCESS_KEY_ID>>
export AWS_SECRET_ACCESS_KEY=<<YOUR_AWS_SECRET_ACCESS_KEY>>
export AWS_DEFAULT_REGION=<<YOUR_AWS_REGION>>
```

- AWS IAM Permission
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "logs:Describe*",
                "logs:Get*",
                "logs:List*",
                "logs:StartQuery",
                "logs:StopQuery",
                "logs:TestMetricFilter",
                "logs:FilterLogEvents"
            ],
            "Effect": "Allow",
            "Resource": "*"
        }
    ]
}
```

- Help 
```
$ ./cloudlogs get --help
$ ./cloudlogs list groups -help
```

- List all groups

```
$ ./cloudlogs list groups
```

- Filter logs groups by prefix

```
$ ./cloudlogs list groups --group "<<PREFIX>>"
```

- List all streams in a group

```
$ ./cloudlogs list streams --group <<GROUP_NAME>>
```

- Filter streams by prefix

```
$ ./cloudlogs list streams --group <<GROUP_NAME>> --prefix <<PREFIX>>
```

- Get logs events

```
$ ./cloudlogs get --group <<GROUP_NAME>> --start "1day ago"
```

- Human-friendly time filtering

  - `--start='13/5/2021 01:23'`
  - `--start='1h ago'`
  - `--start='3d ago'`
  - `--start='6w ago'`

* Filter with pattern

```
$ ./cloudlogs get --group <<GROUP_NAME>> --start "1day ago" --filter "404"
```

- Batch mode - write to file in concurrent

```
$ ./cloudlogs get --group <<GROUP_NAME>> --start "4 weeks ago" --file --interval "7 days"
```

- Watch mode - Get logs in real time

```
$ ./cloudlogs get --group <<GROUP_NAME>> --start "4 weeks ago" --watch
```
