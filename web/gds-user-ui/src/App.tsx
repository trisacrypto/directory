import * as Sentry from '@sentry/react';
import { BrowserRouter } from 'react-router-dom';
import AppRouter from 'application/routes';
import { LanguageProvider } from 'contexts/LanguageContext';
import { useColorMode } from '@chakra-ui/react';

const App: React.FC = () => {
  const { colorMode } = useColorMode();
  console.log('[]', colorMode);
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
