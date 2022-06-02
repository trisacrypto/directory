import VaspVerification from 'components/Section/VaspVerification';

import LandingLayout from 'layouts/LandingLayout';
import Head from 'components/Head/LandingHead';

const StartPage: React.FC = () => {
  return (
    <LandingLayout>
      <Head isStartPage />
      <VaspVerification />
    </LandingLayout>
  );
};

export default StartPage;
