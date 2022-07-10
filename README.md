# Solution

## Part 1

This part completes the verifaction step that can be run as follows.

Generate files first:

```sh
go run generate/main.go -out data/messages.1.data -verify data/verify.1.csv --seed 1560981440 -count 20
go run generate/main.go -out data/messages.2.data -verify data/verify.2.csv --seed 1560980000 -count 10000 -attrs 20 -events 300000 -maxevents 500 -dupes 10
```

Then verify:

```sh
go run verify/main.go --verify-file ./data/verify.1.csv
go run verify/main.go --verify-file ./data/verify.2.csv
```

By default `/main.go` uses `var fileName = "data/messages.1.data"`, change this line to verify against a different `messages` file.

This part also includes all the back-end code necessary for the endpoints.
