import React from 'react';
import { Route, Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';
const PublicOutlet = () => {
  const { isAuthenticated } = useAuth();
  const isLoggedIn = isAuthenticated();
  const { pathname } = useLocation();
  return isLoggedIn ? (
    <Navigate to="/dashboard/overview" state={{ from: pathname }} replace />
  ) : (
    <Outlet />
  );
};
export default PublicOutlet;
