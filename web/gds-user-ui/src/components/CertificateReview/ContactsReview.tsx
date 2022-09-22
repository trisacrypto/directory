import { useSelector } from 'react-redux';
import ContactsReviewDataTable from './ContactsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
import { getCurrentState } from 'application/store/selectors/stepper';

const ContactsReview = () => {
  const currentStateValue = useSelector(getCurrentState);

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={3} title={t`Section 3: Contacts`} />
      <ContactsReviewDataTable data={currentStateValue.data.contacts} />
    </CertificateReviewLayout>
  );
};

export default ContactsReview;
