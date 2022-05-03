import React, { Suspense } from 'react';
import { Routes, Route, Link, Navigate } from 'react-router-dom';
import GoogleAnalyticsWrapper from 'components/GaWrapper';
import useAnalytics from 'hooks/useAnalytics';
import NotFound from 'components/NotFound';
const Home = React.lazy(() => import('modules/home'));
const StartPage = React.lazy(() => import('modules/start'));
const CertificatePage = React.lazy(() => import('modules/dashboard/certificate/registration'));
const VerifyPage = React.lazy(() => import('modules/verify'));

const AppRouter: React.FC = () => {
  const { isInitialized } = useAnalytics();
  return (
    <Suspense fallback="loading...">
      <GoogleAnalyticsWrapper isInitialized={isInitialized}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/getting-started" element={<StartPage />} />
          <Route path="/certificate/registration" element={<CertificatePage />} />
          <Route path="/verify" element={<VerifyPage />} />

          <Route element={<Home />} />

          <Route element={<NotFound />} path="/error-404" />
          <Route path="*" element={<Navigate to="/error-404" />} />
        </Routes>
      </GoogleAnalyticsWrapper>
    </Suspense>
  );
};

export default AppRouter;
