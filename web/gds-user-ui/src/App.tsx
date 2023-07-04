import * as Sentry from '@sentry/react';
import { BrowserRouter } from 'react-router-dom';
import AppRouter from 'application/routes';
import { ErrorBoundary } from 'react-error-boundary';
import ErrorFallback from 'components/ErrorFallback';
import { isMaintenanceMode } from './application/config/index';
import Maintenance from 'components/Maintenance';
import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { isProdEnv } from 'application/config';
import { queryClient } from 'application/config/reactQuery';

const App: React.FC = () => {
  return (
    <ErrorBoundary FallbackComponent={ErrorFallback}>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <div className="App">{isMaintenanceMode() ? <Maintenance /> : <AppRouter />}</div>
        </BrowserRouter>

        <ReactQueryDevtools initialIsOpen={!isProdEnv} />
      </QueryClientProvider>
    </ErrorBoundary>
  );
};

export default Sentry.withProfiler(App);
