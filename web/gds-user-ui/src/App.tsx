import * as Sentry from '@sentry/react';
import { BrowserRouter } from 'react-router-dom';
import AppRouter from 'application/routes';
import { LanguageProvider } from 'contexts/LanguageContext';
import { useColorMode } from '@chakra-ui/react';
import { ErrorBoundary } from 'react-error-boundary';
import ErrorFallback from 'components/ErrorFallback';

const App: React.FC = () => {
  const { colorMode } = useColorMode();
  console.log('[]', colorMode);
  return (
    <ErrorBoundary FallbackComponent={ErrorFallback}>
      <LanguageProvider>
        <BrowserRouter>
          <div className="App">
            <AppRouter />
          </div>
        </BrowserRouter>
      </LanguageProvider>
    </ErrorBoundary>
  );
};

export default Sentry.withProfiler(App);
