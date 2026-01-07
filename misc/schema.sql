-- Database Schema for Dalivim Platform
-- PostgreSQL

-- Users table (Professors and Students)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) NOT NULL CHECK (role IN ('professor', 'student')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Activities table (Created by professors)
CREATE TABLE activities (
    id SERIAL PRIMARY KEY,
    professor_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    language VARCHAR(50) NOT NULL,
    time_limit INTEGER NOT NULL, -- in minutes
    invite_token VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Submissions table (Final student submissions with analysis)
CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    activity_id INTEGER REFERENCES activities(id) ON DELETE CASCADE,
    student_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    student_name VARCHAR(255),
    student_email VARCHAR(255),
    code TEXT,
    
    -- Authorship analysis
    authorship_score DECIMAL(5,4) DEFAULT 0,
    confidence VARCHAR(50),
    signals TEXT[], -- Array of signal strings
    
    -- Behavioral metrics
    avg_keystroke_interval DECIMAL(10,2),
    std_keystroke_interval DECIMAL(10,2),
    paste_events INTEGER DEFAULT 0,
    paste_char_ratio DECIMAL(5,4),
    delete_ratio DECIMAL(5,4),
    focus_loss_count INTEGER DEFAULT 0,
    linear_editing_score DECIMAL(5,4),
    burstiness DECIMAL(10,2),
    time_to_first_run DECIMAL(10,2),
    execution_count INTEGER DEFAULT 0,
    total_time DECIMAL(10,2),
    keystroke_count INTEGER DEFAULT 0,
    
    -- Raw event data
    paste_event_details JSONB,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Telemetry data table (Real-time behavioral data)
CREATE TABLE telemetry_data (
    id SERIAL PRIMARY KEY,
    activity_id INTEGER REFERENCES activities(id) ON DELETE CASCADE,
    student_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    timestamp BIGINT NOT NULL,
    is_final BOOLEAN DEFAULT FALSE,
    
    -- Features and raw events stored as JSONB
    features JSONB,
    raw_events JSONB,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_activities_professor ON activities(professor_id);
CREATE INDEX idx_activities_invite ON activities(invite_token);
CREATE INDEX idx_submissions_activity ON submissions(activity_id);
CREATE INDEX idx_submissions_student ON submissions(student_id);
CREATE INDEX idx_telemetry_activity ON telemetry_data(activity_id);
CREATE INDEX idx_telemetry_student ON telemetry_data(student_id);
CREATE INDEX idx_telemetry_timestamp ON telemetry_data(timestamp);

-- Example data

-- Insert a professor
INSERT INTO users (email, password, name, role) 
VALUES ('professor@dalivim.com', 'hashed_password', 'Prof. Silva', 'professor');

-- Insert an activity
INSERT INTO activities (professor_id, title, description, language, time_limit, invite_token)
VALUES (
    1,
    'Implementar Bubble Sort',
    'Crie uma função que ordena um array usando o algoritmo Bubble Sort. A função deve receber um array de números e retornar o array ordenado.',
    'python',
    60,
    'a1b2c3d4e5f6789012345678'
);

-- Insert a student
INSERT INTO users (email, password, name, role)
VALUES ('student_abc123@anonymous.local', '', 'Anonymous Student', 'student');

-- Insert a submission with analysis
INSERT INTO submissions (
    activity_id,
    student_id,
    student_name,
    student_email,
    code,
    authorship_score,
    confidence,
    signals,
    avg_keystroke_interval,
    paste_events,
    paste_char_ratio,
    delete_ratio,
    linear_editing_score,
    execution_count,
    total_time,
    keystroke_count
) VALUES (
    1,
    2,
    'Anonymous Student',
    'student_abc123@anonymous.local',
    'def bubble_sort(arr):\n    n = len(arr)\n    for i in range(n):\n        for j in range(0, n-i-1):\n            if arr[j] > arr[j+1]:\n                arr[j], arr[j+1] = arr[j+1], arr[j]\n    return arr',
    0.73,
    'medium',
    ARRAY['moderate_paste_ratio', 'low_edit_ratio'],
    150.5,
    2,
    0.35,
    0.08,
    0.75,
    3,
    420.5,
    234
);

-- Query examples

-- Get all activities for a professor
SELECT a.*, COUNT(s.id) as submission_count
FROM activities a
LEFT JOIN submissions s ON a.id = s.activity_id
WHERE a.professor_id = 1
GROUP BY a.id
ORDER BY a.created_at DESC;

-- Get submissions for an activity with suspicion level
SELECT 
    s.*,
    CASE 
        WHEN s.authorship_score > 0.8 THEN 'Very Low'
        WHEN s.authorship_score > 0.6 THEN 'Low'
        WHEN s.authorship_score > 0.4 THEN 'Medium'
        WHEN s.authorship_score > 0.2 THEN 'High'
        ELSE 'Very High'
    END as suspicion_level
FROM submissions s
WHERE s.activity_id = 1
ORDER BY s.authorship_score ASC;

-- Get telemetry timeline for a student
SELECT 
    timestamp,
    features->>'avgKeystrokeInterval' as avg_keystroke,
    features->>'pasteCharRatio' as paste_ratio,
    is_final
FROM telemetry_data
WHERE activity_id = 1 AND student_id = 2
ORDER BY timestamp ASC;

-- Find high-risk submissions (low authorship score + multiple signals)
SELECT 
    s.student_name,
    s.authorship_score,
    s.signals,
    a.title as activity_title
FROM submissions s
JOIN activities a ON s.activity_id = a.id
WHERE s.authorship_score < 0.5
  AND array_length(s.signals, 1) >= 3
ORDER BY s.authorship_score ASC;

-- Average metrics by activity
SELECT 
    a.title,
    COUNT(s.id) as total_submissions,
    AVG(s.authorship_score) as avg_authorship,
    AVG(s.paste_char_ratio) as avg_paste_ratio,
    AVG(s.total_time) as avg_time_seconds
FROM activities a
LEFT JOIN submissions s ON a.id = s.activity_id
GROUP BY a.id, a.title
ORDER BY avg_authorship ASC;
