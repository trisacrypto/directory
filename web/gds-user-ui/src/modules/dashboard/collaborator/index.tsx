import { Suspense } from 'react';
import { Heading } from '@chakra-ui/react';
import CollaboratorsSection from 'components/Collaborators';
// import DashboardLayout from 'layouts/DashboardLayout';

import Loader from 'components/Loader';
const Collaborators: React.FC = () => {
  return (
    <>
      <Heading marginBottom="32px">Collaborators</Heading>
      <Suspense fallback={<Loader />}>
        <CollaboratorsSection />
      </Suspense>
    </>
  );
};

export default Collaborators;
