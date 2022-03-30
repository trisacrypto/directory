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
          <CertificateStepContainer key="1" status="progress" component={<BasicDetails />} />
          <CertificateStepContainer key="2" status="complete" component={<LegalPerson />} />
          <CertificateStepContainer key="3" status="progress" component={<Contacts />} />
          <CertificateStepContainer key="4" status="progress" component={<TrisaImplementation />} />
          <CertificateStepContainer key="5" status="progress" component={<TrixoQuestionnaire />} />
          <CertificateStepContainer key="6" status="progress" component={<ReviewSubmit />} />
        </CertificateSteps>
      </form>
    </>
  );
};

export default ProgressBar;
