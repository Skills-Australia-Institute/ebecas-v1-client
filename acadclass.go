package ebecasv1client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type getAcademicClassesResponse struct {
	PageSize       int             `json:"PageSize"`
	PageCount      int             `json:"PageCount"`
	AcadClassCount int             `json:"AcadClassCount"`
	AcadClasses    []AcademicClass `json:"AcadClasses"`
}

type AcademicClass struct {
	ClassID            int64  `json:"ClassId"`
	LocationCode       string `json:"LocationCode"`
	ClassName          string `json:"ClassName"`
	SubjectCode        string `json:"SubjectCode"`
	SubjectDescription string `json:"SubjectDescription"`
	Rooms              string `json:"Rooms"`
	RoomCapacities     string `json:"RoomCapacities"`
	Teachers           string `json:"Teachers"`
	StudentCount       int    `json:"StudentCount"`
	OnlineClass        bool   `json:"OnlineClass"`
}

type getAcademicClassScheduleResponse struct {
	ScheduleCount int                     `json:"ScheduleCount"`
	ScheduleList  []AcademicClassSchedule `json:"ScheduleList"`
}

type AcademicClassSchedule struct {
	ClassType      string `json:"ClassType"`
	ClassDay       string `json:"ClassDay"`
	StartTime      string `json:"StartTime"`
	FinishTime     string `json:"FinishTime"`
	ClassDate      string `json:"ClassDate"`
	RoomCode       string `json:"RoomCode"`
	RoomCapacity   int    `json:"RoomCapacity"`
	TeacherList    string `json:"TeacherList"`
	MarkAttendance bool   `json:"MarkAttendance"`
}

type GetAcademicClassesParams struct {
	LocationCode string
	AllLocations bool
	ClassDate    string
	PageSize     int
	Page         int
}

type AcademicClassesResult struct {
	Classes    []AcademicClass
	PageSize   int
	Page       int
	TotalCount int
}

// GetAcademicClasses retrieves academic classes for the specified date and location.
func (c *Client) GetAcademicClasses(
	ctx context.Context,
	params GetAcademicClassesParams,
) (AcademicClassesResult, int, error) {
	locationCode := strings.TrimSpace(params.LocationCode)
	classDate := strings.TrimSpace(params.ClassDate)

	if !params.AllLocations && locationCode == "" {
		return AcademicClassesResult{}, http.StatusBadRequest,
			fmt.Errorf("location code is required when all locations is false")
	}

	if params.PageSize < 0 {
		return AcademicClassesResult{}, http.StatusBadRequest,
			fmt.Errorf("page size must not be negative")
	}

	if params.PageSize > maxPageSize {
		return AcademicClassesResult{}, http.StatusBadRequest,
			fmt.Errorf("page size must not exceed %d", maxPageSize)
	}

	if params.Page < 0 {
		return AcademicClassesResult{}, http.StatusBadRequest,
			fmt.Errorf("page must not be negative")
	}

	query := url.Values{}

	if locationCode != "" {
		query.Set("LocationCode", locationCode)
	}

	if params.AllLocations {
		query.Set("AllLocations", "true")
	}

	if classDate != "" {
		query.Set("ClassDate", classDate)
	}

	if params.PageSize > 0 {
		query.Set("PageSize", strconv.Itoa(params.PageSize))
	}

	if params.Page > 0 {
		query.Set("PageCount", strconv.Itoa(params.Page))
	}

	requestURL := c.baseURL + "/acadclass"
	if encodedQuery := query.Encode(); encodedQuery != "" {
		requestURL += "?" + encodedQuery
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return AcademicClassesResult{}, http.StatusInternalServerError, fmt.Errorf(
			"create get academic classes request: %w",
			err,
		)
	}

	data, statusCode, err := c.do(req)
	if err != nil {
		return AcademicClassesResult{}, statusCode, fmt.Errorf(
			"get academic classes, status %d: %w",
			statusCode,
			err,
		)
	}

	if statusCode != http.StatusOK {
		return AcademicClassesResult{}, statusCode, fmt.Errorf(
			"get academic classes returned status %d: %s",
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	var resp getAcademicClassesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return AcademicClassesResult{}, statusCode, fmt.Errorf(
			"decode academic classes response: %w",
			err,
		)
	}

	return AcademicClassesResult{
		Classes:    resp.AcadClasses,
		PageSize:   resp.PageSize,
		Page:       resp.PageCount,
		TotalCount: resp.AcadClassCount,
	}, statusCode, nil
}

// GetAcademicClassSchedule retrieves the schedule for an academic class.
func (c *Client) GetAcademicClassSchedule(
	ctx context.Context,
	classID int64,
) ([]AcademicClassSchedule, int, error) {
	if classID <= 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("class ID must be greater than zero")
	}

	requestURL := fmt.Sprintf("%s/acadclass/%d/schedule", c.baseURL, classID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(
			"create get schedule request for academic class %d: %w",
			classID,
			err,
		)
	}

	data, statusCode, err := c.do(req)
	if err != nil {
		return nil, statusCode, fmt.Errorf(
			"get schedule for academic class %d, status %d: %w",
			classID,
			statusCode,
			err,
		)
	}

	if statusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf(
			"get schedule for academic class %d returned status %d: %s",
			classID,
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	var resp getAcademicClassScheduleResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, statusCode, fmt.Errorf(
			"decode schedule response for academic class %d: %w",
			classID,
			err,
		)
	}

	return resp.ScheduleList, statusCode, nil
}

type getAcademicClassAllocationsResponse struct {
	AllocationCount int                       `json:"AllocationCount"`
	AllocationList  []AcademicClassAllocation `json:"AllocationList"`
}

type AcademicClassAllocation struct {
	StudentID             int64  `json:"StudentId"`
	StudentNo             string `json:"StudentNo"`
	StudentName           string `json:"StudentName"`
	CountryName           string `json:"CountryName"`
	StartDate             string `json:"StartDate"`
	EndDate               string `json:"EndDate"`
	DateOfBirth           string `json:"DateOfBirth"`
	Age                   int    `json:"Age"`
	Gender                string `json:"Gender"`
	WeeksLeft             int    `json:"WeeksLeft"`
	Status                string `json:"Status"`
	CourseCode            string `json:"CourseCode"`
	CourseName            string `json:"CourseName"`
	VisaType              string `json:"VisaType"`
	Arrived               bool   `json:"Arrived"`
	OverallAttendance     int    `json:"OverallAttendance"`
	CurrentAttendance     int    `json:"CurrentAttendance"`
	Holiday               bool   `json:"Holiday"`
	PathwayMonitoringFlag bool   `json:"PathwayMonitoringFlag"`
	PathwayStartDate      string `json:"PathwayStartDate"`
}

// GetAcademicClassAllocations retrieves the student allocations for an academic class.
func (c *Client) GetAcademicClassAllocations(
	ctx context.Context,
	classID int64,
) ([]AcademicClassAllocation, int, error) {
	if classID <= 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("class ID must be greater than zero")
	}

	requestURL := fmt.Sprintf("%s/acadclass/%d/allocation", c.baseURL, classID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(
			"create get allocations request for academic class %d: %w",
			classID,
			err,
		)
	}

	data, statusCode, err := c.do(req)
	if err != nil {
		return nil, statusCode, fmt.Errorf(
			"get allocations for academic class %d, status %d: %w",
			classID,
			statusCode,
			err,
		)
	}

	if statusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf(
			"get allocations for academic class %d returned status %d: %s",
			classID,
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	var resp getAcademicClassAllocationsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, statusCode, fmt.Errorf(
			"decode allocations response for academic class %d: %w",
			classID,
			err,
		)
	}

	return resp.AllocationList, statusCode, nil
}
