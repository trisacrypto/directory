import React from 'react';
import { Route, Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';
import DashboardLayout from 'layouts/DashboardLayout';
const PrivateOutlet = () => {
  console.log('[PrivateOutlet]');
  const { isAuthenticated } = useAuth();
  console.log('[PrivateOutlet] isAuthenticated: ', isAuthenticated());
  const { pathname } = useLocation();
  return isAuthenticated() ? (
    <DashboardLayout>
      <Outlet />
    </DashboardLayout>
  ) : (
    <Navigate to="/" state={{ from: pathname }} replace />
  );
};
export default PrivateOutlet;
