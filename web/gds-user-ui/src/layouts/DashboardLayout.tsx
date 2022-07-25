import React, { useEffect } from 'react';
import { userSelector } from 'modules/auth/login/user.slice';
import { useSelector } from 'react-redux';
import Sidebar from 'components/Sidebar';
import Loader from 'components/Loader';
type DashboardLayoutProp = {
  children: React.ReactNode;
};

const DashboardLayout: React.FC<DashboardLayoutProp> = (props) => {
  const { isFetching } = useSelector(userSelector);
  return <>{isFetching ? <Loader /> : <Sidebar {...props} />}</>;
};
export default DashboardLayout;
