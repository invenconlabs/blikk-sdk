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

The client has one service per API resource. Core: contacts, deleted entities, invoices, logs, material reports, offers, opportunities, payment plans, planning, project collections, projects, supplier invoices, task boards, tasks, time reports and vacation requests. Admin: absence projects, activities, articles, booking projects, company cars, contact tags, cost centers, countries, departments, internal projects, opportunity statuses and tags, project categories, statuses and tags, time codes, units, users and VAT codes.

`ListAll` yields `(item, err)`: on failure it yields a single non-nil error and stops.
