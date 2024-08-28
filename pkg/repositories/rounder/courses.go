package rounder

import "github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"

func (r *repository) CreateCourse(course *models.Course) error {
	course.Id = 0
	return course.Insert(r.db)
}

func (r *repository) CreateCourseDetails(courseDetails *models.CourseDetails) error {
	courseDetails.Id = 0
	return courseDetails.Insert(r.db)
}
