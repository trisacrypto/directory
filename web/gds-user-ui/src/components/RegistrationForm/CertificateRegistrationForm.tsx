import {
  CertificateStepContainer,
  CertificateStepLabel,
  CertificateSteps
} from './CertificateStepper';
import BasicDetails from 'components/BasicDetail';
import LegalPerson from 'components/LegalPerson';
import Contacts from 'components/Contacts';
import TrixoQuestionnaire from 'components/TrixoQuestionnaire';
import TrisaImplementation from 'components/TrisaImplementation';
import CertificateReview from 'components/CertificateReview';
import { useDispatch, useSelector } from 'react-redux';
import { setCurrentStep } from 'application/store/stepper.slice';
import { getCurrentStep } from 'application/store/selectors/stepper';
import { useEffect } from 'react';

const CertificateRegistrationForm = () => {
  const dispatch = useDispatch();
  const currentStep: number = useSelector(getCurrentStep);
  useEffect(() => {
    // if currentStep is 0, then set it to 1
    if (currentStep === 0 || !currentStep) {
      dispatch(setCurrentStep({ currentStep: 1 }));
    }
  }, [currentStep, dispatch]);

  return (
    <>
      <CertificateSteps>
        <CertificateStepLabel />
        <CertificateStepContainer key="1" component={<BasicDetails />} />
        <CertificateStepContainer key="2" component={<LegalPerson />} />
        <CertificateStepContainer key="3" component={<Contacts />} />
        <CertificateStepContainer key="4" component={<TrisaImplementation />} />
        <CertificateStepContainer key="5" component={<TrixoQuestionnaire />} />
        <CertificateStepContainer key="6" isLast component={<CertificateReview />} />
      </CertificateSteps>
    </>
  );
};

export default CertificateRegistrationForm;
