import { detect, fromStorage } from '@lingui/detect-locale';
import { I18nProvider } from '@lingui/react';
import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useContext,
  useEffect,
  useState
} from 'react';
import { DEFAULT_LOCALE, dynamicActivate } from 'utils/i18nLoaderHelper';
import { i18n } from '@lingui/core';
import * as yup from 'yup';
import { t } from '@lingui/macro';
import { LANG_KEY } from 'constants/lang-key';

yup.setLocale({
  string: {
    email: {
      default: t`Email is not valid`,
      required: t``
    }
  }
});

type State = [string | null, Dispatch<SetStateAction<string | null>>];

export const LanguageContext = createContext<State | null>(null);

type LanguageProviderProps = {
  children: ReactNode;
};

const DEFAULT_FALLBACK = () => DEFAULT_LOCALE;
const detectedLanguage = detect(fromStorage(LANG_KEY), DEFAULT_FALLBACK);

const LanguageProvider = ({ children }: LanguageProviderProps) => {
  const [language, setLanguage] = useState<string | null>(detectedLanguage);

  useEffect(() => {
    dynamicActivate(language!);
  }, [language]);

  return (
    <I18nProvider i18n={i18n}>
      <LanguageContext.Provider value={[language, setLanguage]}>
        {children}
      </LanguageContext.Provider>
    </I18nProvider>
  );
};

const useLanguageProvider = () => {
  const context = useContext(LanguageContext);

  if (!context) {
    throw new Error(`useLanguageProvider should be used within a LanguageProvider`);
  }

  return context;
};

export { useLanguageProvider, LanguageProvider };
