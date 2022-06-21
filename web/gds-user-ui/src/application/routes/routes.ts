import HandleAuthCallback from 'modules/auth/callback';
import Overview from 'modules/dashboard/overview';
import NotFound from 'modules/notFound';
import Logout from 'modules/auth/logout';
import ResetPassword from 'modules/auth/reset';
import CertificateRegistration from 'modules/dashboard/certificate/registration';
import Home from 'modules/home';
import Start from 'modules/start';
import VerifyPage from 'modules/verify';
import SuccessPage from 'modules/auth/register/success';
import Login from 'modules/auth/login';
import Register from 'modules/auth/register';

const dashRoutes = [
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Overview,
    layout: 'dashboard'
  },
  {
    path: '/callback',
    name: 'Callback',
    component: HandleAuthCallback,
    layout: 'landing'
  },
  {
    path: '/logout',
    name: 'Logout',
    component: Logout,
    layout: 'landing'
  },
  {
    path: '/reset',
    name: 'Reset',
    component: ResetPassword,
    layout: 'landing'
  },
  {
    path: '/certificate/registration',
    name: 'Certificate Registration',
    component: CertificateRegistration,
    layout: 'dashboard'
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
    path: '/login',
    name: 'Login',
    component: Login,
    layout: 'landing'
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    layout: 'landing'
  },
  {
    path: '/not-found',
    name: 'Not Found',
    component: NotFound,
    layout: 'landing'
  }
];
export default dashRoutes;
