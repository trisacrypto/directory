import {
  // CertificateStepContainer,
  CertificateStepLabel
  // CertificateSteps
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

const renderStep = (step: number) => {
  let stepContent = null;
  switch (step) {
    case 1:
      stepContent = <BasicDetails />;
      break;
    case 2:
      stepContent = <LegalPerson />;
      break;
    case 3:
      stepContent = <Contacts />;
      break;
    case 4:
      stepContent = <TrisaImplementation />;
      break;
    case 5:
      stepContent = <TrixoQuestionnaire />;
      break;
    case 6:
      stepContent = <CertificateReview />;
      break;
    default:
      stepContent = <BasicDetails />;
  }
  return stepContent;
};

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
      <div>
        <CertificateStepLabel />
        {renderStep(currentStep)}
      </div>
    </>
  );
};

export default CertificateRegistrationForm;
