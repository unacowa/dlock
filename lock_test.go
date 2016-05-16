package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"os"
	"testing"
	"time"
)

var (
	table dynamo.Table
)

const (
	region     = "ap-northeast-1"
	table_name = "dynamodb_go_lock"
	domain     = "test-domain"
)

func TestMain(m *testing.M) {
	db := dynamo.New(session.New(), &aws.Config{
		Region:     aws.String(region),
		DisableSSL: aws.Bool(true),
	})
	table = db.Table(table_name)
	exit := m.Run()
	if exit != 0 {
		os.Exit(exit)
	}
}

func TestAcquire(t *testing.T) {
	lock := NewLock(table, domain)
	defer lock.release()
	err := lock.acquire(1)
	if err != nil {
		t.Error("acquire expected to success, but ", err)
	}

	err = lock.acquire(1)
	if err == nil {
		t.Error("acquire twice expected to error, but ", err)
	}

	time.Sleep(2 * time.Second)
	err = lock.acquire(1)
	if err != nil {
		t.Error("acquire after expires expected to success, but ", err)
	}
}

func TestRelease(t *testing.T) {
	lock := NewLock(table, domain)
	lock.acquire(100)
	err := lock.release()
	if err != nil {
		t.Error("acquire then release expected to work, but ", err)
	}

	err = lock.release()
	if err != nil {
		t.Error("release before acquire expects to work, but", err)
	}

	lock.acquire(1)
	time.Sleep(2 * time.Second)
	err = lock.release()
	if err != nil {
		t.Error("acquire then expires and release expected to work, but ", err)
	}
}

func TestConcurrent(t *testing.T) {
	success := make(chan int, 10)
	for i := 0; i < 10; i++ {
		go func() {
			lock := NewLock(table, domain)
			defer lock.release()
			err := lock.acquire(1)
			if err == nil {
				success <- 0
			}
			time.Sleep(1 * time.Second)
		}()
	}

	time.Sleep(2 * time.Second)

	func() {
		lock := NewLock(table, domain)
		defer lock.release()
		err := lock.acquire(1)
		if err != nil {
			t.Error("Slow goroutine expect success, but ", nil)
		}
	}()

	if len(success) != 1 {
		t.Error("Only one transaction must be success, but ", len(success))
	}
}
