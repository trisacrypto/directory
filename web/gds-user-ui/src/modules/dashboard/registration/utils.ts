import { getCookie } from 'utils/cookies';
import { getDefaultValue } from 'components/BasicDetailsForm/validation';
import { getRegistrationDefaultValues } from 'modules/dashboard/certificate/lib';
import { postRegistrationData, getRegistrationData } from 'modules/dashboard/registration/service';
import { handleError } from 'utils/utils';
export const getRegistrationDefaultValue = async () => {
  try {
    const regData = await getRegistrationData();
    console.log('[regData]', regData.data);
    if (regData.status === 200 && Object.keys(regData.data).length > 0) {
      return regData.data;
    }
    const defaultValue: any = localStorage.getItem('certificateForm');
    console.log('defaultCertificateFormValue', defaultValue);
    if (defaultValue) {
      const val = JSON.parse(defaultValue);
      console.log('defaultCertificateFormVal2222', val);
      const postData = await postRegistrationData(val);

      console.log('postData', postData);
      if (postData.status === 204) {
        const getData = await getRegistrationData();
        localStorage.removeItem('certificateForm');
        return getData.data;
      }
    }
    return getDefaultValue();
  } catch (err: any) {
    handleError(err, 'failed to get registration data');
    return getDefaultValue();
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
