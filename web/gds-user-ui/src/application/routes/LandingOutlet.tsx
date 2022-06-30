import React from 'react';
import { Route, Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';

const LandingOutlet = () => {
  const { isUserAuthenticated } = useAuth();
  const { pathname } = useLocation();
  return isUserAuthenticated ? (
    <Navigate to="/dashboard/overview" state={{ from: pathname }} replace />
  ) : (
    <Outlet />
  );
};
export default LandingOutlet;
