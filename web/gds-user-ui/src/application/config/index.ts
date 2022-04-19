export const getAppVersionNumber = () => {
  console.log('getAppVersionNumber', process.env.NODE_ENV);
  if (process.env.NODE_ENV === 'development') {
    return process.env.REACT_APP_VERSION_NUMBER;
  }
};

export const getAppGitVersion = (application: any) => {};
