# enterprise
Atılgan aka Enterprise is a new ORM for postgresql databases. For speed and lightness.

## Features
- A small number of auto generated database model structs and methods
- Create, Get, Update, Delete and List functions on table
- It can load Relation tables and filter them with single query
- Developer friendly helper methods on models
- Simple aggregate functionality
- Simple transaction commit and rollback
- Customizable hooks
- Atlas migrations
- Unique Relation system
- Fastest ORM in golang for postgresql


## Generating Models
First create a package. Create one file for every table. The function name and file name should match.
After write your models. We recommend to create a new package with name generate and put a generate.go for generating models you want.
It will go to create a models package and write files in to it. For examples [Enterprise Example](https://github.com/MrSametBurgazoglu/enterprise_example).

### Migration
For migration we recommend to create a migrate named package and put this script into it.
Execute the script and it will create a directory called migrations and write raw sql for migration.
Further questions please look to Atlas.

### Example
````go
package main

import (
	"context"
	"errors"
	"example/models"
	"github.com/MrSametBurgazoglu/enterprise/client"
	"log"
	"time"
)

func main() {
	dbUrl := "postgresql://testuser:54M3754M37@localhost:5433/testdb?search_path=public"
	db, err := models.NewDB(dbUrl)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	tes := models.NewTest(ctx, db)
	tes.SetName("name")
	tes.SetCreatedAt(time.Now())

	den := models.NewDeneme(ctx, db)
	den.SetCount(20)
	den.SetDenemeType(models.DenemeTypeTestType)
	den.SetTestIDValue(tes.GetID())

	acc := models.NewAccount(ctx, db)
	acc.SetName("name")
	acc.SetSurname("surname")
	acc.SetDenemeIDValue(den.GetID())

	acc2 := models.NewAccount(ctx, db)
	acc2.SetName("name")
	acc2.SetSurname("surname")
	acc2.SetDenemeIDValue(den.GetID())
	println(acc.GetID().String(), "id")
	println(acc2.GetID().String(), "id")

	err = tes.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = den.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = acc.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = acc2.Create()
	if err != nil {
		log.Fatal(err)
	}

	t := models.NewTest(ctx, db)
	t.Where(t.IsIDEqual(tes.GetID()))
	println("test", tes.GetID().String())
	t.WithDenemeList(func(denemeList *models.DenemeList) {
		denemeList.Where(
			denemeList.IsCountEqual(20),
		)
		denemeList.Order(models.DenemeIDField)
		denemeList.WithAccountList()
	})

	err, _ = t.Get()
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range t.DenemeList.Items[0].AccountList.Items {
		println(item.GetID().String())
	}

	var count, maximum, minimum, sum int

	t.DenemeList.Where(
		t.DenemeList.IsIsActiveEqual(true),
	)
	scanner, err := t.DenemeList.Aggregate(func(aggregate *client.Aggregate) {
		aggregate.Count("*", &count)
		aggregate.Max(models.DenemeCountField, &maximum)
		aggregate.Min(models.DenemeCountField, &minimum)
		aggregate.Sum(models.DenemeCountField, &sum)
		aggregate.GroupBy(models.DenemeDenemeTypeField)
	})
	err = scanner()
	for err == nil {
		println(count, maximum, minimum, sum)
		err = scanner()
	}
	if err != nil && !errors.Is(err, client.ErrFinalRow) {
		log.Fatal(err)
	}

	den.SetDenemeType(models.DenemeTypeDenemeType)
	err = den.Update()
	if err != nil {
		log.Fatal(err)
	}
	err = den.Refresh()
	if err != nil {
		log.Fatal(err)
	}

	transaction, err := db.NewTransaction(ctx)
	if err != nil {
		log.Fatal(err)
	}

	test2 := models.NewTest(ctx, transaction)
	test2.SetName("transaction_name")
	test2.SetCreatedAt(time.Now())

	err = test2.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = transaction.Rollback(ctx)
	if err != nil {
		log.Fatal(err)
	}

	transaction2, err := db.NewTransaction(ctx)
	if err != nil {
		log.Fatal(err)
	}

	test3 := models.NewTest(ctx, transaction2)
	test3.SetName("transaction_name")
	test3.SetCreatedAt(time.Now())

	err = test3.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = transaction2.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	testList := models.NewTestList(ctx, db)
	testList.Where(
		testList.IsNameEqual("name"),
	)
	testList.Paging(0, 5)
	testList.WithDenemeList(func(denemeList *models.DenemeList) {
		denemeList.WithAccountList()
	})
	err, found := testList.List()
	if err == nil && found {
		for i, test := range testList.Items {
			println(i, test.GetID().String())
			for i2, deneme := range test.DenemeList.Items {
				println(i2, deneme.GetID().String())
				for i3, account := range deneme.AccountList.Items {
					println(i3, account.GetID().String())
				}
			}
			//println(item.Deneme.GetCount())
		}
	} else {
		println(err.Error())
	}

	acc3 := models.NewAccount(ctx, db)
	acc3.SetName("with_group")
	acc3.SetSurname("surname")

	err = acc3.Create()
	if err != nil {
		log.Fatal(err)
	}

	group := models.NewGroup(ctx, db)
	group.SetName("with_account")
	group.SetSurname("surname")

	err = group.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = acc3.AddIntoGroup(group)
	if err != nil {
		log.Fatal(err)
	}

	acc4 := models.NewAccount(ctx, db)
	acc4.Where(acc4.IsIDEqual(acc3.GetID()))
	acc4.WithGroupList()
	err, ok := acc4.Get()
	if err == nil && ok {
		println(acc4.GroupList.Items[0].GetID().String())
	}

}

````

### Benchmark Results
````text
goos: linux
goarch: amd64
pkg: github.com/FournyP/go-orm-benchmarks/benchmarks
cpu: 12th Gen Intel(R) Core(TM) i7-1255U

BenchmarkEntCreate-4                                 105          11393218 ns/op            7571 B/op        218 allocs/op
BenchmarkEnterpriseCreate-4                          175           6458168 ns/op            9264 B/op        193 allocs/op
BenchmarkGORMCreate-4                                138           8334585 ns/op           22754 B/op        334 allocs/op
BenchmarkSqlxCreate-4                                139           8516631 ns/op            3088 B/op         87 allocs/op

BenchmarkEntUpdate-4                                 100          10051190 ns/op            5302 B/op        137 allocs/op
BenchmarkEnterpriseUpdate-4                          507           2359239 ns/op            1652 B/op         39 allocs/op
BenchmarkGORMUpdate-4                                468           2712004 ns/op            6657 B/op         86 allocs/op
BenchmarkSqlxUpdate-4                                502           2518410 ns/op             535 B/op         16 allocs/op

BenchmarkEntDelete-4                                 367           2897743 ns/op            1904 B/op         45 allocs/op
BenchmarkEnterpriseDelete-4                          513           2089324 ns/op            1157 B/op         21 allocs/op
BenchmarkGORMDelete-4                                492           3073590 ns/op            5634 B/op         85 allocs/op
BenchmarkSqlxDelete-4                                513           2385171 ns/op             304 B/op          9 allocs/op

BenchmarkEntRead-4                                  1107           1179827 ns/op            3744 B/op         93 allocs/op
BenchmarkEnterpriseRead-4                           2575            428796 ns/op            3005 B/op         73 allocs/op
BenchmarkGORMRead-4                                 2296            519182 ns/op            5406 B/op         94 allocs/op
BenchmarkSqlxRead-4                                 1273            943025 ns/op            1200 B/op         32 allocs/op

BenchmarkEntReadWithRelations-4                      424           3185571 ns/op           11888 B/op        295 allocs/op
BenchmarkEnterpriseReadWithRelations-4              2185            544165 ns/op           10781 B/op        219 allocs/op
BenchmarkGORMReadWithRelations-4                     692           2137871 ns/op           34967 B/op        372 allocs/op
BenchmarkSqlxReadWithRelations-4                     456           2687558 ns/op            3856 B/op        105 allocs/op

BenchmarkEntReadSingleField-4                       1153            973295 ns/op            2776 B/op         73 allocs/op
BenchmarkEnterpriseReadSingleField-4                2908            494583 ns/op            2686 B/op         62 allocs/op
BenchmarkGormReadSingleField-4                      2036            515719 ns/op            4738 B/op         73 allocs/op
BenchmarkSqlxReadSingleField-4                      1242           1015917 ns/op             792 B/op         22 allocs/op
````