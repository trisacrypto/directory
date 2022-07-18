import React from 'react';
import { Route, Navigate, useLocation, Outlet } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { userSelector } from 'modules/auth/login/user.slice';
const PrivateOutlet = () => {
  const { isLoggedIn } = useSelector(userSelector);
  const { pathname } = useLocation();
  return isLoggedIn ? (
    <Navigate to="/dashboard/overview" state={{ from: pathname }} replace />
  ) : (
    <Outlet />
  );
};
export default PrivateOutlet;
