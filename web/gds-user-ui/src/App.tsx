import * as Sentry from '@sentry/react';
import { BrowserRouter } from 'react-router-dom';
import AppRouter from 'application/routes';
import { ErrorBoundary } from 'react-error-boundary';
import ErrorFallback from 'components/ErrorFallback';
import { isMaintenanceMode } from './application/config/index';
import Maintenance from 'components/Maintenance';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { isProdEnv } from 'application/config';
const query = new QueryClient();
const App: React.FC = () => {
  return (
    <ErrorBoundary FallbackComponent={ErrorFallback}>
      <QueryClientProvider client={query}>
        <BrowserRouter>
          <div className="App">{isMaintenanceMode() ? <Maintenance /> : <AppRouter />}</div>
        </BrowserRouter>

        <ReactQueryDevtools initialIsOpen={!isProdEnv} />
      </QueryClientProvider>
    </ErrorBoundary>
  );
};

export default Sentry.withProfiler(App);
