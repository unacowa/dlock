
## NAME

dlock - Easy to use, Distributed lock manager, written in Go.

[![Build Status](https://travis-ci.org/unacowa/dlock.svg?branch=master)](https://travis-ci.org/unacowa/dlock)

## Description

'dlock' provides command line tool to use for acquire distributed lock.

Focused on..

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

## Install (on linux amd64)

```
$ curl -L https://github.com/unacowa/dlock/releases/download/pre-release/linux_amd64_dlock -o dlock
$ chmod +x dlock
$ sudo cp dlock /usr/local/bin
```

## Example

### Step 1: setup credential in ENV.
```
$ export AWS_ACCESS_KEY_ID=AKIxxxxxxxxxxxxxxxxx
$ export AWS_SECRET_ACCESS_KEY=Kxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Step 2: create table. (Maybe wait a seconds for created) 
```
$ ./dlock --table dynamodb_go_lock create-table
```

### Step 3: run
```
$ ./dlock --table dynamodb_go_lock run --domain example bash -c 'sleep 3; echo One' & ./dlock --table dynamodb_go_lock run --domain example bash -c 'sleep 3; echo Two'
failed acquire. transaction:fea09f48-c94b-4cb0-8b82-2479dee26564 running and not expires
One
```

## Credential option
Uses `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` for access to AWS DynamoDB.

Both are passed from `ENVIRONMENT` or `GLOBAL OPTION`.


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
