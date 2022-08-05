import { getSteps, getLastStep, resetStepper } from './../application/store/selectors/stepper';
import { getCurrentStep, getHasReachSubmitStep } from 'application/store/selectors/stepper';
import React, { FC, useEffect } from 'react';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import {
  addStep,
  setCurrentStep,
  setStepStatus,
  TStep,
  setStepFormValue,
  setSubmitStep,
  clearStepper,
  setHasReachSubmitStep
} from 'application/store/stepper.slice';
import {
  setStepperFromLocalStorage,
  addStepToLocalStorage,
  clearStepperFromLocalStorage
} from 'utils/localStorageHelper';
import { setRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import { findStepKey } from 'utils/utils';
import { LSTATUS } from 'components/TestnetProgress/CertificateStepLabel';
import { hasStepError } from '../utils/utils';

interface TState {
  status?: boolean;
  isMissed?: boolean;
  step?: number;
  errors?: any;
  isFormCompleted?: boolean;
  formValues?: any;
}

// 'TODO:' this hook should be improve

const useCertificateStepper = () => {
  const dispatch = useDispatch();
  const currentStep: number = useSelector(getCurrentStep);
  const steps: TStep[] = useSelector(getSteps);
  const lastStep: number = useSelector(getLastStep);
  const hasReachSubmitStep: any = useSelector(getHasReachSubmitStep);

  const nextStep = (state?: TState) => {
    // if form value is set then save it to the dedicated step
    if (state?.formValues) {
      dispatch(setStepFormValue({ step: currentStep, formValues: state?.formValues }));
    }
    // only for status update
    if (state?.isFormCompleted || !state?.errors) {
      setStepperFromLocalStorage({ step: currentStep, status: LSTATUS.COMPLETE });
      dispatch(setStepStatus({ status: LSTATUS.COMPLETE, step: currentStep }));
    }
    // if we got an error that means require element are not completed
    if (state?.errors) {
      dispatch(setStepStatus({ status: LSTATUS.ERROR, step: currentStep }));
    }
    // allow manually setting the step status
    if (state?.status) {
      const found = findStepKey(steps, currentStep);
      if (found.length === 1) {
        dispatch(setStepStatus(state));
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
        const foundNext = findStepKey(steps, currentStep + 1);
        if (foundNext.length === 0) {
          if (currentStep === lastStep) {
            return;
          }
          dispatch(addStep({ key: currentStep + 1, status: LSTATUS.PROGRESS }));
        }
      } else {
        if (currentStep === lastStep && state.isFormCompleted) {
          // that mean we move to submit step
          dispatch(setSubmitStep({ submitStep: true }));
        }
        dispatch(addStep({ key: currentStep, status: state.status }));
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
      }
    }
    // if we reach the last step (here review step) , we need to set the submit step
    if (currentStep === lastStep) {
      // that mean we move to submit step
      if (!hasStepError(steps)) {
        setStepperFromLocalStorage({ step: lastStep });
        dispatch(setSubmitStep({ submitStep: true }));
        dispatch(setCurrentStep({ currentStep: lastStep }));
      }
    } else {
      const found = findStepKey(steps, currentStep + 1);

      if (found.length === 0) {
        setStepperFromLocalStorage({ step: currentStep + 1, status: LSTATUS.PROGRESS });
        addStepToLocalStorage({ key: currentStep + 1, status: LSTATUS.PROGRESS });
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
        dispatch(addStep({ key: currentStep + 1, status: LSTATUS.PROGRESS }));
      } else {
        if (found[0].status === LSTATUS.INCOMPLETE) {
          dispatch(setStepStatus({ step: currentStep + 1, status: LSTATUS.PROGRESS }));
        }
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
      }
    }
  };
  const previousStep = (state?: TState) => {
    // if form value is set then save it to the dedicated step
    if (state?.formValues) {
      dispatch(setStepFormValue({ step: currentStep, formValues: state?.formValues }));
    }

    // do not allow to go back for the first step

    const step = currentStep;
    if (currentStep === 1) {
      return;
    }
    dispatch(setCurrentStep({ currentStep: step - 1 }));

    // if the current status is already completed, do not change the status

    const found = findStepKey(steps, currentStep);
    if (found.length > 0 && found[0].status !== LSTATUS.COMPLETE) {
      dispatch(setStepStatus({ step, status: LSTATUS.PROGRESS }));
    }
  };

  const jumpToStep = (step: number) => {
    dispatch(setCurrentStep({ currentStep: step }));
  };

  const jumpToLastStep = () => {
    dispatch(setHasReachSubmitStep({ hasReachSubmitStep: false }));
  };

  const resetForm = () => {
    setRegistrationDefaultValue();
    dispatch(clearStepper());
  };

  return {
    nextStep,
    previousStep,
    jumpToStep,
    resetForm,
    jumpToLastStep
  };
};

export default useCertificateStepper;
