import auth0 from 'auth0-js';
import getAuth0Config from 'application/config/auth0';
import jwt_decode from 'jwt-decode';
import { setCookie } from 'utils/cookies';
import { AUTH0_NAMESPACES } from 'utils/constants';
import dayjs from 'dayjs';

// initialize auth0
const auth0Config = getAuth0Config();
const authWeb = new auth0.WebAuth(auth0Config);

export const auth0SignIn = (options: auth0.CrossOriginLoginOptions) => {
  return new Promise((resolve, reject) => {
    authWeb.login(options, (err: any, authResult: any) => {
      if (err) {
        // console.error('error', err);
        reject(err);
      } else {
        // console.log('authResult', authResult);
        resolve(authResult);
      }
    });
  });
};

export const auth0SignOut = (options: any) => {
  authWeb.logout(options);
};
export const auth0SignUp = (options: auth0.DbSignUpOptions) => {
  return new Promise((resolve, reject) => {
    authWeb.signup(options, (err: any, authResult: any) => {
      if (err) {
        reject(err);
      } else {
        resolve(authResult);
      }
    });
  });
};

export const auth0Logout = (options: auth0.LoginOptions) => {
  authWeb.logout({
    ...options,
    returnTo: process.env.REACT_APP_AUTH0_LOGOUT_REDIRECT_URL || 'localhost:3000/auth/logout'
  });
};

export const auth0ResetPassword = (options: auth0.ChangePasswordOptions) => {
  return new Promise((resolve, reject) => {
    authWeb.changePassword(options, (err: any, authResult: any) => {
      if (err) {
        reject(err);
      } else {
        resolve(authResult);
      }
    });
  });
};

export const auth0Hash = (hash?: any) => {
  return new Promise((resolve, reject) => {
    authWeb.parseHash({ hash: hash || window.location.hash }, (err: any, authResult: any) => {
      if (err) {
        reject(err);
      } else {
        const decodeToken: any = jwt_decode(authResult?.accessToken);
        authResult.idTokenPayload.permissions = decodeToken.permissions;

        resolve(authResult);
      }
    });
  });
};

export const auth0CheckSession = () => {
  return new Promise((resolve, reject) => {
    authWeb.checkSession({}, (err: any, authResult: any) => {
      if (err) {
        reject(err);
      } else {
        resolve(authResult);
      }
    });
  });
};
export const auth0GetUser = (accessToken: any) => {
  return new Promise((resolve, reject) => {
    authWeb.client.userInfo(accessToken, (err: any, user: any) => {
      if (err) {
        reject(err);
      } else {
        resolve(user);
      }
    });
  });
};

// check session and get user info
export const refreshAndFetchUser = () => {
  return new Promise((resolve, reject) => {
    authWeb.checkSession({}, async (err: any, authResult: any) => {
      if (err) {
        reject(err);
      } else {
        const user = await auth0GetUser(authResult.accessToken);
        resolve({
          ...authResult,
          user
        });
      }
    });
  });
};

export const refreshAndSetPermission = async () => {
  const user = (await refreshAndFetchUser()) as any;
  if (user) {
    const decodeToken: any = jwt_decode(user.accessToken);
    user.idTokenPayload.permissions = decodeToken.permissions;
    return user;
  }
};

export const auth0SignWithSocial = (connection: string, options?: auth0.AuthorizeOptions) => {
  return authWeb.authorize({
    ...options,
    connection
  });
};

export const getRefreshToken = async (hasRefreshToken: boolean) => {
  if (hasRefreshToken) {
    // refresh token
    const user = (await refreshAndFetchUser()) as any;
    if (user) {
      setCookie('access_token', user?.accessToken);
    }
  }
};

export const refreshNewToken = async () => {
  const user = (await refreshAndFetchUser()) as any;
  if (user) {
    setCookie('access_token', user?.accessToken);
    return !!user?.accessToken;
  }
};

export const setUserPayload = (userTokenPayload: any, data: Partial<IUserState>) => {
  const { email, name, picture, sub, permissions } = userTokenPayload;
  const { vasp, roles } = data;
  return {
    email,
    name,
    pictureUrl: picture,
    id: sub.split('|')[1],
    permissions,
    roles,
    role: userTokenPayload[AUTH0_NAMESPACES.ROLE],
    lastLogin: dayjs(userTokenPayload[AUTH0_NAMESPACES.LAST_LOGIN]).format('MMM D, YYYY HH:mm:ss'),
    createAt: dayjs(userTokenPayload[AUTH0_NAMESPACES.CREATED_AT]).format('MMM D, YYYY HH:mm:ss'),
    vasp,
    authType: sub.split('|')[0]
  };
};
