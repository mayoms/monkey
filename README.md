# Monkey Programming Language

**Summary** 

Monkey is a toy langauge interpreter, written in Go. It has C-style syntax, and is largely inspured by Ruby and Python

## Overview

This is a project that was inspired by Thorsten Ball's great book, [Writing and Interprer in Go](https://interpreterbook.com/). After completing the book I've gone on to extend the language with a number of additional features including:
  * Reading program from files, as well as 'include statements' to import files.
  * Modulo Operations, Structs and other language additions.
  * Robust string, array and hash methods
  * Expanded standard library including file I/O operations.

There are a number of tasks to complete, as well as a number of bugs. The purpose of this project was to dive deeper into Go, as well as get a better understanding of how programming languages work. It has been successful in those goals. There may or may not be continued work - I do plan on untangling a few messy spots, and there are a few features I'd like to see implemented. This will happen as time and interest allows.

## Installation
```
cd $GOPATH/src
git clone git@github.com:mayoms/monkey.git
cd ./monkey && go install
```

## Basic use
To access the REPL, simply run the following:

```
~ » monkey
Monkey programming language REPL

>>
```

or, to run a program:

```
monkey path/to/file
```

## Contributing

This project welcomes contributions from the community. Contributions are
accepted using GitHub pull requests; for more information, see 
[GitHub documentation - Creating a pull request](https://help.github.com/articles/creating-a-pull-request/).

For a good pull request, we ask you provide the following:

1. Include a clear description of your pull request in the description
   with the basic "what" and "why"s for the request.
2. The tests should pass as best as you can. GitHub will automatically run
   the tests as well, to act as a safety net.
3. The pull request should include tests for the change. A new feature should
   have tests for the new feature and bug fixes should include a test that fails
   without the corresponding code change and passes after they are applied.
   The command `npm run test-cov` will generate a `coverage/` folder that
   contains HTML pages of the code coverage, to better understand if everything
   you're adding is being tested.
4. If the pull request is a new feature, please include appropriate documentation 
   in the `README.md` file as well.
5. To help ensure that your code is similar in style to the existing code,
   run the command `npm run lint` and fix any displayed issues.

## Contributors

Micah Mayo

## License

MIT