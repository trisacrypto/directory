import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import Sidebar from 'components/Sidebar';
import Loader from 'components/Loader';
import { getCookie } from 'utils/cookies';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import { useNavigate } from 'react-router-dom';
import { getAuth0User, userSelector } from 'modules/auth/login/user.slice';
type DashboardLayoutProp = {
  children: React.ReactNode;
};

const DashboardLayout: React.FC<DashboardLayoutProp> = (props) => {
  const { isFetching, isLoggedIn } = useSelector(userSelector);
  //   const { auth0SignIn, auth0SignWithSocial, auth0Hash } = useCustomAuth0();
  //   const navigate = useNavigate();
  //   const dispatch = useDispatch();
  // const getToken = getCookie('access_token');
  // useEffect(() => {
  //   console.log('[getToken]', getToken);
  //   console.log('[isLoggedIn]', isLoggedIn);
  //   if (getToken && !isLoggedIn) {
  //     dispatch(getAuth0User(getToken));
  //   }
  // }, [isLoggedIn, getToken]);

  return <>{isFetching ? <Loader /> : <Sidebar {...props} />}</>;
};
export default DashboardLayout;
