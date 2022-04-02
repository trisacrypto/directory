import React, { FC, useEffect } from 'react';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import {
  addStep,
  setCurrentStep,
  setStepStatus,
  TStep,
  setStepFormValue
} from 'application/store/stepper.slice';

import { findStepKey } from 'utils/utils';
interface TState {
  status?: boolean;
  isMissed?: boolean;
  step?: number;
  errors?: any;
  isFormCompleted?: boolean;
  formValues?: any;
}

// 'todo' this hook should be improve
const useCertificateStepper = () => {
  const dispatch = useDispatch();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const lastStep: number = useSelector((state: RootStateOrAny) => state.stepper.lastStep);

  const nextStep = (state?: TState) => {
    console.log('state is formcompleted', state?.isFormCompleted);
    console.log('state is error', state?.errors);
    console.log('form value from ', state?.formValues);
    if (state?.formValues) {
      dispatch(setStepFormValue({ step: currentStep, formValues: state?.formValues }));
    }
    if (state && state.status) {
      const found = findStepKey(steps, currentStep);
      console.log('found', found);
      console.log('found.length', currentStep);
      if (found.length === 1) {
        dispatch(setStepStatus(state));
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
        const foundNext = findStepKey(steps, currentStep + 1);
        if (foundNext.length === 0) {
          if (currentStep === lastStep) {
            return;
          }
          dispatch(addStep({ key: currentStep + 1, status: 'progress' }));
        }
      } else {
        if (currentStep === lastStep) {
          return;
        }
        dispatch(addStep({ key: currentStep, status: state.status }));
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
      }
    } else {
      if (currentStep === lastStep) {
        return;
      }
      const found = findStepKey(steps, currentStep + 1);

      if (found.length === 0) {
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
        dispatch(addStep({ key: currentStep + 1, status: 'progress' }));
      } else {
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
        dispatch(setStepStatus({ step: currentStep + 1, status: 'progress' }));
      }
    }
  };
  const previousStep = (state?: TState) => {
    // all set the previous state

    if (state) {
    } else {
      const step = currentStep;
      if (currentStep === 1) {
        return;
      }

      dispatch(setCurrentStep({ currentStep: step - 1 }));
      dispatch(setStepStatus({ step, status: 'incomplete' }));
    }
  };

  return {
    nextStep,
    previousStep
  };
};

export default useCertificateStepper;
