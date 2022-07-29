import * as Sentry from '@sentry/react';
import { BrowserTracing } from '@sentry/tracing';
import { getAppVersionNumber } from ".";

const defaultTracingOrigins = ['localhost', /^\//];

const initSentry = () => {
  if (process.env.REACT_APP_SENTRY_DSN) {
    let tracingOrigins = defaultTracingOrigins;
    if (process.env.REACT_APP_TRISA_BASE_URL) {
      const origin = new URL(process.env.REACT_APP_TRISA_BASE_URL);
      tracingOrigins = [origin.host];
    }

    const environment = process.env.REACT_APP_SENTRY_ENVIRONMENT ? process.env.REACT_APP_SENTRY_ENVIRONMENT : process.env.NODE_ENV;

    Sentry.init({
      dsn: process.env.REACT_APP_SENTRY_DSN,
      environment,
      integrations: [
        new BrowserTracing({
          tracingOrigins
        })
      ],

      // Set tracesSampleRate to 1.0 to capture 100%
      // of transactions for performance monitoring.
      // We recommend adjusting this value in production
      tracesSampleRate: 1.0,
      release: getAppVersionNumber()
    });

    // eslint-disable-next-line no-console
    console.log('Sentry tracing initialized');
  } else {
    // eslint-disable-next-line no-console
    console.log('no Sentry configuration available');
  }
};

export default initSentry;
