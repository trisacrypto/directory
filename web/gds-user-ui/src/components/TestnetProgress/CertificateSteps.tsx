import React, { FC, useEffect } from 'react';
import { CertificateStepsProvider, useCertificateSteps } from 'contexts/certificateStepsContext';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { setLastStep } from 'application/store/stepper.slice';
interface StepsProps {
  currentStep?: number;
  children: JSX.Element[];
}

const CertificateSteps: FC<StepsProps> = (props: any): any => {
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const lastStep: number = useSelector((state: RootStateOrAny) => state.stepper.lastStep);

  const dispatch = useDispatch();
  return React.Children.map(props.children, (child: any, index: any) => {
    const isCurrentStep = +child.key === currentStep;
    if (child.props.isLast && !lastStep) {
      dispatch(setLastStep({ lastStep: +child.key }));
    }
    return React.cloneElement(child, {
      ...child.props,
      isCurrentStep
    });
  });
};

export default CertificateSteps;
