# Arquitetura do Sistema Dalivim

## Diagrama de Fluxo Completo

```mermaid
graph TB
    subgraph "Frontend - React"
        A[Professor Dashboard] --> B[Criar Atividade]
        B --> C[Link de Convite]
        D[Aluno] --> E[Student Activity Page]
        E --> F[Monaco Editor]
        F --> G[Telemetry Capture]
        G --> H[Keystroke Dynamics]
        G --> I[Paste Events]
        G --> J[Focus Events]
        G --> K[Edit Patterns]
    end
    
    subgraph "Backend - Go + Gin"
        L[API Gateway] --> M[Auth Service]
        L --> N[Activity Service]
        L --> O[Telemetry Service]
        O --> P[Behavior Analyzer]
        P --> Q[Feature Extraction]
        P --> R[Signal Detection]
        P --> S[Score Calculation]
    end
    
    subgraph "Database - PostgreSQL"
        T[(Users)]
        U[(Activities)]
        V[(Submissions)]
        W[(Telemetry Data)]
    end
    
    subgraph "External Services"
        X[Piston API]
        Y[Code Execution]
    end
    
    C --> D
    F --> X
    X --> Y
    H --> O
    I --> O
    J --> O
    K --> O
    M --> T
    N --> U
    O --> V
    O --> W
    S --> E
    A --> L
    E --> L
```

## Fluxo de Dados - Telemetria

```mermaid
sequenceDiagram
    participant S as Student
    participant E as Editor
    participant T as Telemetry
    participant B as Backend
    participant A as Analyzer
    participant D as Database
    
    S->>E: Digita código
    E->>T: Captura keystroke
    T->>T: Calcula flight time
    T->>T: Calcula dwell time
    
    S->>E: Cola código
    E->>T: Evento de paste
    T->>T: Registra tamanho e conteúdo
    
    S->>E: Sai da aba
    E->>T: Evento blur
    Note over T: Inicia timer
    
    S->>E: Volta à aba
    E->>T: Evento focus
    T->>T: Calcula tempo ausente
    
    loop A cada 10 segundos
        T->>T: Agrega features
        T->>B: POST /api/telemetry
        B->>A: Analisa comportamento
        A->>A: Detecta sinais
        A->>A: Calcula score
        A->>B: Retorna análise
        B->>D: Salva telemetria
        B->>E: Retorna resultado
        E->>S: Exibe feedback (opcional)
    end
    
    S->>E: Envia submissão final
    T->>B: POST /api/telemetry (final=true)
    B->>A: Análise completa
    A->>B: Score final
    B->>D: Salva submission
    B->>E: Confirmação
```

## Fluxo de Criação de Atividade

```mermaid
sequenceDiagram
    participant P as Professor
    participant UI as Dashboard
    participant API as Backend API
    participant DB as Database
    
    P->>UI: Cria nova atividade
    UI->>UI: Preenche formulário
    P->>UI: Submete
    UI->>API: POST /api/activities
    API->>API: Gera invite token
    API->>DB: INSERT activity
    DB->>API: Confirmação
    API->>UI: Retorna atividade + token
    UI->>P: Exibe link de convite
    P->>P: Compartilha link com alunos
```

## Fluxo do Aluno

```mermaid
sequenceDiagram
    participant S as Student
    participant B as Browser
    participant API as Backend API
    participant DB as Database
    participant P as Piston API
    
    S->>B: Acessa link de convite
    B->>API: POST /api/activities/join/:token
    API->>DB: SELECT activity
    API->>DB: INSERT anonymous student
    API->>B: Retorna activity + student
    B->>S: Exibe editor
    
    loop Durante a atividade
        S->>B: Escreve código
        B->>B: Captura telemetria
        B->>API: POST /api/telemetry
        API->>DB: INSERT telemetry_data
        API->>B: Retorna análise
    end
    
    S->>B: Executa código
    B->>P: POST /execute
    P->>B: Retorna output
    B->>S: Exibe resultado
    
    S->>B: Submete final
    B->>API: POST /api/telemetry (final)
    API->>DB: INSERT submission
    API->>B: Confirmação
    B->>S: Sucesso
```

## Arquitetura de Análise

```mermaid
graph LR
    subgraph "Input Features"
        A[Keystroke Dynamics]
        B[Paste Events]
        C[Focus Events]
        D[Edit Patterns]
        E[Execution History]
    end
    
    subgraph "Feature Extraction"
        F[Avg Interval]
        G[Std Interval]
        H[Paste Ratio]
        I[Delete Ratio]
        J[Linear Score]
        K[Burstiness]
        L[Focus Loss Count]
        M[Time to First Run]
    end
    
    subgraph "Signal Detection"
        N{High Paste?}
        O{Low Delete?}
        P{High Linear?}
        Q{No Execution?}
        R{Low Variance?}
    end
    
    subgraph "Scoring"
        S[Suspicion Score]
        T[Authorship Score]
        U[Confidence Level]
    end
    
    A --> F
    A --> G
    A --> K
    B --> H
    D --> I
    D --> J
    C --> L
    E --> M
    
    F --> R
    H --> N
    I --> O
    J --> P
    M --> Q
    
    N --> S
    O --> S
    P --> S
    Q --> S
    R --> S
    
    S --> T
    S --> U
    
    T --> V[0.0 - 1.0]
    U --> W[Low/Medium/High]
```

## Modelo de Dados

```mermaid
erDiagram
    USERS ||--o{ ACTIVITIES : creates
    USERS ||--o{ SUBMISSIONS : submits
    ACTIVITIES ||--o{ SUBMISSIONS : receives
    ACTIVITIES ||--o{ TELEMETRY_DATA : tracks
    USERS ||--o{ TELEMETRY_DATA : generates
    
    USERS {
        int id PK
        string email UK
        string password
        string name
        string role
        timestamp created_at
    }
    
    ACTIVITIES {
        int id PK
        int professor_id FK
        string title
        text description
        string language
        int time_limit
        string invite_token UK
        timestamp created_at
    }
    
    SUBMISSIONS {
        int id PK
        int activity_id FK
        int student_id FK
        text code
        decimal authorship_score
        string confidence
        array signals
        decimal avg_keystroke_interval
        decimal paste_char_ratio
        int execution_count
        timestamp created_at
    }
    
    TELEMETRY_DATA {
        int id PK
        int activity_id FK
        int student_id FK
        bigint timestamp
        bool is_final
        jsonb features
        jsonb raw_events
        timestamp created_at
    }
```

## Stack Tecnológica

```mermaid
graph TB
    subgraph "Frontend Stack"
        A[React 18]
        B[React Router DOM]
        C[Monaco Editor]
        D[Axios/Fetch]
    end
    
    subgraph "Backend Stack"
        E[Go 1.21]
        F[Gin Framework]
        G[GORM]
        H[PostgreSQL Driver]
    end
    
    subgraph "Database"
        I[PostgreSQL 15]
    end
    
    subgraph "External APIs"
        J[Piston API]
    end
    
    subgraph "Future Enhancements"
        K[Redis Cache]
        L[WebSocket]
        M[TensorFlow.js]
        N[S3 Storage]
    end
    
    A --> E
    B --> E
    C --> J
    D --> E
    E --> F
    F --> G
    G --> H
    H --> I
```

## Deployment Architecture

```mermaid
graph TB
    subgraph "Production Environment"
        A[Load Balancer]
        
        subgraph "Frontend Servers"
            B[React App 1]
            C[React App 2]
        end
        
        subgraph "Backend Servers"
            D[Go API 1]
            E[Go API 2]
        end
        
        subgraph "Database Cluster"
            F[(PostgreSQL Primary)]
            G[(PostgreSQL Replica)]
        end
        
        H[Redis Cache]
        I[S3 Storage]
    end
    
    J[Users] --> A
    A --> B
    A --> C
    B --> D
    C --> E
    D --> F
    E --> F
    D --> H
    E --> H
    F --> G
    D --> I
    E --> I
```

## Performance Metrics

### Telemetria
- **Frequência**: A cada 10 segundos
- **Payload médio**: ~2-5 KB
- **Latência esperada**: < 100ms

### Análise
- **Processamento**: Síncrono (< 50ms)
- **Score calculation**: O(n) onde n = features
- **Database writes**: Batch quando possível

### Code Execution
- **Piston API**: Timeout 3s
- **Retry logic**: 2 tentativas
- **Fallback**: Mensagem de erro amigável

## Security Considerations

```mermaid
graph LR
    subgraph "Security Layers"
        A[HTTPS/TLS]
        B[JWT Authentication]
        C[CORS Policy]
        D[Rate Limiting]
        E[SQL Injection Protection]
        F[XSS Prevention]
    end
    
    G[Client Request] --> A
    A --> B
    B --> C
    C --> D
    D --> E
    E --> F
    F --> H[Application Logic]
```

## Análise de Features

### Features Primárias (Alto Peso)
1. **Paste Char Ratio** - Proporção de código colado
2. **Delete Ratio** - Taxa de deleções/correções
3. **Linear Editing Score** - Linearidade do código

### Features Secundárias (Médio Peso)
4. **Keystroke Variance** - Burstiness da digitação
5. **Execution Count** - Número de tentativas
6. **Focus Loss Count** - Saídas da aba

### Features Contextuais (Baixo Peso)
7. **Average Keystroke Interval** - Velocidade média
8. **Total Time** - Tempo total gasto
9. **Time to First Run** - Tempo até testar

## Score Calculation Formula

```
Suspicion Score = Σ(weight_i × signal_i)

where:
- weight_i = peso do sinal detectado
- signal_i = 1 se sinal presente, 0 caso contrário

Authorship Score = 1.0 - min(Suspicion Score, 1.0)

Confidence:
- High: ≥ 4 sinais detectados
- Medium: 2-3 sinais detectados
- Low: 0-1 sinais detectados
```
