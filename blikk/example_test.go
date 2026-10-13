package blikk_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/invenconlabs/blikk-sdk/blikk"
)

func Example_hoursPerProjectLastMonth() {
	ctx := context.Background()
	client, err := blikk.NewClient(os.Getenv("BLIKK_APP_ID"), os.Getenv("BLIKK_APP_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	// time.Date normalises month 0 to the last day of the previous month.
	firstDayOfLastMonth := blikk.NewDate(now.Year(), now.Month()-1, 1)
	lastDayOfLastMonth := blikk.NewDate(now.Year(), now.Month(), 0)

	hoursByProject := map[string]float64{}
	for report, err := range client.TimeReports.ListAll(ctx, blikk.TimeReportListParams{
		UserIDs: []int{1042},
		From:    firstDayOfLastMonth,
		To:      lastDayOfLastMonth,
	}) {
		if err != nil {
			log.Fatal(err)
		}
		projectName := report.ProjectName()
		if projectName == "" {
			projectName = "(no project)"
		}
		hoursByProject[projectName] += report.Hours
	}

	projectNames := make([]string, 0, len(hoursByProject))
	for projectName := range hoursByProject {
		projectNames = append(projectNames, projectName)
	}
	sort.Strings(projectNames)

	fmt.Printf("Hours %s to %s\n", firstDayOfLastMonth, lastDayOfLastMonth)
	for _, projectName := range projectNames {
		fmt.Printf("%-40s %6.1f\n", projectName, hoursByProject[projectName])
	}
}

func Example_activeProjectsPage() {
	ctx := context.Background()
	client, err := blikk.NewClient(os.Getenv("BLIKK_APP_ID"), os.Getenv("BLIKK_APP_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	page, err := client.Projects.List(ctx, blikk.ProjectListParams{
		ListOptions: blikk.ListOptions{
			PageSize:  50,
			SortBy:    blikk.ProjectSortByUpdatedDate,
			SortOrder: blikk.SortDescending,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, project := range page.Items {
		statusName := ""
		if project.Status != nil {
			statusName = project.Status.Name
		}
		fmt.Printf("%-10s %-40s %s\n", project.OrderNumber, project.Title, statusName)
	}
	fmt.Printf("showing %d of %d projects (%d pages)\n", page.ItemCount, page.TotalItemCount, page.TotalPages)
}

func Example_createTimeReport() {
	ctx := context.Background()
	client, err := blikk.NewClient(os.Getenv("BLIKK_APP_ID"), os.Getenv("BLIKK_APP_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	reportID, err := client.TimeReports.Create(ctx, blikk.TimeReportInput{
		Date:             blikk.NewDate(2026, time.September, 14),
		ClockStart:       "08:00",
		ClockEnd:         "16:30",
		BreakMinutes:     30,
		Hours:            8,
		InvoiceableHours: 8,
		Comment:          "Installed new switchboard",
		UserID:           1042,
		ProjectID:        2301,
		ActivityID:       17,
		TimeCodeID:       3,
		TimeArticleID:    9,
	})
	if err != nil {
		log.Fatal(err)
	}

	report, err := client.TimeReports.Get(ctx, reportID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("#%d %s %s %.1fh on %s: %s\n", report.ID, report.Date, report.User.Name, report.Hours, report.ProjectName(), report.Comment)
}

func Example_handleNotFound() {
	ctx := context.Background()
	client, err := blikk.NewClient(os.Getenv("BLIKK_APP_ID"), os.Getenv("BLIKK_APP_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	project, err := client.Projects.Get(ctx, 987654)
	switch {
	case blikk.IsNotFound(err):
		fmt.Println("project does not exist")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("%s %s\n", project.OrderNumber, project.Title)
	}
}

func Example_customClient() {
	client, err := blikk.NewClient(
		os.Getenv("BLIKK_APP_ID"),
		os.Getenv("BLIKK_APP_SECRET"),
		blikk.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
		blikk.WithUserAgent("payroll-export/1.4"),
		blikk.WithRequestsPerSecond(2),
	)
	if err != nil {
		log.Fatal(err)
	}
	_ = client

	clientWithToken, err := blikk.NewClientWithTokenSource(blikk.StaticToken(os.Getenv("BLIKK_TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	_ = clientWithToken
}
