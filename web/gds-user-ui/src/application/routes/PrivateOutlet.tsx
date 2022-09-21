import { Navigate, useLocation, Outlet } from 'react-router-dom';
import useAuth from 'hooks/useAuth';
import DashboardLayout from 'layouts/DashboardLayout';
import { Suspense } from 'react';
import Loader from 'components/Loader';
const PrivateOutlet = () => {
  const { isAuthenticated } = useAuth();
  const isLoggedIn = isAuthenticated();

  const { pathname } = useLocation();

  return isLoggedIn ? (
    <Suspense fallback={<Loader />}>
      <DashboardLayout>
        <Suspense fallback={<Loader />}>
          <Outlet />
        </Suspense>
      </DashboardLayout>
    </Suspense>
  ) : (
    <Navigate to="/" state={{ from: pathname }} replace />
  );
};
export default PrivateOutlet;
