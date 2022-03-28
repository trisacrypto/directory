import useCertificateStepper from 'hooks/useCertificateStepper';
import {
  CertificateStepContainer,
  CertificateStepLabel,
  CertificateSteps
} from './CertificateStepper';
import BasicDetails from 'components/BasicDetail';
import LegalPerson from 'components/LegalPerson';
import { FormProvider, useForm } from 'react-hook-form';
import Contacts from 'components/Contacts';
import TrixoQuestionnaire from 'components/TrixoQuestionnaire';
import ReviewSubmit from 'components/ReviewSubmit';

const ProgressBar = () => {
  const { nextStep, previousStep } = useCertificateStepper();
  const methods = useForm({});

  return (
    <FormProvider {...methods}>
      <form>
        <CertificateSteps>
          <CertificateStepLabel />
          <CertificateStepContainer key="1" status="progress" component={<BasicDetails />} />
          <CertificateStepContainer key="2" status="complete" component={<LegalPerson />} />
          <CertificateStepContainer key="3" status="progress" component={<Contacts />} />
          <CertificateStepContainer key="4" status="progress" component={<TrixoQuestionnaire />} />
          <CertificateStepContainer key="5" status="progress" component={<ReviewSubmit />} />
        </CertificateSteps>
      </form>
    </FormProvider>
    // </CertificateStepsProvider>
  );
};

export default ProgressBar;
