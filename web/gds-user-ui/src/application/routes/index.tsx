import React, { Suspense } from 'react';
import { Routes, Route, Link, Navigate } from 'react-router-dom';
import lazyLoaderPromise from 'components/CustomLazyLoader';
import GoogleAnalyticsWrapper from 'components/GaWrapper';
import useAnalytics from 'hooks/useAnalytics';
import NotFound from 'modules/notFound';
import Logout from 'modules/auth/logout';
import ResetPassword from 'modules/auth/reset';
// const Home = React.lazy(() =>
//   lazyLoaderPromise(() => import('modules/home'))
// );

import HandleAuthCallback from 'modules/auth/callback';
import Overview from 'modules/dashboard/overview';
import CertificateRegistration from 'modules/dashboard/certificate/registration';
const Home = React.lazy(() => import('modules/home'));
const StartPage = React.lazy(() => import('modules/start'));
const CertificatePage = React.lazy(() => import('modules/dashboard/certificate/registration'));
const OverviewPage = React.lazy(() => import('modules/dashboard/overview'));
const VerifyPage = React.lazy(() => import('modules/verify'));
const SuccessAuth = React.lazy(() => import('modules/auth/register/success'));
const LoginPage = React.lazy(() => import('modules/auth/login'));
const RegisterPage = React.lazy(() => import('modules/auth/register'));
import PrivateRoute from './PrivateRoute';
import MembershipGuide from 'components/Section/MembershipGuide';
import IntegrateAndComply from 'components/Section/IntegrateAndComply';
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
          <Route path="/auth/callback" element={<HandleAuthCallback />} />
          <Route path="/auth/logout" element={<Logout />} />
          <Route path="/account/reset" element={<ResetPassword />} />
          <Route path="/auth/register" element={<RegisterPage />} />
          <Route path="/auth/success" element={<SuccessAuth />} />
          <Route path="/guide" element={<MembershipGuide />} />
          <Route path="/comply" element={<IntegrateAndComply />} />

          <Route path="/dashboard/certificate/registration" element={<CertificateRegistration />} />
          <Route path="/dashboard/overview" element={<Overview />} />

          <Route element={<Home />} />

          <Route element={<NotFound />} path="/404" />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </GoogleAnalyticsWrapper>
    </Suspense>
  );
};

export default AppRouter;
