import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import CodeEditor from './CodeEditor';

const StudentActivity = () => {
  const { inviteToken } = useParams();
  const navigate = useNavigate();
  const [activity, setActivity] = useState(null);
  const [student, setStudent] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [telemetryStatus, setTelemetryStatus] = useState(null);

  useEffect(() => {
    loadActivity();
  }, [inviteToken]);

  const loadActivity = async () => {
    try {
      const response = await fetch(`/api/activities/join/${inviteToken}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }
      });

      if (!response.ok) {
        throw new Error('Link inválido ou expirado');
      }

      const data = await response.json();
      setActivity(data.activity);
      setStudent(data.student);
      setLoading(false);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  const handleTelemetryUpdate = (result) => {
    setTelemetryStatus(result);
  };

  if (loading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        color: 'white',
        fontSize: '24px'
      }}>
        Carregando atividade...
      </div>
    );
  }

  if (error) {
    return (
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        color: 'white',
        textAlign: 'center',
        padding: '20px'
      }}>
        <h1 style={{ fontSize: '48px', marginBottom: '20px' }}>⚠️</h1>
        <h2 style={{ fontSize: '24px', marginBottom: '12px' }}>Erro ao carregar atividade</h2>
        <p style={{ fontSize: '16px', opacity: 0.9 }}>{error}</p>
        <button
          onClick={() => navigate('/')}
          style={{
            marginTop: '24px',
            padding: '12px 32px',
            background: 'white',
            color: '#667eea',
            border: 'none',
            borderRadius: '8px',
            fontSize: '16px',
            fontWeight: '600',
            cursor: 'pointer'
          }}
        >
          Voltar ao início
        </button>
      </div>
    );
  }

  return (
    <div style={{ position: 'relative' }}>
      <div style={{
        position: 'absolute',
        top: '20px',
        left: '20px',
        background: 'rgba(255, 255, 255, 0.95)',
        padding: '16px 24px',
        borderRadius: '12px',
        boxShadow: '0 4px 8px rgba(0,0,0,0.15)',
        zIndex: 10,
        maxWidth: '300px'
      }}>
        <h3 style={{ margin: '0 0 8px 0', fontSize: '16px', color: '#667eea' }}>
          {activity.title}
        </h3>
        <p style={{ margin: '0 0 8px 0', fontSize: '13px', color: '#666' }}>
          {activity.description}
        </p>
        <p style={{ margin: 0, fontSize: '12px', color: '#999' }}>
          Aluno: {student.name || student.email}
        </p>
      </div>

      {telemetryStatus && (
        <div style={{
          position: 'absolute',
          bottom: '20px',
          right: '20px',
          background: telemetryStatus.confidence === 'high' && telemetryStatus.authorship_score < 0.5
            ? 'rgba(239, 68, 68, 0.95)'
            : 'rgba(34, 197, 94, 0.95)',
          color: 'white',
          padding: '12px 20px',
          borderRadius: '8px',
          boxShadow: '0 4px 8px rgba(0,0,0,0.15)',
          zIndex: 10,
          fontSize: '13px'
        }}>
          <strong>Score de Autoria:</strong> {(telemetryStatus.authorship_score * 100).toFixed(1)}%
          {telemetryStatus.signals && telemetryStatus.signals.length > 0 && (
            <div style={{ marginTop: '6px', fontSize: '11px', opacity: 0.9 }}>
              🚩 {telemetryStatus.signals.join(', ')}
            </div>
          )}
        </div>
      )}

      <CodeEditor
        activityId={activity.id}
        studentId={student.id}
        onTelemetryUpdate={handleTelemetryUpdate}
      />
    </div>
  );
};

export default StudentActivity;
