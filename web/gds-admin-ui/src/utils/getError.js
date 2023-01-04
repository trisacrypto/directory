import isString from 'lodash/isString';

export default function getErrorMessage(error) {
  if (isString(error)) {
    return error;
  }

  return error?.response && error?.response.data
    ? {
      message: error.response.data.error,
        errorStatus: error.response.status,
      statusText: error.response.statusText,
      }
    : 'Something went wrong';
}
