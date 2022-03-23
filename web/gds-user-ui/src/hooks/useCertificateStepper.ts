import React, { FC, useEffect } from 'react';
import { useCertificateSteps } from 'contexts/certificateStepsContext';
import { useDispatch, useSelector } from 'react-redux';
import { setSteps, setCurrentStep } from 'application/store/stepper.slice';

const useCertificateStepper = () => {
  const [certificateSteps] = useCertificateSteps();
  const next = (step: number, status: string) => {
    // set the localstorage data
    // set the current state  to previous + 1
    const steps = {
      key: step,
      status
    };
  };
  const previous = (step: number) => {
    // all set the previous state
  };
  const saveAndNext = (step: number, datas: any) => {
    // save form data and call next function
  };

  return {
    next,
    previous,
    saveAndNext
  };
};

export default useCertificateStepper;
