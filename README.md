
DynamoDB lock Go (dlock)
==================

Easy to use, Distributed lock manager, written in Go.

## Description

Command line tool, and Go library to use for acquire distributed lock.

Focus on..

* Easy to use. One binary, No deploy, Only cloud.
* Configuration free.
* Almost ZERO pricing. (based on DynamoDB free tier)
* Scalable in the future.

## Motivation
Startup projects, microserivces or etc, Some times we write a small batch program.
In essence, we know it needs `lock`. But most case, it is enough lite and no time to make distributed lock.
In early release, it will not causes problem, because programs takes enough short time. But it causes problem when we forget.
So, I write `Easy to Use` distributed lock.

## Requirement
AWS account only.

## Usage
```
$ ./dlock
NAME:
   dlock - acquire 'lock' to following commands.

USAGE:
   dlock [global options] command [command options] [arguments...]
   
VERSION:
   0.0.0
   
COMMANDS:
     create-table	Create table then exit
     run		run command in distributed lock

GLOBAL OPTIONS:
   --region value, -r value		Specify Region name (default: "ap-northeast-1")
   --table value, -t value		Specify Table name
   --aws_access_key_id value		 [$AWS_ACCESS_KEY_ID]
   --aws_secret_access_key value	 [$AWS_SECRET_ACCESS_KEY]
   --help, -h				show help
   --version, -v			print the version
```

## Example
```
# first create database table.
$ ./lock --table table-name create-table
# run command with distributed lock.
$ ./lock --table table-name run --domain my-batch my-important-exec --some=any
$ ./lock --table table-name run --domain my-batch bash -c 'sleep 10; echo 10' 
```

## Credential
Uses `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` for access to AWS DynamoDB.
Both are passed from `ENV` or `GLOBAL OPTION`.


## Contribution

### Prepare container for library.

```
docker run -itd --name dynamodb_lock_go-gopath -v /go busybox
```

### Install depandancy.
```
docker run --rm --volumes-from dynamodb_lock_go-gopath -v $PWD:/go/src/app -w /go/src/app golang:1.6 go-wrapper download
```

### Build
```
docker run --rm --volumes-from dynamodb_lock_go-gopath -v $PWD:/go/src/app -w /go/src/app golang:1.6 go build -o 'dlock'
```

### Test
Test uses 'REAL' dynamoDB.
Create database as following, Before run.

```
docker run --rm --volumes-from dynamodb_lock_go-gopath -v $PWD:/go/src/app -w /go/src/app -e AWS_ACCESS_KEY_ID=AKIxxxxxxxxx -e AWS_SECRET_ACCESS_KEY=KVxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx golang:1.6 go run lock.go -t dynamodb_go_lock create-table
```

```
docker run --rm --volumes-from dynamodb_lock_go-gopath -v $PWD:/go/src/app -w /go/src/app -e AWS_ACCESS_KEY_ID=AKIxxxxxxxxx -e AWS_SECRET_ACCESS_KEY=KVxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx golang:1.6 go test .
```

## Licence

[MIT](https://github.com/tcnksm/tool/blob/master/LICENCE)
