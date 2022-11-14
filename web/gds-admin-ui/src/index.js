import { StrictMode } from 'react';
import './index.css';
import "nprogress/nprogress.css";
import App from './App';
import reportWebVitals from './reportWebVitals';
import initSentry from 'sentry';
import { createRoot } from 'react-dom/client'


initSentry()

const container = document.getElementById('root');
const root = createRoot(container)

root.render(
  <StrictMode>
    <App />
  </StrictMode>
)

reportWebVitals();
