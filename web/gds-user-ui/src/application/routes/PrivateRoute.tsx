import React from 'react';
import { RouteProps } from 'react-router';
import { Route, Navigate } from 'react-router-dom';
import useAuth from 'hooks/useAuth';

export interface PrivateRouteProps extends RouteProps {
  redirectPath?: string;
  isAuthenticated?: boolean;
}

const PrivateRoute = ({ redirectPath, isAuthenticated, ...props }: PrivateRouteProps) => {
  const { user, isUserAuthenticated } = useAuth();
  if (!isUserAuthenticated) {
    return <Navigate to={'/'} />;
  }
  return <Route {...props} />;
};

export default PrivateRoute;
