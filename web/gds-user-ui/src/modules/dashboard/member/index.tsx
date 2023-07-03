import { Suspense } from 'react';
import { Heading } from '@chakra-ui/react';
import MemberTable from './MemberTable';
import { Trans } from '@lingui/macro';
import Loader from 'components/Loader';
const MemberPage: React.FC = () => {
  return (
    <>
      <Heading marginBottom="69px">
        <Trans>TRISA Member Directory</Trans>
      </Heading>

      <Suspense fallback={<Loader />}>
        <MemberTable />
      </Suspense>
    </>
  );
};

export default MemberPage;
