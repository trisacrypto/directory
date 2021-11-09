import React, { useState, useEffect } from 'react';
import { i18n } from '@lingui/core';
import { I18nProvider } from '@lingui/react';
import { messages as messagesEn } from '../locales/en/messages';
import { messages as messagesDe } from '../locales/de/messages';
import { messages as messagesFr } from '../locales/fr/messages';
import { messages as messagesZh } from '../locales/zh/messages';
import { messages as messagesJa } from '../locales/ja/messages';
import { en, de, fr, zh, ja } from 'make-plural/plurals'


const Context = React.createContext();

export const LanguageStore = ({ children }) => {
  useEffect(() => {
    i18n.load({
      en: messagesEn,
      de: messagesDe,
      fr: messagesFr,
      zh: messagesZh,
      ja: messagesJa,
    });
    i18n.loadLocaleData({
      en: { plurals: en },
      de: { plurals: de },
      fr: { plurals: fr },
      zh: { plurals: zh },
      ja: { plurals: ja}
    })
    i18n.activate('en');
  }, []);

  const [ language, setLanguage ] = useState("en");

  const changeLanguage = (lang) => {
    i18n.activate(lang);
    setLanguage(lang);
  }

  return (
    <Context.Provider value={{ language: language, changeLanguage: changeLanguage}}>
      <I18nProvider i18n={i18n}>
        { children }
      </I18nProvider>
    </Context.Provider>
  );
}

export default Context