[![Build Status](https://travis-ci.org/mantel-digio/covermate.svg?branch=master)](https://travis-ci.org/mantel-digio/covermate)

# Covermate

Covermate is a tool to assist with managing code-coverage metrics.  Go provides builtin tools for assessing code test-coverage.  Covermate extends the capabilities of the Go test-coverage report by allowing coverage-omissions to be flagged, or tagged as exclusions. This allows untested code to be tracked and managed.

Covermate can be integrated into a CI pipeline to verify that untested blocks of code are specifically excluded, and also that the overall coverage meets a specified threshold.

Covermate will exit with non-zero if:
- There is a block of code which is not tested and does not have a `// nocover` comment.
- The overall code coverage does not meet a specified threshold.

*Note:* As per -coverprofile, packages with no tests are ignored.

*Note:* As per -coverprofile, coverage is package-specific.  i.e. code is only considered covered, if it is called from a test within the same package.

## Installion

```bash
go install github.com/mantel-digio/covermate
```

## Usage
```
usage: covermate [<flags>]

Flags:
      --help                     Show context-sensitive help (also try --help-long and --help-man).
  -f, --filename="coverage.out"  coverage report from go test
  -t, --tag="nocover"            comment tag to exclude blocks from mandatory coverage
  -T, --threshold=-1             minimum required overall coverage
```

Covermate depends on the coverage report generated by `go test` with the  `-coverprofile` option.

To execute the commands together:
```bash
go test ./... -coverprofile=coverage.out && covermate
```

By default covermate will not assess the overall level of code coverage.  A threshold can be set by 
providing a value between 0 and 100 with the -T option.  If a theshold is set, and the total coverage is
below the specified level, covermate will exit with an error.

## Example

If a block of code has no test coverage, it must include a specific tag in the comments (`nocover` by default), or it will generate an error.

```go
    // Depending on the nature of `data`, it may not be possible to make `json.Marshal()` return an error.
    data, err := createData()
    if err != nil {
        return err // if this block is not tested, covermate will exit with error
    }
    b, err := json.Marshal(data)
    if err != nil {
        return err // nocover - this block is tagged as being excluded from mandatory code coverage
    }
```