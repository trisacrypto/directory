import React, { useEffect, useState } from 'react';
import { Heading, Stack, Spinner } from '@chakra-ui/react';
import UserEmailVerification from 'components/Section/UserEmailVerification';
import UserEmailConfirmation from 'components/Section/UserEmailConfirmation';
import LandingLayout from 'layouts/LandingLayout';
import AlertMessage from 'components/ui/AlertMessage';
import { t } from '@lingui/macro';
const VerifyPage: React.FC = () => {
  const successRegistrationMessage = t`Your account has been created successfully. Please check your email to verify your account.`;

  return (
    <LandingLayout>
      <AlertMessage
        message={successRegistrationMessage}
        status="success"
        title={t`Account Registration Complete`}
      />
    </LandingLayout>
  );
};

export default VerifyPage;
