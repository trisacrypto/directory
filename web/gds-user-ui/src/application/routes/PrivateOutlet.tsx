import React from 'react';
import { Route, Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';

const PrivateOutlet = () => {
  const { isAuthenticated } = useAuth();
  const { pathname } = useLocation();
  return isAuthenticated() ? <Outlet /> : <Navigate to="/" state={{ from: pathname }} replace />;
};
export default PrivateOutlet;
