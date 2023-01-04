import './index.css';
import 'nprogress/nprogress.css';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import App from './App';
import reportWebVitals from './reportWebVitals';
import initSentry from './sentry';

initSentry();

const container = document.getElementById('root');
const root = createRoot(container as HTMLElement);

root.render(
  <StrictMode>
    <App />
  </StrictMode>
);

reportWebVitals();
