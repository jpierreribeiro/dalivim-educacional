# Dalivim - Plataforma de Estudos com Detecção de IA

Sistema completo para criação de atividades de programação com análise comportamental para detectar uso de IA.

## 🎯 Funcionalidades

### Para Professores
- ✅ Criar atividades de programação
- ✅ Gerar links de convite únicos
- ✅ Visualizar submissões dos alunos
- ✅ Análise comportamental detalhada
- ✅ Detecção de eventos suspeitos (copy/paste, foco, etc.)
- ✅ Score de autoria (0-100%)

### Para Alunos
- ✅ Entrar via link de convite
- ✅ Editor de código Monaco (múltiplas linguagens)
- ✅ Execução de código via Piston API
- ✅ Telemetria invisível e transparente

## 🏗️ Arquitetura

### Frontend (React)
```
frontend/
├── CodeEditor.jsx          # Editor Monaco com telemetria
├── StudentActivity.jsx     # Página do aluno
├── ProfessorDashboard.jsx  # Dashboard do professor
└── ActivityDetails.jsx     # Detalhes e análise de submissões
```

### Backend (Go + Gin + GORM)
```
backend/
├── main.go                 # API completa
└── go.mod                  # Dependências
```

## 📊 Telemetria Capturada

### Keystroke Dynamics
- **Dwell Time**: Tempo entre pressionar e soltar tecla
- **Flight Time**: Tempo entre duas teclas consecutivas
- **Variância**: Desvio padrão dos intervalos (burstiness)

### Eventos de Paste
- Quantidade de eventos
- Tamanho do conteúdo colado
- Proporção de código colado vs. digitado

### Análise de Foco
- Eventos de blur/focus
- Duração fora da aba
- Mudanças suspeitas (saiu 20s, voltou com 30 linhas)

### Padrões de Edição
- Edição linear vs. não-linear
- Taxa de deleção
- Correções e backtracking

### Execução
- Número de tentativas de execução
- Tempo até primeira execução
- Código final vs. intermediários

## 🚩 Sinais de Suspeição

O sistema detecta os seguintes padrões:

| Sinal | Descrição | Peso |
|-------|-----------|------|
| `high_paste_ratio` | >60% do código foi colado | 🔴 Alto |
| `low_edit_ratio` | <2% de deleções (código perfeito) | 🔴 Alto |
| `highly_linear_editing` | >90% edição linear (sem correções) | 🟡 Médio |
| `multiple_paste_events` | Mais de 3 eventos de cola | 🟡 Médio |
| `fast_completion_no_testing` | Código rápido sem execuções | 🟡 Médio |
| `frequent_focus_loss` | >5 saídas da aba | 🟢 Baixo |
| `low_typing_variance` | Digitação robótica | 🟢 Baixo |

### Score de Autoria

```
Score = 1.0 - Σ(pesos_dos_sinais)

- 80-100%: ✅ Muito Baixa Suspeição
- 60-80%:  ✓ Baixa Suspeição
- 40-60%:  ⚠️ Média Suspeição
- 20-40%:  🚨 Alta Suspeição
- 0-20%:   ⛔ Muito Alta Suspeição
```

## 🚀 Setup

### Backend (Go)

1. **Instalar PostgreSQL**
```bash
# Ubuntu/Debian
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql

# Iniciar serviço
sudo systemctl start postgresql
```

2. **Criar Database**
```bash
sudo -u postgres psql
CREATE DATABASE dalivim;
CREATE USER postgres WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE dalivim TO postgres;
\q
```

3. **Instalar Dependências e Rodar**
```bash
cd backend
go mod download
go run main.go
```

O servidor estará em `http://localhost:8080`

### Frontend (React)

1. **Instalar Dependências**
```bash
cd frontend
npm install react react-dom react-router-dom
npm install @monaco-editor/react
```

2. **package.json** necessário:
```json
{
  "name": "dalivim-frontend",
  "version": "1.0.0",
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.20.0",
    "@monaco-editor/react": "^4.6.0"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build"
  }
}
```

3. **Rodar**
```bash
npm start
```

O frontend estará em `http://localhost:3000`

### Configuração de Rotas

**App.jsx**:
```jsx
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import ProfessorDashboard from './ProfessorDashboard';
import StudentActivity from './StudentActivity';
import ActivityDetails from './ActivityDetails';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<ProfessorDashboard />} />
        <Route path="/activity/:inviteToken" element={<StudentActivity />} />
        <Route path="/professor/activity/:activityId" element={<ActivityDetails />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
```

## 🔄 Fluxo Completo

### 1. Professor Cria Atividade
```
POST /api/activities
{
  "title": "Implementar Bubble Sort",
  "description": "Crie um algoritmo de ordenação...",
  "language": "python",
  "timeLimit": 60
}

Response:
{
  "id": 1,
  "inviteToken": "a1b2c3d4e5f6...",
  ...
}
```

### 2. Professor Compartilha Link
```
https://dalivim.com/activity/a1b2c3d4e5f6
```

### 3. Aluno Entra e Faz Atividade
```
POST /api/activities/join/a1b2c3d4e5f6

Response:
{
  "activity": {...},
  "student": {
    "id": 123,
    "email": "student_xyz@anonymous.local"
  }
}
```

### 4. Telemetria em Tempo Real
```
POST /api/telemetry (a cada 10 segundos)
{
  "activityId": 1,
  "studentId": 123,
  "features": {
    "avgKeystrokeInterval": 150.5,
    "pasteCharRatio": 0.35,
    "deleteRatio": 0.08,
    ...
  },
  "rawEvents": {
    "pasteEvents": [...],
    "focusEvents": [...]
  }
}

Response:
{
  "authorship_score": 0.73,
  "confidence": "medium",
  "signals": ["moderate_paste_ratio", "low_edit_ratio"]
}
```

### 5. Submissão Final
```
POST /api/telemetry (isFinal: true)
{
  "isFinal": true,
  "code": "def bubble_sort(arr):\n  ...",
  ...
}
```

### 6. Professor Visualiza Análise
```
GET /api/activities/1/submissions

Response: [
  {
    "studentName": "Anonymous Student",
    "authorshipScore": 0.73,
    "signals": ["moderate_paste_ratio"],
    "pasteEvents": 2,
    "code": "...",
    ...
  }
]
```

## 📡 API Reference

### Endpoints Públicos

#### POST /api/auth/register
Registra novo usuário (professor)

#### POST /api/auth/login
Login de usuário

#### POST /api/activities/join/:inviteToken
Aluno entra na atividade via link

### Endpoints Autenticados

#### POST /api/activities
Cria nova atividade

#### GET /api/activities
Lista atividades do professor

#### GET /api/activities/:id
Detalhes de uma atividade

#### GET /api/activities/:id/submissions
Submissões de uma atividade

### Telemetria (Público)

#### POST /api/telemetry
Recebe dados de telemetria do aluno

## 🎨 UI/UX

### Editor do Aluno
- 🎨 Gradiente roxo moderno
- 📝 Monaco Editor (VS Code)
- ▶️ Execução via Piston API
- 📊 Feedback visual em tempo real
- 🔒 Telemetria transparente

### Dashboard do Professor
- 📚 Cards de atividades
- 📋 Copiar link de convite
- 📊 Contador de submissões
- 🎯 Visualização por atividade

### Análise Detalhada
- 📈 Métricas comportamentais
- 🚩 Sinais de suspeição destacados
- 📋 Eventos de paste detalhados
- 💻 Código fonte completo
- 🎨 Código de cores por nível de suspeição

## 🔐 Segurança

### Considerações Atuais
⚠️ **Este é um MVP educacional. Para produção, implemente:**

1. **Autenticação Forte**
   - JWT com refresh tokens
   - Hash de senhas com bcrypt
   - Rate limiting

2. **Validação**
   - Sanitização de inputs
   - CORS restritivo
   - SQL injection protection (GORM já ajuda)

3. **Privacy**
   - Consentimento explícito de telemetria
   - LGPD compliance
   - Anonimização de dados

## 🔄 Migração para Piston Local

Quando quiser rodar Piston localmente:

```bash
# Clone Piston
git clone https://github.com/engineer-man/piston
cd piston

# Instale linguagens
docker run -v $PWD:'/piston' --tmpfs /piston/jobs -dit -p 2000:2000 --name piston_api ghcr.io/engineer-man/piston

# Atualize URL no frontend
// CodeEditor.jsx
const PISTON_URL = 'http://localhost:2000';
```

## 📊 Análise de Resultados

### Interpretação dos Dados

**Score Alto (>70%)**
- Comportamento natural
- Edições orgânicas
- Múltiplas execuções
- Padrão humano consistente

**Score Médio (40-70%)**
- Alguns sinais mistos
- Pode ser legítimo ou não
- Requer análise contextual

**Score Baixo (<40%)**
- Múltiplos sinais de alerta
- Alta probabilidade de cola
- Investigação recomendada

## 🛠️ Troubleshooting

### Backend não conecta ao Postgres
```bash
# Verificar se PostgreSQL está rodando
sudo systemctl status postgresql

# Verificar conexão
psql -U postgres -d dalivim -h localhost
```

### Frontend não conecta ao Backend
- Verificar CORS em `main.go`
- Confirmar que backend está em `localhost:8080`
- Verificar console do navegador

### Monaco Editor não carrega
```bash
npm install @monaco-editor/react
```

## 🚀 Próximos Passos

### Melhorias Sugeridas

1. **Machine Learning**
   - Treinar modelo de classificação
   - Usar TensorFlow.js no frontend
   - Melhorar detecção com mais features

2. **Features Adicionais**
   - Gravação de sessão (playback)
   - Comparação entre alunos
   - Exportar relatórios PDF
   - Dashboard de analytics

3. **Integração**
   - Google Classroom
   - GitHub Education
   - LMS (Moodle, Canvas)

4. **Escalabilidade**
   - WebSocket para telemetria em tempo real
   - Redis para cache
   - Microservices architecture

## 📝 Licença

MIT License - Use livremente para fins educacionais

## 🤝 Contribuindo

Pull requests são bem-vindos! Para mudanças maiores, abra uma issue primeiro.

---

**Desenvolvido para detectar IA, não para punir aprendizado** 🎓
