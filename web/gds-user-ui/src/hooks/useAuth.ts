import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { userSelector, login, logout, isLoggedInSelector } from 'modules/auth/login/user.slice';
import { useEffect, useState } from 'react';
import { getCookie, setCookie } from 'utils/cookies';
import useCustomAuth0 from './useCustomAuth0';
const useAuth = () => {
  const dispatch = useDispatch();
  const user = useSelector(userSelector);
  const isLoggedIn = useSelector(isLoggedInSelector);
  const { auth0GetUser, auth0CheckSession } = useCustomAuth0();

  const loginUser = (u: TUser) => {
    dispatch(login(u));
  };
  const getToken = getCookie('access_token') || '';

  const logoutUser = () => {
    dispatch(logout());
  };
  const getUser: any = async () => {
    if (getToken) {
      try {
        const userInfo: any = await auth0GetUser(getToken);
        if (userInfo) {
          const u: TUser = {
            isLoggedIn: true,
            user: {
              name: userInfo?.name,
              email: userInfo?.email,
              pictureUrl: userInfo?.picture
            }
          };
          loginUser(u);
        } else {
        }
      } catch (error) {
        console.log('[getUser error]', error);
        throw new Error('401');
      }
    } else {
      logoutUser();
    }
  };
  const isUserAuthenticated = !!isLoggedIn;

  const isAuthenticated = () => {
    console.log('[isLoggedIn]', isLoggedIn);
    if (!isLoggedIn && getToken) {
      getUser();
      return true;
    }
    return isUserAuthenticated;
  };

  return {
    user,
    getUser,
    isLoggedIn,
    loginUser,
    logoutUser,
    isUserAuthenticated,
    isAuthenticated
  };
};

export default useAuth;
