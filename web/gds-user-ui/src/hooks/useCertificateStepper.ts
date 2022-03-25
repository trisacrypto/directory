import React, { FC, useEffect } from 'react';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { addStep, setCurrentStep, setStepStatus, TStep } from 'application/store/stepper.slice';
import {
  addStepToLocalStorage,
  updateStepFromLocalStorage,
  setCurrentStepFromLocalStorage
} from 'utils/localStorageHelper';
import { findStepKey } from 'utils/utils';

// 'todo' comment: this hook should be improve
const useCertificateStepper = () => {
  const dispatch = useDispatch();
  const CurrentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const Steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const nextStep = (state?: any) => {
    if (state) {
      // user can go to the next with doing anything
      const found = findStepKey(Steps, CurrentStep);
      if (found.length === 1) {
        updateStepFromLocalStorage(state, CurrentStep);
        setCurrentStepFromLocalStorage(CurrentStep + 1);
        dispatch(setStepStatus(state));
        dispatch(setCurrentStep({ currentStep: CurrentStep + 1 }));
        const foundNext = findStepKey(Steps, CurrentStep + 1);
        if (foundNext.length === 0) {
          addStepToLocalStorage({ key: CurrentStep + 1, status: 'progress' });
          dispatch(addStep({ key: CurrentStep + 1, status: 'progress' }));
        }
      } else {
        addStepToLocalStorage({ key: CurrentStep, status: 'progress' });
        setCurrentStepFromLocalStorage(CurrentStep);
        dispatch(addStep({ key: CurrentStep, status: state.status }));
        dispatch(setCurrentStep({ currentStep: CurrentStep + 1 }));
      }

      // user can go to the next by updating a status
    } else {
      setCurrentStepFromLocalStorage(CurrentStep + 1);
      dispatch(setCurrentStep({ currentStep: CurrentStep + 1 }));
    }
  };
  const previousStep = () => {
    // all set the previous state

    dispatch(setCurrentStep({ currentStep: CurrentStep - 1 }));
  };

  return {
    nextStep,
    previousStep
  };
};

export default useCertificateStepper;
