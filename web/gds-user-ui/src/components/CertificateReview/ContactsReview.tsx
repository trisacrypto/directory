import { useSelector } from 'react-redux';
import ContactsReviewDataTable from './ContactsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
import { getCurrentState } from 'application/store/selectors/stepper';
import useGetStepStatusByKey from './useGetStepStatusByKey';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
const ContactsReview = () => {
  const currentStateValue = useSelector(getCurrentState);
  const { hasErrorField } = useGetStepStatusByKey(3);

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={3} title={t`Section 3: Contacts`} />
      {hasErrorField ? <RequiredElementMissing elementKey={1} /> : false}
      <ContactsReviewDataTable data={currentStateValue.data.contacts} />
    </CertificateReviewLayout>
  );
};

export default ContactsReview;
