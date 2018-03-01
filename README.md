# go-distributed

[![Build status](https://travis-ci.org/jessedearing/go-distributed.svg?branch=master)](https://travis-ci.org/jessedearing/go-distributed)

## Summary

Provides distributed primitives such as locks using a variety of databases

## Example

```
go-distributed lock --type mongo --db-connection "mongodb://localhost/mydb" && echo "run my job"
```
