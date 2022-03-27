import { Heading, Stack } from '@chakra-ui/react';
import VaspVerification from 'components/Section/VaspVerification';

import LandingLayout from 'layouts/LandingLayout';
import Head from 'components/Head/LandingHead';

const StartPage: React.FC = () => {
  return (
    <LandingLayout>
      <Head
        title="Complete TRISA’s VASP Verfication Process "
        description="All TRISA members must complete TRISA’s VASP verification and due diligence process to become a Verified VASP."
      />
      <VaspVerification />
    </LandingLayout>
  );
};

export default StartPage;
