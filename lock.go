package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/satori/go.uuid"
	"gopkg.in/codegangsta/cli.v2"
	"os"
	"os/exec"
	"time"
)

type Record struct {
	Hkey          string `dynamo:"Hkey,hash"`
	TransactionId string
	Expires       int
}

type Lock struct {
	Table         dynamo.Table
	TransactionId string
	Hkey          string
}

func (record Record) isExpire() bool {
	return record.Expires < int(time.Now().Unix())
}

func isConditionalCheckErr(err error) bool {
	if ae, ok := err.(awserr.RequestFailure); ok {
		return ae.Code() == "ConditionalCheckFailedException"
	}
	return false
}

func (lock Lock) acquire(expires int) error {
	record := Record{
		Hkey:          lock.Hkey,
		TransactionId: lock.TransactionId,
		Expires:       int(time.Now().Unix()) + expires}
	err := lock.Table.Put(record).
		If("attribute_not_exists(Hkey)").
		Run()
	// if attribute exists.
	// Check it's expires time, and begin acquire if expired.
	if !isConditionalCheckErr(err) {
		return err
	}
	lock.Table.Get("Hkey", lock.Hkey).One(&record)
	if record.isExpire() {
		err := lock.Table.Put(record).Run()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("failed acquire. transaction:%s running and not expires", record.TransactionId)
	}
	return nil
}

func (lock Lock) release() error {
	err := lock.Table.Delete("Hkey", lock.Hkey).
		If("begins_with(TransactionId, ?)", lock.TransactionId).
		Run()
	// if same transaction id is not found.
	// Other transaction begin, because this transaction expires.
	if !isConditionalCheckErr(err) {
		return err
	}
	return nil
}

func NewLock(table dynamo.Table, hkey string) Lock {
	return Lock{
		Table:         table,
		TransactionId: uuid.NewV4().String(),
		Hkey:          hkey,
	}
}

func main() {
	var region string
	var table string
	var domain string
	var expires int
	var access_key string
	var secret_key string
	app := cli.NewApp()
	app.Name = "lock"
	app.Usage = "acquire 'lock' to following commands."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "region, r",
			Value:       "ap-northeast-1",
			Usage:       "Specify Region name",
			Destination: &region,
		},
		cli.StringFlag{
			Name:        "table, t",
			Usage:       "Specify Table name",
			Destination: &table,
		},
		cli.StringFlag{
			Name:        "aws_access_key_id",
			Usage:       "",
			Destination: &access_key,
			EnvVar:      "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:        "aws_secret_access_key",
			Usage:       "",
			Destination: &secret_key,
			EnvVar:      "AWS_SECRET_ACCESS_KEY",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "create-table",
			Usage: "Create table then exit",
			Action: func(c *cli.Context) error {
				db := dynamo.New(session.New(), &aws.Config{
					Region:     aws.String(region),
					DisableSSL: aws.Bool(true),
					Credentials: credentials.NewStaticCredentials(
						access_key, secret_key, ""),
				})
				return db.CreateTable(table, Record{}).Run()
			},
		},
		{
			Name:  "run",
			Usage: "run command in distributed lock",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "domain, d",
					Usage:       "Specify Lock resource name.",
					Destination: &domain,
				},
				cli.IntFlag{
					Name:        "expires, e",
					Usage:       "Specify Expire duration as seconds.",
					Destination: &expires,
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return fmt.Errorf("No command to exec.")
				}
				db := dynamo.New(session.New(), &aws.Config{
					Region: aws.String(region),
					Credentials: credentials.NewStaticCredentials(
						access_key, secret_key, ""),
				})
				lock := NewLock(db.Table(table), domain)
				return func() error {
					err := lock.acquire(expires)
					defer lock.release()
					if err != nil {
						return err
					}
					cmd := exec.Command(c.Args().First(), c.Args()[1:]...)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err = cmd.Run()
					if err != nil {
						return err
					}
					return nil
				}()
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
