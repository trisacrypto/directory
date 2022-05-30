import axiosInstance from 'utils/axios';

export const isProdEnv = process.env.NODE_ENV === 'production';

export const getAppVersionNumber = () => {
  if (isProdEnv) {
    return process.env.REACT_APP_VERSION_NUMBER;
  }
};
export const getAppGitVersion = () => {
  if (isProdEnv) {
    return process.env.REACT_APP_GIT_REVISION;
  }
};

export const getBffAndGdsVersion = async () => {
  if (isProdEnv) {
    try {
      const res = await axiosInstance.get('/status');
      return res.data;
    } catch (e) {
      // log error in sentry or console
      console.error('Error while fetching BFF and GDS version', e);
      return false;
    }
  }
};
