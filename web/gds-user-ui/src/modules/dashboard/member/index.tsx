import { Suspense } from 'react';
import { Heading } from '@chakra-ui/react';
import MemberTable from './MemberTable';
import Loader from 'components/Loader';
const MemberPage: React.FC = () => {
  return (
    <>
      <Heading marginBottom="69px">TRISA Member Directory</Heading>

      <Suspense fallback={<Loader />}>
        <MemberTable />
      </Suspense>
    </>
  );
};

export default MemberPage;
