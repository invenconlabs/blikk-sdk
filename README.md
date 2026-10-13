# blikk-sdk

Go client for the Blikk public API.

Module: `github.com/invenconlabs/blikk-sdk`, package `github.com/invenconlabs/blikk-sdk/blikk`.

```bash
go get github.com/invenconlabs/blikk-sdk
```

```go
client, err := blikk.NewClient(appID, appSecret)
if err != nil {
	return err
}

user, err := client.Users.Get(ctx, 42)

page, err := client.Projects.List(ctx, blikk.ProjectListParams{
	Query:       "roof",
	ListOptions: blikk.ListOptions{Page: 2, PageSize: 100},
})

for report, err := range client.TimeReports.ListAll(ctx, blikk.TimeReportListParams{
	From:    blikk.NewDate(2026, 9, 1),
	UserIDs: []int{1, 2},
}) {
	if err != nil {
		return err
	}
	_ = report
}
```

`ListAll` yields `(item, err)`: on failure it yields a single non-nil error and stops.
