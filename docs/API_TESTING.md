# API Testing Guide

## Base URL
```
http://localhost:8080
```

## Authentication Flow

### 1. Register Professor
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "professor@dalivim.com",
    "password": "senha123",
    "name": "Prof. Silva",
    "role": "professor"
  }'
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "email": "professor@dalivim.com",
    "name": "Prof. Silva",
    "role": "professor",
    "createdAt": "2026-01-04T10:00:00Z"
  },
  "token": "abc123def456..."
}
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "professor@dalivim.com",
    "password": "senha123"
  }'
```

## Activity Management

### 3. Create Activity
```bash
curl -X POST http://localhost:8080/api/activities \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Implementar Bubble Sort",
    "description": "Crie uma função que ordena um array usando o algoritmo Bubble Sort.",
    "language": "python",
    "timeLimit": 60
  }'
```

**Response:**
```json
{
  "id": 1,
  "professorId": 1,
  "title": "Implementar Bubble Sort",
  "description": "Crie uma função que ordena um array usando o algoritmo Bubble Sort.",
  "language": "python",
  "timeLimit": 60,
  "inviteToken": "a1b2c3d4e5f6789012345678",
  "createdAt": "2026-01-04T10:05:00Z"
}
```

### 4. List Activities
```bash
curl -X GET http://localhost:8080/api/activities \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
[
  {
    "id": 1,
    "title": "Implementar Bubble Sort",
    "description": "...",
    "language": "python",
    "timeLimit": 60,
    "inviteToken": "a1b2c3d4e5f6789012345678",
    "submissionCount": 3
  }
]
```

### 5. Get Activity Details
```bash
curl -X GET http://localhost:8080/api/activities/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Student Flow

### 6. Join Activity (via invite link)
```bash
curl -X POST http://localhost:8080/api/activities/join/a1b2c3d4e5f6789012345678
```

**Response:**
```json
{
  "activity": {
    "id": 1,
    "title": "Implementar Bubble Sort",
    "description": "...",
    "language": "python",
    "timeLimit": 60
  },
  "student": {
    "id": 2,
    "email": "student_abc123@anonymous.local",
    "name": "Anonymous Student",
    "role": "student"
  }
}
```

### 7. Send Telemetry (every 10 seconds)
```bash
curl -X POST http://localhost:8080/api/telemetry \
  -H "Content-Type: application/json" \
  -d '{
    "activityId": 1,
    "studentId": 2,
    "timestamp": 1704358800000,
    "isFinal": false,
    "features": {
      "avgKeystrokeInterval": 150.5,
      "stdKeystrokeInterval": 45.2,
      "pasteEvents": 1,
      "pasteCharRatio": 0.25,
      "deleteRatio": 0.12,
      "focusLossCount": 2,
      "linearEditingScore": 0.65,
      "burstiness": 0.48,
      "timeToFirstRun": 120.5,
      "executionCount": 2,
      "totalTime": 180.0,
      "totalKeystrokes": 150,
      "codeLength": 200
    },
    "rawEvents": {
      "pasteEvents": [
        {
          "timestamp": 1704358700000,
          "length": 50,
          "content": "def bubble_sort(arr):",
          "linesCount": 1
        }
      ],
      "focusEvents": [
        {
          "type": "blur",
          "timestamp": 1704358650000
        },
        {
          "type": "focus",
          "timestamp": 1704358670000,
          "awayDuration": 20000
        }
      ]
    }
  }'
```

**Response:**
```json
{
  "authorship_score": 0.73,
  "confidence": "medium",
  "signals": [
    "moderate_paste_ratio",
    "low_edit_ratio"
  ]
}
```

### 8. Final Submission
```bash
curl -X POST http://localhost:8080/api/telemetry \
  -H "Content-Type: application/json" \
  -d '{
    "activityId": 1,
    "studentId": 2,
    "timestamp": 1704359000000,
    "isFinal": true,
    "code": "def bubble_sort(arr):\n    n = len(arr)\n    for i in range(n):\n        for j in range(0, n-i-1):\n            if arr[j] > arr[j+1]:\n                arr[j], arr[j+1] = arr[j+1], arr[j]\n    return arr",
    "features": {
      "avgKeystrokeInterval": 155.3,
      "stdKeystrokeInterval": 48.7,
      "pasteEvents": 2,
      "pasteCharRatio": 0.35,
      "deleteRatio": 0.08,
      "focusLossCount": 3,
      "linearEditingScore": 0.75,
      "burstiness": 0.52,
      "timeToFirstRun": 300.0,
      "executionCount": 3,
      "totalTime": 420.5,
      "totalKeystrokes": 234,
      "codeLength": 250
    },
    "rawEvents": {
      "pasteEvents": [
        {
          "timestamp": 1704358700000,
          "length": 50,
          "content": "def bubble_sort(arr):",
          "linesCount": 1
        },
        {
          "timestamp": 1704358900000,
          "length": 80,
          "content": "    for i in range(n):\n        for j in range(0, n-i-1):",
          "linesCount": 2
        }
      ]
    }
  }'
```

## View Submissions

### 9. Get Submissions for Activity
```bash
curl -X GET http://localhost:8080/api/activities/1/submissions \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
[
  {
    "id": 1,
    "activityId": 1,
    "studentId": 2,
    "studentName": "Anonymous Student",
    "studentEmail": "student_abc123@anonymous.local",
    "code": "def bubble_sort(arr):\n    ...",
    "authorshipScore": 0.73,
    "confidence": "medium",
    "signals": [
      "moderate_paste_ratio",
      "low_edit_ratio"
    ],
    "avgKeystrokeInterval": 155.3,
    "stdKeystrokeInterval": 48.7,
    "pasteEvents": 2,
    "pasteCharRatio": 0.35,
    "deleteRatio": 0.08,
    "focusLossCount": 3,
    "linearEditingScore": 0.75,
    "burstiness": 0.52,
    "timeToFirstRun": 300.0,
    "executionCount": 3,
    "totalTime": 420.5,
    "keystrokeCount": 234,
    "pasteEventDetails": "[...]",
    "createdAt": "2026-01-04T10:15:00Z"
  }
]
```

## Piston Code Execution

### 10. Get Available Languages
```bash
curl https://emkc.org/api/v2/piston/runtimes
```

**Response:**
```json
[
  {
    "language": "python",
    "version": "3.10.0",
    "aliases": ["py", "py3", "python3"]
  },
  {
    "language": "javascript",
    "version": "18.15.0",
    "aliases": ["js", "node-javascript"]
  },
  ...
]
```

### 11. Execute Code
```bash
curl -X POST https://emkc.org/api/v2/piston/execute \
  -H "Content-Type: application/json" \
  -d '{
    "language": "python",
    "version": "*",
    "files": [
      {
        "name": "main.py",
        "content": "print(\"Hello, World!\")"
      }
    ]
  }'
```

**Response:**
```json
{
  "run": {
    "stdout": "Hello, World!\n",
    "stderr": "",
    "code": 0,
    "signal": null,
    "output": "Hello, World!\n"
  }
}
```

## Test Scenarios

### Scenario 1: High Suspicion (AI-Generated)
```json
{
  "features": {
    "avgKeystrokeInterval": 50.0,
    "pasteCharRatio": 0.85,
    "deleteRatio": 0.01,
    "linearEditingScore": 0.95,
    "executionCount": 0,
    "totalTime": 45.0
  }
}
```
**Expected Score:** < 0.3 (High suspicion)
**Signals:** `high_paste_ratio`, `low_edit_ratio`, `highly_linear_editing`, `fast_completion_no_testing`

### Scenario 2: Low Suspicion (Human-Written)
```json
{
  "features": {
    "avgKeystrokeInterval": 180.0,
    "stdKeystrokeInterval": 65.0,
    "pasteCharRatio": 0.05,
    "deleteRatio": 0.18,
    "linearEditingScore": 0.45,
    "burstiness": 0.65,
    "executionCount": 8,
    "totalTime": 900.0
  }
}
```
**Expected Score:** > 0.8 (Low suspicion)
**Signals:** None or minimal

### Scenario 3: Medium Suspicion (Mixed)
```json
{
  "features": {
    "avgKeystrokeInterval": 120.0,
    "pasteCharRatio": 0.40,
    "deleteRatio": 0.08,
    "linearEditingScore": 0.70,
    "executionCount": 3,
    "totalTime": 300.0
  }
}
```
**Expected Score:** 0.5-0.7 (Medium suspicion)
**Signals:** `moderate_paste_ratio`, `low_edit_ratio`

## Postman Collection

Import this JSON into Postman for easy testing:

```json
{
  "info": {
    "name": "Dalivim API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Auth",
      "item": [
        {
          "name": "Register",
          "request": {
            "method": "POST",
            "url": "{{base_url}}/api/auth/register",
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"professor@dalivim.com\",\n  \"password\": \"senha123\",\n  \"name\": \"Prof. Silva\",\n  \"role\": \"professor\"\n}"
            }
          }
        },
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "url": "{{base_url}}/api/auth/login",
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"professor@dalivim.com\",\n  \"password\": \"senha123\"\n}"
            }
          }
        }
      ]
    }
  ]
}
```

## Environment Variables for Postman

```
base_url = http://localhost:8080
token = <your-auth-token>
activity_id = 1
student_id = 2
invite_token = a1b2c3d4e5f6789012345678
```
