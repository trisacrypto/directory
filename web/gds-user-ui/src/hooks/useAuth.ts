import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { userSelector, login, logout, isLoggedInSelector } from 'modules/auth/login/user.slice';
import { useEffect, useState } from 'react';
import { getCookie, setCookie } from 'utils/cookies';
import useCustomAuth0 from './useCustomAuth0';
const useAuth = () => {
  const dispatch = useDispatch();
  const user = useSelector(userSelector);
  const isLoggedIn = useSelector(isLoggedInSelector);
  const { auth0GetUser } = useCustomAuth0();

  const loginUser = (u: TUser) => {
    dispatch(login(u));
  };
  const accessToken = getCookie('access_token') || '';

  const logoutUser = () => {
    dispatch(logout());
  };
  const getUser = async () => {
    const getToken = getCookie('access_token');
    if (getToken) {
      try {
        const userInfo: any = await auth0GetUser(getToken);
        const u: TUser = {
          isLoggedIn: true,
          user: {
            name: userInfo?.name,
            email: userInfo?.email,
            pictureUrl: userInfo?.picture
          }
        };
        loginUser(u);
      } catch (error) {
        // log error in sentry
        console.log(error);
        return null;
      }
    } else {
      return accessToken;
    }
  };
  const isUserAuthenticated = !!isLoggedIn;

  return {
    user,
    getUser,
    isLoggedIn,
    loginUser,
    logoutUser,
    isUserAuthenticated
  };
};

export default useAuth;
