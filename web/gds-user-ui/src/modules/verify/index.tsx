import React, { useEffect, useState } from 'react';
import { Heading, Stack, Spinner } from '@chakra-ui/react';
import UserEmailVerification from 'components/Section/UserEmailVerification';
import UserEmailConfirmation from 'components/Section/UserEmailConfirmation';
import LandingLayout from 'layouts/LandingLayout';
import useQuery from 'hooks/useQuery';
import { verifyService } from './verifiy.service';
import AlertMessage from 'components/ui/AlertMessage';
const VerifyPage: React.FC = () => {
  const query = useQuery();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<any>();
  const [result, setResult] = useState<any>(null);
  const vaspID = query.get('vaspID');
  const token = query.get('token');
  const registered_directory = query.get('registered_directory');
  // to-do : should be improve later
  useEffect(() => {
    (async () => {
      try {
        setIsLoading(true);
        if (vaspID && token && registered_directory) {
          const params = { vaspID, token, registered_directory };
          const reponse = await verifyService(params);
          console.log('request', reponse);
          if (!reponse.error) {
            setResult(reponse);
          } else {
            // setError(false)
          }
        } else {
          setError('Invalid params');
        }
      } catch (e: any) {
        console.log('error', e.response);
        if (!e.reponse?.data?.success) {
          setError('could not verify contact');
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
      {result && <UserEmailConfirmation message={result.message} />}
      {error && <AlertMessage message={error} status="error" />}
    </LandingLayout>
  );
};

export default VerifyPage;
