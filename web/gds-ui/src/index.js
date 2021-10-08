import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { i18n } from '@lingui/core'
import { I18nProvider } from '@lingui/react'
import { messages } from './locales/en/messages'

i18n.load('en', messages)
i18n.activate('en')

const TransApp = () => (
  <I18nProvider i18n={i18n}>
    <App />
  </I18nProvider>
)

ReactDOM.render(
  <React.StrictMode>
    <TransApp />
  </React.StrictMode>,
  document.getElementById('root')
)