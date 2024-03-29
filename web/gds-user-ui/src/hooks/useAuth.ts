import { useDispatch, useSelector } from 'react-redux';
import { userSelector, login, logout, isLoggedInSelector } from 'modules/auth/login/user.slice';
import { getCookie, clearCookies, clearLocalStorage } from 'utils/cookies';
import useCustomAuth0 from './useCustomAuth0';
import { persistor } from 'application/store';
const useAuth = () => {
  const dispatch = useDispatch();
  const user = useSelector(userSelector);
  const isLoggedIn = useSelector(isLoggedInSelector);
  const { auth0GetUser } = useCustomAuth0();

  const loginUser = (u: TUser) => {
    dispatch(login(u));
  };
  const getToken = getCookie('access_token') || '';
  // get expiry time from cookie
  const getExpiryTime: any = getCookie('expires_in') || '';

  const logoutUser = () => {
    clearLocalStorage();
    persistor.purge();
    clearCookies();
    dispatch(logout());
  };
  const getUser: any = async () => {
    if (getToken) {
      try {
        const userInfo: any = await auth0GetUser(getToken);
        if (userInfo) {
          clearCookies();
          clearLocalStorage();
          persistor.purge();
          const u: TUser = {
            isLoggedIn: true,
            user: {
              name: userInfo?.name,
              email: userInfo?.email,
              pictureUrl: userInfo?.picture,
              roles: userInfo?.roles || []
            }
          };
          loginUser(u);
        } else {
        }
      } catch (error) {
        throw new Error('401');
      }
    } else {
      logoutUser();
    }
  };
  const isUserAuthenticated = !!isLoggedIn;

  const isAuthenticated = () => {
    // if token is expired then logout
    if (getExpiryTime && isLoggedIn && getToken) {
      const currentTime = new Date().getTime() / 1000;

      if (currentTime > +getExpiryTime) {
        // console.log('token expired');
        clearLocalStorage();
        clearCookies();
        persistor.purge();
        logoutUser();
        return false;
      } else {
        return isUserAuthenticated;
      }
    }
    return false;
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
