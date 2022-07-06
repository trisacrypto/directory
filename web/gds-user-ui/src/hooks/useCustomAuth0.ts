import auth0 from 'auth0-js';
import getAuth0Config from 'application/config/auth0';
import jwt from 'jsonwebtoken';
const useCustomAuth0 = () => {
  // initialize auth0
  const auth0Config = getAuth0Config();
  const authWeb = new auth0.WebAuth(auth0Config);

  const auth0Authorize = (options: any) => {
    authWeb.authorize(options);
  };
  const auth0SignIn = (options: auth0.CrossOriginLoginOptions) => {
    return new Promise((resolve, reject) => {
      authWeb.login(options, (err: any, authResult: any) => {
        if (err) {
          console.error('error', err);
          reject(err);
        } else {
          resolve(authResult);
        }
      });
    });
  };

  const auth0SignOut = (options: any) => {
    authWeb.logout(options);
  };
  const auth0SignUp = (options: auth0.DbSignUpOptions) => {
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

  const auth0Logout = (options: auth0.LoginOptions) => {
    authWeb.logout({
      ...options,
      returnTo: process.env.REACT_APP_AUTH0_LOGOUT_REDIRECT_URL || 'localhost:3000/auth/logout'
    });
  };

  const auth0ResetPassword = (options: auth0.ChangePasswordOptions) => {
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

  const auth0Hash = (hash?: any) => {
    return new Promise((resolve, reject) => {
      authWeb.parseHash({ hash: hash || window.location.hash }, (err: any, authResult: any) => {
        if (err) {
          reject(err);
        } else {
          console.log('[authResult]', authResult);

          resolve(authResult);
        }
      });
    });
  };

  const auth0CheckSession = (options: any) => {
    return new Promise((resolve, reject) => {
      authWeb.checkSession(
        {
          ...options,
          scope: 'read:current_user'
        },
        (err: any, authResult: any) => {
          if (err) {
            reject(err);
          } else {
            resolve(authResult);
          }
        }
      );
    });
  };

  const auth0GetUser = (accessToken: any) => {
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
  const auth0SignWithSocial = (connection: string, options?: auth0.AuthorizeOptions) => {
    return authWeb.authorize({
      ...options,
      connection
    });
  };

  // decode auth access token

  return {
    auth0Authorize,
    auth0SignIn,
    auth0SignOut,
    auth0SignUpWithEmail: auth0SignUp,
    auth0SignWithSocial,
    auth0CheckSession,
    auth0GetUser,
    auth0Logout,
    auth0Hash,
    auth0ResetPassword
  };
};

export default useCustomAuth0;
