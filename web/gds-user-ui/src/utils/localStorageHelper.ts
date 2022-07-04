import { setCurrentStep } from 'application/store/stepper.slice';
import { getRegistrationDefaultValue } from 'modules/dashboard/certificate/lib';
import isEmpty from 'lodash/isEmpty';
export type TStep = {
  status: string;
  key?: number;
  data?: any;
};
export type TPayload = {
  currentStep: number | string;
  steps: TStep[];
};

export const loadStepperFromLocalStorage = () => {
  const getLocalStepper: any = localStorage.getItem('trs_stepper');
  if (getLocalStepper) {
    return JSON.parse(getLocalStepper);
  }
  const defaultValue: any = {
    currentStep: 1,
    steps: [
      {
        key: 1,
        status: 'progress',
        data: {}
      }
    ],
    lastStep: null,
    hasReachSubmitStep: false
  };

  localStorage.setItem('trs_stepper', JSON.stringify(defaultValue));
  return { ...defaultValue, lastStep: null };
};

export const addStepToLocalStorage = (data: any, key?: string) => {
  const getStepper: any = JSON.parse(localStorage.getItem('trs_stepper') || '{"steps":[]}');

  getStepper.steps.push(data);
  localStorage.setItem('trs_stepper', JSON.stringify(getStepper));

  return undefined;
};

export const loadDefaultValueFromLocalStorage = () => {
  const defaultValue: any = localStorage.getItem('certificateForm');
  if (defaultValue) {
    return JSON.parse(defaultValue);
  }

  localStorage.setItem('certificateForm', JSON.stringify(getRegistrationDefaultValue()));
  return getRegistrationDefaultValue();
};

export const setCertificateFormValueToLocalStorage = (values: any) => {
  if (values) {
    localStorage.setItem('certificateForm', JSON.stringify(values));
  }
  return undefined;
};

export const setStepperFromLocalStorage = ({ step, status, data }: any) => {
  try {
    const getStepper: any = JSON.parse(localStorage.getItem('trs_stepper') || '{}');
    if (step && !status && !data) {
      getStepper.currentStep = step;
    }
    if (step && status) {
      getStepper.currentStep = step;
      getStepper.steps.map((s: any, index: any) => {
        if (s.key === step) {
          return (getStepper.steps[index].status = status);
        }
      });
    }
    if (step && data) {
      getStepper.currentStep = step;
      getStepper.steps.map((s: any) => {
        if (s.key === step) {
          return (getStepper.steps = data);
        }
      });
    }
    localStorage.setItem('trs_stepper', JSON.stringify(getStepper));
  } catch {
    return undefined;
  }
};

// load user data from localstorage if exist
export const loadUserDataFromLocalStorage = () => {
  try {
    const getUserData: any = localStorage.getItem('userData');
    if (getUserData) {
      return JSON.parse(getUserData);
    }
    return undefined;
  } catch {
    return undefined;
  }
};

export const clearStepperFromLocalStorage = () => {
  localStorage.removeItem('trs_stepper');
  localStorage.removeItem('certificateForm');
  localStorage.removeItem('isTestNetSent');
  localStorage.removeItem('isMainNetSent');
};
