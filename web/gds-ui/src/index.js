import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { i18n } from '@lingui/core';
import { I18nProvider } from '@lingui/react';
import { messages as messagesEn } from './locales/en/messages';
import { messages as messagesDe } from './locales/de/messages';
import { messages as messagesFr } from './locales/fr/messages';
import { messages as messagesJa } from './locales/ja/messages';
import { messages as messagesZh } from './locales/zh/messages';

i18n.load({
  en: messagesEn,
  de: messagesDe,
  fr: messagesFr,
  ja: messagesJa,
  zh: messagesZh,
});
i18n.activate('en');

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