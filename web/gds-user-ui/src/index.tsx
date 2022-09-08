import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import { ChakraProvider } from '@chakra-ui/react';
import store, { persistor } from 'application/store';
import { Provider } from 'react-redux';
import initSentry from 'application/config/sentry';
import theme from 'theme';
import { PersistGate } from 'redux-persist/integration/react';
import { LanguageProvider } from 'contexts/LanguageContext';
initSentry();

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <PersistGate loading={null} persistor={persistor}>
        <ChakraProvider theme={theme}>
          <LanguageProvider>
            <App />
          </LanguageProvider>
        </ChakraProvider>
      </PersistGate>
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals

// function sendToAnalytics({ id, name, value }) {
//   ga('send', 'event', {
//     eventCategory: 'Web Vitals',
//     eventAction: name,
//     eventValue: Math.round(name === 'CLS' ? value * 1000 : value), // values must be integers
//     eventLabel: id, // id unique to current page load
//     nonInteraction: true // avoids affecting bounce rate
//   });
// }

reportWebVitals();
