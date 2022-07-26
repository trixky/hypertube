# Go Tests

This tester run all go unit tests with a [MageFile](https://magefile.org/).

## Run

You need to have mage installed locally:

```bash
go install github.com/magefile/mage@latest
```

You can then run the tests for each APIs:

```bash
mage TestApiAuth
mage TestApiUser
mage TestApiMedia
```
