import axiosInstance from 'utils/axios';
export const getAppVersionNumber = () => {
  if (process.env.NODE_ENV === 'production') {
    return process.env.REACT_APP_VERSION_NUMBER;
  }
};
export const getAppGitVersion = () => {
  if (process.env.NODE_ENV === 'production') {
    return process.env.REACT_APP_GIT_REVISION;
  }
};

export const getBffAndGdsVersion = async () => {
  if (process.env.NODE_ENV === 'production') {
    const res = await axiosInstance.get('/status');
    return res.data;
  }
};
