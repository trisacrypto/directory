import React, { useEffect, useState } from 'react';
import { Heading, Stack, Spinner } from '@chakra-ui/react';
import UserEmailVerification from 'components/Section/UserEmailVerification';
import UserEmailConfirmation from 'components/Section/UserEmailConfirmation';
import LandingLayout from 'layouts/LandingLayout';
import AlertMessage from 'components/ui/AlertMessage';
const VerifyPage: React.FC = () => {
  //   const [isLoading, setIsLoading] = useState(true);
  //   const [error, setError] = useState<any>();
  //   const [result, setResult] = useState<any>(null);
  const successRegistrationMessage = `Your account has been created successfully. Please check your email to verify your account.`;

  return (
    <LandingLayout>
      <AlertMessage
        message={successRegistrationMessage}
        status="success"
        title="Success Registration"
      />
    </LandingLayout>
  );
};

export default VerifyPage;
