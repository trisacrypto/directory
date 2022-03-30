import React, { FC, useEffect } from 'react';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { addStep, setCurrentStep, setStepStatus, TStep } from 'application/store/stepper.slice';
import {
  addStepToLocalStorage,
  updateStepFromLocalStorage,
  setCurrentStepFromLocalStorage
} from 'utils/localStorageHelper';
import { findStepKey } from 'utils/utils';
interface TState {
  status?: boolean;
  isMissed?: boolean;
  step?: number;
}

// 'todo' comment: this hook should be improve
const useCertificateStepper = () => {
  const dispatch = useDispatch();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);

  const nextStep = (state?: TState) => {
    if (state) {
      // user can go to the next with doing anything
      const found = findStepKey(steps, currentStep);
      if (found.length === 1) {
        dispatch(setStepStatus(state));
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
        const foundNext = findStepKey(steps, currentStep + 1);
        if (foundNext.length === 0) {
          dispatch(addStep({ key: currentStep + 1, status: 'progress' }));
        }
      } else {
        dispatch(addStep({ key: currentStep, status: state.status }));
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
      }

      // user can go to the next by updating a status
    } else {
      const found = findStepKey(steps, currentStep + 1);
      console.log('found', found);
      console.log('nextKey', currentStep + 1);
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
      dispatch(setCurrentStep({ currentStep: step - 1 }));
      console.log('c2', currentStep);
      dispatch(setStepStatus({ step, status: 'incomplete' }));
    }
  };

  return {
    nextStep,
    previousStep
  };
};

export default useCertificateStepper;
