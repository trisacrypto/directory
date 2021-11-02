import _ from 'lodash';
import React, { useContext } from 'react';
import NavDropdown from 'react-bootstrap/NavDropdown';
import LanguageContext from "../../contexts/LanguageContext";

const languages = {
  en: {
    flag: "🇺🇸",
    title: "English",
  },
  fr: {
    flag: "🇫🇷",
    title: "Française",
  },
  de: {
    flag: "🇩🇪",
    title: "Deutsch",
  },
  zh: {
    flag: "🇨🇳",
    title: "中文",
  },
}

const LanguageSelect = () => {
  const context = useContext(LanguageContext);

  const selectLanguage = (lang) => (e) => {
    e.preventDefault();
    context.changeLanguage(lang);
    return false;
  }

  const renderItems = () => {
    return _.map(languages, (value, key) => {
      return (
        <NavDropdown.Item
          key={key}
          href="#"
          onClick={selectLanguage(key)}
        >
          <span className="mr-1">{value.flag}</span> {value.title}
        </NavDropdown.Item>
      );
    })
  }

  const currentLanguage = () => {
    return <><span className="mr-1">{languages[context.language].flag}</span> {context.language.toUpperCase()}</>
  }

  return (
    <NavDropdown title={currentLanguage()} id="select-language-dropdown">
      {renderItems()}
    </NavDropdown>
  );
};

export default LanguageSelect;