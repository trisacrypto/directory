import React from 'react';
import CallbackPage from 'modules/auth/callback';
import Overview from 'modules/dashboard/overview';
import NotFound from 'modules/notFound';
import Logout from 'modules/auth/logout';
import ResetPassword from 'modules/auth/reset';
import CertificateRegistrationPage from 'modules/dashboard/certificate/registration';
import Home from 'modules/home';
import Start from 'modules/start';
import SuccessPage from 'modules/auth/register/success';
import Login from 'modules/auth/login';
import Register from 'modules/auth/register';
import Maintenance from 'components/Maintenance';
const CertificateRegistration = React.lazy(
  () => import('modules/dashboard/certificate/registration')
);

import MembershipGuide from 'components/Section/MembershipGuide';
import IntegrateAndComply from 'components/Section/IntegrateAndComply';
import CertificateManagement from '../../components/CertificateManagement';

const appRoutes = [
  // -------LANDING  ROUTES-------
  {
    path: '/',
    name: 'Home',
    component: Home,
    layout: 'landing'
  },
  {
    path: '/',
    name: 'Success Auth',
    component: Home,
    layout: 'landing'
  },
  {
    path: '/getting-started',
    name: 'getting-started',
    component: Start,
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
  {
    path: '/maintenance',
    name: 'Maintenance',
    component: Maintenance,
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

  {
    path: '/auth/success',
    name: 'Success',
    component: SuccessPage,
    layout: 'landing'
  },
  {
    path: '/auth/callback',
    name: 'Callback',
    component: CallbackPage,
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
  {
    path: '/dashboard/logout',
    name: 'Logout',
    component: Logout,
    layout: 'dashboard'
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
