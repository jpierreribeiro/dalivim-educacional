# 🎓 Sistema de Semestres e Detecção de Similaridade

## 📋 Visão Geral

Este sistema adiciona:
1. **Gestão de Semestres Acadêmicos**
2. **Atualização Automática do Semestre dos Alunos**
3. **Detecção Automática de Similaridade entre Códigos**
4. **Clusterização de Submissões Similares**

## 🏗️ Novos Modelos

### 1. Semester (Semestre Acadêmico)
```go
type Semester struct {
    ID        uint
    Year      int       // 2024, 2025, etc
    Period    int       // 1 (primeiro semestre) ou 2 (segundo semestre)
    StartDate time.Time // Data de início
    EndDate   time.Time // Data de término
}
```

**Exemplo**: 2024.1 (primeiro semestre de 2024)

### 2. User (Atualizado)
```go
type User struct {
    // ... campos anteriores
    CurrentSemester  int // Semestre atual do aluno (1-10)
    EnrollmentYear   int // Ano de matrícula
    EnrollmentPeriod int // Período de matrícula (1 ou 2)
}
```

### 3. Activity (Atualizada)
```go
type Activity struct {
    // ... campos anteriores
    SemesterID     uint // Semestre acadêmico (2024.1, 2024.2, etc)
    TargetSemester int  // Semestre do curso (1º, 2º, 3º, etc)
}
```

### 4. SimilarityDetection (Novo)
```go
type SimilarityDetection struct {
    ActivityID      uint
    SubmissionID1   uint
    SubmissionID2   uint
    StudentID1      uint
    StudentID2      uint
    SimilarityScore float64 // 0.0 a 1.0
    Algorithm       string  // "levenshtein_normalized"
    IsSuspicious    bool    // true se score > 0.75
    ClusterID       *uint   // Grupo de códigos similares
}
```

### 5. SimilarityCluster (Novo)
```go
type SimilarityCluster struct {
    ActivityID     uint
    ClusterSize    int     // Quantas submissões no cluster
    AvgSimilarity  float64 // Similaridade média
    SuspicionLevel string  // "low", "medium", "high"
}
```

## 🔄 Fluxo de Funcionamento

### 1. Criação de Semestre
```
Professor/Admin cria semestre → Sistema armazena
Exemplo: 2024.1 (01/02/2024 a 30/06/2024)
```

### 2. Cadastro de Aluno
```
Aluno se cadastra → Informa:
- Ano de matrícula: 2024
- Período de matrícula: 1 (2024.1)
- Sistema calcula: CurrentSemester = 1
```

### 3. Atualização Automática
```
Quando semestre muda (2024.1 → 2024.2):
Sistema recalcula automaticamente:
- Aluno de 2024.1 → CurrentSemester = 2
- Aluno de 2023.2 → CurrentSemester = 4
```

### 4. Criação de Atividade
```
Professor cria atividade:
- Semestre acadêmico: 2024.1
- Semestre alvo: 3 (atividade para alunos do 3º semestre)
```

### 5. Detecção Automática de Similaridade
```
Quando alunos submetem → Sistema automaticamente:
1. Compara todos os códigos par a par
2. Calcula score de similaridade (0-100%)
3. Identifica submissões suspeitas (>75%)
4. Agrupa códigos similares em clusters
```

## 📊 Algoritmo de Similaridade

### Levenshtein Distance Normalizado
```
1. Normaliza código (remove espaços, lowercase)
2. Calcula distância de Levenshtein
3. Normaliza para score 0-1:
   
   similarity = 1 - (distance / max_length)
```

### Thresholds
- **> 90%**: 🔴 Alta suspeição (provavelmente copiado)
- **75-90%**: 🟡 Média suspeição (muito similar)
- **< 75%**: 🟢 Baixa suspeição (normal)

### Clusterização
```
Algoritmo de componentes conectados:
1. Cria grafo onde nodes = submissões
2. Edges = similaridade > 75%
3. Encontra componentes conectados = clusters
4. Cluster com 2+ submissões = grupo de cola
```

## 🎯 Casos de Uso

### Caso 1: Detectar Grupos de Cola

**Cenário**: 5 alunos do 3º semestre fazem atividade de Bubble Sort

```
Submissões:
- Aluno A: código único
- Aluno B: código 92% similar ao C
- Aluno C: código 92% similar ao B
- Aluno D: código 88% similar ao B e C
- Aluno E: código único

Sistema detecta:
Cluster 1: [B, C, D]
- Avg Similarity: 90.6%
- Suspicion Level: high
- Conclusão: Provável cola entre B, C e D
```

### Caso 2: Progressão de Semestre

**Cenário**: Aluno se matriculou em 2024.1

```
2024.1 (início): CurrentSemester = 1
2024.2 (automático): CurrentSemester = 2
2025.1 (automático): CurrentSemester = 3
2025.2 (automático): CurrentSemester = 4
```

### Caso 3: Atividade por Semestre

**Cenário**: Professor cria atividade de Estrutura de Dados

```
Configuração:
- Semestre acadêmico: 2024.2
- Semestre alvo: 3 (alunos do 3º semestre)

Sistema permite:
- Apenas alunos do 3º semestre acessam
- Comparações apenas dentro deste grupo
- Histórico por semestre acadêmico
```

## 🔌 API Endpoints

### Semestres

```bash
# Criar semestre
POST /api/semesters
{
  "year": 2024,
  "period": 1,
  "startDate": "2024-02-01T00:00:00Z",
  "endDate": "2024-06-30T23:59:59Z"
}

# Listar semestres
GET /api/semesters

# Semestre ativo
GET /api/semesters/active

# Atualizar semestres dos alunos (cron job)
POST /api/semesters/update-students
```

### Atividades (Atualizado)

```bash
# Criar atividade com semestre
POST /api/activities
{
  "title": "Bubble Sort",
  "description": "...",
  "language": "python",
  "timeLimit": 60,
  "semesterId": 1,
  "targetSemester": 3
}
```

### Similaridade

```bash
# Detectar similaridades (automático ou manual)
POST /api/activities/:id/detect-similarities

# Ver similaridades de uma atividade
GET /api/activities/:id/similarities

# Ver clusters de cola
GET /api/activities/:id/clusters
```

### Cadastro de Aluno (Atualizado)

```bash
POST /api/auth/register
{
  "email": "aluno@email.com",
  "password": "senha123",
  "name": "João Silva",
  "role": "student",
  "enrollmentYear": 2024,
  "enrollmentPeriod": 1
}
```

## 📈 Dashboard do Professor

### Nova View: Similaridades

```
┌─────────────────────────────────────────────┐
│ Atividade: Bubble Sort (2024.1 - 3º Sem)  │
├─────────────────────────────────────────────┤
│                                             │
│ 📊 Estatísticas de Similaridade:           │
│   • Total de submissões: 45                 │
│   • Clusters detectados: 3                  │
│   • Submissões suspeitas: 12 (26.7%)       │
│                                             │
│ 🔴 Cluster 1 (Alta Suspeição)              │
│   • Tamanho: 5 alunos                       │
│   • Similaridade média: 94.2%               │
│   • Alunos: João, Maria, Pedro, Ana, Lucas │
│   └─ [Ver Detalhes] [Comparar Códigos]     │
│                                             │
│ 🟡 Cluster 2 (Média Suspeição)             │
│   • Tamanho: 3 alunos                       │
│   • Similaridade média: 82.5%               │
│   • Alunos: Carlos, Fernanda, Bruno        │
│   └─ [Ver Detalhes]                         │
│                                             │
│ 🟢 Cluster 3 (Baixa Suspeição)             │
│   • Tamanho: 4 alunos                       │
│   • Similaridade média: 68.1%               │
│   • Alunos: Ricardo, Juliana, Paula, Diego │
│   └─ [Ver Detalhes]                         │
│                                             │
└─────────────────────────────────────────────┘
```

## ⚙️ Configuração

### 1. Criar Semestre Inicial

```bash
curl -X POST http://localhost:8080/api/semesters \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "year": 2024,
    "period": 1,
    "startDate": "2024-02-01T00:00:00Z",
    "endDate": "2024-06-30T23:59:59Z"
  }'
```

### 2. Configurar Cron Job

```bash
# Adicionar no crontab
# Atualiza semestres dos alunos toda segunda-feira às 2h
0 2 * * 1 curl -X POST http://localhost:8080/api/semesters/update-students \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### 3. Configurar Detecção Automática

```go
// Opção 1: Detecção ao finalizar submissão
// No telemetry_service.go
if isFinal {
    // ... criar submission ...
    
    // Trigger similarity detection
    go s.similarityService.DetectSimilarities(activityID)
}

// Opção 2: Detecção manual pelo professor
// Via endpoint POST /api/activities/:id/detect-similarities
```

## 🔐 Considerações de Privacidade

1. **LGPD**: Dados de similaridade devem ser pseudonimizados
2. **Acesso**: Apenas professor da disciplina vê comparações
3. **Retenção**: Deletar após período acadêmico
4. **Transparência**: Informar alunos sobre análise

## 🚀 Melhorias Futuras

### 1. Algoritmos Avançados
- **AST Comparison**: Comparar árvore sintática abstrata
- **Cosine Similarity**: TF-IDF em tokens de código
- **MOSS**: Measure of Software Similarity

### 2. Machine Learning
```python
# Treinar modelo para detectar padrões de cola
features = [
    similarity_score,
    temporal_proximity,  # Submissões próximas no tempo
    behavioral_patterns, # Telemetria similar
    network_analysis     # Relações sociais
]
```

### 3. Visualizações
- **Grafo de Similaridade**: Nodes = alunos, edges = similaridade
- **Heatmap**: Matriz de similaridade entre todos os pares
- **Timeline**: Quando cada aluno submeteu

### 4. Alertas Automáticos
```
Sistema detecta cluster suspeito →
Email automático para professor →
"🚨 Possível cola detectada na atividade X"
```

## 📚 Referências

- [Levenshtein Distance](https://en.wikipedia.org/wiki/Levenshtein_distance)
- [MOSS - Plagiarism Detection](https://theory.stanford.edu/~aiken/moss/)
- [Graph Clustering Algorithms](https://en.wikipedia.org/wiki/Graph_clustering)