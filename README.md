# go-quiz

A simple quiz app in Go.

## Features

- Takes an input CSV, reads questions from it and prompts user for answers.
- Tally ups the score for the correct answers.

## Usage

```sh
go run .
```

## Example

```
> cat ./data/problems.csv | head

question,answer
5+5,10
7+3,10
1+1,2

› go run .

Enter file path [./data/problems.csv]:
Using filepath: ./data/problems.csv
Number of records: 5
5+5?
10
7+3?
10
...
You got 5 (100.0%) correct!
```
