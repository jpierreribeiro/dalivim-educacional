import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom'; // Importe isto
import CodeEditor from './CodeEditor';

// Se quiser importar os outros também para testar rotas futuras:
// import StudentActivity from './StudentActivity';
// import CodeEditor from './CodeEditor';

function App() {
  return (
    // O BrowserRouter deve ser o pai de todos que usam navegação
    <BrowserRouter>
      <Routes>
        {/* Define que o dashboard é a página inicial (/) */}
        <Route path="/" element={<CodeEditor />} />
        
        {/* Se quiser adicionar as outras páginas depois: */}
        {/* <Route path="/student" element={<StudentActivity />} /> */}
      </Routes>
    </BrowserRouter>
  );
}

export default App;