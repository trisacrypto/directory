import LegalPersonReviewDataTable from './LegalPersonReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';

const LegalPersonReview = () => (
  <CertificateReviewLayout>
    <CertificateReviewHeader step={2} title="Section 2: Legal Person" />
    <LegalPersonReviewDataTable />
  </CertificateReviewLayout>
);

export default LegalPersonReview;
