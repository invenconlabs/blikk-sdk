package blikk

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"math"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	DefaultBaseURL           = "https://publicapi.blikk.com"
	DefaultRequestsPerSecond = 4.0

	defaultUserAgent    = "blikk-sdk-go"
	defaultHTTPTimeout  = 30 * time.Second
	maxRateLimitRetries = 5
	maxRetryAfter       = 30 * time.Second
	defaultRetryAfter   = time.Second
	maxErrorBodyBytes   = 64 << 10
	jsonContentType     = "application/json"
	anyContentType      = "*/*"
)

type Client struct {
	Contacts           *ContactsService
	DeletedEntities    *DeletedEntitiesService
	Invoices           *InvoicesService
	Logs               *LogsService
	MaterialReports    *MaterialReportsService
	Offers             *OffersService
	Opportunities      *OpportunitiesService
	PaymentPlans       *PaymentPlansService
	Planning           *PlanningService
	ProjectCollections *ProjectCollectionsService
	Projects           *ProjectsService
	SupplierInvoices   *SupplierInvoicesService
	TaskBoards         *TaskBoardsService
	Tasks              *TasksService
	TimeReports        *TimeReportsService
	VacationRequests   *VacationRequestsService

	AbsenceProjects     *AbsenceProjectsService
	Activities          *ActivitiesService
	Articles            *ArticlesService
	BookingProjects     *BookingProjectsService
	CompanyCars         *CompanyCarsService
	ContactTags         *ContactTagsService
	CostCenters         *CostCentersService
	Countries           *CountriesService
	Departments         *DepartmentsService
	InternalProjects    *InternalProjectsService
	OpportunityStatuses *OpportunityStatusesService
	OpportunityTags     *OpportunityTagsService
	ProjectCategories   *ProjectCategoriesService
	ProjectStatuses     *ProjectStatusesService
	ProjectTags         *ProjectTagsService
	TimeCodes           *TimeCodesService
	Units               *UnitsService
	Users               *UsersService
	VatCodes            *VatCodesService

	baseURL     string
	httpClient  *http.Client
	tokenSource TokenSource
	limiter     *rate.Limiter
	userAgent   string
	sleep       func(ctx context.Context, d time.Duration) error
}

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = baseURL }
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}

// WithTokenSource replaces the app ID and secret exchange; the credentials passed to NewClient are then ignored.
func WithTokenSource(tokenSource TokenSource) Option {
	return func(c *Client) { c.tokenSource = tokenSource }
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) { c.userAgent = userAgent }
}

// WithRequestsPerSecond sets the client-side request rate; the API allows 4 per second. Zero or less disables limiting.
func WithRequestsPerSecond(requestsPerSecond float64) Option {
	return func(c *Client) { c.limiter = newLimiter(requestsPerSecond) }
}

func newLimiter(requestsPerSecond float64) *rate.Limiter {
	if requestsPerSecond <= 0 {
		return rate.NewLimiter(rate.Inf, 1)
	}
	return rate.NewLimiter(rate.Limit(requestsPerSecond), 1)
}

func NewClient(appID, appSecret string, opts ...Option) (*Client, error) {
	c := &Client{
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: defaultHTTPTimeout},
		limiter:    newLimiter(DefaultRequestsPerSecond),
		userAgent:  defaultUserAgent,
		sleep:      sleepContext,
	}
	for _, opt := range opts {
		opt(c)
	}
	if err := validateBaseURL(c.baseURL); err != nil {
		return nil, err
	}
	if c.tokenSource == nil {
		tokenSource, err := c.newAppTokenSource(appID, appSecret)
		if err != nil {
			return nil, err
		}
		c.tokenSource = tokenSource
	}
	c.registerServices()
	return c, nil
}

func (c *Client) registerServices() {
	c.Contacts = &ContactsService{client: c}
	c.DeletedEntities = &DeletedEntitiesService{client: c}
	c.Invoices = &InvoicesService{client: c}
	c.Logs = &LogsService{client: c}
	c.MaterialReports = &MaterialReportsService{client: c}
	c.Offers = &OffersService{client: c}
	c.Opportunities = &OpportunitiesService{client: c}
	c.PaymentPlans = &PaymentPlansService{client: c}
	c.Planning = &PlanningService{client: c}
	c.ProjectCollections = &ProjectCollectionsService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.SupplierInvoices = &SupplierInvoicesService{client: c}
	c.TaskBoards = &TaskBoardsService{client: c}
	c.Tasks = &TasksService{client: c}
	c.TimeReports = &TimeReportsService{client: c}
	c.VacationRequests = &VacationRequestsService{client: c}

	c.AbsenceProjects = &AbsenceProjectsService{client: c}
	c.Activities = &ActivitiesService{client: c}
	c.Articles = &ArticlesService{client: c}
	c.BookingProjects = &BookingProjectsService{client: c}
	c.CompanyCars = &CompanyCarsService{client: c}
	c.ContactTags = &ContactTagsService{client: c}
	c.CostCenters = &CostCentersService{client: c}
	c.Countries = &CountriesService{client: c}
	c.Departments = &DepartmentsService{client: c}
	c.InternalProjects = &InternalProjectsService{client: c}
	c.OpportunityStatuses = &OpportunityStatusesService{client: c}
	c.OpportunityTags = &OpportunityTagsService{client: c}
	c.ProjectCategories = &ProjectCategoriesService{client: c}
	c.ProjectStatuses = &ProjectStatusesService{client: c}
	c.ProjectTags = &ProjectTagsService{client: c}
	c.TimeCodes = &TimeCodesService{client: c}
	c.Units = &UnitsService{client: c}
	c.Users = &UsersService{client: c}
	c.VatCodes = &VatCodesService{client: c}
}

func NewClientWithTokenSource(source TokenSource, opts ...Option) (*Client, error) {
	if source == nil {
		return nil, errors.New("blikk: token source is nil")
	}
	return NewClient("", "", append([]Option{WithTokenSource(source)}, opts...)...)
}

func validateBaseURL(baseURL string) error {
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("blikk: invalid base URL %q", baseURL)
	}
	if !isAllowedScheme(parsed) {
		return fmt.Errorf("blikk: base URL %q must use https unless the host is loopback", baseURL)
	}
	return nil
}

func isAllowedScheme(parsed *url.URL) bool {
	switch parsed.Scheme {
	case "https":
		return true
	case "http":
		return isLoopbackHost(parsed.Hostname())
	default:
		return false
	}
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func (c *Client) newAppTokenSource(appID, appSecret string) (*appTokenSource, error) {
	if appID == "" || appSecret == "" {
		return nil, errors.New("blikk: app ID and app secret are required unless WithTokenSource is used")
	}
	tokenURL, err := url.JoinPath(c.baseURL, tokenPath)
	if err != nil {
		return nil, fmt.Errorf("blikk: invalid base URL: %w", err)
	}
	return &appTokenSource{
		appID:      appID,
		appSecret:  appSecret,
		tokenURL:   tokenURL,
		userAgent:  c.userAgent,
		httpClient: c.httpClient,
		now:        time.Now,
		sleep:      c.sleep,
	}, nil
}

func (c *Client) do(ctx context.Context, method, resourcePath string, query url.Values, body any, out any) error {
	encoded, err := encodeJSONBody(body)
	if err != nil {
		return err
	}
	return c.doEncoded(ctx, method, resourcePath, query, encoded, out)
}

type UploadFile struct {
	FieldName   string
	Filename    string
	ContentType string
	Content     io.Reader
}

func (c *Client) upload(ctx context.Context, resourcePath string, fields map[string]string, files []UploadFile, out any) error {
	encoded, err := encodeMultipartBody(fields, files)
	if err != nil {
		return err
	}
	return c.doEncoded(ctx, http.MethodPost, resourcePath, nil, encoded, out)
}

func (c *Client) doEncoded(ctx context.Context, method, resourcePath string, query url.Values, body requestBody, out any) error {
	resp, err := c.execute(ctx, method, resourcePath, query, body, jsonContentType)
	if err != nil {
		return err
	}
	responseBody, err := readBody(resp)
	if err != nil {
		return fmt.Errorf("blikk: reading response of %s %s: %w", method, resp.Request.URL.Path, err)
	}
	return decodeBody(method, resourcePath, responseBody, out)
}

type Download struct {
	Content     io.ReadCloser
	ContentType string
	Filename    string
}

// download streams the body; the caller must close Content.
func (c *Client) download(ctx context.Context, resourcePath string, query url.Values) (*Download, error) {
	resp, err := c.execute(ctx, http.MethodGet, resourcePath, query, requestBody{}, anyContentType)
	if err != nil {
		return nil, err
	}
	return &Download{
		Content:     resp.Body,
		ContentType: resp.Header.Get("Content-Type"),
		Filename:    filenameFromContentDisposition(resp.Header.Get("Content-Disposition")),
	}, nil
}

func filenameFromContentDisposition(header string) string {
	_, params, err := mime.ParseMediaType(header)
	if err != nil {
		return ""
	}
	return baseName(params["filename"])
}

// The name comes from the server, so drop any directory part before a caller writes it to disk.
func baseName(filename string) string {
	if i := strings.LastIndexAny(filename, `/\`); i >= 0 {
		filename = filename[i+1:]
	}
	if filename == "." || filename == ".." {
		return ""
	}
	return filename
}

// execute returns only successful responses, with the body still open for the caller to read and close.
func (c *Client) execute(ctx context.Context, method, resourcePath string, query url.Values, body requestBody, accept string) (*http.Response, error) {
	requestURL, err := c.buildURL(resourcePath, query)
	if err != nil {
		return nil, err
	}

	rateLimitRetries := 0
	tokenInvalidated := false
	for {
		token, err := c.tokenSource.Token(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := c.send(ctx, method, requestURL, body, token, accept)
		if err != nil {
			return nil, err
		}
		if isSuccess(resp.StatusCode) {
			return resp, nil
		}
		errorBody, err := readBody(resp)
		if err != nil {
			return nil, fmt.Errorf("blikk: reading response of %s %s: %w", method, resp.Request.URL.Path, err)
		}
		switch {
		case resp.StatusCode == http.StatusTooManyRequests && rateLimitRetries < maxRateLimitRetries:
			rateLimitRetries++
			if err := c.sleep(ctx, retryAfter(resp.Header, errorBody)); err != nil {
				return nil, err
			}
		case resp.StatusCode == http.StatusUnauthorized && !tokenInvalidated:
			tokenInvalidated = true
			c.tokenSource.Invalidate(token)
		default:
			return nil, newError(method, "/"+resourcePath, resp.StatusCode, resp.Header, errorBody)
		}
	}
}

func (c *Client) buildURL(resourcePath string, query url.Values) (string, error) {
	requestURL, err := url.JoinPath(c.baseURL, resourcePath)
	if err != nil {
		return "", fmt.Errorf("blikk: building URL for %s: %w", resourcePath, err)
	}
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}
	return requestURL, nil
}

func (c *Client) send(ctx context.Context, method, requestURL string, body requestBody, token, accept string) (*http.Response, error) {
	if err := c.waitForRateLimit(ctx); err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, method, requestURL, body, token, accept)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("blikk: %s %s: %w", method, req.URL.Path, err)
	}
	return resp, nil
}

// rate.Limiter fails early with a plain error when the wait would outlast the deadline, before ctx is done.
func (c *Client) waitForRateLimit(ctx context.Context) error {
	err := c.limiter.Wait(ctx)
	switch {
	case err == nil:
		return nil
	case ctx.Err() != nil:
		return fmt.Errorf("blikk: waiting for rate limit: %w", ctx.Err())
	default:
		return fmt.Errorf("blikk: waiting for rate limit: %w (%v)", context.DeadlineExceeded, err)
	}
}

func (c *Client) newRequest(ctx context.Context, method, requestURL string, body requestBody, token, accept string) (*http.Request, error) {
	var reader io.Reader
	if body.payload != nil {
		reader = bytes.NewReader(body.payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, requestURL, reader)
	if err != nil {
		return nil, fmt.Errorf("blikk: creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", c.userAgent)
	if body.payload != nil {
		req.Header.Set("Content-Type", body.contentType)
	}
	return req, nil
}

func readResponse(resp *http.Response) (int, []byte, error) {
	body, err := readBody(resp)
	return resp.StatusCode, body, err
}

func readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	reader := io.Reader(resp.Body)
	if !isSuccess(resp.StatusCode) {
		reader = io.LimitReader(resp.Body, maxErrorBodyBytes)
	}
	return io.ReadAll(reader)
}

// The payload is held in memory so a retry after 429 or 401 can send it again.
type requestBody struct {
	payload     []byte
	contentType string
}

func encodeJSONBody(body any) (requestBody, error) {
	if body == nil {
		return requestBody{}, nil
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return requestBody{}, fmt.Errorf("blikk: encoding request body: %w", err)
	}
	return requestBody{payload: payload, contentType: jsonContentType}, nil
}

func encodeMultipartBody(fields map[string]string, files []UploadFile) (requestBody, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	for _, name := range slices.Sorted(maps.Keys(fields)) {
		if err := writer.WriteField(name, fields[name]); err != nil {
			return requestBody{}, fmt.Errorf("blikk: encoding form field %q: %w", name, err)
		}
	}
	for _, file := range files {
		if err := writeFilePart(writer, file); err != nil {
			return requestBody{}, err
		}
	}
	if err := writer.Close(); err != nil {
		return requestBody{}, fmt.Errorf("blikk: encoding multipart body: %w", err)
	}
	return requestBody{payload: buffer.Bytes(), contentType: writer.FormDataContentType()}, nil
}

func writeFilePart(writer *multipart.Writer, file UploadFile) error {
	if file.Content == nil {
		return fmt.Errorf("blikk: upload file %q has no content", file.Filename)
	}
	part, err := writer.CreatePart(filePartHeader(file))
	if err != nil {
		return fmt.Errorf("blikk: encoding file %q: %w", file.Filename, err)
	}
	if _, err := io.Copy(part, file.Content); err != nil {
		return fmt.Errorf("blikk: reading file %q: %w", file.Filename, err)
	}
	return nil
}

var quoteEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

// Quoted like multipart.Writer.CreateFormFile, which cannot set the part's Content-Type.
func filePartHeader(file UploadFile) textproto.MIMEHeader {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
		quoteEscaper.Replace(file.FieldName), quoteEscaper.Replace(file.Filename)))
	header.Set("Content-Type", cmp.Or(file.ContentType, "application/octet-stream"))
	return header
}

func setFieldIfNotEmpty(fields map[string]string, name, value string) {
	if value != "" {
		fields[name] = value
	}
}

func withDefaultFieldName(file UploadFile, fieldName string) UploadFile {
	file.FieldName = cmp.Or(file.FieldName, fieldName)
	return file
}

func decodeBody(method, resourcePath string, body []byte, out any) error {
	if out == nil || len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("blikk: decoding response of %s %s: %w", method, resourcePath, err)
	}
	return nil
}

func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

type rateLimitBody struct {
	WaitInSeconds float64 `json:"waitInSeconds"`
}

func retryAfter(header http.Header, body []byte) time.Duration {
	wait, ok := retryAfterFromHeader(header)
	if !ok {
		wait, ok = retryAfterFromBody(body)
	}
	if !ok {
		wait = defaultRetryAfter
	}
	return wait
}

func retryAfterFromHeader(header http.Header) (time.Duration, bool) {
	seconds, err := strconv.ParseFloat(header.Get("Retry-After"), 64)
	if err != nil || !isUsableWait(seconds) {
		return 0, false
	}
	return secondsToDuration(seconds), true
}

func retryAfterFromBody(body []byte) (time.Duration, bool) {
	var decoded rateLimitBody
	if err := json.Unmarshal(body, &decoded); err != nil || !isUsableWait(decoded.WaitInSeconds) {
		return 0, false
	}
	return secondsToDuration(decoded.WaitInSeconds), true
}

func isUsableWait(seconds float64) bool {
	return !math.IsNaN(seconds) && !math.IsInf(seconds, 0) && seconds > 0
}

// Clamping before the conversion matters: a huge float overflows time.Duration.
func secondsToDuration(seconds float64) time.Duration {
	return time.Duration(min(seconds, maxRetryAfter.Seconds()) * float64(time.Second))
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func resourcePath(base string, pathParams ...string) (string, error) {
	segments := []string{base}
	for _, param := range pathParams {
		if err := validatePathParam(param); err != nil {
			return "", err
		}
		segments = append(segments, url.PathEscape(param))
	}
	return strings.Join(segments, "/"), nil
}

// url.JoinPath cleans the path, so these would silently address a different resource.
func validatePathParam(param string) error {
	if param == "" || param == "." || param == ".." {
		return fmt.Errorf("blikk: invalid path parameter %q", param)
	}
	return nil
}
