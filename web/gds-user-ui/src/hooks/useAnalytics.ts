import React, { FC, useEffect } from 'react';
import { BrowserRouter as Router, Link, useLocation } from 'react-router-dom';
import ReactGa from 'react-ga';
import { isProdEnv } from 'application/config';
const useAnalytics = () => {
  const [isInitialized, setIsInitialized] = React.useState<boolean>(false);
  const trackingID: string | undefined = process.env.REACT_APP_ANALYTICS_ID;
  useEffect(() => {
    //  initialize google analytics only in production environment
    if (isProdEnv && trackingID) {
      ReactGa.initialize(trackingID);
    }
    setIsInitialized(true);
  }, [trackingID]);
  return {
    isInitialized
  };
};

export default useAnalytics;
