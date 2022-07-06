import * as Sentry from '@sentry/react'
import { BrowserTracing } from '@sentry/tracing';
import config from './config';

const defaultTracingOrigins = ['localhost', /^\//];

const initSentry = () => {
    console.log(config.appVersion)

    if (process.env.REACT_APP_SENTRY_DSN) {
        let tracingOrigins = defaultTracingOrigins;
        if (process.env.REACT_APP_GDS_API_ENDPOINT) {
            const origin = new URL(process.env.REACT_APP_GDS_API_ENDPOINT);
            tracingOrigins = [origin.host];
        }

        Sentry.init({
            dsn: process.env.REACT_APP_SENTRY_DSN,
            integrations: [new BrowserTracing({ tracingOrigins })],
            environment: process.env.NODE_ENV,
            tracesSampleRate: 1.0,
            release: config.appVersion
        });

        // eslint-disable-next-line no-console
        console.log('Sentry tracing initialized');
    } else {
        // eslint-disable-next-line no-console
        console.log('no Sentry configuration available');
    }
}

export default initSentry;