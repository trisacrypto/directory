import React from 'react';
import ContactsReviewDataTable from './ContactsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
import { StepEnum } from 'types/enums';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';

const ContactsReview = () => {
  const { certificateStep } = useFetchCertificateStep({
    key: StepEnum.CONTACTS
  });

  const hasErrors = certificateStep?.errors;

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={3} title={t`Section 3: Contacts`} />
      {hasErrors ? <RequiredElementMissing elementKey={3} errorFields={hasErrors} /> : false}
      <ContactsReviewDataTable data={certificateStep?.form?.contacts} />
    </CertificateReviewLayout>
  );
};

export default ContactsReview;
