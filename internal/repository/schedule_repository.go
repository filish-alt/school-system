package repository

import (
	"context"
	"database/sql"
	"school-exam/internal/domain"
)

type ScheduleRepository struct {
	DB *sql.DB
}

func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{DB: db}
}

func (r *ScheduleRepository) Create(ctx context.Context, s domain.TeacherSchedule) error {
	query := `INSERT INTO teacher_schedules (id, tenant_id, teacher_id, subject_id, section_id, day_of_week, start_time, end_time, room, activity_type)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, s.ID, s.TenantID, s.TeacherID, s.SubjectID, s.SectionID, s.DayOfWeek, s.StartTime, s.EndTime, s.Room, s.ActivityType)
	return err
}

func (r *ScheduleRepository) ListByTeacher(ctx context.Context, teacherID string) ([]domain.TeacherSchedule, error) {
	query := `SELECT sc.id, sc.tenant_id, sc.teacher_id, (t.first_name || ' ' || t.last_name) as teacher_name,
	                 sc.subject_id, sub.name as subject_name, sc.section_id, sec.name as section_name,
	                 sc.day_of_week, sc.start_time, sc.end_time, sc.room, sc.activity_type
	          FROM teacher_schedules sc
	          JOIN teachers t ON sc.teacher_id = t.id
	          LEFT JOIN subjects sub ON sc.subject_id = sub.id
	          LEFT JOIN sections sec ON sc.section_id = sec.id
	          WHERE sc.teacher_id = ?
	          ORDER BY CASE sc.day_of_week
	            WHEN 'Monday' THEN 1
	            WHEN 'Tuesday' THEN 2
	            WHEN 'Wednesday' THEN 3
	            WHEN 'Thursday' THEN 4
	            WHEN 'Friday' THEN 5
	            WHEN 'Saturday' THEN 6
	            WHEN 'Sunday' THEN 7
	            ELSE 8
	          END, sc.start_time ASC`
	rows, err := r.DB.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.TeacherSchedule
	for rows.Next() {
		var s domain.TeacherSchedule
		var teaName, subID, subName, secID, secName, room sql.NullString
		if err := rows.Scan(&s.ID, &s.TenantID, &s.TeacherID, &teaName,
			&subID, &subName, &secID, &secName,
			&s.DayOfWeek, &s.StartTime, &s.EndTime, &room, &s.ActivityType); err != nil {
			return nil, err
		}
		if teaName.Valid { s.TeacherName = &teaName.String }
		if subID.Valid { s.SubjectID = &subID.String }
		if subName.Valid { s.SubjectName = &subName.String }
		if secID.Valid { s.SectionID = &secID.String }
		if secName.Valid { s.SectionName = &secName.String }
		if room.Valid { s.Room = &room.String }
		list = append(list, s)
	}
	return list, nil
}

func (r *ScheduleRepository) ListBySection(ctx context.Context, sectionID string) ([]domain.TeacherSchedule, error) {
	query := `SELECT sc.id, sc.tenant_id, sc.teacher_id, (t.first_name || ' ' || t.last_name) as teacher_name,
	                 sc.subject_id, sub.name as subject_name, sc.section_id, sec.name as section_name,
	                 sc.day_of_week, sc.start_time, sc.end_time, sc.room, sc.activity_type
	          FROM teacher_schedules sc
	          JOIN teachers t ON sc.teacher_id = t.id
	          LEFT JOIN subjects sub ON sc.subject_id = sub.id
	          LEFT JOIN sections sec ON sc.section_id = sec.id
	          WHERE sc.section_id = ?
	          ORDER BY sc.day_of_week, sc.start_time ASC`
	rows, err := r.DB.QueryContext(ctx, query, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.TeacherSchedule
	for rows.Next() {
		var s domain.TeacherSchedule
		var teaName, subID, subName, secID, secName, room sql.NullString
		if err := rows.Scan(&s.ID, &s.TenantID, &s.TeacherID, &teaName,
			&subID, &subName, &secID, &secName,
			&s.DayOfWeek, &s.StartTime, &s.EndTime, &room, &s.ActivityType); err != nil {
			return nil, err
		}
		if teaName.Valid { s.TeacherName = &teaName.String }
		if subID.Valid { s.SubjectID = &subID.String }
		if subName.Valid { s.SubjectName = &subName.String }
		if secID.Valid { s.SectionID = &secID.String }
		if secName.Valid { s.SectionName = &secName.String }
		if room.Valid { s.Room = &room.String }
		list = append(list, s)
	}
	return list, nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM teacher_schedules WHERE id = ?", id)
	return err
}
