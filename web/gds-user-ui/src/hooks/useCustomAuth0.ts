import auth0 from 'auth0-js';
import getAuth0Config from 'application/config/auth0';
const useCustomAuth0 = () => {
  // initialize auth0
  const auth0Config = getAuth0Config();
  const authWeb = new auth0.WebAuth(auth0Config);

  const auth0Authorize = (options: any) => {
    authWeb.authorize(options);
  };
  const auth0SignIn = (options: any, callback: any) => {
    return new Promise((resolve, reject) => {
      authWeb.login(options, (err: any, authResult: any) => {
        if (err) {
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

  const auth0Logout = (options: any) => {
    authWeb.logout(options);
  };

  const auth0CheckSession = (options: any) => {
    return new Promise((resolve, reject) => {
      authWeb.checkSession(options, (err: any, authResult: any) => {
        if (err) {
          reject(err);
        } else {
          resolve(authResult);
        }
      });
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

  return {
    authWeb,
    auth0Authorize,
    auth0SignIn,
    auth0SignOut,
    auth0SignUpWithEmail: auth0SignUp,
    auth0SignWithSocial,
    auth0CheckSession,
    auth0GetUser,
    auth0Logout
  };
};

export default useCustomAuth0;
