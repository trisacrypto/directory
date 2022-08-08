import { getDefaultValue } from 'components/BasicDetailsForm/validation';
import { getRegistrationDefaultValues } from 'modules/dashboard/certificate/lib';
import { postRegistrationData, getRegistrationData } from 'modules/dashboard/registration/service';
import { handleError } from 'utils/utils';
export const getRegistrationDefaultValue = async () => {
  try {
    const regData = await getRegistrationData();
    if (regData.status === 200 && Object.keys(regData.data).length > 0) {
      return regData.data;
    }
    const defaultValue: any = localStorage.getItem('certificateForm');
    if (defaultValue) {
      const val = JSON.parse(defaultValue);
      const postData = await postRegistrationData(val);
      if (postData.status === 204) {
        const getData = await getRegistrationData();
        localStorage.removeItem('certificateForm');
        return getData.data;
      }
    }
    return getRegistrationDefaultValues();
  } catch (err: any) {
    handleError(err, 'failed to get registration data');
    return getRegistrationDefaultValues();
  }
};

export const postRegistrationValue = (data: any) => {
  return new Promise((resolve, reject) => {
    postRegistrationData(data)
      .then((res) => {
        console.log('[postRegistrationData]', res);
        if (res.status === 204) {
          resolve(res);
        } else {
          reject(res);
        }
      })
      .catch((err) => {
        console.log('[postRegistrationData]', err);
        reject(err);
      });
  });
};

export const setRegistrationDefaultValue = () => {
  console.log('[setRegistrationDefaultValue]');
  const defaultValue: any = getRegistrationDefaultValues();
  return new Promise((resolve, reject) => {
    postRegistrationData(defaultValue)
      .then((res) => {
        console.log('[default postRegistration value]', res);
        if (res.status === 204) {
          resolve(res);
        } else {
          reject(res);
        }
      })
      .catch((err) => {
        console.log('[postRegistrationData]', err);
        reject(err);
      });
  });
};

// get registration data from backend and download the file
export const downloadRegistrationData = async () => {
  try {
    const regData = await getRegistrationData();
    if (regData.status === 200 && Object.keys(regData.data).length > 0) {
      const blob = new Blob([JSON.stringify(regData.data)], { type: 'application/json' });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'registration.json');
      document.body.appendChild(link);
      link.click();
    }
  } catch (err: any) {
    handleError(err, 'failed to get registration data');
  }
};
