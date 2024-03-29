# Summary

This was a takehome test in which I spent around 20 hours and got nothing in return from the entity that sent it to me.
To look at the original problem write-up please look at [./README_PROBLEM.md](https://github.com/Sleepful/customer-data-test/blob/main/README_PROBLEM.md)

To look at progress, check out the merged PULL REQUESTS.

To look at other code challenges, look here: https://github.com/Sleepful/CodingChallenges

# Solution

## Part 1

This part completes the verification step that can be run as follows.

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

## Part 2

Completes the following with TailwindCSS styling:

- Displays the list of users with pagination
- Displays attributes for a user
- Uses EmberJS utilities for doing the operations mentioned above

to run this do:

```
cd ./ember-quickstart
npm install
ember serve
```

and then visit `http://localhost:4200` on your browser

## Missing

This second part is missing a couple of things, namely:

- edit/delete/create for records
