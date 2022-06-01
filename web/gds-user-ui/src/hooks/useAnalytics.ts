import React, { FC, useEffect } from 'react';
import { BrowserRouter as Router, Link, useLocation } from 'react-router-dom';
import ReactGA from 'react-ga4';
import { isProdEnv } from 'application/config';

const useAnalytics = () => {
  const [isInitialized, setIsInitialized] = React.useState<boolean>(false);
  const trackingID: any = process.env.REACT_APP_ANALYTICS_ID;
  useEffect(() => {
    //  initialize google analytics only in production environment
    if (isProdEnv && trackingID) {
      // eslint-disable-next-line no-console
      console.log('initializing google analytics');
      ReactGA.initialize(trackingID, {
        gaOptions: {
          siteSpeedSampleRate: 100
        }
      });
    }
    setIsInitialized(true);
  }, [trackingID]);
  return {
    isInitialized
  };
};

export default useAnalytics;
