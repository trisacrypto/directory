import React from 'react';
import { Route, Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';
import DashboardLayout from 'layouts/DashboardLayout';
const PrivateOutlet = () => {
  const { isAuthenticated } = useAuth();

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
