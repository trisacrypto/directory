import React, { FC } from 'react';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { setLastStep } from 'application/store/stepper.slice';
interface StepsProps {
  currentStep?: number;
  children: JSX.Element[];
}

const CertificateSteps: FC<StepsProps> = (props: any): any => {
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const lastStep: number = useSelector((state: RootStateOrAny) => state.stepper.lastStep);
  const [toggle, setToggle] = React.useState(false);
  // check if react children is already mounted

  React.useEffect(() => {
    setToggle((prev) => !prev);
  }, [currentStep]);

  const dispatch = useDispatch();
  return React.Children.map(props.children, (child: any) => {
    const isCurrentStep = +child.key === currentStep;
    if (child.props.isLast && !lastStep) {
      dispatch(setLastStep({ lastStep: +child.key }));
    }
    return React.cloneElement(child, {
      ...child.props,
      key: toggle ? +child.key + 1 : +child.key,
      isCurrentStep
    });
  });
};

export default CertificateSteps;
