import React from 'react';
import ReactGA from 'react-ga4';
import { useLocation } from 'react-router-dom';

interface IProps {
  children: React.ReactNode;
  isInitialized: boolean;
}

const GoogleAnalyticsWrapper: React.FC<IProps> = ({ children, isInitialized }) => {
  const location = useLocation();

  React.useEffect(() => {
    if (isInitialized) {
      console.log('location path', location.pathname);
      ReactGA.set({ page: location.pathname });
      ReactGA.send(location.pathname + location.search);
    }
  }, [isInitialized, location]);

  return <>{children}</>;
};

export default GoogleAnalyticsWrapper;
