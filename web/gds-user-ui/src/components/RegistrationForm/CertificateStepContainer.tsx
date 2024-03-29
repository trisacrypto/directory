import { FC } from 'react';
import { Collapse } from '@chakra-ui/transition';
import { ScrollToTop } from '../ScrollTop/index';
interface StepLabelProps {
  key: string;
  component: JSX.Element;
  isCurrentStep?: boolean;
  isLast?: boolean;
}

const CertificateStepContainer: FC<StepLabelProps> = (props) => {
  return (
    <>
      <ScrollToTop />
      <Collapse in={props.isCurrentStep} unmountOnExit>
        {props.component}
      </Collapse>
    </>
  );
}; // ProgressBar

export default CertificateStepContainer;
