import React, { Suspense } from 'react';
import { Routes, Route, Link, Navigate } from 'react-router-dom';
import PrivateOutlet from 'application/routes/PrivateOutlet';
import LandingOutlet from 'application/routes/LandingOutlet';
import GoogleAnalyticsWrapper from 'components/GaWrapper';
import useAnalytics from 'hooks/useAnalytics';
import VerifyPage from 'modules/verify';
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
        const dashPath =
          prop.path === '/verify' ? `/${prop.path}${prop.route}` : `/${prop.layout}${prop.route}`;

        return <Route key={key} path={dashPath} element={<prop.component />} />;
      } else {
        return null;
      }
    });
  };

  // get current route from pathname
  const currentRoute = window.location.pathname.split('/')[1];
  console.log('currentRoute', currentRoute);

  const { isInitialized } = useAnalytics();
  return (
    <Suspense fallback="Loading">
      <GoogleAnalyticsWrapper isInitialized={isInitialized}>
        <Routes>
          {currentRoute === 'verify' ? (
            <Route path="/verify" element={<VerifyPage />} />
          ) : (
            <>
              <Route path="/" element={<LandingOutlet />}>
                {getLandingRoutes()}
              </Route>
              <Route path="/dashboard" element={<PrivateOutlet />}>
                {getProtectedRoutes()}
              </Route>
            </>
          )}
        </Routes>
      </GoogleAnalyticsWrapper>
    </Suspense>
  );
};

export default AppRouter;
