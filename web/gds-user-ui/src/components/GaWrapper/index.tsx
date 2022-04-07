import React from 'react';
import ReactGa from 'react-ga';
import { useLocation } from 'react-router-dom';

interface IProps {
  children: React.ReactNode;
  isInitialized: boolean;
}

const GoogleAnalyticsWrapper: React.FC<IProps> = ({ children, isInitialized }) => {
  const location = useLocation();

  React.useEffect(() => {
    if (isInitialized) {
      ReactGa.pageview(location.pathname);
    }
  }, [isInitialized, location]);

  return <>{children}</>;
};

export default GoogleAnalyticsWrapper;
