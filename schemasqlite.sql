
CREATE TABLE tenants (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT,
    phone TEXT,
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);


INSERT INTO roles (name) VALUES
('super_admin'),
('school_admin'),
('teacher'),
('student'),
('department_head');

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email TEXT,
    role_id INTEGER,
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE departments (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    name TEXT NOT NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE TABLE sections (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    name TEXT NOT NULL,
    department_id TEXT,
    grade_level TEXT,
    academic_year TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (department_id) REFERENCES departments(id)
);

CREATE TABLE students (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    student_code TEXT UNIQUE,
    first_name TEXT,
    last_name TEXT,
    year TEXT,
    section_id TEXT,
    department_id TEXT,
    user_id TEXT,
    status TEXT DEFAULT 'active',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (section_id) REFERENCES sections(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE teachers (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    first_name TEXT,
    last_name TEXT,
    department_id TEXT,
    user_id TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE subjects (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    name TEXT NOT NULL,
    department_id TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (department_id) REFERENCES departments(id)
);

CREATE TABLE teacher_subjects (
    id TEXT PRIMARY KEY,
    teacher_id TEXT,
    subject_id TEXT,
    section_id TEXT,
    FOREIGN KEY (teacher_id) REFERENCES teachers(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (section_id) REFERENCES sections(id)
);

CREATE TABLE question_banks (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    subject_id TEXT,
    created_by_teacher_id TEXT,
    title TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (created_by_teacher_id) REFERENCES teachers(id)
);

CREATE TABLE questions (
    id TEXT PRIMARY KEY,
    question_bank_id TEXT,
    type TEXT CHECK (type IN ('mcq','true_false','short','essay')),
    question_text TEXT,
    marks INTEGER,
    difficulty_level TEXT,
    FOREIGN KEY (question_bank_id) REFERENCES question_banks(id)
);

CREATE TABLE question_options (
    id TEXT PRIMARY KEY,
    question_id TEXT,
    option_text TEXT,
    is_correct INTEGER DEFAULT 0,
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

CREATE TABLE exams (
    id TEXT PRIMARY KEY,
    tenant_id TEXT,
    title TEXT,
    subject_id TEXT,
    section_id TEXT,
    created_by_teacher_id TEXT,
    duration_minutes INTEGER,
    start_time DATETIME,
    end_time DATETIME,
    status TEXT DEFAULT 'draft',
    total_marks INTEGER,
    shuffle_options INTEGER DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (section_id) REFERENCES sections(id),
    FOREIGN KEY (created_by_teacher_id) REFERENCES teachers(id)
);

CREATE TABLE exam_questions (
    id TEXT PRIMARY KEY,
    exam_id TEXT,
    question_id TEXT,
    marks INTEGER,
    order_index INTEGER,
    FOREIGN KEY (exam_id) REFERENCES exams(id),
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

CREATE TABLE exam_sessions (
    id TEXT PRIMARY KEY,
    exam_id TEXT,
    student_id TEXT,
    start_time DATETIME,
    end_time DATETIME,
    status TEXT DEFAULT 'in_progress',
    total_score INTEGER,
    FOREIGN KEY (exam_id) REFERENCES exams(id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE TABLE student_answers (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    question_id TEXT,
    answer_text TEXT,
    selected_option_id TEXT,
    score INTEGER,
    FOREIGN KEY (session_id) REFERENCES exam_sessions(id),
    FOREIGN KEY (question_id) REFERENCES questions(id),
    FOREIGN KEY (selected_option_id) REFERENCES question_options(id)
);


CREATE TABLE results (
    id TEXT PRIMARY KEY,
    exam_id TEXT,
    student_id TEXT,
    total_score INTEGER,
    grade TEXT,
    published_at DATETIME,
    FOREIGN KEY (exam_id) REFERENCES exams(id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE TABLE audit_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    action TEXT,
    entity_type TEXT,
    entity_id TEXT,
    ip_address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE exam_violations (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    violation_type TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES exam_sessions(id)
);


CREATE INDEX idx_users_tenant ON users(tenant_id);
CREATE INDEX idx_students_section ON students(section_id);
CREATE INDEX idx_exam_sessions_exam ON exam_sessions(exam_id);



