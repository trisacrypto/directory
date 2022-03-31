import {
  CertificateStepContainer,
  CertificateStepLabel,
  CertificateSteps
} from './CertificateStepper';
import BasicDetails from 'components/BasicDetail';
import LegalPerson from 'components/LegalPerson';
import Contacts from 'components/Contacts';
import TrixoQuestionnaire from 'components/TrixoQuestionnaire';
import ReviewSubmit from 'components/ReviewSubmit';
import TrisaImplementation from 'components/TrisaImplementation';

const ProgressBar = () => {
  return (
    <>
      <form>
        <CertificateSteps>
          <CertificateStepLabel />
          <CertificateStepContainer key="1" component={<BasicDetails />} />
          <CertificateStepContainer key="2" component={<LegalPerson />} />
          <CertificateStepContainer key="3" component={<Contacts />} />
          <CertificateStepContainer key="4" component={<TrisaImplementation />} />
          <CertificateStepContainer key="5" component={<TrixoQuestionnaire />} />
          <CertificateStepContainer key="6" isLast component={<ReviewSubmit />} />
        </CertificateSteps>
      </form>
    </>
  );
};

export default ProgressBar;
