import React, { Suspense } from 'react';
import { Routes, Route, Link } from 'react-router-dom';

const Home = React.lazy(() => import('modules/home'));
const StartPage = React.lazy(() => import('modules/start'));
const CertificatePage = React.lazy(() => import('modules/dashboard/certificate/registration'));
const SubmitPage = React.lazy(() => import('modules/dashboard/certificate/submit'));

const AppRouter: React.FC = () => {
  return (
    <Suspense fallback="loading...">
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/getting-started" element={<StartPage />} />
        <Route path="/certificate/registration" element={<CertificatePage />} />
        <Route path="/certificate/submit" element={<SubmitPage />} />
        <Route element={<Home />} />
      </Routes>
    </Suspense>
  );
};

export default AppRouter;
