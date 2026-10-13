
-- 1. Tenants (Schools)
CREATE TABLE IF NOT EXISTS tenants (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT,
    phone TEXT,
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Roles
CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

INSERT OR IGNORE INTO roles (id, name) VALUES
(1, 'super_admin'),
(2, 'school_director'),
(3, 'vice_director'),
(4, 'department_head'),
(5, 'teacher'),
(6, 'student'),
(7, 'school_admin');

-- 3. Academic Years
CREATE TABLE IF NOT EXISTS academic_years (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    year_name TEXT NOT NULL,
    start_date DATE,
    end_date DATE,
    is_active INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- 4. System Settings (Configuration-driven rules)
CREATE TABLE IF NOT EXISTS system_settings (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    setting_key TEXT NOT NULL,
    setting_value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    UNIQUE(tenant_id, setting_key)
);

-- 5. Users
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    profile_photo TEXT,
    role_id INTEGER,
    status TEXT DEFAULT 'active',
    must_change_password INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- 6. Departments
CREATE TABLE IF NOT EXISTS departments (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    name TEXT NOT NULL,
    head_teacher_id TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- 7. Sections / Classes
CREATE TABLE IF NOT EXISTS sections (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    name TEXT NOT NULL,
    department_id TEXT,
    grade_level TEXT,
    academic_year TEXT,
    academic_year_id TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (academic_year_id) REFERENCES academic_years(id)
);

-- 8. Students
CREATE TABLE IF NOT EXISTS students (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    student_code TEXT UNIQUE,
    first_name TEXT,
    last_name TEXT,
    year TEXT,
    phone TEXT,
    photo_url TEXT,
    section_id TEXT,
    department_id TEXT,
    academic_year_id TEXT,
    user_id TEXT,
    status TEXT DEFAULT 'active',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (section_id) REFERENCES sections(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (academic_year_id) REFERENCES academic_years(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 9. Teachers
CREATE TABLE IF NOT EXISTS teachers (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    teacher_code TEXT UNIQUE,
    first_name TEXT,
    last_name TEXT,
    phone TEXT,
    photo_url TEXT,
    department_id TEXT,
    user_id TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 10. Subjects
CREATE TABLE IF NOT EXISTS subjects (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    name TEXT NOT NULL,
    code TEXT,
    department_id TEXT,
    grade_level TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (department_id) REFERENCES departments(id)
);

-- 11. Teacher Assignments (Subjects + Sections)
CREATE TABLE IF NOT EXISTS teacher_subjects (
    id TEXT PRIMARY KEY,
    teacher_id TEXT,
    subject_id TEXT,
    section_id TEXT,
    academic_year_id TEXT,
    FOREIGN KEY (teacher_id) REFERENCES teachers(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (section_id) REFERENCES sections(id),
    FOREIGN KEY (academic_year_id) REFERENCES academic_years(id)
);

-- 12. Teacher Timetables / Schedules
CREATE TABLE IF NOT EXISTS teacher_schedules (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    teacher_id TEXT NOT NULL,
    subject_id TEXT,
    section_id TEXT,
    day_of_week TEXT NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    room TEXT,
    activity_type TEXT DEFAULT 'lecture',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (teacher_id) REFERENCES teachers(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (section_id) REFERENCES sections(id)
);

-- 13. Question Banks
CREATE TABLE IF NOT EXISTS question_banks (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    subject_id TEXT,
    department_id TEXT,
    created_by_teacher_id TEXT,
    title TEXT,
    grade_level TEXT,
    academic_year_id TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (academic_year_id) REFERENCES academic_years(id),
    FOREIGN KEY (created_by_teacher_id) REFERENCES teachers(id)
);

-- 14. Questions
CREATE TABLE IF NOT EXISTS questions (
    id TEXT PRIMARY KEY,
    question_bank_id TEXT,
    type TEXT CHECK (type IN ('mcq','true_false','fill_blank','matching','short','essay')),
    question_text TEXT,
    image_url TEXT,
    marks REAL DEFAULT 1.0,
    difficulty_level TEXT DEFAULT 'medium',
    topic TEXT,
    grade_level TEXT,
    academic_year_id TEXT,
    creator_teacher_id TEXT,
    is_private INTEGER DEFAULT 1,
    status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'submitted_for_review', 'approved_for_final', 'archived')),
    answer_key TEXT,
    case_sensitive INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (question_bank_id) REFERENCES question_banks(id),
    FOREIGN KEY (academic_year_id) REFERENCES academic_years(id),
    FOREIGN KEY (creator_teacher_id) REFERENCES teachers(id)
);

-- 15. Question Options (MCQ & True/False)
CREATE TABLE IF NOT EXISTS question_options (
    id TEXT PRIMARY KEY,
    question_id TEXT,
    option_text TEXT,
    is_correct INTEGER DEFAULT 0,
    order_index INTEGER DEFAULT 0,
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

-- 16. Question Match Pairs (Matching Type)
CREATE TABLE IF NOT EXISTS question_match_pairs (
    id TEXT PRIMARY KEY,
    question_id TEXT NOT NULL,
    left_prompt TEXT NOT NULL,
    right_match TEXT NOT NULL,
    order_index INTEGER DEFAULT 0,
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

-- 17. Exams
CREATE TABLE IF NOT EXISTS exams (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    title TEXT NOT NULL,
    exam_type TEXT DEFAULT 'mid' CHECK (exam_type IN ('quiz', 'mid', 'final')),
    subject_id TEXT NOT NULL,
    section_id TEXT NOT NULL,
    department_id TEXT,
    academic_year_id TEXT,
    created_by_teacher_id TEXT,
    duration_minutes INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'under_review', 'approved', 'rejected', 'scheduled', 'in_progress', 'completed', 'published', 'archived')),
    approval_status TEXT DEFAULT 'draft' CHECK (approval_status IN ('draft', 'under_review', 'approved', 'rejected')),
    is_locked INTEGER DEFAULT 0,
    is_visible INTEGER DEFAULT 0,
    is_started INTEGER DEFAULT 0,
    access_code TEXT,
    instructions TEXT,
    total_marks REAL,
    pass_percentage REAL DEFAULT 50.0,
    shuffle_options INTEGER DEFAULT 0,
    version INTEGER DEFAULT 1,
    approved_by TEXT,
    approved_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (section_id) REFERENCES sections(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (academic_year_id) REFERENCES academic_years(id),
    FOREIGN KEY (created_by_teacher_id) REFERENCES teachers(id)
);

-- 18. Exam Questions
CREATE TABLE IF NOT EXISTS exam_questions (
    id TEXT PRIMARY KEY,
    exam_id TEXT,
    question_id TEXT,
    marks REAL,
    order_index INTEGER,
    is_dropped INTEGER DEFAULT 0,
    FOREIGN KEY (exam_id) REFERENCES exams(id),
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

-- 19. Exam Sessions (Student Exam Instances)
CREATE TABLE IF NOT EXISTS exam_sessions (
    id TEXT PRIMARY KEY,
    exam_id TEXT,
    student_id TEXT,
    start_time DATETIME,
    end_time DATETIME,
    status TEXT DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'submitted', 'timed_out', 'cancelled')),
    total_score REAL,
    is_remote INTEGER DEFAULT 0,
    FOREIGN KEY (exam_id) REFERENCES exams(id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

-- 20. Student Answers
CREATE TABLE IF NOT EXISTS student_answers (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    question_id TEXT,
    answer_text TEXT,
    selected_option_id TEXT,
    matching_answers TEXT,
    score REAL,
    is_manually_graded INTEGER DEFAULT 0,
    FOREIGN KEY (session_id) REFERENCES exam_sessions(id),
    FOREIGN KEY (question_id) REFERENCES questions(id),
    FOREIGN KEY (selected_option_id) REFERENCES question_options(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_student_answer_unique
ON student_answers(session_id, question_id);

-- 21. Exam Time Extensions
CREATE TABLE IF NOT EXISTS exam_extensions (
    id TEXT PRIMARY KEY,
    exam_id TEXT NOT NULL,
    student_id TEXT,
    extended_minutes INTEGER NOT NULL,
    reason TEXT NOT NULL,
    authorized_by TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (exam_id) REFERENCES exams(id),
    FOREIGN KEY (authorized_by) REFERENCES users(id)
);

-- 22. Exam Remote Access Authorizations
CREATE TABLE IF NOT EXISTS exam_remote_authorizations (
    id TEXT PRIMARY KEY,
    exam_id TEXT NOT NULL,
    student_id TEXT NOT NULL,
    valid_until DATETIME NOT NULL,
    reason TEXT NOT NULL,
    approved_by TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (exam_id) REFERENCES exams(id),
    FOREIGN KEY (student_id) REFERENCES students(id),
    FOREIGN KEY (approved_by) REFERENCES users(id)
);

-- 23. Continuous Assessment / Coursework (Quizzes, Mid, Assignments, Projects)
CREATE TABLE IF NOT EXISTS coursework_assessments (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    academic_year_id TEXT,
    subject_id TEXT NOT NULL,
    section_id TEXT NOT NULL,
    teacher_id TEXT NOT NULL,
    title TEXT NOT NULL,
    type TEXT CHECK (type IN ('quiz', 'mid_exam', 'assignment', 'project')),
    max_marks REAL NOT NULL,
    weight_percentage REAL NOT NULL,
    due_date DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (academic_year_id) REFERENCES academic_years(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (section_id) REFERENCES sections(id),
    FOREIGN KEY (teacher_id) REFERENCES teachers(id)
);

-- 24. Student Coursework Marks
CREATE TABLE IF NOT EXISTS student_coursework_marks (
    id TEXT PRIMARY KEY,
    coursework_id TEXT NOT NULL,
    student_id TEXT NOT NULL,
    marks_obtained REAL NOT NULL,
    remarks TEXT,
    graded_by TEXT,
    graded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (coursework_id) REFERENCES coursework_assessments(id),
    FOREIGN KEY (student_id) REFERENCES students(id),
    UNIQUE(coursework_id, student_id)
);

-- 25. Dropped Questions Record
CREATE TABLE IF NOT EXISTS dropped_questions (
    id TEXT PRIMARY KEY,
    exam_id TEXT NOT NULL,
    question_id TEXT NOT NULL,
    reason TEXT NOT NULL,
    credit_type TEXT DEFAULT 'full_credit' CHECK (credit_type IN ('full_credit', 'reduce_total')),
    dropped_by TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (exam_id) REFERENCES exams(id),
    FOREIGN KEY (question_id) REFERENCES questions(id),
    FOREIGN KEY (dropped_by) REFERENCES users(id)
);

-- 26. Answer Key Corrections Record
CREATE TABLE IF NOT EXISTS answer_key_corrections (
    id TEXT PRIMARY KEY,
    exam_id TEXT NOT NULL,
    question_id TEXT NOT NULL,
    old_answer TEXT,
    new_answer TEXT,
    reason TEXT NOT NULL,
    corrected_by TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (exam_id) REFERENCES exams(id),
    FOREIGN KEY (question_id) REFERENCES questions(id),
    FOREIGN KEY (corrected_by) REFERENCES users(id)
);

-- 27. Student Complaints
CREATE TABLE IF NOT EXISTS student_complaints (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    student_id TEXT NOT NULL,
    subject_id TEXT,
    exam_id TEXT,
    coursework_id TEXT,
    question_id TEXT,
    category TEXT CHECK (category IN ('exam_result', 'question_key', 'assignment_mark', 'project_mark', 'other')),
    description TEXT NOT NULL,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'escalated_to_dept_head', 'resolved', 'rejected')),
    response TEXT,
    resolved_by TEXT,
    resolved_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

-- 28. Official Results (CA + Final aggregation)
CREATE TABLE IF NOT EXISTS results (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    academic_year_id TEXT,
    student_id TEXT NOT NULL,
    subject_id TEXT NOT NULL,
    section_id TEXT,
    quiz_score REAL DEFAULT 0,
    mid_score REAL DEFAULT 0,
    assignment_score REAL DEFAULT 0,
    project_score REAL DEFAULT 0,
    ca_score REAL DEFAULT 0,
    final_score REAL DEFAULT 0,
    total_score REAL DEFAULT 0,
    grade TEXT,
    is_published INTEGER DEFAULT 0,
    published_at DATETIME,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (academic_year_id) REFERENCES academic_years(id),
    FOREIGN KEY (student_id) REFERENCES students(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (section_id) REFERENCES sections(id),
    UNIQUE(student_id, subject_id, academic_year_id)
);

-- 29. Audit Logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    user_id TEXT,
    action TEXT NOT NULL,
    entity_type TEXT,
    entity_id TEXT,
    details TEXT,
    ip_address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 30. Exam Violations
CREATE TABLE IF NOT EXISTS exam_violations (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    violation_type TEXT,
    details TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES exam_sessions(id)
);

-- 31. Database Performance Indexes
CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_students_section ON students(section_id);
CREATE INDEX IF NOT EXISTS idx_students_tenant ON students(tenant_id);
CREATE INDEX IF NOT EXISTS idx_teachers_dept ON teachers(department_id);
CREATE INDEX IF NOT EXISTS idx_exam_sessions_exam ON exam_sessions(exam_id);
CREATE INDEX IF NOT EXISTS idx_exam_sessions_student ON exam_sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_coursework_section ON coursework_assessments(section_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant ON audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_complaints_student ON student_complaints(student_id);
