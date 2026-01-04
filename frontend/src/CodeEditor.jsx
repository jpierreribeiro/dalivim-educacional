import React, { useEffect, useRef, useState } from 'react';
import Editor from '@monaco-editor/react';

const CodeEditor = ({ activityId, studentId, onTelemetryUpdate }) => {
  const editorRef = useRef(null);
  const telemetryRef = useRef({
    keystrokes: [],
    pasteEvents: [],
    focusEvents: [],
    edits: [],
    executions: [],
    sessionStart: Date.now()
  });
  
  const [code, setCode] = useState('');
  const [output, setOutput] = useState('');
  const [isRunning, setIsRunning] = useState(false);
  const [language, setLanguage] = useState('javascript');
  const lastFocusTime = useRef(Date.now());
  const lastKeystrokeTime = useRef(Date.now());

  useEffect(() => {
    // Send telemetry every 10 seconds
    const interval = setInterval(() => {
      sendTelemetry();
    }, 10000);

    // Track window focus/blur
    const handleFocus = () => {
      const now = Date.now();
      const awayTime = now - lastFocusTime.current;
      telemetryRef.current.focusEvents.push({
        type: 'focus',
        timestamp: now,
        awayDuration: awayTime
      });
      lastFocusTime.current = now;
    };

    const handleBlur = () => {
      const now = Date.now();
      telemetryRef.current.focusEvents.push({
        type: 'blur',
        timestamp: now
      });
      lastFocusTime.current = now;
    };

    window.addEventListener('focus', handleFocus);
    window.addEventListener('blur', handleBlur);

    return () => {
      clearInterval(interval);
      window.removeEventListener('focus', handleFocus);
      window.removeEventListener('blur', handleBlur);
      sendTelemetry(true); // Final telemetry
    };
  }, []);

  const handleEditorDidMount = (editor, monaco) => {
    editorRef.current = editor;

    // Track keystrokes
    editor.onKeyDown((e) => {
      const now = Date.now();
      const timeSinceLastKey = now - lastKeystrokeTime.current;
      
      telemetryRef.current.keystrokes.push({
        key: e.browserEvent.key,
        keyCode: e.keyCode,
        timestamp: now,
        dwellTime: 0, // Will be updated on keyup
        flightTime: timeSinceLastKey,
        position: editor.getPosition()
      });
      
      lastKeystrokeTime.current = now;
    });

    editor.onKeyUp((e) => {
      const now = Date.now();
      const lastKeystroke = telemetryRef.current.keystrokes[telemetryRef.current.keystrokes.length - 1];
      if (lastKeystroke && lastKeystroke.keyCode === e.keyCode) {
        lastKeystroke.dwellTime = now - lastKeystroke.timestamp;
      }
    });

    // Track paste events
    editor.onDidPaste((e) => {
      const now = Date.now();
      const pastedText = editor.getModel().getValueInRange(e.range);
      
      telemetryRef.current.pasteEvents.push({
        timestamp: now,
        range: e.range,
        length: pastedText.length,
        content: pastedText.substring(0, 200), // First 200 chars for analysis
        linesCount: pastedText.split('\n').length
      });
    });

    // Track content changes for edit analysis
    editor.onDidChangeModelContent((e) => {
      const now = Date.now();
      
      e.changes.forEach(change => {
        const isDelete = change.text === '';
        const isLinear = isLinearEdit(change, editor);
        
        telemetryRef.current.edits.push({
          timestamp: now,
          range: change.range,
          text: change.text,
          isDelete,
          isLinear,
          rangeLength: change.rangeLength,
          position: editor.getPosition()
        });
      });
    });
  };

  const isLinearEdit = (change, editor) => {
    // Check if edit is at the end of content (linear)
    const model = editor.getModel();
    const totalLines = model.getLineCount();
    const lastLineLength = model.getLineLength(totalLines);
    
    const isAtEnd = 
      change.range.startLineNumber === totalLines &&
      change.range.startColumn >= lastLineLength - 2;
    
    return isAtEnd;
  };

  const calculateTelemetryFeatures = () => {
    const telemetry = telemetryRef.current;
    const totalTime = Date.now() - telemetry.sessionStart;
    
    // Keystroke dynamics
    const flightTimes = telemetry.keystrokes
      .map(k => k.flightTime)
      .filter(t => t > 0 && t < 5000); // Filter outliers
    
    const avgKeystrokeInterval = flightTimes.length > 0
      ? flightTimes.reduce((a, b) => a + b, 0) / flightTimes.length
      : 0;
    
    const stdKeystrokeInterval = flightTimes.length > 1
      ? Math.sqrt(
          flightTimes
            .map(t => Math.pow(t - avgKeystrokeInterval, 2))
            .reduce((a, b) => a + b, 0) / flightTimes.length
        )
      : 0;
    
    // Paste analysis
    const totalPastedChars = telemetry.pasteEvents.reduce((sum, p) => sum + p.length, 0);
    const totalChars = code.length || 1;
    const pasteCharRatio = totalPastedChars / totalChars;
    
    // Delete ratio
    const deleteEvents = telemetry.edits.filter(e => e.isDelete).length;
    const totalEdits = telemetry.edits.length || 1;
    const deleteRatio = deleteEvents / totalEdits;
    
    // Linear editing score
    const linearEdits = telemetry.edits.filter(e => e.isLinear).length;
    const linearEditingScore = linearEdits / totalEdits;
    
    // Focus loss analysis
    const focusLossEvents = telemetry.focusEvents.filter(e => e.type === 'blur');
    const suspiciousFocusLoss = focusLossEvents.filter(e => {
      const nextFocus = telemetry.focusEvents.find(
        f => f.type === 'focus' && f.timestamp > e.timestamp
      );
      return nextFocus && nextFocus.awayDuration > 10000; // Away for more than 10s
    }).length;
    
    // Burstiness (variance in typing speed)
    const burstiness = stdKeystrokeInterval / (avgKeystrokeInterval || 1);
    
    // Time to first run
    const timeToFirstRun = telemetry.executions.length > 0
      ? (telemetry.executions[0].timestamp - telemetry.sessionStart) / 1000
      : totalTime / 1000;
    
    return {
      avgKeystrokeInterval,
      stdKeystrokeInterval,
      pasteEvents: telemetry.pasteEvents.length,
      pasteCharRatio,
      deleteRatio,
      focusLossCount: suspiciousFocusLoss,
      linearEditingScore,
      burstiness,
      timeToFirstRun,
      executionCount: telemetry.executions.length,
      totalTime: totalTime / 1000,
      totalKeystrokes: telemetry.keystrokes.length,
      codeLength: totalChars
    };
  };

  const sendTelemetry = async (isFinal = false) => {
    const features = calculateTelemetryFeatures();
    
    const payload = {
      activityId,
      studentId,
      timestamp: Date.now(),
      isFinal,
      code: isFinal ? code : null, // Only send code on final submission
      features,
      rawEvents: {
        pasteEvents: telemetryRef.current.pasteEvents,
        focusEvents: telemetryRef.current.focusEvents.slice(-20), // Last 20 events
        keystrokeSample: telemetryRef.current.keystrokes.slice(-100) // Last 100 keystrokes
      }
    };

    try {
      const response = await fetch('/api/telemetry', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      
      const result = await response.json();
      
      if (onTelemetryUpdate) {
        onTelemetryUpdate(result);
      }
    } catch (error) {
      console.error('Failed to send telemetry:', error);
    }
  };

  const runCode = async () => {
    setIsRunning(true);
    const now = Date.now();
    
    telemetryRef.current.executions.push({
      timestamp: now,
      codeSnapshot: code
    });

    try {
      const response = await fetch('https://emkc.org/api/v2/piston/execute', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          language: language,
          version: '*',
          files: [{
            name: 'main',
            content: code
          }]
        })
      });

      const result = await response.json();
      
      if (result.run) {
        setOutput(result.run.output || result.run.stderr || 'No output');
      } else {
        setOutput('Execution failed');
      }
    } catch (error) {
      setOutput(`Error: ${error.message}`);
    } finally {
      setIsRunning(false);
    }
  };

  return (
    <div style={{ 
      display: 'flex', 
      flexDirection: 'column', 
      height: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      fontFamily: "'Inter', sans-serif"
    }}>
      <div style={{
        padding: '20px',
        background: 'rgba(255, 255, 255, 0.95)',
        backdropFilter: 'blur(10px)',
        borderBottom: '2px solid #667eea',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        boxShadow: '0 4px 6px rgba(0,0,0,0.1)'
      }}>
        <h2 style={{ margin: 0, color: '#667eea', fontSize: '24px', fontWeight: '700' }}>
          Dalivim Code Editor
        </h2>
        <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
          <select
            value={language}
            onChange={(e) => setLanguage(e.target.value)}
            style={{
              padding: '10px 16px',
              borderRadius: '8px',
              border: '2px solid #667eea',
              fontSize: '14px',
              fontWeight: '500',
              cursor: 'pointer',
              background: 'white'
            }}
          >
            <option value="javascript">JavaScript</option>
            <option value="python">Python</option>
            <option value="java">Java</option>
            <option value="cpp">C++</option>
            <option value="go">Go</option>
            <option value="rust">Rust</option>
          </select>
          <button
            onClick={runCode}
            disabled={isRunning}
            style={{
              padding: '10px 24px',
              background: isRunning ? '#ccc' : 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              fontSize: '14px',
              fontWeight: '600',
              cursor: isRunning ? 'not-allowed' : 'pointer',
              boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
              transition: 'all 0.3s ease'
            }}
          >
            {isRunning ? '⏳ Executando...' : '▶ Executar'}
          </button>
          <button
            onClick={() => sendTelemetry(true)}
            style={{
              padding: '10px 24px',
              background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              fontSize: '14px',
              fontWeight: '600',
              cursor: 'pointer',
              boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
              transition: 'all 0.3s ease'
            }}
          >
            📤 Enviar
          </button>
        </div>
      </div>
      
      <div style={{ flex: 1, display: 'flex', padding: '20px', gap: '20px' }}>
        <div style={{ 
          flex: 1, 
          background: 'white', 
          borderRadius: '12px',
          overflow: 'hidden',
          boxShadow: '0 8px 16px rgba(0,0,0,0.15)'
        }}>
          <Editor
            height="100%"
            language={language}
            value={code}
            onChange={(value) => setCode(value || '')}
            onMount={handleEditorDidMount}
            theme="vs-dark"
            options={{
              minimap: { enabled: true },
              fontSize: 14,
              lineNumbers: 'on',
              roundedSelection: true,
              scrollBeyondLastLine: false,
              automaticLayout: true,
            }}
          />
        </div>
        
        <div style={{ 
          width: '400px',
          background: 'rgba(255, 255, 255, 0.95)',
          borderRadius: '12px',
          padding: '20px',
          boxShadow: '0 8px 16px rgba(0,0,0,0.15)',
          overflowY: 'auto'
        }}>
          <h3 style={{ 
            margin: '0 0 16px 0', 
            color: '#667eea',
            fontSize: '18px',
            fontWeight: '700'
          }}>
            📊 Output
          </h3>
          <pre style={{
            background: '#1e1e1e',
            color: '#d4d4d4',
            padding: '16px',
            borderRadius: '8px',
            fontSize: '13px',
            fontFamily: "'Fira Code', monospace",
            whiteSpace: 'pre-wrap',
            wordWrap: 'break-word',
            minHeight: '200px',
            maxHeight: '400px',
            overflow: 'auto'
          }}>
            {output || 'Execute o código para ver o resultado...'}
          </pre>
        </div>
      </div>
    </div>
  );
};

export default CodeEditor;
