import { Suspense } from 'react';
import { Heading } from '@chakra-ui/react';

import Loader from 'components/Loader';
const MemberPage: React.FC = () => {
  return (
    <>
      <Heading marginBottom="69px">TRISA Member Directory</Heading>
      <Suspense fallback={<Loader />}></Suspense>
    </>
  );
};

export default MemberPage;
