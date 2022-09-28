import React, { Suspense } from 'react';
import { Routes, Route } from 'react-router-dom';
import PrivateOutlet from 'application/routes/PrivateOutlet';
import LandingOutlet from 'application/routes/LandingOutlet';
import GoogleAnalyticsWrapper from 'components/GaWrapper';
import useAnalytics from 'hooks/useAnalytics';
import appRoutes from 'application/routes/routes';

const AppRouter: React.FC = () => {
  const getLandingRoutes = () => {
    return appRoutes.map((prop, key) => {
      if (prop.layout === 'landing' || prop.layout === 'dash-landing') {
        return <Route key={key} path={prop.path} element={<prop.component />} />;
      } else {
        return null;
      }
    });
  };
  const getProtectedRoutes = () => {
    return appRoutes.map((prop, key) => {
      if (prop.route && (prop.layout === 'dashboard' || prop.layout === 'dash-landing')) {
        const dashPath = `/${prop.layout}${prop.route}`;
        return <Route key={key} path={dashPath} element={<prop.component />} />;
      } else {
        return null;
      }
    });
  };

  // get current route from pathname

  const { isInitialized } = useAnalytics();

  return (
    <Suspense fallback="Loading">
      <GoogleAnalyticsWrapper isInitialized={isInitialized}>
        <Routes>
          <>
            <Route path="/" element={<LandingOutlet />}>
              {getLandingRoutes()}
            </Route>
            <Route path="/dashboard" element={<PrivateOutlet />}>
              {getProtectedRoutes()}
            </Route>
          </>
        </Routes>
      </GoogleAnalyticsWrapper>
    </Suspense>
  );
};

export default AppRouter;
