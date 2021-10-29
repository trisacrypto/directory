import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { LanguageStore } from './contexts/LanguageContext';


ReactDOM.render(
  <React.StrictMode>
    <LanguageStore>
      <App />
    </LanguageStore>
  </React.StrictMode>,
  document.getElementById('root')
)