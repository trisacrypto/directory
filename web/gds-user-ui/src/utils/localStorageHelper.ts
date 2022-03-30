import { setCurrentStep } from 'application/store/stepper.slice';
export type TStep = {
  status: string;
  key?: number;
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
        status: 'progress'
      }
    ]
  };

  localStorage.setItem('trs_stepper', JSON.stringify(defaultValue));
  return { ...defaultValue, lastStep: null };
};

export const addStepToLocalStorage = (data: any, key?: string) => {
  const hasLocalStepper: any = localStorage.getItem('trs_stepper');
  if (hasLocalStepper) {
    const getStepper = JSON.parse(hasLocalStepper);
    getStepper.steps.push(data);
    localStorage.setItem('trs_stepper', JSON.stringify(getStepper));
  }
};

export const updateStepFromLocalStorage = (data: any, stepKey: number) => {
  const hasLocalStepper: any = localStorage.getItem('trs_stepper');
  if (hasLocalStepper) {
    const getStepper: any = JSON.parse(hasLocalStepper);
    getStepper.steps.map((step: any) => {
      if (step.key === stepKey - 1) {
        return (getStepper.steps[step.key].status = data.status);
      }
    });
    console.log(getStepper);
  }
};
export const setCurrentStepFromLocalStorage = (currentStep: number) => {
  const hasLocalStepper: any = localStorage.getItem('trs_stepper');
  if (hasLocalStepper) {
    const getStepper: any = JSON.parse(hasLocalStepper);
    getStepper.currentStep = currentStep;

    localStorage.setItem('trs_stepper', JSON.stringify(getStepper));
  }
};

export const setStepFormValue = () => {
  // add each step form value to localstorage
};
