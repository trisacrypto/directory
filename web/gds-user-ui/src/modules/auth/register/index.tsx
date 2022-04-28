import { Heading, Stack } from '@chakra-ui/react';
import Register from 'components/Section/CreateAccount';

import LandingLayout from 'layouts/LandingLayout';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import { useForm } from 'react-hook-form';

const StartPage: React.FC = () => {
  const { auth0SignUpWithEmail, auth0SignWithSocial } = useCustomAuth0();
  const handleAuth0Register = (evt: any, type: any) => {
    evt.preventDefault();
    console.log('type', type);
    if (type === 'google') {
      auth0SignWithSocial('google-oauth2');
    }
    if (type === 'email') {
      (async () => {
        const response = await auth0SignUpWithEmail({
          email: 'masskoder+007787@gmail.com',
          password: 'Local123@!',
          connection: 'Username-Password-Authentication'
        });
        console.log('response', response);
      })();
    }
  };
  return (
    <LandingLayout>
      <form>
        <Register handleSignUp={handleAuth0Register} />
      </form>
    </LandingLayout>
  );
};

export default StartPage;
