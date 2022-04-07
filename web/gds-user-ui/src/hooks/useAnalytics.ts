import React, { FC, useEffect } from 'react';
import { BrowserRouter as Router, Link, useLocation } from 'react-router-dom';
import ReactGa from 'react-ga';

const useAnalytics = () => {
  const [isInitialized, setIsInitialized] = React.useState<boolean>(false);
  const trackingID: any = process.env.REACT_APP_GA_TRACKING_ID;
  useEffect(() => {
    //  initialize google analytics only in production environment
    if (!window.location.href.includes('localhost')) {
      ReactGa.initialize(trackingID);
    }
    setIsInitialized(true);
  }, [trackingID]);
  return {
    isInitialized
  };
};

export default useAnalytics;
