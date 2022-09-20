import * as Sentry from '@sentry/react';
import { BrowserRouter } from 'react-router-dom';
import AppRouter from 'application/routes';
import { ErrorBoundary } from 'react-error-boundary';
import ErrorFallback from 'components/ErrorFallback';
import { isMaintenanceMode } from './application/config/index';
import Maintenance from 'components/Maintenance';
const App: React.FC = () => {
  return (
    <ErrorBoundary FallbackComponent={ErrorFallback}>
      <BrowserRouter>
        <div className="App">{isMaintenanceMode() ? <Maintenance /> : <AppRouter />}</div>
      </BrowserRouter>
    </ErrorBoundary>
  );
};

export default Sentry.withProfiler(App);
