import * as Sentry from '@sentry/react';
import { BrowserTracing } from '@sentry/tracing';

const defaultTracingOrigins = ['localhost', /^\//];

const initSentry = () => {
  if (process.env.REACT_APP_SENTRY_DSN) {
    let tracingOrigins = defaultTracingOrigins;
    if (process.env.REACT_APP_TRISA_BASE_URL) {
      const origin = new URL(process.env.REACT_APP_TRISA_BASE_URL);
      tracingOrigins = [origin.host];
    }

    Sentry.init({
      dsn: process.env.REACT_APP_SENTRY_DSN,
      environment: process.env.NODE_ENV,
      integrations: [
        new BrowserTracing({
          tracingOrigins
        })
      ],

      // Set tracesSampleRate to 1.0 to capture 100%
      // of transactions for performance monitoring.
      // We recommend adjusting this value in production
      tracesSampleRate: 1.0
    });

    console.log("Sentry tracing initialized");
  } else {
    console.log("no Sentry configuration available");
}
};

export default initSentry;
