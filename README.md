# ted

**Test Execution Dashboard** - a facility to collect &amp; present test results, and provide a list of tests to be executed (TODO).

Demo version : https://arcane-ravine-69473.herokuapp.com/  
This is very much a WIP!

# TODO :

- Provide a list of tests (that can be dynamically updated) that a test runner can query, so that the test runner's next test(s) can be controlled
- UI dashboard for test results
- Store test results in file/DB

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
