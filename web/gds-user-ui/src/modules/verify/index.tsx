import React, { useEffect, useState } from 'react';
import { Heading, Stack, Spinner } from '@chakra-ui/react';
import UserEmailVerification from 'components/Section/UserEmailVerification';
import UserEmailConfirmation from 'components/Section/UserEmailConfirmation';
import LandingLayout from 'layouts/LandingLayout';
import useQuery from 'hooks/useQuery';
import { verifyService } from './verify.service';
import AlertMessage from 'components/ui/AlertMessage';
const VerifyPage: React.FC = () => {
  const query = useQuery();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<any>();
  const [result, setResult] = useState<any>(null);
  const vaspID = query.get('vaspID');
  const token = query.get('token');
  const registered_directory = query.get('registered_directory');
  // to-do : should be improve later
  useEffect(() => {
    (async () => {
      try {
        if (vaspID && token && registered_directory) {
          const params = { vaspID, token, registered_directory };
          const reponse = await verifyService(params);

          if (!reponse.error) {
            setResult(reponse);
          } else {
            console.log('Something went wrong');
            // setError(false)
          }
        } else {
          setError('Invalid params');
        }
      } catch (e: any) {
        if (!e.response?.data?.success) {
          setError(e.response?.data?.error);
        } else {
          // log error
          console.log('sorry something went wrong , please try again');
        }
      } finally {
        setIsLoading(false);
      }
    })();
  }, [vaspID, token, registered_directory]);
  return (
    <LandingLayout>
      {isLoading && <Spinner size={'2xl'} />}
      {result && (
        <AlertMessage message={result.message} status="success" title="Contact Verified " />
      )}
      {error && <AlertMessage message={error} status="error" />}
    </LandingLayout>
  );
};

export default VerifyPage;
