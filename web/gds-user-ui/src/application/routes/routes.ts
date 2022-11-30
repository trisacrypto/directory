import { lazy } from 'react';
import CallbackPage from 'modules/auth/callback';
import NotFound from 'modules/notFound';
import Logout from 'modules/auth/logout';
import ResetPassword from 'modules/auth/reset';
import Home from 'modules/home';
import Start from 'modules/start';
import SuccessPage from 'modules/auth/register/success';
import Login from 'modules/auth/login';
import Register from 'modules/auth/register';
import Maintenance from 'components/Maintenance';

import MembershipGuide from 'components/Section/MembershipGuide';
import IntegrateAndComply from 'components/Section/IntegrateAndComply';
// import CertificateManagement from 'components/CertificateManagement';
import VerifyPage from 'modules/verify';
import Collaborators from 'modules/dashboard/collaborator';
import Profile from 'modules/dashboard/profile';
import SwitchOrganization from 'modules/dashboard/organization/SwitchOrganization';
const Overview = lazy(() => import('modules/dashboard/overview'));
const CertificateRegistrationPage = lazy(
  () => import('modules/dashboard/certificate/registration')
);

const CertificateInventory = lazy(() => import('components/CertificateInventory'));

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
    path: '/certificate/registration',
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
  {
    path: '/verify',
    name: 'verify',
    component: VerifyPage,
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
    path: '/dashboard/switch',
    name: 'Switch Organization',
    component: SwitchOrganization,
    layout: 'dashboard',
    route: '/switch'
  },
  {
    path: '/dashboard/certificate/inventory',
    name: 'Certificate Inventory',
    component: CertificateInventory,
    layout: 'dashboard',
    route: '/certificate/inventory'
  },
  {
    path: '/dashboard/logout',
    name: 'Logout',
    component: Logout,
    layout: 'dashboard'
  },
  {
    path: '/dashboard/profile',
    route: '/profile',
    name: 'Profile',
    component: Profile,
    layout: 'dashboard'
  },

  {
    path: '/dashboard/collaborators',
    route: '/collaborators',
    name: 'Collaborators',
    component: Collaborators,
    layout: 'dashboard'
  },
  {
    path: '/dashboard/organization/switch/:id',
    name: 'Switch Organization',
    component: SwitchOrganization,
    layout: 'dashboard',
    route: '/organization/switch/:id'
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
