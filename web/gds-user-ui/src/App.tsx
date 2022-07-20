import LandingLayout from './layouts/LandingLayout';
import * as Sentry from '@sentry/react';
import { BrowserRouter } from 'react-router-dom';
import AppRouter from 'application/routes';
import { LanguageProvider } from 'contexts/LanguageContext';

const App: React.FC = () => {
  return (
    <LanguageProvider>
      <BrowserRouter>
        <div className="App">
          <AppRouter />
        </div>
      </BrowserRouter>
    </LanguageProvider>
  );
};

export default Sentry.withProfiler(App);
