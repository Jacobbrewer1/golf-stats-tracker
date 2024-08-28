// Package rounder provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package rounder

import (
	"time"

	externalRef0 "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/common"
)

const (
	BasicAuthScopes = "basicAuth.Scopes"
)

// AverageType defines the model for average_type.
type AverageType = string

// List of AverageType
const (
	AverageType_fairway_hit AverageType = "fairway_hit"
	AverageType_green_hit   AverageType = "green_hit"
	AverageType_par_3       AverageType = "par_3"
	AverageType_par_4       AverageType = "par_4"
	AverageType_par_5       AverageType = "par_5"
	AverageType_penalties   AverageType = "penalties"
	AverageType_putts       AverageType = "putts"
)

// ChartDataPoint defines the model for chart_data_point.
type ChartDataPoint struct {
	// X The x-axis label
	X *string `json:"x,omitempty"`

	// Y The y-axis label
	Y *float32 `json:"y,omitempty"`
}

// Course defines the model for course.
type Course struct {
	Details *[]CourseDetails `json:"details,omitempty"`
	Id      *int64           `json:"id,omitempty"`
	Name    *string          `json:"name,omitempty"`
}

// CourseDetails defines the model for course_details.
type CourseDetails struct {
	Holes            *[]Hole  `json:"holes,omitempty"`
	Id               *int64   `json:"id,omitempty"`
	Marker           *string  `json:"marker,omitempty"`
	MetersBackNine   *int64   `json:"meters_back_nine,omitempty"`
	MetersFrontNine  *int64   `json:"meters_front_nine,omitempty"`
	MetersTotal      *int64   `json:"meters_total,omitempty"`
	ParBackNine      *int64   `json:"par_back_nine,omitempty"`
	ParFrontNine     *int64   `json:"par_front_nine,omitempty"`
	ParTotal         *int64   `json:"par_total,omitempty"`
	Rating           *float64 `json:"rating,omitempty"`
	Slope            *int64   `json:"slope,omitempty"`
	YardageBackNine  *int64   `json:"yardage_back_nine,omitempty"`
	YardageFrontNine *int64   `json:"yardage_front_nine,omitempty"`
	YardageTotal     *int64   `json:"yardage_total,omitempty"`
}

// HitInRegulation defines the model for hit_in_regulation.
type HitInRegulation = string

// List of HitInRegulation
const (
	HitInRegulation_HIT            HitInRegulation = "HIT"
	HitInRegulation_LEFT           HitInRegulation = "LEFT"
	HitInRegulation_LONG           HitInRegulation = "LONG"
	HitInRegulation_NOT_APPLICABLE HitInRegulation = "NOT_APPLICABLE"
	HitInRegulation_RIGHT          HitInRegulation = "RIGHT"
	HitInRegulation_SHORT          HitInRegulation = "SHORT"
)

// Hole defines the model for hole.
type Hole struct {
	Id          *int64 `json:"id,omitempty"`
	Meters      *int64 `json:"meters,omitempty"`
	Number      *int64 `json:"number,omitempty"`
	Par         *int64 `json:"par,omitempty"`
	StrokeIndex *int64 `json:"stroke_index,omitempty"`
	Yardage     *int64 `json:"yardage,omitempty"`
}

// HoleStats defines the model for hole_stats.
type HoleStats struct {
	FairwayHit *HitInRegulation `json:"fairway_hit,omitempty"`
	GreenHit   *HitInRegulation `json:"green_hit,omitempty"`

	// Penalties The number of penalty strokes
	Penalties *int64 `json:"penalties,omitempty"`

	// PinLocation The pin position
	PinLocation *string `json:"pin_location,omitempty"`

	// Putts The number of putts
	Putts *int64 `json:"putts,omitempty"`

	// Score The number of strokes
	Score *int64 `json:"score,omitempty"`
}

// Round defines the model for round.
type Round struct {
	// CourseName The course name
	CourseName *string `json:"course_name,omitempty"`

	// Id The round id
	Id *int64 `json:"id,omitempty"`

	// Marker The marker
	Marker *string `json:"marker,omitempty"`

	// TeeTime The tee time
	TeeTime *time.Time `json:"tee_time,omitempty"`
}

// RoundCreate defines the model for round_create.
type RoundCreate struct {
	// CourseId The course id
	CourseId *int64 `json:"course_id,omitempty"`

	// Id The round id
	Id *int64 `json:"id,omitempty"`

	// MarkerId The marker id
	MarkerId *int64 `json:"marker_id,omitempty"`

	// TeeTime The tee time
	TeeTime *time.Time `json:"tee_time,omitempty"`
}

// Token defines the model for token.
type Token struct {
	// Token The token
	Token *string `json:"token,omitempty"`
}

// User defines the model for user.
type User struct {
	// Id The user id
	Id *int64 `json:"id,omitempty"`

	// Name The name of the user
	Name *string `json:"name,omitempty"`

	// Password The password
	Password *string `json:"password,omitempty"`

	// Username The username
	Username *string `json:"username,omitempty"`
}

// PathCourseId defines the model for path_course_id.
type PathCourseId = int64

// PathHoleId defines the model for path_hole_id.
type PathHoleId = int64

// PathRoundId defines the model for path_round_id.
type PathRoundId = int64

// QueryAverageType defines the model for query_average_type.
type QueryAverageType = AverageType

// QueryNameParam defines the model for query_name_param.
type QueryNameParam = string

// LoginJSONBody defines parameters for Login.
type LoginJSONBody struct {
	// Password The password
	Password *string `json:"password,omitempty"`

	// Username The username
	Username *string `json:"username,omitempty"`
}

// GetNewRoundCoursesParams defines parameters for GetNewRoundCourses.
type GetNewRoundCoursesParams struct {
	// Name The name of the club
	Name *QueryNameParam `form:"name,omitempty" json:"name,omitempty"`
}

// GetLineChartAveragesParams defines parameters for GetLineChartAverages.
type GetLineChartAveragesParams struct {
	// AverageType The type of average
	AverageType QueryAverageType `form:"average_type" json:"average_type"`

	// FromDate Filter by date, from date.
	FromDate *externalRef0.FromDate `form:"from_date,omitempty" json:"from_date,omitempty"`

	// Since Filter by the duration, since the current date. (E.g. 1d, 1w, 1m, 1y)
	Since *externalRef0.Since `form:"since,omitempty" json:"since,omitempty"`
}

// GetPieChartAveragesParams defines parameters for GetPieChartAverages.
type GetPieChartAveragesParams struct {
	// AverageType The type of average
	AverageType QueryAverageType `form:"average_type" json:"average_type"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody LoginJSONBody

// CreateRoundJSONRequestBody defines body for CreateRound for application/json ContentType.
type CreateRoundJSONRequestBody = RoundCreate

// UpdateHoleStatsJSONRequestBody defines body for UpdateHoleStats for application/json ContentType.
type UpdateHoleStatsJSONRequestBody = HoleStats

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = User
