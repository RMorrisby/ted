# ted

**Test Execution Dashboard** - a facility to collect &amp; present test results, and provide a list of tests to be executed (TODO).

Demo version : https://arcane-ravine-69473.herokuapp.com/  
This is very much a WIP!

# TODO :

TODO Data page - add links for previous test runs

TODO Data page - add link to all results (most recent X (e.g. 100)) regardless of test run

TODO Data page - remove TED Status column - colour-code TED Notes

TODO Data page - auto-scroll page down when received new result

# How to use

## Is TED alive?

`GET <hostname>/is-alive`

## Register a test

```
POST <hostname>/test
{


}
```

## Register a test suite

```
POST <hostname>/suite
{


}
```

## Send in a test result

```
POST <hostname>/result
{


}

OR send in a placeholder, which will be updated later
POST <hostname>/result
{

}
```

## Claim a test result, so that other concurretnly-running test agents don't try to run the test

```
POST <hostname>/claimtest
{




}
```

## Update a test result

```
PUT <hostname>/result
{


}
```

## Get latest test run ID

```
GET <hostname>/testrunid/latest
```

## Get next test run ID

If testrun IDs are being incremented, what is the next testrun ID after the latest one?

```
GET <hostname>/testrunid/next
```

## Get all failed & not-run tests in a test run

```
GET <hostname>/reruns?testrun=<testrun>
```
