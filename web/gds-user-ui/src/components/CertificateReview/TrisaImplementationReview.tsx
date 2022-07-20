import TrisaImplementationReviewDataTable from './TrisaImplementationReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';

const TrisaImplementationReview = () => (
  <CertificateReviewLayout>
    <CertificateReviewHeader step={4} title="Section 4: TRISA Implementation" />
    <TrisaImplementationReviewDataTable />
  </CertificateReviewLayout>
);

export default TrisaImplementationReview;
