import React, { Suspense } from 'react';
import { Routes, Route, Link } from 'react-router-dom';

const Home = React.lazy(() => import('modules/home'));
const StartPage = React.lazy(() => import('modules/start'));
const CertifacatePage = React.lazy(() => import('modules/dashboard/Certificate'));

const AppRouter: React.FC = () => {
  return (
    <Suspense fallback="loading...">
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/start" element={<StartPage />} />
        <Route path="/certificate" element={<CertifacatePage />} />
      </Routes>
    </Suspense>
  );
};

export default AppRouter;
