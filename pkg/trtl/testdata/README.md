# testdata

The `testdata` directory contains the required fixtures for running the trtl test suite. There are two files which must be maintained here.

`db.json`

The "source of truth" which declares the objects that should exist in the trtl database used by the tests. This file is used to both regenerate the test database and verify the results of operations against the test database. The structure of the file is a JSON map where each key is an identifier the tests can reference and each value represents the object stored in the database. E.g.,

```json
{
    "alice": {
        "namespace": "people",
        "key": "45abc34d-f8f9-4f0e-b8e3-f8f9f8f8f8f8",
        "value": {
            "name": "Alice",
            "email": "alice@example.com"
        }
    },
}
```
Defines a trtl database consisting of one entry in the namespace `people` with the key `45abc34d-f8f9-4f0e-b8e3-f8f9f8f8f8f8` and the JSON object:

```json
{
    "name": "Alice",
    "email": "alice@example.com"
}
```

The `value` field is a JSON map which can contain any kind of JSON-compatible fields. It's a generic representation of a data object which the test suite marshals into binary data when generating the test database.

`db.tar.gz`

The compressed representation of the trtl test database. Before running the tests, the test suite unzips this file into a temporary directory in `testdata` so that the tests have access to a fresh database. If the file does exist, then it is automatically generated from `db.json`. The test runner can also force an update to the gzipped file by supplying the `-update` flag (e.g., `go test -update`), which should be done every time the `db.json` file is updated.

An important consideration is that the test suite wraps the defined objects in `db.json` into Honu objects when generating the database. This means that an update to a newer version of Honu may introduce breaking changes and `db.tar.gz` should be regenerated in that case.