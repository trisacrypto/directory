import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import "nprogress/nprogress.css";
import App from './App';
import { Provider } from 'react-redux';
import { configureStore } from './redux/store';

import reportWebVitals from './reportWebVitals';
import initSentry from 'sentry';


initSentry()

ReactDOM.render(
  <React.StrictMode>
    <Provider store={configureStore({})}>
      <App />
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);

reportWebVitals();
