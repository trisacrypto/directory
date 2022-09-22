import { useSelector } from 'react-redux';
import TrisaImplementationReviewDataTable from './TrisaImplementationReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
import { getCurrentState } from 'application/store/selectors/stepper';

const TrisaImplementationReview = () => {
  const currentStateValue = useSelector(getCurrentState);
  const { data: trisaData } = currentStateValue;

  const trisa = {
    mainnet: trisaData.mainnet,
    testnet: trisaData.testnet
  };

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={4} title={t`Section 4: TRISA Implementation`} />
      <TrisaImplementationReviewDataTable mainnet={trisa.mainnet} testnet={trisa.testnet} />
    </CertificateReviewLayout>
  );
};

export default TrisaImplementationReview;
