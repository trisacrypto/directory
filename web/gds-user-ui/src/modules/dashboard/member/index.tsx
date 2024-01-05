import { Heading } from '@chakra-ui/react';
import { Suspense } from 'react';
import MemberTable from './components/MemberTable';
import { Trans } from '@lingui/macro';
import Loader from 'components/Loader';
import DirectoryNotification from './components/DirectoryNotification';
import FormLayout from 'layouts/FormLayout';
import MemberSelectNetwork from './components/MemberNetworkSelect';
import MemberHeader from './components/MemberHeader';
const MemberPage: React.FC = () => {
  return (
    <>
      <Heading marginBottom="69px">
        <Trans>TRISA Member Directory</Trans>
      </Heading>
      <Suspense fallback={<Loader />}>
        <DirectoryNotification />
        <FormLayout overflowX={'auto'}>
          <MemberHeader />
          <MemberSelectNetwork />
          <MemberTable />
        </FormLayout>
      </Suspense>
    </>
  );
};

export default MemberPage;
