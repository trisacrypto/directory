import React from 'react';
import { RouteProps } from 'react-router';
import { Route, Navigate } from 'react-router-dom';
//import { useAuthState } from '../../';

export interface PrivateRouteProps extends RouteProps {
    redirectPath: string;
    isAuthenticated: boolean;
}

export const PrivateRoute = ({ redirectPath, isAuthenticated,  ...props }: PrivateRouteProps) => {
  //const { user } = useAuthState()//;
  if (!isAuthenticated) {
    return <Navigate to={redirectPath} />;
  }
  return <Route {...props} />;
};