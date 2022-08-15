import { Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';
import DashboardLayout from 'layouts/DashboardLayout';
const PrivateOutlet = () => {
  const { isAuthenticated } = useAuth();
  const isLoggedIn = isAuthenticated();

  const { pathname } = useLocation();
  return isLoggedIn ? (
    <DashboardLayout>
      <Outlet />
    </DashboardLayout>
  ) : (
    <Navigate to="/" state={{ from: pathname }} replace />
  );
};
export default PrivateOutlet;
