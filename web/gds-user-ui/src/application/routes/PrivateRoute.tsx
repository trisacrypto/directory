import React from 'react';
import { Route, Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';

const PrivateRoute = ({ children }: any) => {
  const { isUserAuthenticated } = useAuth();
  const { pathname } = useLocation();
  return isUserAuthenticated ? (
    children
  ) : (
    <Navigate to="/login" state={{ from: pathname }} replace />
  );
};
export default PrivateRoute;
