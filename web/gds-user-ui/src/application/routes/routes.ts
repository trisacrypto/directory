import React from 'react';
import HandleAuthCallback from 'modules/auth/callback';
import Overview from 'modules/dashboard/overview';
import NotFound from 'modules/notFound';
import Logout from 'modules/auth/logout';
import ResetPassword from 'modules/auth/reset';
import CertificateRegistrationPage from 'modules/dashboard/certificate/registration';
import { Routes, Route, Link, Outlet } from 'react-router-dom';
import Home from 'modules/home';
import Start from 'modules/start';
import VerifyPage from 'modules/verify';
import SuccessPage from 'modules/auth/register/success';
import Login from 'modules/auth/login';
import Register from 'modules/auth/register';
const CertificateRegistration = React.lazy(
  () => import('modules/dashboard/certificate/registration')
);

// const Home = React.lazy(() => import('modules/home'));
// const StartPage = React.lazy(() => import('modules/start'));
// const CertificatePage = React.lazy(() => import('modules/dashboard/certificate/registration'));
// const VerifyPage = React.lazy(() => import('modules/verify'));
// const SuccessAuth = React.lazy(() => import('modules/auth/register/success'));
// const LoginPage = React.lazy(() => import('modules/auth/login'));
// const RegisterPage = React.lazy(() => import('modules/auth/register'));
// const Home = React.lazy(() => import('modules/home'));
// const StartPage = React.lazy(() => import('modules/start'));
// const CertificatePage = React.lazy(() => import('modules/dashboard/certificate/registration'));
// const VerifyPage = React.lazy(() => import('modules/verify'));
// const SuccessAuth = React.lazy(() => import('modules/auth/register/success'));
// const LoginPage = React.lazy(() => import('modules/auth/login'));
// const RegisterPage = React.lazy(() => import('modules/auth/register'));

import MembershipGuide from 'components/Section/MembershipGuide';
import IntegrateAndComply from 'components/Section/IntegrateAndComply';
import CertificateManagement from "../../components/CertificateManagement";

const appRoutes = [
  // -------LANDING  ROUTES-------
  {
    path: '/',
    name: 'Home',
    component: Home,
    layout: 'landing'
  },
  {
    path: '/start',
    name: 'Start',
    component: Start,
    layout: 'landing'
  },
  {
    path: '/certificate/registration',
    name: 'Certificate Registration',
    component: CertificateRegistration,
    layout: 'landing'
  },
  {
    path: '/verify',
    name: 'Verify',
    component: VerifyPage,
    layout: 'landing'
  },
  {
    path: '/success',
    name: 'Success',
    component: SuccessPage,
    layout: 'landing'
  },
  {
    path: '/auth/callback',
    name: 'Callback',
    component: HandleAuthCallback,
    layout: 'landing'
  },
  {
    path: '/comply',
    name: 'Comply and Integrate',
    component: IntegrateAndComply,
    layout: 'landing'
  },
  {
    path: '/guide',
    name: 'Membership Guide',
    component: MembershipGuide,
    layout: 'landing'
  },

  // -------AUTH ROUTES-------
  {
    path: '/auth/login',
    name: 'Login',
    component: Login,
    layout: 'landing'
  },
  {
    path: '/auth/register',
    name: 'Register',
    component: Register,
    layout: 'landing'
  },
  {
    path: '/auth/logout',
    name: 'Logout',
    component: Logout,
    layout: 'landing'
  },
  {
    path: '/auth/reset',
    name: 'Reset',
    component: ResetPassword,
    layout: 'landing'
  },

  // ------- DASHBOARD ROUTES-------
  {
    path: '/dashboard/overview',
    name: 'Dashboard',
    component: Overview,
    layout: 'dashboard',
    route: '/overview'
  },
  {
    path: '/dashboard/certificate/registration',
    name: 'Certificate Registration',
    component: CertificateRegistrationPage,
    layout: 'dashboard',
    route: '/certificate/registration'
  },
  {
    path: '/dashboard/certificate-management',
    name: 'Certificate Management',
    component: CertificateManagement,
    layout: 'dashboard',
    route: '/dashboard/certificate-management'
  },

  //  -------ERROR ROUTES-------
  {
    path: '/not-found',
    name: 'Not Found',
    component: NotFound,
    layout: 'landing'
  },
  {
    path: '*',
    name: 'Not Found',
    component: NotFound,
    layout: 'dash-landing'
  }
];

export default appRoutes;
