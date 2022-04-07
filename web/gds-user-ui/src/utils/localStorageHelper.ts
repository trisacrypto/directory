import { setCurrentStep } from 'application/store/stepper.slice';
import { getRegistrationDefaultValue } from 'modules/dashboard/Certificate/lib';
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
  try {
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
  } catch {
    return undefined;
  }
};

export const addStepToLocalStorage = (data: any, key?: string) => {
  try {
    const getStepper: any = JSON.parse(localStorage.getItem('trs_stepper') || '{}');
    if (getStepper) {
      getStepper.steps.push(data);
      localStorage.setItem('trs_stepper', JSON.stringify(getStepper));
    }
  } catch {
    return undefined;
  }
};

export const updateStepFromLocalStorage = (data: any, stepKey: number) => {
  try {
    const hasLocalStepper: any = localStorage.getItem('trs_stepper');
    if (hasLocalStepper) {
      const getStepper: any = JSON.parse(hasLocalStepper);
      getStepper.steps.map((step: any) => {
        if (step.key === stepKey - 1) {
          return (getStepper.steps[step.key].status = data.status);
        }
      });
    }
  } catch {
    return undefined;
  }
};
export const setCurrentStepFromLocalStorage = (currentStep: number) => {
  try {
    const hasLocalStepper: any = localStorage.getItem('trs_stepper');
    if (hasLocalStepper) {
      const getStepper: any = JSON.parse(hasLocalStepper);
      getStepper.currentStep = currentStep;

      localStorage.setItem('trs_stepper', JSON.stringify(getStepper));
    }
  } catch {
    return undefined;
  }
};

export const setStepFormValueToLocalStorage = (currentStep: number, values: any) => {
  // add each step form value to localstorage
  try {
    const hasLocalStepper: any = localStorage.getItem('trs_stepper');
    if (hasLocalStepper) {
      const getStepper: any = JSON.parse(hasLocalStepper);
      getStepper.steps.map((step: any) => {
        if (step.key === currentStep) {
          return (getStepper.steps[step.key].datas = values);
        }
      });
    }
  } catch {
    return undefined;
  }
};

export const loadDefaultValueFromLocalStorage = () => {
  try {
    const defaultValue: any = localStorage.getItem('certificateForm');
    if (defaultValue) {
      return JSON.parse(defaultValue);
    } else {
      localStorage.setItem('certificateForm', JSON.stringify(getRegistrationDefaultValue()));
      return getRegistrationDefaultValue();
    }
  } catch {
    return undefined;
  }
};

export const setCertificateFormValueToLocalStorage = (values: any) => {
  try {
    localStorage.setItem('certificateForm', JSON.stringify(values));
  } catch {
    return undefined;
  }
};
export const setStepperFromLocalStorage = ({ step, status, data }: any) => {
  try {
    const getStepper: any = JSON.parse(localStorage.getItem('trs_stepper') || '{}');
    if (getStepper) {
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
        getStepper.steps.map((s: any) => {
          if (s.key === step) {
          }
        });
      }
      localStorage.setItem('trs_stepper', JSON.stringify(getStepper));
    }
  } catch {
    return undefined;
  }
};
