import BasicDetailsReviewDataTable from './BasicDetailsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';

const BasicDetailsReview = () => {
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={1} title="Section 1: Basic Details" />
      <BasicDetailsReviewDataTable />
    </CertificateReviewLayout>
  );
};

export default BasicDetailsReview;
