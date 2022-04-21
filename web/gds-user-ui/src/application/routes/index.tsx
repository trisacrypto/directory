import React, { Suspense } from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import GoogleAnalyticsWrapper from 'components/GaWrapper';
import useAnalytics from 'hooks/useAnalytics';
const Home = React.lazy(() => import('modules/home'));
const StartPage = React.lazy(() => import('modules/start'));
const CertificatePage = React.lazy(() => import('modules/dashboard/certificate/registration'));
const VerifyPage = React.lazy(() => import('modules/verify'));
const LoginPage = React.lazy(() => import('modules/auth/login'));
const RegisterPage = React.lazy(() => import('modules/auth/register'));

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
          <Route path="/auth/login" element={<LoginPage />} />
          <Route path="/auth/register" element={<RegisterPage />} />
          <Route path="*" element={<Home />} />

          <Route element={<Home />} />
        </Routes>
      </GoogleAnalyticsWrapper>
    </Suspense>
  );
};

export default AppRouter;
