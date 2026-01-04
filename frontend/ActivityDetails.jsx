import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

const ActivityDetails = () => {
  const { activityId } = useParams();
  const [activity, setActivity] = useState(null);
  const [submissions, setSubmissions] = useState([]);
  const [selectedSubmission, setSelectedSubmission] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadActivityDetails();
  }, [activityId]);

  const loadActivityDetails = async () => {
    try {
      const [activityRes, submissionsRes] = await Promise.all([
        fetch(`/api/activities/${activityId}`, {
          headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
        }),
        fetch(`/api/activities/${activityId}/submissions`, {
          headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
        })
      ]);

      const activityData = await activityRes.json();
      const submissionsData = await submissionsRes.json();

      setActivity(activityData);
      setSubmissions(submissionsData);
      setLoading(false);
    } catch (error) {
      console.error('Failed to load activity details:', error);
      setLoading(false);
    }
  };

  const getSuspicionLevel = (score) => {
    if (score > 0.8) return { level: 'Muito Baixa', color: '#10b981', emoji: '✅' };
    if (score > 0.6) return { level: 'Baixa', color: '#3b82f6', emoji: '✓' };
    if (score > 0.4) return { level: 'Média', color: '#f59e0b', emoji: '⚠️' };
    if (score > 0.2) return { level: 'Alta', color: '#ef4444', emoji: '🚨' };
    return { level: 'Muito Alta', color: '#dc2626', emoji: '⛔' };
  };

  const formatDuration = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}m ${secs}s`;
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
        Carregando...
      </div>
    );
  }

  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      padding: '40px'
    }}>
      <div style={{
        maxWidth: '1400px',
        margin: '0 auto'
      }}>
        <div style={{
          background: 'white',
          borderRadius: '16px',
          padding: '32px',
          marginBottom: '24px',
          boxShadow: '0 8px 16px rgba(0,0,0,0.15)'
        }}>
          <h1 style={{
            margin: '0 0 12px 0',
            fontSize: '32px',
            color: '#667eea',
            fontWeight: '700'
          }}>
            {activity?.title}
          </h1>
          <p style={{
            margin: '0 0 16px 0',
            fontSize: '16px',
            color: '#666',
            lineHeight: '1.6'
          }}>
            {activity?.description}
          </p>
          <div style={{
            display: 'flex',
            gap: '12px'
          }}>
            <span style={{
              padding: '6px 16px',
              background: '#f3f4f6',
              borderRadius: '8px',
              fontSize: '14px',
              color: '#667eea',
              fontWeight: '600'
            }}>
              📝 {submissions.length} submissões
            </span>
            <span style={{
              padding: '6px 16px',
              background: '#f3f4f6',
              borderRadius: '8px',
              fontSize: '14px',
              color: '#666',
              fontWeight: '600'
            }}>
              💻 {activity?.language}
            </span>
            <span style={{
              padding: '6px 16px',
              background: '#f3f4f6',
              borderRadius: '8px',
              fontSize: '14px',
              color: '#666',
              fontWeight: '600'
            }}>
              ⏱ {activity?.timeLimit} min
            </span>
          </div>
        </div>

        <div style={{
          display: 'grid',
          gridTemplateColumns: selectedSubmission ? '400px 1fr' : '1fr',
          gap: '24px'
        }}>
          <div style={{
            display: 'flex',
            flexDirection: 'column',
            gap: '16px'
          }}>
            {submissions.map(submission => {
              const suspicion = getSuspicionLevel(submission.authorshipScore);
              return (
                <div
                  key={submission.id}
                  onClick={() => setSelectedSubmission(submission)}
                  style={{
                    background: 'white',
                    borderRadius: '12px',
                    padding: '20px',
                    cursor: 'pointer',
                    boxShadow: selectedSubmission?.id === submission.id
                      ? '0 0 0 3px #667eea'
                      : '0 4px 8px rgba(0,0,0,0.1)',
                    transition: 'all 0.3s ease',
                    border: selectedSubmission?.id === submission.id
                      ? '2px solid #667eea'
                      : '2px solid transparent'
                  }}
                >
                  <div style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    marginBottom: '12px'
                  }}>
                    <h3 style={{
                      margin: 0,
                      fontSize: '16px',
                      color: '#333',
                      fontWeight: '700'
                    }}>
                      {submission.studentName || submission.studentEmail}
                    </h3>
                    <span style={{
                      fontSize: '20px'
                    }}>
                      {suspicion.emoji}
                    </span>
                  </div>

                  <div style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '12px',
                    marginBottom: '12px'
                  }}>
                    <div style={{
                      flex: 1,
                      height: '8px',
                      background: '#f3f4f6',
                      borderRadius: '4px',
                      overflow: 'hidden'
                    }}>
                      <div style={{
                        width: `${submission.authorshipScore * 100}%`,
                        height: '100%',
                        background: suspicion.color,
                        borderRadius: '4px',
                        transition: 'width 0.3s ease'
                      }} />
                    </div>
                    <span style={{
                      fontSize: '14px',
                      fontWeight: '700',
                      color: suspicion.color
                    }}>
                      {(submission.authorshipScore * 100).toFixed(0)}%
                    </span>
                  </div>

                  <div style={{
                    display: 'flex',
                    gap: '8px',
                    flexWrap: 'wrap'
                  }}>
                    <span style={{
                      padding: '4px 10px',
                      background: '#f3f4f6',
                      borderRadius: '6px',
                      fontSize: '11px',
                      color: '#666',
                      fontWeight: '600'
                    }}>
                      ⏱ {formatDuration(submission.totalTime)}
                    </span>
                    <span style={{
                      padding: '4px 10px',
                      background: '#f3f4f6',
                      borderRadius: '6px',
                      fontSize: '11px',
                      color: '#666',
                      fontWeight: '600'
                    }}>
                      ⌨️ {submission.keystrokeCount}
                    </span>
                    {submission.pasteEvents > 0 && (
                      <span style={{
                        padding: '4px 10px',
                        background: '#fee2e2',
                        borderRadius: '6px',
                        fontSize: '11px',
                        color: '#dc2626',
                        fontWeight: '600'
                      }}>
                        📋 {submission.pasteEvents} colas
                      </span>
                    )}
                  </div>

                  {submission.signals && submission.signals.length > 0 && (
                    <div style={{
                      marginTop: '12px',
                      padding: '8px',
                      background: '#fef3c7',
                      borderRadius: '6px',
                      fontSize: '11px',
                      color: '#92400e',
                      fontWeight: '600'
                    }}>
                      🚩 {submission.signals.join(', ')}
                    </div>
                  )}
                </div>
              );
            })}
          </div>

          {selectedSubmission && (
            <div style={{
              background: 'white',
              borderRadius: '12px',
              padding: '24px',
              boxShadow: '0 8px 16px rgba(0,0,0,0.15)',
              maxHeight: 'calc(100vh - 200px)',
              overflow: 'auto'
            }}>
              <div style={{
                marginBottom: '24px',
                paddingBottom: '20px',
                borderBottom: '2px solid #f3f4f6'
              }}>
                <h2 style={{
                  margin: '0 0 16px 0',
                  fontSize: '24px',
                  color: '#667eea',
                  fontWeight: '700'
                }}>
                  Análise Detalhada
                </h2>

                <div style={{
                  display: 'grid',
                  gridTemplateColumns: 'repeat(2, 1fr)',
                  gap: '16px',
                  marginBottom: '16px'
                }}>
                  <MetricCard
                    label="Intervalo Médio (ms)"
                    value={selectedSubmission.avgKeystrokeInterval?.toFixed(0) || 'N/A'}
                    icon="⌨️"
                  />
                  <MetricCard
                    label="Desvio Padrão (ms)"
                    value={selectedSubmission.stdKeystrokeInterval?.toFixed(0) || 'N/A'}
                    icon="📊"
                  />
                  <MetricCard
                    label="Taxa de Deleção"
                    value={`${((selectedSubmission.deleteRatio || 0) * 100).toFixed(1)}%`}
                    icon="⌫"
                  />
                  <MetricCard
                    label="Edição Linear"
                    value={`${((selectedSubmission.linearEditingScore || 0) * 100).toFixed(1)}%`}
                    icon="📝"
                    warning={selectedSubmission.linearEditingScore > 0.8}
                  />
                  <MetricCard
                    label="Taxa de Cola"
                    value={`${((selectedSubmission.pasteCharRatio || 0) * 100).toFixed(1)}%`}
                    icon="📋"
                    warning={selectedSubmission.pasteCharRatio > 0.3}
                  />
                  <MetricCard
                    label="Perdas de Foco"
                    value={selectedSubmission.focusLossCount || 0}
                    icon="👁️"
                    warning={selectedSubmission.focusLossCount > 5}
                  />
                  <MetricCard
                    label="Execuções"
                    value={selectedSubmission.executionCount || 0}
                    icon="▶️"
                    warning={selectedSubmission.executionCount === 0}
                  />
                  <MetricCard
                    label="Burstiness"
                    value={selectedSubmission.burstiness?.toFixed(2) || 'N/A'}
                    icon="⚡"
                  />
                </div>

                {selectedSubmission.pasteEventDetails && selectedSubmission.pasteEventDetails.length > 0 && (
                  <div style={{
                    marginTop: '20px',
                    padding: '16px',
                    background: '#fef3c7',
                    borderRadius: '8px'
                  }}>
                    <h4 style={{
                      margin: '0 0 12px 0',
                      fontSize: '14px',
                      color: '#92400e',
                      fontWeight: '700'
                    }}>
                      🚩 Eventos de Cola Detectados
                    </h4>
                    {selectedSubmission.pasteEventDetails.map((paste, idx) => (
                      <div
                        key={idx}
                        style={{
                          marginBottom: '8px',
                          padding: '8px',
                          background: 'white',
                          borderRadius: '6px',
                          fontSize: '12px'
                        }}
                      >
                        <div style={{ color: '#92400e', fontWeight: '600', marginBottom: '4px' }}>
                          {paste.length} caracteres • {paste.linesCount} linhas
                        </div>
                        <pre style={{
                          margin: 0,
                          fontSize: '11px',
                          color: '#666',
                          whiteSpace: 'pre-wrap',
                          wordBreak: 'break-word'
                        }}>
                          {paste.content.substring(0, 100)}...
                        </pre>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div>
                <h3 style={{
                  margin: '0 0 12px 0',
                  fontSize: '18px',
                  color: '#333',
                  fontWeight: '700'
                }}>
                  Código Submetido
                </h3>
                <pre style={{
                  background: '#1e1e1e',
                  color: '#d4d4d4',
                  padding: '16px',
                  borderRadius: '8px',
                  fontSize: '13px',
                  fontFamily: "'Fira Code', monospace",
                  overflow: 'auto',
                  maxHeight: '500px'
                }}>
                  {selectedSubmission.code || 'Código não disponível'}
                </pre>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

const MetricCard = ({ label, value, icon, warning }) => (
  <div style={{
    padding: '12px',
    background: warning ? '#fee2e2' : '#f9fafb',
    borderRadius: '8px',
    border: warning ? '2px solid #fca5a5' : '2px solid #e5e7eb'
  }}>
    <div style={{
      fontSize: '11px',
      color: warning ? '#991b1b' : '#666',
      fontWeight: '600',
      marginBottom: '4px'
    }}>
      {icon} {label}
    </div>
    <div style={{
      fontSize: '18px',
      color: warning ? '#dc2626' : '#333',
      fontWeight: '700'
    }}>
      {value}
    </div>
  </div>
);

export default ActivityDetails;
