import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import ProfessorDashboard from './ProfessorDashboard';
import StudentActivity from './StudentActivity';
import ActivityDetails from './ActivityDetails';
import CodeEditor from './CodeEditor';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<ProfessorDashboard />} />
        <Route path="/activity/:inviteToken" element={<StudentActivity />} />
        <Route path="/professor/activity/:activityId" element={<ActivityDetails />} />
        <Route path="/editor" element={<CodeEditor activityId={1} studentId={1} />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
