import { Navigate, useLocation, Outlet } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { userSelector } from 'modules/auth/login/user.slice';
import DashboardLayout from 'layouts/DashboardLayout';
const PrivateOutlet = () => {
  const { isLoggedIn } = useSelector(userSelector);

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
