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
interface ProgressBarProps {
  onSetRegistrationState?: any;
  registrationState?: any;
}
const ProgressBar = ({ onSetRegistrationState }: ProgressBarProps) => {
  return (
    <>
      <CertificateSteps>
        <CertificateStepLabel />
        <CertificateStepContainer
          key="1"
          component={<BasicDetails onChangeRegistrationState={onSetRegistrationState} />}
        />
        <CertificateStepContainer key="2" component={<LegalPerson />} />
        <CertificateStepContainer key="3" component={<Contacts />} />
        <CertificateStepContainer key="4" component={<TrisaImplementation />} />
        <CertificateStepContainer key="5" component={<TrixoQuestionnaire />} />
        <CertificateStepContainer key="6" isLast component={<CertificateReview />} />
      </CertificateSteps>
    </>
  );
};

export default ProgressBar;
