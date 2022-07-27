import React, { useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import Sidebar from 'components/Sidebar';
import Loader from 'components/Loader';
import { getCookie } from 'utils/cookies';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import { useNavigate } from 'react-router-dom';
import { getAuth0User, userSelector } from 'modules/auth/login/user.slice';
import { getRegistrationData } from 'modules/dashboard/registration/service';
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

  // const [registrationData, setRegistrationData] = useState<any>();
  // const [isLoading, setIsLoading] = useState(true);
  // useEffect(() => {
  //   const fetchData = async () => {
  //     const data = await getRegistrationData();
  //     console.log('[getRegistrationData]', data);
  //     setRegistrationData(data);
  //     setIsLoading(false);
  //   };
  //   fetchData();
  // }, []);

  // return <Loader />;

  return <>{isFetching ? <Loader /> : <Sidebar {...props} />}</>;
};
export default DashboardLayout;
