# ted

**Test Execution Dashboard** - a facility to collect &amp; present test results, and provide a list of tests to be executed (TODO).

Demo version : https://arcane-ravine-69473.herokuapp.com/  
This is very much a WIP!

# How to use

## Is TED alive?

`GET <hostname>/is-alive`

```
200
{"is-alive": true}
```

## Is TED paused?

`GET <hostname>/pause`

```
200
false
```

## Register a test

```
POST <hostname>/test
{
  "Name": "test_1_foo_bar",
  "Dir": "read_ABC",
  "Priority": 4,
  "Categories": "foo|bar",
  "Description": "Basic test to read from the ABC interface",
  "Notes": "",
  "Owner": "XYZ",
  "IsKnownIssue": false,
  "KnownIssueDescription": ""
}
```

```
201
```

## Register a test suite

```
POST <hostname>/suite
{
  "Name": "Suite ABC",
  "Description": "Suite for the ABC interface",
  "Owner": "XYZ",
  "Notes": "",
}
```

```
201
```

## Send in a test result

```
POST <hostname>/result
{
  "SuiteName": "Suite ABC",
  "TestName": "test_1_foo_bar",
  "TestRunIdentifier": "0.1.15",
  "Status": "FAILED",
  "StartTimestamp": "2022-01-01 07:01:37 UTC",
  "EndTimestamp": "2022-01-01 07:01:59 UTC",
  "RanBy": "XYZ",
  "Overwrite": false
}

OR send in a placeholder, which will be updated later
POST <hostname>/result
{
  "SuiteName": "Suite ABC",
  "TestName": "test_1_foo_bar",
  "TestRunIdentifier": "0.1.15",
  "Status": "NOT RUN",
  "StartTimestamp": null,
  "EndTimestamp": null,
  "RanBy": null
}
```

```
201    (if the result was accepted)
```

```
400    (if the result was rejected because a result exists for this test & test-run )
"Result received on POST, but there was an existing result in the DB for test test_1_foo_bar for testrun 0.1.15"
```

## Claim a test result, so that other concurrently-running test agents don't try to run the test

```
POST <hostname>/claimtest
{
  "TestName": "test_1_foo_bar",
  "TestRunIdentifier": "0.1.15",
  "IsRerun": false
}
```

```
200
true
```

## Update a test result

```
PUT <hostname>/result
{
  "SuiteName": "Suite ABC",
  "TestName": "test_1_foo_bar",
  "TestRunIdentifier": "0.1.15",
  "Status": "FAILED",
  "StartTimestamp": "2022-01-01 07:01:37 UTC",
  "EndTimestamp": "2022-01-01 07:01:59 UTC",
  "RanBy": "XYZ",
  "Overwrite": true
}
```

```
200
```

## Get latest test run ID

```
GET <hostname>/testrunid/latest
```

```
200
0.1.14
```

## Get next test run ID

If testrun IDs are being incremented, what is the next testrun ID after the latest one?

```
GET <hostname>/testrunid/next
```

```
200
0.1.15
```

## Get all failed & not-run tests in a test run

```
GET <hostname>/reruns?testrun=<testrun>
```

```
200
[
  {
  "Categories": "foo|bar",
  "Dir": "read_ABC",
  "TestName": "test_1_foo_bar",
  "TestRunIdentifier": "0.1.15",
  "Status": "FAILED",
  "Priority": 4,
  "StartTimestamp": "2022-01-01 07:01:37 UTC",
  "EndTimestamp": "2022-01-01 07:01:59 UTC",
  "RanBy": "XYZ",
  "TedStatus": "FAILED"
  },
  {
  "Categories": "foo|bar",
  "Dir": "read_ABC",
  "TestName": "test_2_foo_bar_baz",
  "TestRunIdentifier": "0.1.15",
  "Status": "FAILED",
  "Priority": 3,
  "StartTimestamp": "2022-01-01 07:02:01 UTC",
  "EndTimestamp": "2022-01-01 07:02:20 UTC",
  "RanBy": "XYZ",
  "TedStatus": "FAILED"
  }
]
```

## Get testrun stats

```
GET <hostname>/stats?testrun=<testrun>
```

```
200
{
  "TestRunName": "0.1.15",
  "Total": 3,
  "Passed": 1,
  "PassedOnRerun": 0,
  "Failed": 2,
  "KnownIssue": 0,
  "NotRun": 0,
  "LastRun": ""
}
```

# TODO :

TODO Data page - add links for previous test runs

TODO Data page - add link to all results (most recent X (e.g. 100)) regardless of test run

TODO Data page - remove TED Status column - colour-code TED Notes

TODO Data page - auto-scroll page down when received new result
