import { HStack, Select, Tooltip } from '@chakra-ui/react';
import { t } from '@lingui/macro';
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
    <HStack>
      <Tooltip label={t`Select language`} hasArrow>
        <Select w="100%" ml={2} value={language as string} onChange={handleLanguageClick}>
          <LanguageOptions />
        </Select>
      </Tooltip>
    </HStack>
  );
};

export default LanguagesDropdown;
