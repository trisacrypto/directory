import { FC } from 'react';
import { Collapse } from '@chakra-ui/transition';
interface StepLabelProps {
  key: string;
  component: JSX.Element;
  isCurrentStep?: boolean;
  isLast?: boolean;
}

const CertificateStepContainer: FC<StepLabelProps> = (props) => {
  console.log('[Called] CertificateStepContainer.tsx');
  return (
    <>
      <Collapse in={props.isCurrentStep} unmountOnExit>
        {props.component}
      </Collapse>
    </>
  );
}; // ProgressBar

export default CertificateStepContainer;
