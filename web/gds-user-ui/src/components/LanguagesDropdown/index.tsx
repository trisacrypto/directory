import { Select } from '@chakra-ui/react';
import { isDashLocale } from 'application/config';
import { useLanguageProvider } from 'contexts/LanguageContext';

const languages = {
  en: {
    flag: '🇺🇸',
    title: 'English'
  },
  fr: {
    flag: '🇫🇷',
    title: 'Française'
  },
  de: {
    flag: '🇩🇪',
    title: 'Deutsch'
  },
  zh: {
    flag: '🇨🇳',
    title: '中文'
  },
  ja: {
    flag: '🇯🇵',
    title: '日本語'
  }
};

const LanguageOptions = () => {
  return (
    <>
      {Object.entries(languages).map(([k, v]) => (
        <option key={k} value={k}>
          {v.flag} {k.toUpperCase() as string}
        </option>
      ))}
      {isDashLocale() && <option value={'en-dh'}>- DH</option>}
    </>
  );
};

const LanguagesDropdown: React.FC = () => {
  const [language, setLanguage] = useLanguageProvider();

  const handleLanguageClick = (e: any) => {
    localStorage.setItem('gds_lang', e.target.value);
    setLanguage(e.target.value);
  };
  return (
    <Select w="100%" maxW="100" ml={3} value={language as string} onChange={handleLanguageClick}>
      <LanguageOptions />
    </Select>
  );
};

export default LanguagesDropdown;
