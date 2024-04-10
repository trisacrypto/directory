import VaspVerification from 'components/Section/VaspVerification';

import LandingLayout from 'layouts/LandingLayout';
import Head from 'components/Head/LandingHead';
import LandingBanner from 'components/Banner/LandingBanner';

const StartPage: React.FC = () => {
  return (
    <LandingLayout>
      <LandingBanner />
      <Head isStartPage />
      <VaspVerification />
    </LandingLayout>
  );
};

export default StartPage;
