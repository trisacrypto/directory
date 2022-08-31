import React, { useEffect } from 'react';
import { colors } from 'utils/theme';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import { useFormContext } from 'react-hook-form';
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
