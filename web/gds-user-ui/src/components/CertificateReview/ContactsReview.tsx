import React, { useEffect } from 'react';
import { colors } from 'utils/theme';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

import ContactsReviewDataTable from './ContactsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
interface ContactsReviewProps {
  data: any;
}
const ContactsReview = ({ data }: ContactsReviewProps) => {
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={3} title={t`Section 3: Contacts`} />
      <ContactsReviewDataTable data={data} />
    </CertificateReviewLayout>
  );
};

export default ContactsReview;
