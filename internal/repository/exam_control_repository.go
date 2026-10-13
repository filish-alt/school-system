package repository

import (
	"context"
	"database/sql"
	"school-exam/internal/domain"
	"time"
)

type ExamControlRepository struct {
	DB *sql.DB
}

func NewExamControlRepository(db *sql.DB) *ExamControlRepository {
	return &ExamControlRepository{DB: db}
}

// 1. Extensions
func (r *ExamControlRepository) AddExtension(ctx context.Context, ext domain.ExamExtension) error {
	query := `INSERT INTO exam_extensions (id, exam_id, student_id, extended_minutes, reason, authorized_by, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, ext.ID, ext.ExamID, ext.StudentID, ext.ExtendedMinutes, ext.Reason, ext.AuthorizedBy, time.Now())
	return err
}

func (r *ExamControlRepository) ListExtensions(ctx context.Context, examID string) ([]domain.ExamExtension, error) {
	query := `SELECT e.id, e.exam_id, e.student_id, (s.first_name || ' ' || s.last_name) as student_name,
	                 e.extended_minutes, e.reason, e.authorized_by, u.username as authorizer_name, e.created_at
	          FROM exam_extensions e
	          LEFT JOIN students s ON e.student_id = s.id
	          LEFT JOIN users u ON e.authorized_by = u.id
	          WHERE e.exam_id = ?
	          ORDER BY e.created_at DESC`
	rows, err := r.DB.QueryContext(ctx, query, examID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.ExamExtension
	for rows.Next() {
		var ext domain.ExamExtension
		var stuID, stuName, authName sql.NullString
		if err := rows.Scan(&ext.ID, &ext.ExamID, &stuID, &stuName, &ext.ExtendedMinutes, &ext.Reason, &ext.AuthorizedBy, &authName, &ext.CreatedAt); err != nil {
			return nil, err
		}
		if stuID.Valid { ext.StudentID = &stuID.String }
		if stuName.Valid { ext.StudentName = &stuName.String }
		if authName.Valid { ext.AuthorizerName = &authName.String }
		list = append(list, ext)
	}
	return list, nil
}

// 2. Remote Access
func (r *ExamControlRepository) AddRemoteAuthorization(ctx context.Context, auth domain.ExamRemoteAuthorization) error {
	query := `INSERT INTO exam_remote_authorizations (id, exam_id, student_id, valid_until, reason, approved_by, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, auth.ID, auth.ExamID, auth.StudentID, auth.ValidUntil, auth.Reason, auth.ApprovedBy, time.Now())
	return err
}

func (r *ExamControlRepository) ListRemoteAuthorizations(ctx context.Context, examID string) ([]domain.ExamRemoteAuthorization, error) {
	query := `SELECT r.id, r.exam_id, r.student_id, s.student_code, (s.first_name || ' ' || s.last_name) as student_name,
	                 r.valid_until, r.reason, r.approved_by, u.username as authorizer_name, r.created_at
	          FROM exam_remote_authorizations r
	          JOIN students s ON r.student_id = s.id
	          LEFT JOIN users u ON r.approved_by = u.id
	          WHERE r.exam_id = ?
	          ORDER BY r.created_at DESC`
	rows, err := r.DB.QueryContext(ctx, query, examID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.ExamRemoteAuthorization
	for rows.Next() {
		var a domain.ExamRemoteAuthorization
		var code, name, authName sql.NullString
		if err := rows.Scan(&a.ID, &a.ExamID, &a.StudentID, &code, &name, &a.ValidUntil, &a.Reason, &a.ApprovedBy, &authName, &a.CreatedAt); err != nil {
			return nil, err
		}
		if code.Valid { a.StudentCode = &code.String }
		if name.Valid { a.StudentName = &name.String }
		if authName.Valid { a.AuthorizerName = &authName.String }
		list = append(list, a)
	}
	return list, nil
}

// 3. Dropped Questions & Recalculation
func (r *ExamControlRepository) DropQuestion(ctx context.Context, dq domain.DroppedQuestion) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert dropped_question record
	if _, err := tx.ExecContext(ctx, `INSERT INTO dropped_questions (id, exam_id, question_id, reason, credit_type, dropped_by, created_at)
	                                 VALUES (?, ?, ?, ?, ?, ?, ?)`, dq.ID, dq.ExamID, dq.QuestionID, dq.Reason, dq.CreditType, dq.DroppedBy, time.Now()); err != nil {
		return err
	}

	// Mark in exam_questions
	if _, err := tx.ExecContext(ctx, "UPDATE exam_questions SET is_dropped = 1 WHERE exam_id = ? AND question_id = ?", dq.ExamID, dq.QuestionID); err != nil {
		return err
	}

	// Recalculate score for all sessions of this exam
	// If full_credit: grant full question marks to all student answers
	// If reduce_total: set score to 0 and reduce exam total
	var qMarks float64
	_ = tx.QueryRowContext(ctx, "SELECT marks FROM exam_questions WHERE exam_id = ? AND question_id = ?", dq.ExamID, dq.QuestionID).Scan(&qMarks)

	if dq.CreditType == "full_credit" {
		_, _ = tx.ExecContext(ctx, `UPDATE student_answers SET score = ?
		                            WHERE question_id = ? AND session_id IN (SELECT id FROM exam_sessions WHERE exam_id = ?)`,
			qMarks, dq.QuestionID, dq.ExamID)
	} else {
		_, _ = tx.ExecContext(ctx, `UPDATE student_answers SET score = 0
		                            WHERE question_id = ? AND session_id IN (SELECT id FROM exam_sessions WHERE exam_id = ?)`,
			dq.QuestionID, dq.ExamID)
	}

	// Recompute total_score in exam_sessions
	_, _ = tx.ExecContext(ctx, `UPDATE exam_sessions SET total_score = (
		SELECT COALESCE(SUM(score), 0) FROM student_answers WHERE student_answers.session_id = exam_sessions.id
	) WHERE exam_id = ?`, dq.ExamID)

	return tx.Commit()
}

// 4. Answer Key Corrections & Automatic Recalculation
func (r *ExamControlRepository) CorrectAnswerKey(ctx context.Context, corr domain.AnswerKeyCorrection) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `INSERT INTO answer_key_corrections (id, exam_id, question_id, old_answer, new_answer, reason, corrected_by, created_at)
	                                 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, corr.ID, corr.ExamID, corr.QuestionID, corr.OldAnswer, corr.NewAnswer, corr.Reason, corr.CorrectedBy, time.Now()); err != nil {
		return err
	}

	// Update question_options: clear is_correct for this question and set for the new option
	if corr.NewAnswer != nil && *corr.NewAnswer != "" {
		_, _ = tx.ExecContext(ctx, "UPDATE question_options SET is_correct = 0 WHERE question_id = ?", corr.QuestionID)
		_, _ = tx.ExecContext(ctx, "UPDATE question_options SET is_correct = 1 WHERE id = ?", *corr.NewAnswer)

		// Get marks for this question
		var qMarks float64 = 1.0
		_ = tx.QueryRowContext(ctx, "SELECT marks FROM exam_questions WHERE exam_id = ? AND question_id = ?", corr.ExamID, corr.QuestionID).Scan(&qMarks)

		// Re-score answers for this question
		_, _ = tx.ExecContext(ctx, `UPDATE student_answers
		                            SET score = CASE WHEN selected_option_id = ? THEN ? ELSE 0 END
		                            WHERE question_id = ? AND session_id IN (SELECT id FROM exam_sessions WHERE exam_id = ?)`,
			*corr.NewAnswer, qMarks, corr.QuestionID, corr.ExamID)

		// Recalculate session totals
		_, _ = tx.ExecContext(ctx, `UPDATE exam_sessions SET total_score = (
			SELECT COALESCE(SUM(score), 0) FROM student_answers WHERE student_answers.session_id = exam_sessions.id
		) WHERE exam_id = ?`, corr.ExamID)
	}

	return tx.Commit()
}

// 5. Exam Workflow State (Start, Lock, Approve, Visibility)
func (r *ExamControlRepository) UpdateExamControls(ctx context.Context, examID string, isStarted, isVisible, isLocked *bool, status, approvalStatus, accessCode *string, approvedBy *string) error {
	query := `UPDATE exams SET
	          is_started = COALESCE(?, is_started),
	          is_visible = COALESCE(?, is_visible),
	          is_locked = COALESCE(?, is_locked),
	          status = COALESCE(?, status),
	          approval_status = COALESCE(?, approval_status),
	          access_code = COALESCE(?, access_code),
	          approved_by = COALESCE(?, approved_by),
	          approved_at = CASE WHEN ? IS NOT NULL THEN CURRENT_TIMESTAMP ELSE approved_at END
	          WHERE id = ?`
	var startVal, visVal, lockVal *int
	if isStarted != nil { v := 0; if *isStarted { v = 1 }; startVal = &v }
	if isVisible != nil { v := 0; if *isVisible { v = 1 }; visVal = &v }
	if isLocked != nil { v := 0; if *isLocked { v = 1 }; lockVal = &v }

	_, err := r.DB.ExecContext(ctx, query, startVal, visVal, lockVal, status, approvalStatus, accessCode, approvedBy, approvedBy, examID)
	return err
}

// 6. Matching Pairs
func (r *ExamControlRepository) AddMatchPairs(ctx context.Context, pairs []domain.QuestionMatchPair) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO question_match_pairs (id, question_id, left_prompt, right_match, order_index) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range pairs {
		if _, err := stmt.ExecContext(ctx, p.ID, p.QuestionID, p.LeftPrompt, p.RightMatch, p.OrderIndex); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *ExamControlRepository) ListMatchPairs(ctx context.Context, questionID string) ([]domain.QuestionMatchPair, error) {
	query := `SELECT id, question_id, left_prompt, right_match, order_index
	          FROM question_match_pairs WHERE question_id = ? ORDER BY order_index ASC`
	rows, err := r.DB.QueryContext(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.QuestionMatchPair
	for rows.Next() {
		var p domain.QuestionMatchPair
		if err := rows.Scan(&p.ID, &p.QuestionID, &p.LeftPrompt, &p.RightMatch, &p.OrderIndex); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
