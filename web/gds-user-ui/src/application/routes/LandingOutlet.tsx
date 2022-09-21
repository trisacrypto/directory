import React from 'react';
import { Route, Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';
const PublicOutlet = () => {
  const { isAuthenticated } = useAuth();
  const isLoggedIn = isAuthenticated();
  const { pathname } = useLocation();
  return <Outlet />;
};
export default PublicOutlet;
