import { Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';
import DashboardLayout from 'layouts/DashboardLayout';
const PrivateOutlet = () => {
  const { isAuthenticated } = useAuth();
  const isLoggedIn = isAuthenticated();

  const { pathname } = useLocation();
  // get current route from pathname
  const currentRoute = pathname.split('/')[1];
  console.log('currentRoute', currentRoute);
  return isLoggedIn ? (
    <DashboardLayout>
      <Outlet />
    </DashboardLayout>
  ) : (
    <Navigate to="/" state={{ from: pathname }} replace />
  );
};
export default PrivateOutlet;
