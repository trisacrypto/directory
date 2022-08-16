import { getDefaultValue } from 'components/BasicDetailsForm/validation';
import { getRegistrationDefaultValues } from 'modules/dashboard/certificate/lib';
import {
  postRegistrationData,
  getRegistrationData,
  getSubmissionStatus
} from 'modules/dashboard/registration/service';
import { handleError } from 'utils/utils';

export const postRegistrationValue = (data: any) => {
  console.log('[postRegistrationValue]', data);
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
    const v = getRegistrationDefaultValues();
    await postRegistrationValue(v);
    return getRegistrationDefaultValues();
  } catch (err: any) {
    handleError(err, 'failed to get registration data');
  }
};

export const setRegistrationDefaultValue = () => {
  console.log('[setRegistrationDefaultValue]');
  const defaultValue: any = getRegistrationDefaultValues();
  return new Promise((resolve, reject) => {
    postRegistrationData(defaultValue)
      .then((res) => {
        // console.log('[default postRegistration value]', res);
        if (res.status === 204) {
          resolve(res);
        } else {
          reject(res);
        }
      })
      .catch((err) => {
        handleError(err, 'failed to post registration value');
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

// load default stepper

export const getDefaultStepper = async () => {
  try {
    const [regData, regStatus] = await Promise.all([getRegistrationData(), getSubmissionStatus()]);
    console.log('[regStatus]', regStatus.data);
    if (regData.status === 200 && Object.keys(regData.data).length > 0) {
      return {
        currentStep: regData.data.state.current,
        steps: regData.data.state.steps,
        lastStep: null,
        hasReachSubmitStep: regData.data.state.ready_to_submit,
        testnetSubmitted: regStatus?.data?.testnetSubmitted || false,
        mainnetSubmitted: regStatus?.data?.mainnetSubmitted || false
      };
    }
    const defaultValue: any = {
      current: 1,
      steps: [
        {
          key: 1,
          status: 'progress'
        }
      ]
    };
    // update registrations state object
    const postData = await postRegistrationData({ state: defaultValue });
    if (postData.status === 204) {
      const getData = await getRegistrationData();
      return {
        currentStep: getData.data.state.current,
        steps: getData.data.state.steps,
        lastStep: null,
        hasReachSubmitStep: false,
        testnetSubmitted: false,
        mainnetSubmitted: false
      };
    }
  } catch (err: any) {
    handleError(err, 'failed to get stepper data');
  }
};

// // load default stepper without async call
// export const loadDefaultStepperSync = () => {
//   return Promise.resolve({
//     currentStep: 1,
//     steps: [
//       {
//         key: 1,
//         status: 'progress'
//       }
//     ],
//     lastStep: null,
//     hasReachSubmitStep: false
//   });
// };

// set stepper data
export const getRegistrationAndStepperData = async () => {
  try {
    const [regData, regStatus] = await Promise.all([
      getRegistrationDefaultValue(),
      getSubmissionStatus()
    ]);
    if (regData) {
      const response: any = {
        registrationData: regData,
        stepperData: {
          currentStep: regData?.state?.current || 1,
          steps: regData?.state?.steps || [{ key: 1, status: 'progress' }],
          lastStep: null,
          hasReachSubmitStep: regData?.state?.ready_to_submit || false,
          testnetSubmitted: !!regStatus?.data?.testnet_submitted,
          mainnetSubmitted: !!regStatus?.data?.mainnet_submitted
        }
      };
      return response;
    }
  } catch (err: any) {
    handleError(err, 'failed to get stepper data');
  }
};
