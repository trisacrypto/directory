import React, { FC, useEffect } from 'react';
import { Collapse } from '@chakra-ui/transition';
interface StepLabelProps {
  status: string;
  key: string;
  component: JSX.Element;
  isCurrentStep?: boolean;
}

const CertificateStepContainer: FC<StepLabelProps> = (props) => {
  console.log('stepcontainerprops', props);
  return (
    <>
      <Collapse in={props.isCurrentStep}>{props.component}</Collapse>
    </>
  );
}; // ProgressBar

export default CertificateStepContainer;
