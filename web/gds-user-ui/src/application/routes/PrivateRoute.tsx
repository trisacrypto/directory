import React from 'react';
import { RouteProps } from 'react-router';
import { Route, Navigate } from 'react-router-dom';
import useAuth from 'hooks/useAuth';

export interface PrivateRouteProps extends RouteProps {
  redirectPath: string;
  isAuthenticated: boolean;
}

export const PrivateRoute = ({ redirectPath, isAuthenticated, ...props }: PrivateRouteProps) => {
  const { user } = useAuth();
  if (!user) {
    return <Navigate to={redirectPath} />;
  }
  return <Route {...props} />;
};
