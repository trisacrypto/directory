import { Heading, Stack } from '@chakra-ui/react';
import Login from 'components/Section/Login';

import LandingLayout from 'layouts/LandingLayout';
import Head from 'components/Head/LandingHead';

const StartPage: React.FC = () => {
  return (
    <LandingLayout>
      <Login />
    </LandingLayout>
  );
};

export default StartPage;
