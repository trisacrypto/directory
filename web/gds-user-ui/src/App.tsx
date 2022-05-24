import LandingLayout from './layouts/LandingLayout';
import * as Sentry from '@sentry/react';
import { BrowserRouter } from 'react-router-dom';
import AppRouter from 'application/routes';
import { I18nProvider } from '@lingui/react';
import { i18n } from '@lingui/core';
import { useEffect } from 'react';
import { DEFAULT_LOCALE, dynamicActivate } from 'utils/i18nLoaderHelper';

const App: React.FC = () => {
  useEffect(() => {
    dynamicActivate(DEFAULT_LOCALE);
  }, []);

  return (
    <I18nProvider i18n={i18n}>
      <BrowserRouter>
        <div className="App">
          <AppRouter />
        </div>
      </BrowserRouter>
    </I18nProvider>
  );
};

export default Sentry.withProfiler(App);
