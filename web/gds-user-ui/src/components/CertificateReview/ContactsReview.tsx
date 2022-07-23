import ContactsReviewDataTable from './ContactsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';

const ContactsReview = () => (
  <CertificateReviewLayout>
    <CertificateReviewHeader step={3} title="Section 3: Contacts" />
    <ContactsReviewDataTable />
  </CertificateReviewLayout>
);

export default ContactsReview;
